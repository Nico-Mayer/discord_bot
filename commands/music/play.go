package music

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/ffmpeg-audio"
	mybot "github.com/nico-mayer/discordbot/bot"
)

var (
	urlPattern = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
)

var PlayCommand = discord.SlashCommandCreate{
	Name:        "play",
	Description: "Startet die Wiedergabe eines Songs",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "query",
			Description: "url oder suche",
			Required:    true,
		},
	},
}

func PlayCommandHandler(event *events.ApplicationCommandInteractionCreate, bot *mybot.Bot) error {
	data := event.SlashCommandInteractionData()
	author := event.User()
	query := data.String("query")

	if !urlPattern.MatchString(query) {
		query = "ytsearch:" + strings.TrimSpace(query)
	}

	voiceState, ok := event.Client().Caches().VoiceState(*event.GuildID(), author.ID)
	if !ok {
		return event.CreateMessage(discord.MessageCreate{
			Flags:   discord.MessageFlagEphemeral,
			Content: "Du musst in einem Voice Channel sein um diesen command zu nutzen",
		})
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
		slog.Error("creating stdout pipe", err)
		return err
	}

	if err = event.DeferCreateMessage(false); err != nil {
		slog.Error("creating defer message", err)
		return err
	}

	go func() {
		conn := event.Client().VoiceManager().CreateConn(*event.GuildID())
		if err = conn.Open(context.TODO(), *voiceState.ChannelID, false, false); err != nil {
			slog.Error("connecting to voice", err)
		}
		defer conn.Close(context.TODO())

		if err = conn.SetSpeaking(context.TODO(), voice.SpeakingFlagMicrophone); err != nil {
			slog.Error("set speaking", err)
		}

		if err = cmd.Start(); err != nil {
			slog.Error("starting yt-dlp", err)
		}

		opusProvider, err := ffmpeg.New(context.TODO(), bufio.NewReader(stdout))
		if err != nil {
			slog.Error("ffmpeg start", err)
		}
		defer opusProvider.Close()

		conn.SetOpusFrameProvider(opusProvider)

		event.Client().Rest().CreateFollowupMessage(event.ApplicationID(), event.Token(), discord.MessageCreate{
			Content: fmt.Sprintf("Loaded song: [%s]", query),
		})
		if err = cmd.Wait(); err != nil {
			slog.Error("error waiting for yt-dlp: ", err)
		}
	}()

	return nil
}
