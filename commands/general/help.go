package general

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/nico-mayer/discordbot/config"
)

var HelpCommand = discord.SlashCommandCreate{
	Name:        "help",
	Description: "Zeigt liste aller Commands",
}

func HelpCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	event.DeferCreateMessage(true)

	var slashCommands []discord.SlashCommand

	commands, err := event.Client().Rest().GetGuildCommands(config.APP_ID, config.GUILD_ID, false)
	if err != nil {
		slog.Error("receiving guild commands", err)
	}

	for _, command := range commands {
		if slashCommand, ok := command.(discord.SlashCommand); ok {
			slashCommands = append(slashCommands, slashCommand)
		}
	}

	_, err = event.Client().Rest().CreateFollowupMessage(config.APP_ID, event.Token(), discord.MessageCreate{
		Embeds: []discord.Embed{
			{
				Title:       "ℹ️    Help",
				Description: generateList(slashCommands),
			},
		},
	})
	if err != nil {
		slog.Error("error on sending response", slog.Any("err", err))
	}

}

func generateList(slashCommands []discord.SlashCommand) (desc string) {
	for _, command := range slashCommands {
		line := fmt.Sprintf("- `/%s` - (%s)\n", command.Name(), command.Description)
		desc = desc + line
	}
	return desc
}
