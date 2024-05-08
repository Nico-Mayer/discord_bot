package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nico-mayer/discordbot/commands"
	"github.com/nico-mayer/discordbot/config"
)

func main() {
	osSignals := make(chan os.Signal, 1)

	// Setup bot
	bot := NewBot()
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
