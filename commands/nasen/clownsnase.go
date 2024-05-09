package nasen

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	mybot "github.com/nico-mayer/discordbot/bot"
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

func ClownsnaseCommandHandler(event *events.ApplicationCommandInteractionCreate, b *mybot.Bot) error {
	data := event.SlashCommandInteractionData()

	author := event.User()
	target := data.User("user")
	reason := data.String("reason")

	if target.Bot {
		return event.CreateMessage(discord.MessageCreate{
			Flags:   discord.MessageFlagEphemeral,
			Content: "Du kannst mir keine clownsnase geben ich bin fucking " + fmt.Sprintf("<@%s>", target.ID),
		})
	}

	event.DeferCreateMessage(false)
	if !db.UserInDatabase(target.ID) {
		err := db.InsertDBUser(target.ID, target.Username)
		if err != nil {
			return err
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
		return err
	}

	nasenCount, err := db.GetNasenCountForUser(target.ID)
	if err != nil {
		return err
	}

	_, err = event.Client().Rest().CreateFollowupMessage(config.APP_ID, event.Token(), discord.MessageCreate{
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
	return err
}
