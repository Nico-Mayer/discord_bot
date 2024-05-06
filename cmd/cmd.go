package cmd

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/nico-mayer/discordbot/cmd/general"
	"github.com/nico-mayer/discordbot/cmd/music"
)

var commands = []discord.ApplicationCommandCreate{
	general.SayCommand,
	general.PingCommand,
	general.HelpCommand,
	general.UserCommand,
	music.PlayCommand,
}

func RegisterSlashCommands(client bot.Client, guildID snowflake.ID) {
	if _, err := client.Rest().SetGuildCommands(client.ApplicationID(), guildID, commands); err != nil {
		slog.Error("error while registering commands", slog.Any("err", err))
	}
}
