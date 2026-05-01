package internal

import (
	"bytes"
	"net"
	"os"
	"strings"
)

// formatBytes splits the raw header bytes into individual request lines.
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

func writeToFile(request []string) {
	logfile, err := os.OpenFile("serverlogs.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	check(err)
	defer logfile.Close()

	logger.Info("writing request to log file", "line_count", len(request), "path", "serverlogs.txt")
	for _, line := range request {
		_, err := logfile.WriteString(line + "\n")
		check(err)
	}
}

func ReadConnection(file net.Conn) {
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
		//Do not stop the function that reads the bytes 
		//Instead read until header end, parse the then if the body is there start reading else stop
		if index != -1 {
			stringrequest := formatBytes(buff[:index])
			writeToFile(stringrequest)
			startline, err, a := HeaderParser(stringrequest)
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
			// we are going to check the content legth,
			//if the reminder of the bytes in buf[index:] is smaller than the content length this mean there is more bytes to be read
			// so we go again through the loop and read the byte
			if len(buff[index+4:]) < contentlength {
				logger.Info("waiting for remaining body bytes",
					"content_length", contentlength,
					"bytes_received", len(buff[index+4:]),
				)

				continue
			}
			body := BodyParser(buff[index+4 : index+4+contentlength])
			logger.Info("parsed request body", "body_length", len(body))
			break
		} else {
			break
		}
	}
}

