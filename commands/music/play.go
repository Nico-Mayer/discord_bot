package music

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/ffmpeg-audio"
	mybot "github.com/nico-mayer/discordbot/bot"
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
	author := event.User()
	query := data.String("query")

	if bot.BotStatus == mybot.Playing {
		return handleError(event, "I am already playing a song", errors.New("already playing a song"))
	}

	if !urlPattern.MatchString(query) {
		query = "ytsearch:" + strings.TrimSpace(query)
	}

	voiceState, ok := event.Client().Caches().VoiceState(*event.GuildID(), author.ID)
	if !ok {
		return handleError(event, "Du musst in einem Voice Channel sein um diesen command zu nutzen", errors.New("author of play command not in a voice channel"))
	}

	if err := event.DeferCreateMessage(false); err != nil {
		return handleError(event, "Error creating defer message", err)
	}

	cmd := exec.Command(
		"yt-dlp", query,
		"--extract-audio",
		"--audio-format", "opus",
		"--no-playlist",
		"-o", "-",
		"--quiet",
		"--ignore-errors",
		"--no-warnings",
	)
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return handleError(event, "Error creating stdout pipe for song", err)
	}

	go func() {
		conn := event.Client().VoiceManager().CreateConn(*event.GuildID())
		if err = conn.Open(context.TODO(), *voiceState.ChannelID, false, false); err != nil {
			handleDeferError(event, "Error connecting to voice channel", err)
			return
		}

		if err = conn.SetSpeaking(context.TODO(), voice.SpeakingFlagMicrophone); err != nil {
			handleDeferError(event, "Error set bot to speaking", err)
			return
		}

		bot.BotStatus = mybot.Playing

		if err = cmd.Start(); err != nil {
			handleDeferError(event, "Error starting yt-dlp", err)
			return
		}

		opusProvider, err := ffmpeg.New(context.TODO(), bufio.NewReader(stdout))
		if err != nil {
			handleDeferError(event, "Error starting ffmpeg", err)
			return
		}

		conn.SetOpusFrameProvider(opusProvider)

		event.Client().Rest().CreateFollowupMessage(event.ApplicationID(), event.Token(), discord.MessageCreate{
			Content: fmt.Sprintf("Loaded song: [%s]", query),
		})
		if err = cmd.Wait(); err != nil {
			handleDeferError(event, "ERROR waiting for yt-dlp", err)
		}

		time.Sleep(5 * time.Second)

		defer func() {
			bot.BotStatus = mybot.Resting
			opusProvider.Close()
			conn.Close(context.TODO())
		}()
	}()
	return nil
}

func handleError(event *events.ApplicationCommandInteractionCreate, message string, err error) error {
	slog.Error(message, err)
	return event.CreateMessage(discord.MessageCreate{
		Flags:   discord.MessageFlagEphemeral,
		Content: message,
	})
}

func handleDeferError(event *events.ApplicationCommandInteractionCreate, message string, err error) error {
	slog.Error(message, err)
	_, error := event.Client().Rest().CreateFollowupMessage(event.ApplicationID(), event.Token(), discord.MessageCreate{
		Content: message,
	})
	return error
}
