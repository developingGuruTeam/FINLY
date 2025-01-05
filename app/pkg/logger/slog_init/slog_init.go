package slog_init

import (
	"cachManagerApp/app/pkg/logger/slogpretty"
	"log/slog"
	"os"
)

func Init() *slog.Logger {
	log := setupPrettySlog()
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
