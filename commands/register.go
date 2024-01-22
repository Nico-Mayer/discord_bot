package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/nico-mayer/go_discordbot/config"
)

type Collection struct {
	Name     string
	Icon     rune
	Commands []*discordgo.ApplicationCommand
}

var Collections []Collection

func RegisterCommands(session *discordgo.Session) {
	var err error

	general := Collection{
		Name: "General",
		Icon: '📜',
		Commands: []*discordgo.ApplicationCommand{
			{
				Name:        "help",
				Description: "Zeigt liste aller Commands",
			},
			{
				Name:        "user",
				Description: "Zeigt Informationen über einen User an.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "bre",
						Description: "Wähle den bre aus.",
						Required:    true,
					},
				},
			},
		},
	}

	nasen := Collection{
		Name: "Nasen",
		Icon: '👃',
		Commands: []*discordgo.ApplicationCommand{
			{
				Name:        "clownsnase",
				Description: "Gib einem bre eine Clownsnase.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "bre",
						Description: "Wähle den bre aus der unfassbar müffelt.",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "grund",
						Description: "Warum müffelt derjenige so hart?",
						MaxLength:   64,
						Required:    false,
					},
				},
			},
			{
				Name:        "clownfiesta",
				Description: "Gib allen bres im Channel eine Clownsnase!",
			},
			{
				Name:        "leaderboard",
				Description: "Zeigt Clownsnasen Leaderboard!",
			},
			{
				Name:        "nasen",
				Description: "Zeigt alle Nasen eines Users an.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionUser,
						Name:        "bre",
						Description: "Wähle den bre dessen clownsnasen du sehen willst",
						Required:    true,
					},
				},
			},
		},
	}

	music := Collection{
		Name: "Music",
		Icon: '🎵',
		Commands: []*discordgo.ApplicationCommand{
			{
				Name:        "play",
				Description: "Spielt einen Song ab.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "url",
						Description: "Suche nach einem Song.",
						Required:    true,
					},
				},
			}, {
				Name:        "stop",
				Description: "Stoppt DJ Rosine",
			},
		},
	}

	lol := Collection{
		Name: "League of Legends",
		Icon: '🎮',
		Commands: []*discordgo.ApplicationCommand{
			{
				Name:        "live_game",
				Description: "Zeigt Informationen über das aktuelle Spiel an.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "summoner",
						Description: "Gib den Summoner Namen an.",
						Required:    true,
					},
				},
			},
		},
	}

	Collections = append(Collections, general)
	Collections = append(Collections, nasen)
	Collections = append(Collections, music)
	Collections = append(Collections, lol)

	var commands []*discordgo.ApplicationCommand
	for _, collection := range Collections {
		commands = append(commands, collection.Commands...)
	}

	_, err = session.ApplicationCommandBulkOverwrite(config.APP_ID, config.GUILD_ID, commands)

	if err != nil {
		log.Fatal(err)
	}
}
