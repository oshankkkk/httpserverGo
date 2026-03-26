package main
import(
"net"
)
func SendResponse(file net.Conn) {
	msg := "HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\nHello\r\n"
	n, err := file.Write([]byte(msg))
	check(err)
	logger.Info("response sent", "bytes_written", n, "status_code", 200)
}
