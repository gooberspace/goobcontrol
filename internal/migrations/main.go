package migrations

import (
	"os"

	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

////go:embed migrations/*
//var migrationsFS embed.FS

func init() {
	if err := Migrations.Discover(os.DirFS("migrationfiles")); err != nil {
		panic(err)
	}
}
