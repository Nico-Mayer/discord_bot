package music

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/voice"
	mybot "github.com/nico-mayer/discordbot/bot"
	"github.com/nico-mayer/discordbot/config"
)

var (
	urlPattern = regexp.MustCompile(`^https?://(?:www\.)?(?:youtube\.com/watch\?v=|youtu\.be/)[a-zA-Z0-9_-]{11}`)
)

var PlayCommand = discord.SlashCommandCreate{
	Name:        "play",
	Description: "Startet die Wiedergabe eines Songs",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "query",
			Description: "youtube url oder suche",
			Required:    true,
		},
	},
}

func PlayCommandHandler(event *events.ApplicationCommandInteractionCreate, bot *mybot.Bot) error {
	data := event.SlashCommandInteractionData()
	query := data.String("query")
	if !urlPattern.MatchString(query) {
		query = "ytsearch:" + strings.TrimSpace(query)
	}
	voiceState, ok := event.Client().Caches().VoiceState(*event.GuildID(), event.User().ID)
	if !ok {
		return event.CreateMessage(discord.MessageCreate{
			Flags:   discord.MessageFlagEphemeral,
			Content: "Du musst in einem voice channel sein um diesen command zu benutzen.",
		})
	}

	if err := event.DeferCreateMessage(false); err != nil {
		return err
	}

	song, err := bot.Enqueue(query)
	if err != nil {
		event.Client().Rest().CreateFollowupMessage(event.ApplicationID(), event.Token(), discord.MessageCreate{
			Content: "Fehler beim laden der songdaten",
		})
		return err
	}

	// ADD TO QUEUE CASE
	if bot.BotStatus == mybot.Playing {
		_, err = event.Client().Rest().CreateFollowupMessage(event.ApplicationID(), event.Token(), discord.MessageCreate{
			Embeds: []discord.Embed{
				{
					Author: &discord.EmbedAuthor{
						Name:    event.User().Username,
						IconURL: *event.User().AvatarURL(),
					},
					Title:       "📃 Warteschlange",
					Description: fmt.Sprintf("Added song: [`%s`]", song.Title),
					Thumbnail: &discord.EmbedResource{
						URL: song.ThumbnailURL,
					},
				},
			},
		})
		return err
	}

	// Play SONG CASE
	// CONNECT TO VOICE, IS A BLOCKING CALL SO RUN IN GO ROUTINE
	go func() {
		conn := bot.Client.VoiceManager().CreateConn(config.GUILD_ID)
		if err = conn.Open(context.TODO(), *voiceState.ChannelID, false, false); err != nil {
			slog.Error("connecting to voice channel", err)
		}
		if err = conn.SetSpeaking(context.TODO(), voice.SpeakingFlagMicrophone); err != nil {
			slog.Error("setting bot to speaking", err)
		}
	}()

	go bot.PlaySong()
	event.Client().Rest().CreateFollowupMessage(event.ApplicationID(), event.Token(), discord.MessageCreate{
		Embeds: []discord.Embed{
			{
				Author: &discord.EmbedAuthor{
					Name:    event.User().Username,
					IconURL: *event.User().AvatarURL(),
				},
				URL:         "https://www.youtube.com/watch?v=" + song.ID,
				Title:       song.Title,
				Description: fmt.Sprintf("▶️ Playing [`%s`]", song.Duration),
				Thumbnail: &discord.EmbedResource{
					URL: song.ThumbnailURL,
				},
				Footer: &discord.EmbedFooter{
					Text: "source: youtube",
				},
			},
		},
	})

	return nil
}
