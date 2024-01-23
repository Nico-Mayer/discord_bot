package music

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nico-mayer/go_discordbot/player"
	"github.com/nico-mayer/go_discordbot/utils"
)

func Pause(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {

	err := p.Pause()
	if err != nil {
		utils.ReplyError(s, i, err, "DJ-Rosine spielt aktuell keinen banger der pausiert werden kann")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Type:        discordgo.EmbedTypeRich,
					Title:       "⏸️ - Paused",
					Description: "`/resume` to continue playing",
				},
			},
		},
	})
	utils.Check(err)

}
