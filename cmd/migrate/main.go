package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/config"
	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/adapters/postgres/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.Load()
	db, err := sql.Open("pgx", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	goose.SetBaseFS(migrations.EmbedFS)

	if err := goose.SetDialect("postgres"); err != nil {
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
