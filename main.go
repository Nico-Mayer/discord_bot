package main

import (
	"log"

	"github.com/nico-mayer/go_discordbot/bot"
	"github.com/nico-mayer/go_discordbot/db"
)

func main() {
	err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	bot.Run()
}
