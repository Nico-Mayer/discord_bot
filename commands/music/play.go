package music

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nico-mayer/go_discordbot/player"
	"github.com/nico-mayer/go_discordbot/utils"
)

func Play(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	user := i.Member.User
	var url string
	for _, option := range i.ApplicationCommandData().Options {
		if option.Type == discordgo.ApplicationCommandOptionString {
			url = option.StringValue()
		}
	}

	voiceState, err := s.State.VoiceState(i.GuildID, user.ID)
	if err != nil {
		utils.ReplyDeferredError(s, i, err, "Du musst in einem Sprachkanal sein, um Musik abspielen zu können")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	utils.Check(err)

	song, err := player.GetSongInfo(url)
	if err != nil {
		utils.ReplyDeferredError(s, i, err, "Songdaten nicht gefunden. Gib bitte eine gültige URL an.")
		return
	}

	p.JoinChannel(voiceState)
	p.Enqueue(song)
	go p.Play()

	time.Sleep(500 * time.Millisecond)

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Type:        discordgo.EmbedTypeRich,
				Title:       "▶️ - Playing",
				Description: formatDesc(song, p),
				Color:       0xff0001,
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL:    song.Thumbnail.URL,
					Width:  int(song.Thumbnail.Width),
					Height: int(song.Thumbnail.Height),
				},
			},
		},
	})
	utils.Check(err)
}

func formatDesc(song *player.Song, player *player.Player) string {
	var sb strings.Builder

	heading := fmt.Sprintf("[%s](%s) \n", song.Name, song.FullUrl)
	sb.WriteString(heading)
	sb.WriteString("\n📃 Warteschlange: \n \n")

	for i, s := range player.QueueList {
		line := fmt.Sprintf("%d. `%s` \n", i, s)
		sb.WriteString(line)
	}

	return sb.String()
}
