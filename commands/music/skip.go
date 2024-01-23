package music

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/nico-mayer/go_discordbot/player"
	"github.com/nico-mayer/go_discordbot/utils"
)

func Skip(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	msg := "Skip ⏭️"

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	utils.Check(err)

	err = p.Skip()
	if err != nil {
		log.Println(err)
		msg = "Kein Song in der Warteschlange"
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: msg,
	})
	utils.Check(err)

}
