package mybot

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/nico-mayer/discordbot/config"
)

type Bot struct {
	Client   bot.Client
	Lavalink disgolink.Client
	Handlers map[string]func(event *events.ApplicationCommandInteractionCreate, b *Bot) error
}

func NewBot() *Bot {
	return &Bot{}
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
		// Lavalink event handler
		bot.WithEventListenerFunc(b.onVoiceStateUpdate),
		bot.WithEventListenerFunc(b.onVoiceServerUpdate),
	)
	if err != nil {
		log.Fatal("FATAL: failed to setup bot client", err)
	}
	defer b.Client.Close(context.TODO())

	// Initialize lavalink client
	b.Lavalink = disgolink.New(b.Client.ApplicationID())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = b.Lavalink.AddNode(ctx, disgolink.NodeConfig{
		Address:  config.NODE_ADDRESS,
		Password: config.NODE_PW,
		Secure:   true,
	})
	if err != nil {
		slog.Error("failed to add node", slog.Any("err", err))
		os.Exit(1)
	}

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

func (b *Bot) onVoiceStateUpdate(event *events.GuildVoiceStateUpdate) {
	if event.VoiceState.UserID != b.Client.ApplicationID() {
		return
	}
	b.Lavalink.OnVoiceStateUpdate(context.TODO(), event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
}

func (b *Bot) onVoiceServerUpdate(event *events.VoiceServerUpdate) {
	b.Lavalink.OnVoiceServerUpdate(context.TODO(), event.GuildID, event.Token, *event.Endpoint)
}
