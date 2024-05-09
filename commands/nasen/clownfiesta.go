package nasen

import (
	"log/slog"
	"sync"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/nico-mayer/discordbot/config"
	"github.com/nico-mayer/discordbot/db"
)

var ClownfiestaCommand = discord.SlashCommandCreate{
	Name:        "clownfiesta",
	Description: "Gib allen bres im Voice Channel eine Clownsnase",
}

func ClownfiestaCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	//data := event.SlashCommandInteractionData()

	voiceState, ok := event.Client().Caches().VoiceState(config.GUILD_ID, event.User().ID)
	if !ok {
		event.CreateMessage(discord.MessageCreate{
			Content: "You need to be in a voice channel to use this command",
		})
		return
	}

	voiceChannel, ok := event.Client().Caches().GuildAudioChannel(*voiceState.ChannelID)
	if !ok {
		event.CreateMessage(discord.MessageCreate{
			Content: "Voice channel is not existing",
		})
		return
	}

	usersInChannel := event.Client().Caches().AudioChannelMembers(voiceChannel)
	event.DeferCreateMessage(false)

	author := event.User()

	var wg sync.WaitGroup
	errChan := make(chan error, len(usersInChannel))

	for _, user := range usersInChannel {
		wg.Add(1)

		go func(user discord.Member) {
			defer wg.Done()

			var nase db.Nase = db.Nase{
				ID:       snowflake.New(time.Now()),
				UserID:   user.User.ID,
				AuthorID: author.ID,
				Reason:   "Clownfiesta 🤡",
				Created:  time.Now(),
			}

			errChan <- db.InsertNase(nase)
		}(user)

	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			slog.Error("inserting nase for user in clownfiesta loop", err)
		}
	}

	event.Client().Rest().CreateFollowupMessage(config.APP_ID, event.Token(), discord.MessageCreate{
		Content: "Clownfiesta",
	})

}
