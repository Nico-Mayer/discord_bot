package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo/events"
	mybot "github.com/nico-mayer/discordbot/bot"
	"github.com/nico-mayer/discordbot/commands"
	"github.com/nico-mayer/discordbot/commands/general"
	"github.com/nico-mayer/discordbot/commands/music"
	"github.com/nico-mayer/discordbot/commands/nasen"
	"github.com/nico-mayer/discordbot/config"
)

func main() {
	osSignals := make(chan os.Signal, 1)

	// Setup bot
	bot := mybot.NewBot()

	// Populate slash handler slice
	bot.Handlers = map[string]func(event *events.ApplicationCommandInteractionCreate, b *mybot.Bot) error{
		general.HelpCommand.Name:      general.HelpCommandHandler,
		general.UserCommand.Name:      general.UserCommandHandler,
		nasen.ClownsnaseCommand.Name:  nasen.ClownsnaseCommandHandler,
		nasen.ClownfiestaCommand.Name: nasen.ClownfiestaCommandHandler,
		nasen.NasenCommand.Name:       nasen.NasenCommandHandler,
		nasen.LeaderboardCommand.Name: nasen.LeaderboardCommandHandler,
		music.TestCommand.Name:        music.TestCommandHandler,
	}

	bot.SetupBot()

	// Register slash commands
	commands.RegisterSlashCommands(bot.Client, config.GUILD_ID)

	// Open Gateway
	if err := bot.Client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while connecting to gateway", slog.Any("err", err))
	}
	slog.Info("bot is now running. Press CTRL-C to exit.")

	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-osSignals
}
