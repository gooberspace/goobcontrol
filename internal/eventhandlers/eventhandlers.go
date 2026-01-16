package eventhandlers

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/gooberspace/goobcontrol/internal/commands"
)

func HandleDiscordEvent(e bot.Event) {
	switch e := e.(type) {
	case *events.Ready:
		botReadyMessage(e)
	case *events.HeartbeatAck:
		//handleHeartBeat(e)
	case *events.GuildReady:
		//guildReadyMessage(e)
	case *events.GuildsReady:
		guildsReadyMessage(e)
	case *events.ApplicationCommandInteractionCreate:
		commands.HandleCommand(e)
	default:
		slog.Debug(fmt.Sprint(reflect.TypeOf(e)))
	}
}

func guildReadyMessage(e *events.GuildReady) {
	slog.Info(
		"Guild Ready",
		"Guild Name", e.Guild.Name,
		"Guild ID", e.Guild.ID,
	)
}

func guildsReadyMessage(e *events.GuildsReady) {
	guilds, error := e.Client().Rest().GetCurrentUserGuilds("", 0, 0, 10, false)
	if error != nil {
		slog.Error("Error getting guilds")
	} else {
		slog.Info(
			"Guilds Ready",
			"Guilds Connected", len(guilds),
		)
	}
}

func botReadyMessage(e *events.Ready) {
	slog.Info(
		"Bot connected",
		"Username", e.User.Username,
	)
}

func handleHeartBeat(e *events.HeartbeatAck) {
	slog.Debug("Gateway Heartbeat Ack",
		slog.Time("Last", e.LastHeartbeat),
		slog.Time("New", e.NewHeartbeat),
	)
}
