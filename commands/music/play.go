package music

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nico-mayer/go_discordbot/player"
)

func Play(s *discordgo.Session, i *discordgo.InteractionCreate, player *player.Player) {
	user := i.Member.User

	voiceState, err := s.State.VoiceState(i.GuildID, user.ID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Du musst in einem Voice Channel sein um Musik abspielen zu können.",
			},
		})
		return
	}

	var url string
	for _, option := range i.ApplicationCommandData().Options {
		if option.Type == discordgo.ApplicationCommandOptionString {
			url = option.StringValue()
		}
	}

	player.Play(url, voiceState)
}
