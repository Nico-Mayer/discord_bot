package general

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var PingCommand = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "answers with pong",
}

func PingCommandExecute(event *events.ApplicationCommandInteractionCreate) {
	err := event.CreateMessage(discord.NewMessageCreateBuilder().
		SetContent("pong").
		Build(),
	)
	if err != nil {
		slog.Error("error on sending response", slog.Any("err", err))
	}
}
