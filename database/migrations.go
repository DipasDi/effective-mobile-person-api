package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func RunMigrations() error {
	Dbconnect()
	defer Db.Close()

	driver, err := postgres.WithInstance(Db, &postgres.Config{})
	if err != nil {
		log.Warn().Err(err).Msg("postgres driver creation failed")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver)
	if err != nil {
		log.Info().Msg("Successfully enriched person data")
	}

	defer m.Close()

	if err := m.Up(); err != nil {
		fmt.Println(err)
	}
	return err
}
