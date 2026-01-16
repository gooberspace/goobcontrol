package commands

import (
	"log/slog"
	"os"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/gooberspace/goobcontrol/internal/common"
)

func handleKick(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	user := data.Member("member")
	reason, reasonSet := data.OptString("reason")
	var embed *discord.EmbedBuilder
	if !reasonSet {
		reason = "n/a"
	}
	if kickErr := event.Client().Rest().RemoveMember(*event.GuildID(), user.User.ID); kickErr != nil {
		embed = discord.NewEmbedBuilder().
			SetTitle("Error kicking member").
			SetDescriptionf("Couldn't kick %s, check if the user is still in the server, if %s has enough permissions to kick them and if its role is above the targeted user", user.EffectiveName(), os.Getenv("BOT_NAME")).
			SetColor(common.ColourError)
	} else {
		embed = discord.NewEmbedBuilder().
			SetTitle("Member kicked").
			AddField("Name:", user.EffectiveName(), false).
			AddField("Reason", reason, false).
			AddField("Kicked by", event.User().Username, false).
			SetColor(common.ColourSuccess)
	}
	if err := event.CreateMessage(
		discord.NewMessageCreateBuilder().AddEmbeds(embed.Build()).Build(),
	); err != nil {
		slog.Error("Error responding to interaction", slog.Any("err", err))
	}

}
