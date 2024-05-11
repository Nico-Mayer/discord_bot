package mybot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/nico-mayer/discordbot/config"
)

func (b *Bot) onTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
	fmt.Println("Track ended handler")
	if err := b.Client.UpdateVoiceState(context.TODO(), config.GUILD_ID, nil, false, false); err != nil {
		slog.Error("err", err)
	}
}
