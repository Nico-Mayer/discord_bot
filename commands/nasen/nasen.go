package nasen

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/jedib0t/go-pretty/table"
	"github.com/nico-mayer/discordbot/config"
	"github.com/nico-mayer/discordbot/db"
)

var NasenCommand = discord.SlashCommandCreate{
	Name:        "nasen",
	Description: "Liste aller Nasen eines Users an",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "Wähle einen User aus",
			Required:    true,
		},
	},
}

func NasenCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	target := data.User("user")

	nasen, err := db.GetNasenForUser(target.ID)
	if err != nil {
		slog.Error("fetching nasen array for user", err)
	}

	if len(nasen) == 0 {
		event.CreateMessage(discord.MessageCreate{
			Flags:   discord.MessageFlagEphemeral,
			Content: "Diser user hat noch keine Clownsnase, du kannst ihm eine mit `/clownsnase` geben.",
		})
		return
	}

	event.DeferCreateMessage(false)
	event.Client().Rest().CreateFollowupMessage(config.APP_ID, event.Token(), discord.MessageCreate{
		Content: getDesc(target, nasen),
	})

}

func getDesc(user discord.User, nasen []db.Nase) string {
	var sb strings.Builder

	heading := fmt.Sprintf("Alle Clownsnasen von <@%s> \n", user.ID)
	sb.WriteString(heading)
	sb.WriteString("```\n")
	sb.WriteString(generateTable(nasen))
	sb.WriteString("```")

	return sb.String()
}

func generateTable(nasen []db.Nase) string {
	var tableString strings.Builder
	t := table.NewWriter()
	t.SetOutputMirror(&tableString)
	t.AppendHeader(table.Row{"Datum", "Von", "Grund"})
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true
	t.SetAutoIndex(true)

	for i := len(nasen) - 1; i >= 0; i-- {
		nase := nasen[i]
		author, err := db.GetUser(nase.AuthorID)
		if err != nil {
			slog.Error("error fetching user data in table generate", err)
		}

		date := fmt.Sprintf("%v", nase.Created.Format("02-Jan-06"))
		t.AppendRow(table.Row{date, author.Name, nase.Reason})
	}

	t.Render()

	return tableString.String()
}
