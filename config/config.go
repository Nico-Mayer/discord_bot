package config

import (
	"os"

	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
)

var (
	TOKEN    string
	APP_ID   snowflake.ID
	GUILD_ID snowflake.ID
)

func init() {
	godotenv.Load()

	TOKEN = os.Getenv("TOKEN")
	APP_ID = snowflake.GetEnv("APP_ID")
	GUILD_ID = snowflake.GetEnv("GUILD_ID")
}
