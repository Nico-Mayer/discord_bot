package mybot

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/nico-mayer/discordbot/db"
	"github.com/nico-mayer/discordbot/levels"
)

type BotStatus int32

const (
	Resting BotStatus = 0
	Playing BotStatus = 1
)

type Bot struct {
	Client        bot.Client
	BotStatus     BotStatus
	Handlers      map[string]func(event *events.ApplicationCommandInteractionCreate, b *Bot) error
	Queue         []Song
	SkipInterrupt chan bool
}

func NewBot() *Bot {
	return &Bot{
		BotStatus:     Resting,
		SkipInterrupt: make(chan bool, 1),
	}
}

func (b *Bot) SetupBot() {
	var err error

	// Initialize bot client
	b.Client, err = disgo.New(
		os.Getenv("TOKEN"),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates, cache.FlagMembers, cache.FlagChannels),
		),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentGuildVoiceStates,
				gateway.IntentsAll,
			),
		),
		// Slash command listener
		bot.WithEventListenerFunc(b.onApplicationCommand),

		// Message create listener
		bot.WithEventListenerFunc(b.onMessageCreate),

		// Voice join listener
		bot.WithEventListenerFunc(b.onVoiceJoin),
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

func (b *Bot) onMessageCreate(event *events.MessageCreate) {
	author := event.Message.Member.User
	if author.Bot {
		return
	}

	dbUser, err := db.ValidateAndFetchUser(author.ID, author.Username)
	if err != nil {
		slog.Error("validating and fetching user on message create")
	}

	level, levelUp, err := levels.GrantExpToUser(dbUser, levels.EXP_PER_MESSAGE)
	if err != nil {
		slog.Error("granting exp to user on message create")
	}
	if levelUp {
		levels.HandleLevelUp(event.Client(), author.ID, level)
	}
}

func (b *Bot) onVoiceJoin(event *events.GuildVoiceJoin) {
	author := event.Member.User
	if event.Member.User.Bot {
		return
	}
	dbUser, err := db.ValidateAndFetchUser(author.ID, author.Username)
	if err != nil {
		slog.Error("validating and fetching user on voice join")
	}

	level, levelUp, err := levels.GrantExpToUser(dbUser, levels.EXP_PER_VOICE_JOIN)
	if err != nil {
		slog.Error("granting exp to user on voice join")
	}
	if levelUp {
		levels.HandleLevelUp(event.Client(), author.ID, level)
	}
}

func (b *Bot) SetStatus(status string) error {
	return b.Client.SetPresence(context.TODO(), gateway.WithCustomActivity(status))
}
