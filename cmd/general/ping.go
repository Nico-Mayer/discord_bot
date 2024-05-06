package general

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/nico-mayer/discordbot/db"
)

var PingCommand = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "answers with pong",
}

func PingCommandExecute(event *events.ApplicationCommandInteractionCreate) {

	user, _ := db.GetUser("488328811063148554")

	slog.Info(user.Name)

	err := event.CreateMessage(discord.NewMessageCreateBuilder().
		SetContent("pong").
		Build(),
	)
	if err != nil {
		slog.Error("error on sending response", slog.Any("err", err))
	}
}
