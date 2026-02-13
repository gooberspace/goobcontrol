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
	dbMigrator := goobcontrol.SetupDbMigrator(database)
	commandHandler := commands.HandleCommand

	// This function creates a new instance of our bot with a shared logger, config, database etc.
	gc := goobcontrol.New(logger, config, version, commandHandler, database, dbMigrator)

	gc.Logger.Info("Starting the bot named " + gc.Config.GetString("bot.name"))

	if dBerr := gc.TestDatabase(); dBerr != nil {
		gc.Logger.Error("Failed to connect to database", "Error", dBerr)
	} else if dBerr = gc.RunDbMigrations(); dBerr != nil {
		gc.Logger.Error("Error running Database Migrations", "Error", dBerr)
	}

	// This function sets up the actual connection to Discord
	gc.SetupBot()

	// Here we set the commands we want to register with the Discord API, the global commands work for everyone on every server
	// while the private commands only work in Goober Space or my private servers
	GlobalApplicationCommands = append(GlobalApplicationCommands,
		commands.GoobCommand,
	)

	PrivateApplicationCommands = append(PrivateApplicationCommands,
		commands.KickCommand,
		commands.BanCommand,
	)

	privateGuilds := gc.Config.GetStringSlice("discord.privateGuilds")
	gc.RegisterCommands(GlobalApplicationCommands, PrivateApplicationCommands, privateGuilds)

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s

	gc.Logger.Info("Attempting graceful shutdown")
	gc.Client.Close(context.TODO())
}
