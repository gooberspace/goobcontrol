package commands

import (
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/json"
	"github.com/gooberspace/goobcontrol/internal/common"
	"github.com/gooberspace/goobcontrol/internal/goobcontrol"
)

var BanCommand = discord.SlashCommandCreate{
	Name:                     "ban",
	Description:              "Ban a member",
	DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionAdministrator),
	Contexts:                 []discord.InteractionContextType{discord.InteractionContextTypeGuild},
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "member",
			Description: "The member you would like to ban",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "reason",
			Description: "reason for banning the user",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "purge_duration",
			Description: "purge last x duration of messages (ex. 1h, 24h, 1h30m max 168h)",
			Required:    false,
		},
	},
}

func handleBan(gc *goobcontrol.GoobControl, event *events.ApplicationCommandInteractionCreate) {

	// Handle ban and return embed
	embed := checkDurationValidityAndBan(gc, event)

	if err := event.CreateMessage(
		discord.NewMessageCreateBuilder().AddEmbeds(embed.Build()).Build(),
	); err != nil {
		slog.Error("Error responding to interaction", slog.Any("err", err))
	}
}

func checkDurationValidityAndBan(gc *goobcontrol.GoobControl, event *events.ApplicationCommandInteractionCreate) *discord.EmbedBuilder {
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
	parsedDuration, parseErr := time.ParseDuration(purgeDurationString)
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

	if event.Client().Rest().AddBan(*event.GuildID(), user.User.ID, parsedDuration) != nil {
		return discord.NewEmbedBuilder().
			SetTitle("Error banning member").
			SetDescriptionf("Couldn't ban %s, check if the user is still in the server, if %s has enough permissions to ban them and if its role is above the targeted user.", user.EffectiveName(), gc.Config.GetString("bot.name")).
			SetColor(common.ColourError)
	}

	return discord.NewEmbedBuilder().
		SetTitle("Member banned").
		AddField("Name:", user.EffectiveName(), false).
		AddField("Reason", reason, false).
		AddField("Banned by", event.User().Username, false).
		SetColor(common.ColourSuccess)
}
