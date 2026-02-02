package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/gooberspace/goobcontrol/internal/goobcontrol"
)

var GoobCommand = discord.SlashCommandCreate{
	Name:        "goob",
	Description: "Goob Control utility commands",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionSubCommand{
			Name:        "info",
			Description: "Get info and runtime stats",
		},
	},
}

func handleGoob(gb *goobcontrol.GoobControl, event *events.ApplicationCommandInteractionCreate) {
	switch *event.SlashCommandInteractionData().SubCommandName {
	case "info":
		handleInfo(gb, event)
	}
}
