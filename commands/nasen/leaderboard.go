package nasen

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/nico-mayer/discordbot/db"
)

var LeaderboardCommand = discord.SlashCommandCreate{
	Name:        "leaderboard",
	Description: "Zeigt Clownsnasen Leaderboard an",
}

func LeaderboardCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	leaderboard, err := db.GetLeaderboard()
	if err != nil {
		slog.Error("feting leaderboard data from database", err)
	}

	event.CreateMessage(discord.MessageCreate{
		Content: generateLeaderboard(leaderboard),
	})
}

func generateLeaderboard(leaderboard []db.LeaderboardEntry) string {
	var sb strings.Builder

	for _, entry := range leaderboard {
		sb.WriteString(fmt.Sprintf("Key: <@%s> - Value: %d\n", entry.UserID, entry.NasenCount))
	}

	return sb.String()
}
