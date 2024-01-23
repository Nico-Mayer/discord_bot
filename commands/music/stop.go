package music

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nico-mayer/go_discordbot/player"
	"github.com/nico-mayer/go_discordbot/utils"
)

func Stop(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {

	p.Stop()

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Tschüss, euer DJ Rosine",
		},
	})
	utils.Check(err)
}
