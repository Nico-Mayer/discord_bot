package general

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var HelpCommand = discord.SlashCommandCreate{
	Name:        "help",
	Description: "helps",
}

func HelpCommandExecute(event *events.ApplicationCommandInteractionCreate) {

	err := event.CreateMessage(discord.NewMessageCreateBuilder().
		SetContent("Pls help me i am under the water").
		Build(),
	)
	if err != nil {
		slog.Error("error on sending response", slog.Any("err", err))
	}

}
