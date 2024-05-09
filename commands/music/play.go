package music

import (
	"context"
	"errors"
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
	mybot "github.com/nico-mayer/discordbot/bot"
	"github.com/nico-mayer/discordbot/config"
)

var PlayCommand = discord.SlashCommandCreate{
	Name:        "play",
	Description: "Startet die Wiedergabe eines Songs",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "query",
			Description: "search query",
			Required:    true,
		},
	},
}

func PlayCommandHandler(event *events.ApplicationCommandInteractionCreate, bot *mybot.Bot) error {
	data := event.SlashCommandInteractionData()
	query := data.String("query")
	author := event.User()

	event.DeferCreateMessage(false)

	var toPlay *lavalink.Track
	bot.Lavalink.BestNode().LoadTracksHandler(context.TODO(), query, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			_, _ = bot.Client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					buildPlayingEmbed(track, author),
				},
			})
			toPlay = &track
		},
		func(playlist lavalink.Playlist) {
			_, _ = bot.Client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
				Content: json.Ptr(fmt.Sprintf("Loaded playlist: `%s` with `%d` tracks", playlist.Info.Name, len(playlist.Tracks))),
			})
			toPlay = &playlist.Tracks[0]
		},
		func(tracks []lavalink.Track) {
			_, _ = bot.Client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
				Content: json.Ptr(fmt.Sprintf("Loaded search result: [`%s`](<%s>)", tracks[0].Info.Title, *tracks[0].Info.URI)),
			})
			toPlay = &tracks[0]
		},
		func() {
			_, _ = bot.Client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
				Content: json.Ptr(fmt.Sprintf("Nothing found for: `%s`", query)),
			})
		},
		func(err error) {
			_, _ = bot.Client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
				Content: json.Ptr(fmt.Sprintf("Error while looking up query: `%s`", err)),
			})
		},
	))
	if toPlay == nil {
		return errors.New("error fetching song data")
	}

	channelID := snowflake.MustParse("1082979754312994880")

	if err := bot.Client.UpdateVoiceState(context.TODO(), config.GUILD_ID, &channelID, false, false); err != nil {
		return err
	}

	return bot.Lavalink.Player(*event.GuildID()).Update(context.TODO(), lavalink.WithTrack(*toPlay))

}

func buildPlayingEmbed(track lavalink.Track, author discord.User) discord.Embed {
	return discord.Embed{
		Author: &discord.EmbedAuthor{
			Name:    author.Username,
			IconURL: *author.AvatarURL(),
		},
		Title:       "▶️ - Playing:",
		Description: fmt.Sprintf("Loaded track: [`%s`](<%s>) [03:11] 🔊", track.Info.Title, *track.Info.URI),
		Thumbnail:   &discord.EmbedResource{URL: *track.Info.ArtworkURL},
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("Source: %s", track.Info.SourceName),
		},
	}
}
