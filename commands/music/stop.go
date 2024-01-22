package music

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nico-mayer/go_discordbot/player"
)

func Stop(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	fmt.Println("Stop player")

	p.Stop()

}
