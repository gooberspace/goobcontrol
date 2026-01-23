package commands

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func HandleCommand(event *events.ApplicationCommandInteractionCreate) {
	switch event.Data.Type() {
	case discord.ApplicationCommandTypeSlash:
		HandleSlashCommand(event)
	}
}

func HandleSlashCommand(event *events.ApplicationCommandInteractionCreate) {
	slog.Info("Slash command ran", "command", event.SlashCommandInteractionData().CommandName(), "username", event.User().Username, "userID", event.User().ID, "Options", event.SlashCommandInteractionData().Options)
	switch event.ApplicationCommandInteraction.Data.CommandName() {
	case "kick":
		handleKick(event)
	case "ban":
		handleBan(event)
	case "goob":
		handleGoob(event)
	}
}
