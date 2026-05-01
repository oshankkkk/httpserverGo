package internal

import (
	"log/slog"
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
