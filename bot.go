package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nico-mayer/discordbot/cmd"
	"github.com/nico-mayer/discordbot/cmd/general"
	"github.com/nico-mayer/discordbot/cmd/music"
	"github.com/nico-mayer/discordbot/config"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func main() {
	osSignals := make(chan os.Signal, 1)

	// Initialize bot client
	client, err := disgo.New(config.TOKEN,
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuildVoiceStates)),
		// Add listener for slash commands
		bot.WithEventListenerFunc(func(event *events.ApplicationCommandInteractionCreate) {
			data := event.SlashCommandInteractionData()

			switch data.CommandName() {
			case "help":
				general.HelpCommandExecute(event)
			case "ping":
				general.PingCommandExecute(event)
			case "say":
				general.SayCommandExecute(event)
			case "play":
				go music.PlayCommandExecute(event)
			}
		}),
	)
	if err != nil {
		slog.Error("error while building disgo bot instance", slog.Any("err", err))
		return
	}
	defer client.Close(context.TODO())

	// Register slash commands
	cmd.RegisterSlashCommands(client, config.GUILD_ID)

	// Open Gateway
	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while connecting to gateway", slog.Any("err", err))
	}
	slog.Info("bot is now running. Press CTRL-C to exit.")

	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-osSignals
}
