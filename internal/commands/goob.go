package commands

import "github.com/disgoorg/disgo/events"

func handleGoob(event *events.ApplicationCommandInteractionCreate) {
	switch *event.SlashCommandInteractionData().SubCommandName {
	case "info":
		handleInfo(event)
	}
}
