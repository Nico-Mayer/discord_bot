package nasen

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nico-mayer/go_discordbot/config"
	"github.com/nico-mayer/go_discordbot/db"
)

func Clownfiesta(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var channelID string
	var channelName string
	var usersInChannel []db.User
	author := i.Member.User
	inVoice := false

	guild, err := s.State.Guild(config.GUILD_ID)
	if err != nil {
		fmt.Println(err)
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == author.ID {
			inVoice = true
			channelID = vs.ChannelID
			channel, _ := s.Channel(channelID)
			channelName = channel.Name

			if err != nil {
				log.Println(err)
			}
		}
	}

	if !inVoice {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Du musst in einem Sprachkanal sein!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Println(err)
		}
		return
	}

	for _, vs := range guild.VoiceStates {
		if vs.ChannelID == channelID {
			user, err := s.User(vs.UserID)
			if err != nil {
				log.Println(err)
			}
			usersInChannel = append(usersInChannel, db.User{
				Name: user.Username,
				ID:   user.ID,
			})
		}
	}

	err = db.GiveNase(usersInChannel, author.ID, "Clownfiesta!")

	if err != nil {
		log.Println(err)
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{

			Embeds: []*discordgo.MessageEmbed{
				{
					Type:  discordgo.EmbedTypeRich,
					Title: "Clownfiesta! 🤡",
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://i.kym-cdn.com/photos/images/newsfeed/001/480/336/e0a.gif",
					},
					Description: fmt.Sprintf("**Komplette Clownfiesta in `%s`!**", channelName) + buildList(usersInChannel),
				},
			},
		},
	})

	if err != nil {
		log.Println(err)
	}
}

func buildList(users []db.User) string {
	var sb strings.Builder
	for _, user := range users {
		row := fmt.Sprintf("\n- <@%s> +1 Clownsnase", user.ID)
		sb.WriteString(row)
	}
	return sb.String()
}
