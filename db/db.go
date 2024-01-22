package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/nico-mayer/go_discordbot/config"
)

var DB *sql.DB

func Connect() error {
	connStr := "user=" + config.PGUSER + " password=" + config.PGPASSWORD + " dbname=" + config.PGDATABASE + " host=" + config.PGHOST + " sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	log.Println("Successfully connected to database")
	return nil
}
