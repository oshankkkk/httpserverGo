package main

import (
	"httptcp/main/internal"
	"log/slog"
	"net"
	"os"
)



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
		internal.SendResponse(conn)
	}
}



