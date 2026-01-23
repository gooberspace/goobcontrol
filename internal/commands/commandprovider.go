package commands

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

var (
	GlobalApplicationCommands  []discord.ApplicationCommandCreate
	PrivateApplicationCommands []discord.ApplicationCommandCreate
)

func init() {
	var goobCommand = discord.SlashCommandCreate{
		Name:        "goob",
		Description: "Goob Control utility commands",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "info",
				Description: "Get info and runtime stats",
			},
		},
	}

	var kickCommand = discord.SlashCommandCreate{
		Name:                     "kick",
		Description:              "Kick a member",
		DefaultMemberPermissions: json.NewNullablePtr(discord.PermissionAdministrator),
		Contexts:                 []discord.InteractionContextType{discord.InteractionContextTypeGuild},
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionUser{
				Name:        "member",
				Description: "The member you would like to kick",
				Required:    true,
			},
			discord.ApplicationCommandOptionString{
				Name:        "reason",
				Description: "reason for kicking the user",
				Required:    true,
			},
		},
	}
	var banCommand = discord.SlashCommandCreate{
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

	GlobalApplicationCommands = append(GlobalApplicationCommands,
		goobCommand,
	)

	PrivateApplicationCommands = append(PrivateApplicationCommands,
		kickCommand,
		banCommand,
	)
}

func RegisterCommands(client bot.Client) {
	if _, err := client.Rest().SetGlobalCommands(client.ApplicationID(), GlobalApplicationCommands); err != nil {
		slog.Error("error while registering commands", slog.Any("err", err))
	}
}

func RegisterGuildCommands(client bot.Client, guilds []string) {
	for _, g := range guilds {
		if _, err := client.Rest().SetGuildCommands(client.ApplicationID(), snowflake.MustParse(g), PrivateApplicationCommands); err != nil {
			slog.Error("error while registering commands", slog.Any("err", err))
		}
	}

}
