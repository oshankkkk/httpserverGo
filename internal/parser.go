package internal

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type StartLine struct {
	method  string
	path    string
	version string
}

func HeaderParser(buff []string) (StartLine, error, int) {
	startlinestring := buff[0]
	templist := strings.Fields(startlinestring)
	if len(templist) != 3 {
		return StartLine{}, errors.New("malformed method"), 0
	}
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
	headermap, contentlength, err := HeaderfieldParser(buff[1:])
	if err != nil {
		return StartLine{}, err, 0
	}
	logger.Info("parsed headers", "headers", headermap, "content_length", contentlength)
	return startline, nil, contentlength
}

func fieldNameValidation(fieldName string) error {
	for _, value := range fieldName {
		if !(value >= 'a' && value <= 'z' || value >= 'A' && value <= 'Z' || strings.ContainsRune("!#$%&'*+-.^_`|~", value) || value >= '0' && value <= '9') {
			logger.Warn("invalid header field name", "field_name", fieldName)
			return errors.New("invalid fieldName")
		}
	}

	return nil
}

func HeaderfieldParser(buff []string) (map[string]string, int, error) {
	headermap := make(map[string]string)
	seperator := ":"

	var contentlength int
	for _, value := range buff {
		index := strings.Index(value, seperator)
		if index == -1 {
			return nil, 0, errors.New("malformed method")
		}

		fieldName := value[:index]
		err := fieldNameValidation(fieldName)
		if err != nil {
			return nil, 0, err
		}

		fieldValue := value[index+1:]

		if strings.ToLower(fieldName) == "content-length" {
			contentlength, err = strconv.Atoi(strings.TrimSpace(fieldValue))
			if err != nil {
				return nil, 0, err
			}
		}
		_, ok := headermap[fieldName]
		if !ok {
			headermap[fieldName] = strings.TrimSpace(fieldValue)
		} else {
			headermap[fieldName] += "," + strings.TrimSpace(fieldValue)
		}
	}
	return headermap, contentlength, nil
}

func BodyParser(buff []byte) string {
	fullbodystring := string(buff)
	return fullbodystring
}
