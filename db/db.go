package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/nico-mayer/go_discordbot/config"
)

var DB *sql.DB

func Connect() error {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.PGUSER, config.PGPASSWORD, config.PGHOST, config.PGPORT, config.PGDATABASE)
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
