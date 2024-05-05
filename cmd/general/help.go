package general

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var HelpCommand = discord.SlashCommandCreate{
	Name:        "help",
	Description: "answers with pong",
}

func HelpCommandListener(event *events.ApplicationCommandInteractionCreate) {

	data := event.SlashCommandInteractionData()
	if data.CommandName() == HelpCommand.Name {

		// test, _ := rest.Applications.GetGuildCommands(config.APP_ID, config.GUILD_ID, true, nil)

		/* for _, cmd := range test {
			fmt.Println(cmd.Name())
		} */

		err := event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Pls help me i am under the water").
			Build(),
		)
		if err != nil {
			slog.Error("error on sending response", slog.Any("err", err))
		}
	}
}
