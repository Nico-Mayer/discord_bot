package general

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nico-mayer/go_discordbot/commands"
	"github.com/nico-mayer/go_discordbot/utils"
)

func Help(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var sb strings.Builder

	for _, collection := range commands.Collections {
		icon := string(collection.Icon)
		name := collection.Name
		str := fmt.Sprintf("\n**%s - %s:**\n", icon, name)
		sb.WriteString(str)
		for _, command := range collection.Commands {
			name := command.Name
			desc := command.Description
			str := fmt.Sprintf("- `/%s` - %s\n", name, desc)
			sb.WriteString(str)
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{

			Embeds: []*discordgo.MessageEmbed{
				{
					Type:        discordgo.EmbedTypeRich,
					Title:       "Help:",
					Description: sb.String(),
				},
			},
		},
	})
	utils.Check(err)
}
