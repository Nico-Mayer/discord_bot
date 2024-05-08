package nasen

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var ClownfiestaCommand = discord.SlashCommandCreate{
	Name:        "clownfiesta",
	Description: "Gib allen bres im Voice Channel eine Clownsnase",
}

func ClownfiestaCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	//data := event.SlashCommandInteractionData()

	_, ok := event.Client().Caches().VoiceState(*event.GuildID(), event.User().ID)
	if !ok {
		event.CreateMessage(discord.MessageCreate{
			Content: "You need to be in a voice channel to use this command",
		})
	}

	event.DeferCreateMessage(false)
}
