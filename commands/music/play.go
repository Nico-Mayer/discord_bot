package music

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"
	mybot "github.com/nico-mayer/discordbot/bot"
	"github.com/nico-mayer/discordbot/config"
)

var (
	urlPattern = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
)

var PlayCommand = discord.SlashCommandCreate{
	Name:        "play",
	Description: "Startet die Wiedergabe eines Songs",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "identifier",
			Description: "url oder suche",
			Required:    true,
		},
	},
}

func PlayCommandHandler(event *events.ApplicationCommandInteractionCreate, bot *mybot.Bot) error {
	data := event.SlashCommandInteractionData()
	author := event.User()

	identifier := data.String("identifier")
	if !urlPattern.MatchString(identifier) {
		identifier = lavalink.SearchTypeYouTube.Apply(identifier)
	}

	voiceState, ok := event.Client().Caches().VoiceState(*event.GuildID(), author.ID)
	if !ok {
		return event.CreateMessage(discord.MessageCreate{
			Flags:   discord.MessageFlagEphemeral,
			Content: "Du musst in einem Voice Channel sein um diesen command zu nutzen",
		})
	}

	event.DeferCreateMessage(false)

	var toPlay *lavalink.Track

	bot.Lavalink.BestNode().LoadTracksHandler(context.TODO(), identifier, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			bot.Client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					buildPlayingEmbed(track, author),
				},
			})
			toPlay = &track
		},
		func(playlist lavalink.Playlist) {
			bot.Client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
				Content: json.Ptr(fmt.Sprintf("Loaded playlist: `%s` with `%d` tracks", playlist.Info.Name, len(playlist.Tracks))),
			})
			toPlay = &playlist.Tracks[0]
		},
		func(tracks []lavalink.Track) {
			bot.Client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
				Embeds: &[]discord.Embed{
					buildPlayingEmbed(tracks[0], author),
				},
			})
			toPlay = &tracks[0]
		},
		func() {
			bot.Client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
				Content: json.Ptr(fmt.Sprintf("Nothing found for: `%s`", identifier)),
			})
		},
		func(err error) {
			bot.Client.Rest().UpdateInteractionResponse(event.ApplicationID(), event.Token(), discord.MessageUpdate{
				Content: json.Ptr(fmt.Sprintf("Error while looking up query: `%s`", err)),
			})
		},
	))
	if toPlay == nil {
		return errors.New("error fetching song data")
	}

	if err := bot.Client.UpdateVoiceState(context.TODO(), config.GUILD_ID, voiceState.ChannelID, false, false); err != nil {
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
