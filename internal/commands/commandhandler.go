package commands

import (
	"github.com/disgoorg/disgo/events"
)

func HandleCommand(event *events.ApplicationCommandInteractionCreate) {
	switch event.ApplicationCommandInteraction.Data.CommandName() {
	case "stats":
		handleStats(event)
	case "kick":
		handleKick(event)
	case "ban":
		handleBan(event)
	}
}
