package levels

import (
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/snowflake/v2"
	"github.com/nico-mayer/discordbot/db"
)

const (
	EXP_PER_MESSAGE         int = 10
	EXP_PER_VOICE_JOIN      int = 25
	EXP_NEEDED_FOR_LEVEL_UP int = 250
)

var levelMapping = map[int]string{
	5:  "1196937224634257479",
	15: "1196936971776438302",
	25: "1196936450306998322",
	40: "696672218117177437",
	60: "696671960507351100",
}

func GrantExpToUser(userId snowflake.ID, username string, exp int) (int, bool, error) {
	if !db.UserInDatabase(userId) {
		err := db.InsertDBUser(userId, username)
		if err != nil {
			return 0, false, err
		}
	}

	dbUser, err := db.GetUser(userId)
	if err != nil {
		return 0, false, err
	}

	err = dbUser.GrantExp(exp)
	if err != nil {
		return 0, false, err
	}

	oldLevel := CalcUserLevel(dbUser.Exp)
	newLevel := CalcUserLevel(dbUser.Exp + exp)
	levelUp := oldLevel != newLevel

	return newLevel, levelUp, err
}

func CalcUserLevel(exp int) int {
	if exp < EXP_NEEDED_FOR_LEVEL_UP {
		return 1
	}

	return exp / EXP_NEEDED_FOR_LEVEL_UP
}

func HandleLevelUp(botClient bot.Client, userId snowflake.ID, level int) {

	newRank := levelMapping[level]
	if newRank != "" {
		err := botClient.Rest().AddMemberRole(snowflake.GetEnv("GUILD_ID"), userId, snowflake.MustParse(newRank))
		if err != nil {
			slog.Error("giving role to user", err)
		}
		return
	}
}
