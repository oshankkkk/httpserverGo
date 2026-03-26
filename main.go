/*
-- sample req --
*/
package main

import (
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
		//add the server impl for sending errors back
		panic(err)
	}
}

func main() {
	listner, err := net.Listen("tcp", ":8080")
	check(err)
	logger.Info("server started", "addr", ":8080")
	for {
		file, err := listner.Accept()
		check(err)
		ReadConnection(file)
		SendResponse(file)

	}
}



//func writingStatusCode(code int) {
//	switch code {
//	case 200:
//		logger.Info("response status", "status_code", code)
//	case 400:
//		logger.Warn("response status", "status_code", code)
//	case 500:
//		logger.Error("response status", "status_code", code)
//	default:
//		logger.Info("response status", "status_code", code)
//	}
//
//}
//


