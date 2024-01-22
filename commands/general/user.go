package general

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nico-mayer/go_discordbot/db"
	"github.com/nico-mayer/go_discordbot/utils"
)

func User(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var target *discordgo.User

	for _, option := range i.ApplicationCommandData().Options {
		if option.Type == discordgo.ApplicationCommandOptionUser {
			target = option.UserValue(s)
		}
	}

	var dbUser db.User

	dbUser, err := db.GetUser(target.ID)
	utils.Check(err)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Type:      discordgo.EmbedTypeRich,
					Title:     target.Username,
					Color:     0x00ff00,
					Thumbnail: &discordgo.MessageEmbedThumbnail{URL: target.AvatarURL("")},
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Level",
							Value: fmt.Sprintf("```%d```", dbUser.Level),
						}, {
							Name:  "Exp",
							Value: fmt.Sprintf("```%d```", dbUser.Exp),
						}, {
							Name:  "Nasen",
							Value: fmt.Sprintf("```%d```", len(dbUser.GetNasen())),
						},
					},
				},
			},
		},
	})
	utils.Check(err)
}
