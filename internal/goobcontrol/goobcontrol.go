package goobcontrol

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
)

func New(logger *slog.Logger, config *viper.Viper, version string, commandHandler func(*GoobControl, *events.ApplicationCommandInteractionCreate), db *bun.DB) *GoobControl {
	return &GoobControl{
		Config:         config,
		Logger:         logger,
		Version:        version,
		CommandHandler: commandHandler,
		DB:             db,
	}
}

type GoobControl struct {
	Client         bot.Client
	Logger         *slog.Logger
	Config         *viper.Viper
	Version        string
	DB             *bun.DB
	CommandHandler func(*GoobControl, *events.ApplicationCommandInteractionCreate)
}

func (gc *GoobControl) SetupBot() {
	var err error
	if gc.Client, err = disgo.New(gc.Config.GetString("discord.token"),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentDirectMessages,
			),
		),
		bot.WithEventListenerFunc(gc.HandleDiscordEvent),
	); err != nil {
		panic(err)
	}
	if err = gc.Client.OpenGateway(context.TODO()); err != nil {
		panic(err)
	}
}

func (gc *GoobControl) RegisterCommands(publicCommands []discord.ApplicationCommandCreate, privateCommands []discord.ApplicationCommandCreate, privateGuilds []string) {
	if _, err := gc.Client.Rest().SetGlobalCommands(gc.Client.ApplicationID(), publicCommands); err != nil {
		slog.Error("error while registering commands", slog.Any("err", err))
	}

	for _, g := range privateGuilds {
		if _, err := gc.Client.Rest().SetGuildCommands(gc.Client.ApplicationID(), snowflake.MustParse(g), privateCommands); err != nil {
			slog.Error("error while registering commands", slog.Any("err", err))
		}
	}

}

func (gc *GoobControl) HandleDiscordEvent(e bot.Event) {
	switch e := e.(type) {
	case *events.Ready:
		botReadyMessage(gc, e)
	case *events.HeartbeatAck:
		//do nothing
	case *events.GuildReady:
		//do nothing
	case *events.GuildsReady:
		guildsReadyMessage(gc, e)
	case *events.ApplicationCommandInteractionCreate:
		gc.CommandHandler(gc, e)
	default:
		slog.Debug(fmt.Sprint(reflect.TypeOf(e)))
	}
}

func guildsReadyMessage(gc *GoobControl, e *events.GuildsReady) {
	guilds, error := e.Client().Rest().GetCurrentUserGuilds("", 0, 0, 10, false)
	if error != nil {
		gc.Logger.Error("Error getting guilds")
	} else {
		gc.Logger.Info(
			"Guilds Ready",
			"Guilds Connected", len(guilds),
		)
	}
}

func botReadyMessage(gc *GoobControl, e *events.Ready) {
	gc.Logger.Info(
		"Bot connected",
		"Username", e.User.Username,
	)
}
