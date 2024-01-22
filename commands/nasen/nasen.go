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
		Content: "```" + generateTable(user, user.GetNasen()) + "```",
	})
	utils.Check(err)
}

func generateTable(user db.User, nasen []db.Nase) string {
	var tableString strings.Builder
	t := table.NewWriter()
	t.SetOutputMirror(&tableString)
	t.AppendHeader(table.Row{"Datum", "Von", "Grund"})
	t.SetStyle(table.StyleLight)

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
