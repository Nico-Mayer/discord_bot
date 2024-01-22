package nasen

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nico-mayer/go_discordbot/db"
	"github.com/nico-mayer/go_discordbot/utils"
)

func Clownsnase(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var target db.User
	var thumbnail *discordgo.MessageEmbedThumbnail
	authorID := i.Member.User.ID
	reason := "einfach schlecht"

	if len(i.ApplicationCommandData().Options) > 0 {
		for _, option := range i.ApplicationCommandData().Options {
			if option.Type == discordgo.ApplicationCommandOptionUser {
				target.ID = option.UserValue(s).ID
				target.Name = option.UserValue(s).Username
				thumbnail = &discordgo.MessageEmbedThumbnail{
					URL: option.UserValue(s).AvatarURL(""),
				}
			}

			if option.Type == discordgo.ApplicationCommandOptionString {
				reason = option.StringValue()
			}
		}
	} else {
		return
	}

	nasenCount, err := db.GetNasenCount(target.ID)
	utils.Check(err)
	err = db.GiveNase([]db.User{target}, authorID, reason)
	utils.Check(err)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{

			Embeds: []*discordgo.MessageEmbed{
				{
					Type:        discordgo.EmbedTypeRich,
					Title:       "Clownsnase Kassiert  🤡",
					Thumbnail:   thumbnail,
					Description: fmt.Sprintf("`Grund: %s`", reason),
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Von: ",
							Value:  fmt.Sprintf("<@%s>", authorID),
							Inline: true,
						},
						{
							Name:   "An: ",
							Value:  fmt.Sprintf("<@%s>", target.ID),
							Inline: true,
						},
						{
							Name:   "Total: ",
							Value:  fmt.Sprintf("%d", nasenCount+1),
							Inline: true,
						},
					},
				},
			},
		},
	})
	utils.Check(err)
}
