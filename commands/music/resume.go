package music

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nico-mayer/go_discordbot/player"
	"github.com/nico-mayer/go_discordbot/utils"
)

func Resume(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	err := p.Resume()
	if err != nil {
		utils.ReplyError(s, i, err, "DJ Rosine ist aktuell nicht pausiert")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Type:  discordgo.EmbedTypeRich,
					Title: "▶️ - Resume",
				},
			},
		},
	})
	utils.Check(err)
}
