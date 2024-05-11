package mybot

import "github.com/disgoorg/disgolink/v3/lavalink"

type Queue struct {
	Tracks []lavalink.Track
}

func (m *Queue) Add(track ...lavalink.Track) {
	m.Tracks = append(m.Tracks, track...)
}
