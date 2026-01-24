package main

import (
	"context"
	"errors"
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
	"github.com/spf13/viper"
)

var (
	logDebug = flag.Bool("debug", false, "Enable debugging")
)

func init() {
	flag.Parse()
}

func main() {
	viper.SetEnvPrefix("goob")
	viper.SetTypeByDefaultValue(true)
	viper.SetConfigName("goobconfig")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotfoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotfoundError) {
			//Ignore these errors as we might be reading from env
		} else {
			panic("Couldn't load config file: " + err.Error())
		}
	}
	viper.SetDefault("bot.name", "Goob Control")
	viper.SetDefault("bot.debug", false)
	viper.SetDefault("discord.privateGuilds", []string{})
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// All this to set up some logging huh
	logLevel := slog.LevelInfo
	if *logDebug || viper.GetString("bot.debug") == "true" {
		logLevel = slog.LevelDebug
	}
	slogOpts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, slogOpts))
	slog.SetDefault(logger)

	// Setting up a basic disgo client with some sane defaults
	// We're doing all event handling elsewhere so this file can stay small
	slog.Info("Starting the bot named " + viper.GetString("bot.name"))
	client, err := disgo.New(viper.GetString("discord.token"),
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
	privateGuilds := viper.GetStringSlice("discord.privateGuilds")
	commands.RegisterGuildCommands(client, privateGuilds)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s

	slog.Info("Attempting graceful shutdown")
	client.Close(context.TODO())
}
