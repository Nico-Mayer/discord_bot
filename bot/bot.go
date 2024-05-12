package mybot

import (
	"context"
	"log"
	"log/slog"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/nico-mayer/discordbot/config"
)

type BotStatus int32

const (
	Resting BotStatus = 0
	Playing BotStatus = 1
)

type Bot struct {
	Client    bot.Client
	BotStatus BotStatus
	Handlers  map[string]func(event *events.ApplicationCommandInteractionCreate, b *Bot) error
}

func NewBot() *Bot {
	return &Bot{
		BotStatus: Resting,
	}
}

func (b *Bot) SetupBot() {
	var err error

	// Initialize bot client
	b.Client, err = disgo.New(config.TOKEN,
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates, cache.FlagMembers, cache.FlagChannels),
		),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentGuildVoiceStates,
			),
		),

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
	err := handler(event, b)
	if err != nil {
		slog.Error("executing slash command", slog.String("command", data.CommandName()), err)
	}
}
