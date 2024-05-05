package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nico-mayer/discordbot/cmd"
	"github.com/nico-mayer/discordbot/cmd/general"
	"github.com/nico-mayer/discordbot/config"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
)

func main() {
	// Initialize bot client
	client, err := disgo.New(config.TOKEN,
		bot.WithDefaultGateway(),
		// Event Listeners
		// General Collection
		bot.WithEventListenerFunc(general.SayCommandListener),
		bot.WithEventListenerFunc(general.PingCommandListener),
		bot.WithEventListenerFunc(general.HelpCommandListener),
	)
	if err != nil {
		slog.Error("error while building disgo bot instance", slog.Any("err", err))
		return
	}
	defer client.Close(context.TODO())

	// Register Slash Commands
	cmd.RegisterSlashCommands(client, config.GUILD_ID)

	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while connecting to gateway", slog.Any("err", err))
	}

	slog.Info("bot is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
