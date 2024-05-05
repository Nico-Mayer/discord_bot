package general

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var SayCommand = discord.SlashCommandCreate{
	Name:        "say",
	Description: "says what you say",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "message",
			Description: "What to say",
			Required:    true,
		},
		discord.ApplicationCommandOptionBool{
			Name:        "ephemeral",
			Description: "If the response should only be visible to you",
			Required:    true,
		},
	},
}

func SayCommandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	if data.CommandName() == "say" {
		err := event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent(data.String("message")).
			SetEphemeral(data.Bool("ephemeral")).
			Build(),
		)
		if err != nil {
			slog.Error("error on sending response", slog.Any("err", err))
		}
	}
}
