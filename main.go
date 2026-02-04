package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo/discord"
	"github.com/gooberspace/goobcontrol/internal/commands"
	"github.com/gooberspace/goobcontrol/internal/goobcontrol"
)

var (
	version                    = "localdevelopment"
	GlobalApplicationCommands  []discord.ApplicationCommandCreate
	PrivateApplicationCommands []discord.ApplicationCommandCreate
)

func main() {
	config := goobcontrol.CreateConfig()
	logger := goobcontrol.CreateLogger(config)
	database := goobcontrol.SetupDatabase(config)
	commandHandler := commands.HandleCommand

	// This function creates a new instance of our bot with a shared logger, config, database etc.
	gb := goobcontrol.New(logger, config, version, commandHandler, database)

	gb.Logger.Info("Starting the bot named " + gb.Config.GetString("bot.name"))

	// This function sets up the actual connection to Discord
	gb.SetupBot()

	gb.TestDatabase()

	// Here we set the commands we want to register with the Discord API, the global commands work for everyone on every server
	// while the private commands only work in Goober Space or my private servers
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

	gb.Logger.Info("Attempting graceful shutdown")
	gb.Client.Close(context.TODO())
}
