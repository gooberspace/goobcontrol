package commands

import (
	"log/slog"
	"os"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/gooberspace/goobcontrol/internal/common"
)

func handleBan(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	user := data.Member("member")
	reason, reasonSet := data.OptString("reason")
	purgeDurationString, purgeDurationSet := data.OptString("purge_duration")

	if !reasonSet {
		reason = "n/a"
	}
	if !purgeDurationSet {
		purgeDurationString = "0s"
	}

	// Validate and parse duration
	parsedDuration, err := time.ParseDuration(purgeDurationString)
	embed := checkDurationValidityAndBan(event, user, reason, parsedDuration, err)

	if err := event.CreateMessage(
		discord.NewMessageCreateBuilder().AddEmbeds(embed.Build()).Build(),
	); err != nil {
		slog.Error("Error responding to interaction", slog.Any("err", err))
	}
}

func checkDurationValidityAndBan(event *events.ApplicationCommandInteractionCreate, user discord.ResolvedMember, reason string, parsedDuration time.Duration, parseErr error) *discord.EmbedBuilder {
	if parseErr != nil {
		return discord.NewEmbedBuilder().
			SetTitle("Parameter error").
			SetDescription("Failed to parse purge duration, check your formatting. you can use any combination of h, m and s (not days!). with a maximum of 604800 seconds (or 7 days in other units)")
	}

	if parsedDuration.Seconds() > 604800 {
		return discord.NewEmbedBuilder().
			SetTitle("Parameter error").
			SetDescription("Parsed purge duration correctly but it's more than 604800 seconds. Please pick a lower duration")
	}

	if banErr := event.Client().Rest().AddBan(*event.GuildID(), user.User.ID, parsedDuration); banErr != nil {
		return discord.NewEmbedBuilder().
			SetTitle("Error banning member").
			SetDescriptionf("Couldn't ban %s, check if the user is still in the server, if %s has enough permissions to ban them and if its role is above the targeted user.", user.EffectiveName(), os.Getenv("BOT_NAME")).
			SetColor(common.ColourError)
	}

	return discord.NewEmbedBuilder().
		SetTitle("Member banned").
		AddField("Name:", user.EffectiveName(), false).
		AddField("Reason", reason, false).
		AddField("Banned by", event.User().Username, false).
		SetColor(common.ColourSuccess)
}
