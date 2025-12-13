package main

import (
	"context"
	"database/sql"

	"log"
	"os"

	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/config"
	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/adapters/sqlite/migrations"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func main() {
	cfg := config.Load()
	db, err := sql.Open("sqlite", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.EmbedFS)

	if err := goose.SetDialect("sqlite"); err != nil {
		log.Fatal(err)
	}

	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	if err := goose.RunContext(context.Background(), command, db, "."); err != nil {
		log.Fatal(err)
	}
}
