package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo/discord"
	"github.com/gooberspace/goobcontrol/internal/commands"
	"github.com/gooberspace/goobcontrol/internal/goobcontrol"
	_ "github.com/joho/godotenv/autoload"
)

var (
	version                    = "localdevelopment"
	GlobalApplicationCommands  []discord.ApplicationCommandCreate
	PrivateApplicationCommands []discord.ApplicationCommandCreate
)

func init() {

}

func main() {
	config := goobcontrol.CreateConfig()
	logger := goobcontrol.CreateLogger(config)
	commandHandler := commands.HandleCommand

	gb := goobcontrol.New(*logger, config, version, commandHandler)

	gb.Logger.Info("Starting the bot named " + gb.Config.GetString("bot.name"))

	gb.SetupBot()

	// Setting up a basic disgo client with some sane defaults
	// We're doing all event handling elsewhere so this file can stay small

	GlobalApplicationCommands = append(GlobalApplicationCommands,
		commands.GoobCommand,
	)

	PrivateApplicationCommands = append(PrivateApplicationCommands,
		commands.KickCommand,
		commands.BanCommand,
	)

	privateGuilds := gb.Config.GetStringSlice("discord.privateGuilds")
	gb.RegisterCommands(GlobalApplicationCommands, PrivateApplicationCommands, privateGuilds)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s

	slog.Info("Attempting graceful shutdown")
	gb.Client.Close(context.TODO())
}
