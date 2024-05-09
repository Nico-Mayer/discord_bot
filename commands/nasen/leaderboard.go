package nasen

import (
	"fmt"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	mybot "github.com/nico-mayer/discordbot/bot"
	"github.com/nico-mayer/discordbot/config"
	"github.com/nico-mayer/discordbot/db"
)

var LeaderboardCommand = discord.SlashCommandCreate{
	Name:        "leaderboard",
	Description: "Zeigt Clownsnasen Leaderboard an",
}

func LeaderboardCommandHandler(event *events.ApplicationCommandInteractionCreate, b *mybot.Bot) error {
	leaderboard, err := db.GetLeaderboard()
	if err != nil {
		return err
	}

	event.DeferCreateMessage(false)

	_, err = event.Client().Rest().CreateFollowupMessage(config.APP_ID, event.Token(), discord.MessageCreate{Embeds: []discord.Embed{
		{
			Title:       "Clownsnasen Leaderboard  🤡",
			Thumbnail:   &discord.EmbedResource{URL: "https://media.tenor.com/81u64lUzA_QAAAAi/clown-peepo.gif"},
			Description: formatLeaderboard(leaderboard),
		},
	}})

	return err
}

func formatLeaderboard(leaderboard []db.LeaderboardEntry) string {
	var sb strings.Builder

	if len(leaderboard) == 0 {
		sb.WriteString("Bis jetzt hat noch niemand eine Clownsnase, verteile clownsnasen mit `/clownsnase`")
	}

	for i, entry := range leaderboard {
		n := "n"

		if entry.NasenCount == 1 {
			n = ""
		} else if entry.NasenCount == 0 {
			continue
		}

		str := fmt.Sprintf("**%d.** <@%s> = %d Nase%s\n", i+1, entry.UserID, entry.NasenCount, n)
		sb.WriteString(str)
	}

	return sb.String()
}
