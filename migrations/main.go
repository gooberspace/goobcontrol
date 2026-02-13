package migrations

import (
	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

////go:embed migrations/*
//var migrationsFS embed.FS

func init() {
	//this is fine
}
