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

func PingCommandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	if data.CommandName() == PingCommand.Name {
		err := event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("pong").
			Build(),
		)
		if err != nil {
			slog.Error("error on sending response", slog.Any("err", err))
		}
	}
}
