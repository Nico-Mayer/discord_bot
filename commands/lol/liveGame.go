package lol

import (
	"fmt"
	"strings"
	"time"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/datadragon"
	"github.com/KnutZuidema/golio/riot/lol"
	"github.com/bwmarrin/discordgo"
	"github.com/nico-mayer/go_discordbot/utils"
)

const (
	RED_SIDE  = 100
	BLUE_SIDE = 200
)

func LiveGame(s *discordgo.Session, i *discordgo.InteractionCreate, golio *golio.Client) {
	var summonerName string

	for _, option := range i.ApplicationCommandData().Options {
		summonerName = option.StringValue()
	}

	summoner, err := golio.Riot.LoL.Summoner.GetByName(summonerName)
	if err != nil {
		fmt.Println(err)
		utils.ReplyError(s, i, fmt.Sprintf("Summoner `%s` not found!", summonerName))
		return
	}

	liveGame, err := golio.Riot.LoL.Spectator.GetCurrent(summoner.ID)
	if err != nil {
		fmt.Println(err)
		utils.ReplyError(s, i, fmt.Sprintf("`%s` is currently not in a game!", summonerName))
		return
	}

	replyLiveGame(s, i, golio, liveGame, summoner)

}

func replyLiveGame(s *discordgo.Session, i *discordgo.InteractionCreate, golio *golio.Client, liveGame *lol.GameInfo, summoner *lol.Summoner) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		fmt.Println(err)
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Type:  discordgo.EmbedTypeRich,
				Title: fmt.Sprintf("🔎 Live Game - %s", summoner.Name),
				Color: 0x90456f,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "🔵 Blue Team",
						Value:  getTeam(golio.DataDragon, liveGame, BLUE_SIDE),
						Inline: true,
					}, {
						Name:   "🔴 Red Team",
						Value:  getTeam(golio.DataDragon, liveGame, RED_SIDE),
						Inline: true,
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
}

func getTeam(dd *datadragon.Client, liveGame *lol.GameInfo, teamID int) string {
	var res string

	for _, p := range liveGame.Participants {
		if p.TeamID == teamID {
			var champName string

			champ, err := p.GetChampion(dd)
			if err != nil {
				fmt.Println(err)
				champName = "unknown"
			} else {
				champName = champ.Name
			}

			res += fmt.Sprintf("`%s` - [%s](https://www.op.gg/summoners/euw/%s-EUW) \n", champName, p.SummonerName, strings.ReplaceAll(p.SummonerName, " ", ""))
		}
	}
	return res
}

func getIngameTime(startTime int64) string {
	elapsedMilliseconds := time.Now().UnixNano()/int64(time.Millisecond) - startTime

	totalMinutes := elapsedMilliseconds / (1000 * 60)
	remainingSeconds := (elapsedMilliseconds / 1000) % 60

	ingameTime := fmt.Sprintf("%02d:%02d", totalMinutes, remainingSeconds)

	return ingameTime
}
