package db

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/nico-mayer/discordbot/config"
)

var DB *sql.DB

func init() {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", config.PGUSER, config.PGPASSWORD, config.PGHOST, config.PGPORT, config.PGDATABASE)
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	slog.Info("successfully connected to database.")
}
