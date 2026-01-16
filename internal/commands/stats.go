package commands

import (
	"log/slog"
	"runtime"
	"strconv"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
)

var appStart time.Time

func init() {
	appStart = time.Now()
}

func handleStats(event *events.ApplicationCommandInteractionCreate) {
	var gatewayPing string
	if event.Client().HasGateway() {
		gatewayPing = event.Client().Gateway().Latency().String()
	}

	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	embed := discord.NewEmbedBuilder().
		SetTitle("Bot Statistics").
		AddField("Name", "Loading...", false).
		AddField("REST Latency", "Loading...", false).
		AddField("Gateway Latency", gatewayPing, false).
		AddField("Started", appStart.Format("2006-01-02 15:04:05 MST"), false).
		AddField("Uptime", time.Since(appStart).Truncate(time.Second).String(), false).
		AddField("Memory used", bToMB(m.Sys), false).
		AddField("Guilds Joined", "Loading...", false)

	defer func() {
		var start int64

		//guilds, _ := event.Client().Rest().GetCurrentUserGuilds("", 0, 0, 0, false)
		botInfo, _ := event.Client().Rest().GetBotApplicationInfo(func(config *rest.RequestConfig) {
			start = time.Now().UnixNano()
		})
		duration := time.Now().UnixNano() - start

		embed.SetField(0, "Name", botInfo.Bot.Username, false)
		embed.SetField(1, "REST Latency", time.Duration(duration).String(), false)
		embed.SetField(6, "Guilds Joined", strconv.FormatInt(int64(*botInfo.ApproximateGuildCount), 10), false)
		if _, err := event.Client().Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{Embeds: &[]discord.Embed{embed.Build()}}); err != nil {
			slog.Error("Failed to update ping embed: ", slog.Any("err", err))
		}
	}()

	if err := event.CreateMessage(discord.NewMessageCreateBuilder().SetEmbeds(embed.Build()).SetEphemeral(true).Build()); err != nil {
		slog.Error("Error responding to interaction", slog.Any("err", err))
	}
}

func bToMB(b uint64) string {
	return strconv.FormatUint(b/1048576, 10) + " MB"
}
