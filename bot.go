package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/nico-mayer/discordbot/commands/general"
	"github.com/nico-mayer/discordbot/commands/nasen"
	"github.com/nico-mayer/discordbot/config"
)

type Bot struct {
	Client   bot.Client
	Handlers map[string]func(event *events.ApplicationCommandInteractionCreate)
}

func NewBot() *Bot {
	return &Bot{}
}

func (b *Bot) SetupBot() {
	var err error

	// Setup Slash Handlers
	b.Handlers = map[string]func(event *events.ApplicationCommandInteractionCreate){
		"ping":        general.PingCommandHandler,
		"user":        general.UserCommandHandler,
		"help":        general.HelpCommandHandler,
		"say":         general.SayCommandHandler,
		"clownsnase":  nasen.ClownsnaseCommandHandler,
		"clownfiesta": nasen.ClownfiestaCommandHandler,
	}

	// Initialize bot client
	b.Client, err = disgo.New(config.TOKEN,
		bot.WithDefaultGateway(),
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuildVoiceStates, gateway.IntentGuildMessages)),

		// Slash command listener
		bot.WithEventListenerFunc(b.onApplicationCommand),
	)
	if err != nil {
		log.Fatal("FATAL: failed to setup bot client", err)
	}
	defer b.Client.Close(context.TODO())

}

func (b *Bot) onApplicationCommand(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()

	handler, ok := b.Handlers[data.CommandName()]
	if !ok {
		slog.Info("unknown command", slog.String("command", data.CommandName()))
		return
	}
	handler(event)
}
