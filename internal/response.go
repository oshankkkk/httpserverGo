package internal

import (
	"fmt"
	"io"
	"net"
	"sort"
	"strconv"
	"strings"
)

const defaultResponseBody = "Hello"

type StatusLine struct {
	version string
	code    int
	status  string
}

var header map[string]string

func writeStatusCode() string {
	statusLine := StatusLine{
		version: "HTTP/1.1",
		code:    200,
		status:  "OK",
	}

	return fmt.Sprintf("%s %d %s\r\n", statusLine.version, statusLine.code, statusLine.status)
}

func getDefaultHeaders() {
	header = map[string]string{
		"Connection":     "close",
		"Content-Length": strconv.Itoa(len(defaultResponseBody)),
		"Content-Type":   "text/plain; charset=utf-8",
	}
}

func writeHeaders(headers map[string]string) string {
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for _, key := range keys {
		builder.WriteString(fmt.Sprintf("%s: %s\r\n", key, headers[key]))
	}
	builder.WriteString("\r\n")

	return builder.String()
}

func SendResponse(conn net.Conn) {
	defer conn.Close()

	getDefaultHeaders()
	msg := writeStatusCode() + writeHeaders(header) + defaultResponseBody

	n, err := io.WriteString(conn, msg)
	check(err)
	logger.Info("response sent", "bytes_written", n, "status_code", 200)
}
