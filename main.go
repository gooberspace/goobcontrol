package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/gooberspace/goobcontrol/internal/commands"
	"github.com/gooberspace/goobcontrol/internal/eventhandlers"
	_ "github.com/joho/godotenv/autoload"
)

var (
	logDebug = flag.Bool("debug", false, "Enable debugging")
)

func init() {
	flag.Parse()
}

func main() {
	// All this to set up some logging huh
	logLevel := slog.LevelInfo
	if *logDebug {
		logLevel = slog.LevelDebug
	}
	slogOpts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, slogOpts))
	slog.SetDefault(logger)

	// Setting up a basic disgo client with some sane defaults
	// We're doing all event handling elsewhere so this file can stay small
	slog.Info("Starting the bot named " + os.Getenv("BOT_NAME"))
	client, err := disgo.New(os.Getenv("BOT_DISCORD_TOKEN"),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentDirectMessages,
			),
		),
		bot.WithEventListenerFunc(eventhandlers.HandleDiscordEvent),
	)
	if err != nil {
		panic(err)
	}
	if err = client.OpenGateway(context.TODO()); err != nil {
		panic(err)
	}

	commands.RegisterCommands(client)
	privateGuilds := strings.Split(os.Getenv("BOT_PRIVATE_GUILDS"), ",")
	commands.RegisterGuildCommands(client, privateGuilds)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s

	slog.Info("Attempting graceful shutdown")
	client.Close(context.TODO())
}
