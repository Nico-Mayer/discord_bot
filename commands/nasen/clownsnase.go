package nasen

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/nico-mayer/discordbot/config"
	"github.com/nico-mayer/discordbot/db"
)

var ClownsnaseCommand = discord.SlashCommandCreate{
	Name:        "clownsnase",
	Description: "Gib einem bre eine Clownsnase",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "Wähle einen user aus",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "reason",
			Description: "Grund für die clownsnase",
			Required:    false,
		},
	},
}

func ClownsnaseCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()

	author := event.User()
	target := data.User("user")
	reason := data.String("reason")

	event.DeferCreateMessage(false)

	if target.Bot {
		event.Client().Rest().CreateFollowupMessage(config.APP_ID, event.Token(), discord.MessageCreate{
			Content: "Du kannst mir keine clownsnase geben ich bin fucking " + fmt.Sprintf("<@%s>", target.ID),
		})
		return
	}

	if !db.UserInDatabase(target.ID) {
		err := db.InsertDBUser(target.ID, target.Username)
		if err != nil {
			slog.Error("inserting user into database", err)
		}
	}

	nase := db.Nase{
		ID:       snowflake.New(time.Now()),
		UserID:   target.ID,
		AuthorID: author.ID,
		Reason:   reason,
		Created:  time.Now(),
	}

	err := db.InsertNase(nase)
	if err != nil {
		slog.Error("inserting nase into database", err)
	}

	nasenCount, err := db.GetNasenCountForUser(target.ID)
	if err != nil {
		slog.Error("fetching nasen count for target user", err)
	}

	event.Client().Rest().CreateFollowupMessage(config.APP_ID, event.Token(), discord.MessageCreate{
		Embeds: []discord.Embed{
			{
				Title:       "Clownsnase Kassiert  🤡",
				Description: fmt.Sprintf("`Grund: %s`", reason),
				Color:       0x00ff00,
				Thumbnail: &discord.EmbedResource{
					URL: *target.AvatarURL(),
				},
				Fields: []discord.EmbedField{
					{
						Name:  "Von",
						Value: fmt.Sprintf("<@%s>", author.ID),
					}, {
						Name:  "An: ",
						Value: fmt.Sprintf("<@%s>", target.ID),
					}, {
						Name:  "Total: ",
						Value: fmt.Sprintf("%d", nasenCount),
					},
				},
			},
		},
	})

}
