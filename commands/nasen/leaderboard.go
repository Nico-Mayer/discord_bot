package nasen

import (
	"fmt"
	"log"

	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nico-mayer/go_discordbot/db"
)

func Leaderboard(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Println(err)
	}

	leaderboard, err := db.GetLeaderboard()
	if err != nil {
		log.Println(err)
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Type:  discordgo.EmbedTypeRich,
				Title: "Clownsnasen Leaderboard  🤡",
				Thumbnail: &discordgo.MessageEmbedThumbnail{
					URL: "https://media.tenor.com/81u64lUzA_QAAAAi/clown-peepo.gif",
				},
				Description: formatLeaderboard(leaderboard),
			},
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func formatLeaderboard(leaderboard []db.User) string {
	var sb strings.Builder

	if len(leaderboard) == 0 {
		sb.WriteString("Bis jetzt hat noch niemand eine Clownsnase, verteile clownsnasen mit `/clownsnase`")
	}

	for i, user := range leaderboard {
		n := "n"

		if user.NasenCount == 1 {
			n = ""
		} else if user.NasenCount == 0 {
			continue
		}

		str := fmt.Sprintf("**%d.** <@%s> = %d Nase%s\n", i+1, user.ID, user.NasenCount, n)
		sb.WriteString(str)
	}

	return sb.String()
}
