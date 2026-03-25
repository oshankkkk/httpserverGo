package main
import(
"strings"
"unicode"
"errors"
"strconv"
)

/*
The normal procedure for parsing an HTTP message is to read the start-line into a structure,
read each header field line into a hash table by field name until the empty line, and then use the parsed data to determine if a message body is expected.
If a message body has been indicated, then it is read as a stream until an amount of octets equal to the message body length is read or the connection is closed.
*/

type StartLine struct {
	method  string
	path    string
	version string
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

func parse(startLine string) {
	headermap := make(map[string]string)
	complist := strings.Split(startLine, " ")
	logger.Info("parsed start line", "method", complist[0], "path", complist[1], "version", complist[2])
	headermap["method"] = complist[0]
	headermap["path"] = complist[1]
	headermap["version"] = complist[2]
}


func bodyParser(buff []byte) string {
	fullbodystring := string(buff)
	return fullbodystring
}

