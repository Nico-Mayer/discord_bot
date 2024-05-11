package bot

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"

	"github.com/nico-mayer/go_discordbot/commands"
	"github.com/nico-mayer/go_discordbot/commands/general"
	"github.com/nico-mayer/go_discordbot/commands/music"
	"github.com/nico-mayer/go_discordbot/commands/nasen"
	"github.com/nico-mayer/go_discordbot/config"
	"github.com/nico-mayer/go_discordbot/levels"
	"github.com/nico-mayer/go_discordbot/player"
	"github.com/nico-mayer/go_discordbot/utils"
)

var session *discordgo.Session

func Run() {
	session, _ = discordgo.New("Bot " + config.TOKEN)

	// Register Commands
	commands.RegisterCommands(session)

	// Init needed Clients
	player := player.NewPlayer(session)

	session.AddHandler(func(
		s *discordgo.Session,
		i *discordgo.InteractionCreate,
	) {
		data := i.ApplicationCommandData()

		switch data.Name {
		case "help":
			general.Help(s, i)
		case "user":
			general.User(s, i)
		case "clownsnase":
			nasen.Clownsnase(s, i)
		case "clownfiesta":
			nasen.Clownfiesta(s, i)
		case "leaderboard":
			nasen.Leaderboard(s, i)
		case "nasen":
			nasen.Nasen(s, i)
		case "play":
			music.Play(s, i, player)
		case "stop":
			music.Stop(s, i, player)
		case "skip":
			music.Skip(s, i, player)
		case "pause":
			music.Pause(s, i, player)
		case "resume":
			music.Resume(s, i, player)
		default:
			return
		}

	})

	levels.Init(session)

	err := session.Open()
	if err != nil {
		log.Fatal(err)
	}

	err = session.UpdateGameStatus(0, "mit seinem Zipfel")
	utils.Check(err)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	err = session.Close()
	if err != nil {
		log.Fatal(err)
	}
}
