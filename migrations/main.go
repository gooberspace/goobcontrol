package migrations

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

//go:embed migrations/*
var migrationsFS embed.FS

func init() {
	if err := Migrations.Discover(migrationsFS); err != nil {
		panic(err)
	}
}
