package nasen

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nico-mayer/go_discordbot/db"
	"github.com/nico-mayer/go_discordbot/utils"
)

func Nasen(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var target *discordgo.User

	for _, option := range i.ApplicationCommandData().Options {
		if option.Type == discordgo.ApplicationCommandOptionUser {
			target = option.UserValue(s)
		}
	}

	user, err := db.GetUser(target.ID)
	if err != nil {
		utils.ReplyError(s, i, err, "Bre hat noch keine Clownsnasen! gib ihm eine mit `/clownsnase`")
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	utils.Check(err)

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: getDesc(user),
	})
	utils.Check(err)
}

func getDesc(user db.User) string {
	var sb strings.Builder

	heading := fmt.Sprintf("Alle Clownsnasen von <@%s> \n", user.ID)
	sb.WriteString(heading)
	sb.WriteString("```\n")
	sb.WriteString(generateTable(user, user.GetNasen()))
	sb.WriteString("```")

	return sb.String()
}

func generateTable(user db.User, nasen []db.Nase) string {
	var tableString strings.Builder
	t := table.NewWriter()
	t.SetOutputMirror(&tableString)
	t.AppendHeader(table.Row{"Datum", "Von", "Grund"})
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true
	t.SetAutoIndex(true)

	for _, nase := range nasen {
		author, err := db.GetUser(nase.AuthorID)
		utils.Check(err)

		date := fmt.Sprintf("%v", nase.Created.Format("02-Jan-06"))
		t.AppendRow(table.Row{date, author.Name, nase.Reason})
	}

	t.AppendSeparator()

	t.Render()

	return tableString.String()
}
