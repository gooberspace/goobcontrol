package commands

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/gooberspace/goobcontrol/internal/common"
	"github.com/spf13/viper"
)

func handleKick(event *events.ApplicationCommandInteractionCreate) {

	embed := tryKickOrFail(event)

	if err := event.CreateMessage(
		discord.NewMessageCreateBuilder().AddEmbeds(embed.Build()).Build(),
	); err != nil {
		slog.Error("Error responding to interaction", slog.Any("err", err))
	}

}

func tryKickOrFail(event *events.ApplicationCommandInteractionCreate) *discord.EmbedBuilder {
	data := event.SlashCommandInteractionData()
	user := data.Member("member")
	reason, reasonSet := data.OptString("reason")
	if !reasonSet {
		reason = "n/a"
	}

	if kickErr := event.Client().Rest().RemoveMember(*event.GuildID(), user.User.ID); kickErr != nil {
		return discord.NewEmbedBuilder().
			SetTitle("Error kicking member").
			SetDescriptionf("Couldn't kick %s, check if the user is still in the server, if %s has enough permissions to kick them and if its role is above the targeted user", user.EffectiveName(), viper.GetString("bot.name")).
			SetColor(common.ColourError)
	}

	return discord.NewEmbedBuilder().
		SetTitle("Member kicked").
		AddField("Name:", user.EffectiveName(), false).
		AddField("Reason", reason, false).
		AddField("Kicked by", event.User().Username, false).
		SetColor(common.ColourSuccess)
}
