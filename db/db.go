package db

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func init() {
	if os.Getenv("ENV") != "PROD" {
		slog.Info("running on stage enviroment, loading env variables from .env file")
		godotenv.Load()
	}

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGDATABASE"))
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
