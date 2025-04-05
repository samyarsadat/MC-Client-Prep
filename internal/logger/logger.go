package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func Init(opts *slog.HandlerOptions) {
	Log = slog.New(slog.NewTextHandler(os.Stderr, opts))
}
