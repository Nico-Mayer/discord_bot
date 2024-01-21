package player

import (
	"time"

	"github.com/kkdai/youtube/v2"
)

type Song struct {
	name        string
	fullUrl     string
	downloadUrl string
	duration    time.Duration
}

func GetSongInfo(url string) (*Song, error) {
	client := youtube.Client{}
	sng, err := client.GetVideo(url)

	if err != nil {
		return nil, err
	}

	return &Song{
		name:        sng.Title,
		fullUrl:     url,
		downloadUrl: sng.Formats.WithAudioChannels()[0].URL,
		duration:    sng.Duration,
	}, nil
}
