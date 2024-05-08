package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nico-mayer/discordbot/commands"
	"github.com/nico-mayer/discordbot/commands/general"
	"github.com/nico-mayer/discordbot/commands/music"
	"github.com/nico-mayer/discordbot/config"
	"github.com/nico-mayer/discordbot/db"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func main() {
	osSignals := make(chan os.Signal, 1)

	// Initialize bot client
	client, err := disgo.New(config.TOKEN,
		bot.WithDefaultGateway(),
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuildVoiceStates, gateway.IntentGuildMessages)),

		// Voice join listener
		bot.WithEventListenerFunc(func(event *events.GuildVoiceJoin) {}),

		// Message created listener
		bot.WithEventListenerFunc(func(event *events.MessageCreate) {
			author := event.Message.Author

			if author.Bot {
				return
			}

			userInDatabase := db.UserInDatabase(author.ID)

			if userInDatabase {
				dbuser, err := db.GetUser(author.ID)
				if err != nil {
					slog.Error("fetching user from database")
				}
				fmt.Println(dbuser.Level)
			}
		}),

		// Slash command listener
		bot.WithEventListenerFunc(func(event *events.ApplicationCommandInteractionCreate) {
			data := event.SlashCommandInteractionData()

			switch data.CommandName() {
			case "help":
				general.HelpCommandExecute(event)
			case "user":
				general.UserCommandExecute(event)
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
	commands.RegisterSlashCommands(client, config.GUILD_ID)

	// Open Gateway
	if err = client.OpenGateway(context.TODO()); err != nil {
		slog.Error("error while connecting to gateway", slog.Any("err", err))
	}
	slog.Info("bot is now running. Press CTRL-C to exit.")

	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-osSignals
}
