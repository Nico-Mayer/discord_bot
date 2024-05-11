package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	TOKEN        string
	APP_ID       string
	GUILD_ID     string
	PGUSER       string
	PGPASSWORD   string
	PGDATABASE   string
	PGHOST       string
	PGPORT       string
	NODE_ADDRESS string
	NODE_PW      string
)

func init() {
	godotenv.Load()

	TOKEN = os.Getenv("TOKEN")
	APP_ID = os.Getenv("APP_ID")
	GUILD_ID = os.Getenv("GUILD_ID")
	PGUSER = os.Getenv("PGUSER")
	PGPASSWORD = os.Getenv("PGPASSWORD")
	PGDATABASE = os.Getenv("PGDATABASE")
	PGHOST = os.Getenv("PGHOST")
	PGPORT = os.Getenv("PGPORT")
	NODE_ADDRESS = os.Getenv("NODE_ADDRESS")
	NODE_PW = os.Getenv("NODE_PW")
}
