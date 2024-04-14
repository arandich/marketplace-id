package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/arandich/marketplace-id/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Config{}
	err := config.New(ctx, &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	connString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=720",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)

	db, err := sql.Open("postgres", connString)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrator/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalln(err)
	}
	err = m.Up()
	if err != nil {
		return
	}
}
