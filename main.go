/*
-- sample req --
POST /coffee HTTP/1.1
Host: localhost:42069
User-Agent: curl/7.81.0
Accept: */ /*
Content-Length: 21

{"flavor":"dark mode"}
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
		readConnection(file)
		sendResponse(file)

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


