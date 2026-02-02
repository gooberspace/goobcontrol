package commands

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/gooberspace/goobcontrol/internal/goobcontrol"
)

func HandleCommand(gb *goobcontrol.GoobControl, event *events.ApplicationCommandInteractionCreate) {
	switch event.Data.Type() {
	case discord.ApplicationCommandTypeSlash:
		handleSlashCommand(gb, event)
	}
}

func handleSlashCommand(gc *goobcontrol.GoobControl, event *events.ApplicationCommandInteractionCreate) {
	slog.Info("Slash command ran", "command", event.SlashCommandInteractionData().CommandName(), "username", event.User().Username, "userID", event.User().ID, "Options", event.SlashCommandInteractionData().Options)
	switch event.ApplicationCommandInteraction.Data.CommandName() {
	case "kick":
		handleKick(gc, event)
	case "ban":
		handleBan(gc, event)
	case "goob":
		handleGoob(gc, event)
	}
}
