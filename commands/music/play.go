package music

import (
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
				Type:  discordgo.EmbedTypeRich,
				Title: song.Name,
			},
		},
	})
	utils.Check(err)

	err = p.Play(song, voiceState)
	if err != nil {
		utils.ReplyError(s, i, err, "Something went wrong playing your song.")
		return
	}
}
