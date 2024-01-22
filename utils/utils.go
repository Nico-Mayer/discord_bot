package utils

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func ReplyError(s *discordgo.Session, i *discordgo.InteractionCreate, e error, msg string) {
	log.Println(e)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: msg,
		},
	})
	Check(err)
}

func Check(e error) {
	if e != nil {
		log.Println(e)
	}
}
