package general

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	mybot "github.com/nico-mayer/discordbot/bot"
	"github.com/nico-mayer/discordbot/config"
)

var HelpCommand = discord.SlashCommandCreate{
	Name:        "help",
	Description: "Zeigt liste aller Commands",
}

func HelpCommandHandler(event *events.ApplicationCommandInteractionCreate, b *mybot.Bot) error {
	event.DeferCreateMessage(true)

	var slashCommands []discord.SlashCommand

	commands, err := event.Client().Rest().GetGuildCommands(config.APP_ID, config.GUILD_ID, false)
	if err != nil {
		return err
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
		return err
	}

	return nil

}

func generateList(slashCommands []discord.SlashCommand) (desc string) {
	for _, command := range slashCommands {
		line := fmt.Sprintf("- `/%s` - (%s)\n", command.Name(), command.Description)
		desc = desc + line
	}
	return desc
}
