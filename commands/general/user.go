package general

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/nico-mayer/discordbot/config"
	"github.com/nico-mayer/discordbot/db"
)

var UserCommand = discord.SlashCommandCreate{
	Name:        "user",
	Description: "Zeigt Informationen über einen User an",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "Wähle einen User aus",
			Required:    true,
		},
	},
}

func UserCommandHandler(event *events.ApplicationCommandInteractionCreate) error {
	data := event.SlashCommandInteractionData()
	targetUser := data.User("user")

	if targetUser.Bot {
		return event.CreateMessage(discord.MessageCreate{
			Flags:   discord.MessageFlagEphemeral,
			Content: "Du kannst keine infos von Bots abrufen.",
		})
	}

	event.DeferCreateMessage(false)
	if !db.UserInDatabase(targetUser.ID) {
		err := db.InsertDBUser(targetUser.ID, targetUser.Username)
		if err != nil {
			return err
		}
	}

	dbUser, err := db.GetUser(targetUser.ID)
	if err != nil {
		return err
	}

	userNasenCount, err := db.GetNasenCountForUser(dbUser.ID)
	if err != nil {
		return err
	}

	_, err = event.Client().Rest().CreateFollowupMessage(config.APP_ID, event.Token(), discord.MessageCreate{
		Embeds: []discord.Embed{
			{
				Title:       targetUser.Username,
				Description: "User stats:",
				Color:       0x00ff00,
				Thumbnail: &discord.EmbedResource{
					URL: *targetUser.AvatarURL(),
				},
				Fields: []discord.EmbedField{
					{
						Name: "Level",
						// Todo add level calculation
						Value: fmt.Sprintf("```%d```", 1),
					}, {
						Name:  "Exp",
						Value: fmt.Sprintf("```%d```", dbUser.Exp),
					}, {
						Name:  "Nasen",
						Value: fmt.Sprintf("```%d```", userNasenCount),
					},
				},
				Footer: &discord.EmbedFooter{
					Text: "Um eine Liste aller Clownsnasen des Benutzers zu sehen, benutze /nasen.",
				},
			},
		},
	})
	return err
}
