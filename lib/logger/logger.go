package logger

import (
	"log/slog"
	"os"
)

func InitLogger(env string) {
	var opts PrettyHandlerOptions

	if env == "local" {
		opts = PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
	} else {
		opts = PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		}
	}
	handler := newPrettyHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

}
