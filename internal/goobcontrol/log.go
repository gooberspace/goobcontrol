package goobcontrol

import (
	"flag"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

var (
	logDebug = flag.Bool("debug", false, "Enable debugging")
)

func init() {
	flag.Parse()
}

func CreateLogger(config *viper.Viper) *slog.Logger {
	// All this to set up some logging huh
	logLevel := slog.LevelInfo
	if *logDebug || config.GetString("bot.debug") == "true" {
		logLevel = slog.LevelDebug
	}
	slogOpts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, slogOpts))
	return logger
}
