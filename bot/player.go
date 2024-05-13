package mybot

import (
	"bufio"
	"context"
	"log/slog"
	"os"
	"time"

	"os/exec"
	"strings"

	"github.com/disgoorg/ffmpeg-audio"
	"github.com/nico-mayer/discordbot/config"
)

type Song struct {
	Title        string
	ID           string
	ThumbnailURL string
	Duration     string
	Query        string
}

func (b *Bot) Enqueue(query string) (Song, error) {
	song, err := getSongData(query)
	if err != nil {
		return Song{}, err
	}

	b.Queue = append(b.Queue, song)

	return song, nil
}

func (b *Bot) Dequeue() Song {
	song := b.Queue[0]
	b.Queue = b.Queue[1:]
	return song
}

func getSongData(query string) (Song, error) {
	cmd := exec.Command("yt-dlp",
		"--get-title",
		"--get-id",
		"--get-thumbnail",
		"--get-duration",
		"--ignore-errors",
		"--no-warnings",
		"--skip-download",
		query,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("loading songdata with yt-dlp command", err)
		return Song{}, err
	}

	metadata := strings.Split(string(output), "\n")

	var song Song
	song.Title = metadata[0]
	song.ID = metadata[1]
	song.ThumbnailURL = metadata[2]
	song.Duration = metadata[3]
	song.Query = query
	return song, nil
}

func (b *Bot) PlaySong() error {
	cmd := exec.Command(
		"yt-dlp", b.Queue[0].Query,
		"--extract-audio",
		"--audio-format", "opus",
		"--no-playlist",
		"-o", "-",
		"--quiet",
		"--ignore-errors",
		"--no-warnings",
	)
	cmd.Stderr = os.Stderr
	b.Dequeue()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("creating stdout pipe", err)
		return err
	}

	b.BotStatus = Playing

	if err = cmd.Start(); err != nil {
		slog.Error("stating yt-dlp command", err)
		return err
	}

	opusProvider, err := ffmpeg.New(context.TODO(), bufio.NewReader(stdout))
	if err != nil {
		slog.Error("creating opus provider", err)
		return err
	}

	conn := b.Client.VoiceManager().GetConn(config.GUILD_ID)
	conn.SetOpusFrameProvider(opusProvider)

	if err = cmd.Wait(); err != nil {
		slog.Error("waiting for yt-dlp command", err)
		return err
	}

	// 10 SEC TIMEOUT TO ENSURE SONG NOT GETS CUT AT THE END
	time.Sleep(time.Second * 10)

	defer func() {
		opusProvider.Close()
		if len(b.Queue) > 0 {
			go b.PlaySong()
		} else {
			conn.Close(context.TODO())
			b.BotStatus = Resting
		}
	}()

	return nil
}
