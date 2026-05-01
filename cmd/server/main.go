
package main

import (
	"httptcp/main/internal"
	"log/slog"
	"net"
	"os"
)

func SendResponse(conn net.Conn) {
	msg := "HTTP/1.1 200 OK\r\nContent-Length: 5\r\n\r\nHello\r\n"
	n, err := conn.Write([]byte(msg))
	check(err)
	logger.Info("response sent", "bytes_written", n, "status_code", 200)
}

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelInfo,
}))

func check(err error) {
	if err != nil {
		logger.Error("unexpected error", "err", err)
		// TODO: send an HTTP error response instead of panicking.
		panic(err)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	check(err)
	logger.Info("server started", "addr", ":8080")
	for {
		conn, err := listener.Accept()
		check(err)
		internal.ReadConnection(conn)
		SendResponse(conn)
	}
}



