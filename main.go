package main

import (
	"bytes"
	"errors"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelInfo,
}))

type StartLine struct {
	method  string
	path    string
	version string
}

func check(err error) {
	if err != nil {
		logger.Error("unexpected error", "err", err)
		panic(err)
	}
}

/*
POST /coffee HTTP/1.1
Host: localhost:42069
User-Agent: curl/7.81.0
Accept: */ /*
Content-Length: 21

{"flavor":"dark mode"}
*/

func readConnection(file net.Conn) {
	stream := make([]byte, 1024)
	buff := []byte{}
	var contentlength int
	contentlength = 0
	logger.Info("accepted connection", "remote_addr", file.RemoteAddr().String())
	for {
		count, err := file.Read(stream)
		check(err)
		buff = append(buff, stream[:count]...)
		index := bytes.Index(buff, []byte("\r\n\r\n"))
		//Do not stop the function that reads the bytes instead read until header end, parse the then if the body is there start reading else stop
		if index != -1 {
			stringrequest := formatBytes(buff[:index])
			writeTofile(stringrequest)
			startline, err, a := headerParser(stringrequest)
			check(err)
			contentlength = a
			logger.Info("parsed request",
				"method", startline.method,
				"path", startline.path,
				"version", startline.version,
				"content_length", contentlength,
			)

		}
		//parse the body if it is here by reading the test of the remainig bytes from the connection
		//body is basically buff[index:]
		if contentlength != 0 {
			// we are going to check the content legth, if the reminder of the bytes in buf[index:] is smaller than the content length this mean there is more bytes to be read
			// so we go again through the loop and read the byte
			if len(buff[index+4:]) < contentlength {
				logger.Info("waiting for remaining body bytes",
					"content_length", contentlength,
					"bytes_received", len(buff[index+4:]),
				)
				continue
			}
			body := bodyParser(buff[index+4:])
			logger.Info("parsed request body", "body_length", len(body))
			break
		} else {
			break
		}
	}
}
func bodyParser(buff []byte) string {
	fullbodystring := string(buff)
	return fullbodystring
}
func formatBytes(header []byte) []string {
	buff := []string{}
	var temp string
	for _, value := range string(header) {
		if string(value) == "\n" {
			temp = strings.TrimSpace(temp)
			buff = append(buff, temp)
			temp = ""
		} else {
			temp += string(value)
		}

	}
	temp = strings.TrimSpace(temp)
	if temp != "" {
		buff = append(buff, temp)
	}
	return buff
}

func headerParser(buff []string) (StartLine, error, int) {
	startlinestring := buff[0]
	templist := strings.Split(startlinestring, " ")
	var startline StartLine
	startline.method = templist[0]
	startline.path = templist[1]
	startline.version = templist[2]
	for i := 0; i < len(startline.method); i++ {
		if !unicode.IsUpper(rune(startline.method[i])) {
			return StartLine{}, errors.New("incorrect http method"), 0
		}
	}
	if startline.version != "HTTP/1.1" {
		return StartLine{}, errors.New("wrong http version"), 0
	}
	// buf[1:] are headers
	headermap, contentlength := headerfieldParser(buff[1:])

	logger.Info("parsed headers", "headers", headermap, "content_length", contentlength)
	return startline, nil, contentlength
}

func fieldNameValidation(fieldName string) error {
	for _, value := range fieldName {
		if !(value >= 'a' && value <= 'z' || value >= 'A' && value <= 'Z' || strings.ContainsRune("!#$%&'*+-.^_`|~", value) || value >= 0 && value <= 9) {
			logger.Warn("invalid header field name", "field_name", fieldName)
			return errors.New("invalid fieldName")
		}
	}

	return nil
}
func headerfieldParser(buff []string) (map[string]string, int) {
	headermap := make(map[string]string)
	seperator := ":"
	var contentlength int
	for _, value := range buff {
		value = string(value)
		index := strings.Index(value, seperator)
		fieldName := value[:index]
		err := fieldNameValidation(fieldName)
		check(err)

		fieldValue := value[index+1:]

		if strings.ToLower(fieldName) == "content-length" {

			contentlength, err = strconv.Atoi(strings.TrimSpace(fieldValue))
			check(err)

		}
		_, ok := headermap[fieldName]
		if ok == false {
			headermap[fieldName] = strings.TrimSpace(fieldValue)
		} else {

			headermap[fieldName] += "," + strings.TrimSpace(fieldValue)
		}
	}
	return headermap, contentlength

}

/*
The normal procedure for parsing an HTTP message is to read the start-line into a structure,
read each header field line into a hash table by field name until the empty line, and then use the parsed data to determine if a message body is expected.
If a message body has been indicated, then it is read as a stream until an amount of octets equal to the message body length is read or the connection is closed.
*/

func writeTofile(request []string) {
	logfile, err := os.OpenFile("serverlogs.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	check(err)
	defer logfile.Close()

	logger.Info("writing request to log file", "line_count", len(request), "path", "serverlogs.txt")
	for _, line := range request {
		_, err := logfile.WriteString(line + "\n")
		check(err)
	}
}

func parse(startLine string) {
	headermap := make(map[string]string)
	complist := strings.Split(startLine, " ")
	logger.Info("parsed start line", "method", complist[0], "path", complist[1], "version", complist[2])
	headermap["method"] = complist[0]
	headermap["path"] = complist[1]
	headermap["version"] = complist[2]
}
func writingStatusCode(code int) {
	switch code {
	case 200:
		logger.Info("response status", "status_code", code)
	case 400:
		logger.Warn("response status", "status_code", code)
	case 500:
		logger.Error("response status", "status_code", code)
	default:
		logger.Info("response status", "status_code", code)
	}

}

func sendResponse(file net.Conn) {
	msg := "HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\nHello\r\n"
	n, err := file.Write([]byte(msg))
	check(err)
	logger.Info("response sent", "bytes_written", n, "status_code", 200)
}

func main() {
	listner, err := net.Listen("tcp", ":8080")
	check(err)
	logger.Info("server started", "addr", ":8080")
	for {
		file, err := listner.Accept()
		check(err)
		readConnection(file)
		sendResponse(file)

	}
}
