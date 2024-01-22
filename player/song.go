package player

import (
	"time"

	"github.com/kkdai/youtube/v2"
)

type Song struct {
	Name        string
	Description string
	FullUrl     string
	downloadUrl string
	Duration    time.Duration
	Thumbnail   *youtube.Thumbnail
}

func GetSongInfo(url string) (*Song, error) {
	client := youtube.Client{}
	sng, err := client.GetVideo(url)

	if err != nil {
		return nil, err
	}

	return &Song{
		Name:        sng.Title,
		Description: sng.Description,
		FullUrl:     url,
		downloadUrl: sng.Formats.WithAudioChannels()[0].URL,
		Duration:    sng.Duration,
		Thumbnail:   &sng.Thumbnails[0],
	}, nil
}
