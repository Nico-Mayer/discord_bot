package music

import (
	"fmt"
	"strings"

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
		utils.ReplyError(s, i, err, "Du musst in einem Voice Channel sein um Musik abspielen zu können.")
		return
	}
	p.JoinChannel(voiceState)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	utils.Check(err)

	song, err := player.GetSongInfo(url)
	if err != nil {
		utils.ReplyError(s, i, err, "Error Fetching Song Data")
		return
	}

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
	p.Enqueue(song)
	p.Play()
}

func formatDesc(song *player.Song, player *player.Player) string {
	var sb strings.Builder

	heading := fmt.Sprintf("[%s](%s) \n", song.Name, song.FullUrl)

	sb.WriteString(heading)

	sb.WriteString("\n📃 Warteschlange: \n \n")
	sb.WriteString(fmt.Sprintf("%s \n", song.Name))

	for _, s := range player.QueueList {
		line := fmt.Sprintf("%s \n", s)
		sb.WriteString(line)
	}

	return sb.String()
}
