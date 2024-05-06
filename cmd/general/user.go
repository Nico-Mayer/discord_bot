package general

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/nico-mayer/discordbot/db"
)

var UserCommand = discord.SlashCommandCreate{
	Name:        "user",
	Description: "Zeigt Informationen über einen User an.",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "Waehle einen user aus.",
			Required:    true,
		},
	},
}

func UserCommandExecute(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()

	event.DeferCreateMessage(true)
	targetUser := data.User("user")

	if targetUser.Bot {
		return
	}

	if !db.UserInDatabase(targetUser.ID.String()) {
		err := db.InsertDBUser(targetUser.ID.String(), targetUser.Username)
		if err != nil {
			slog.Error("inserting user to database", err)
		}
	}

	/* event.MessageCommandInteractionData().TargetID()
	event.Client().Rest().UpdateInteractionResponse(config.APP_ID)
	*/
}
