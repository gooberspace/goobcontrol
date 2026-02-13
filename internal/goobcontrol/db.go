package goobcontrol

import (
	"context"
	"database/sql"
	"time"

	"github.com/gooberspace/goobcontrol/migrations"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
)

func SetupDatabase(config *viper.Viper) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithAddr(config.GetString("database.host")),
		pgdriver.WithUser(config.GetString("database.user")),
		pgdriver.WithPassword(config.GetString("database.password")),
		pgdriver.WithDatabase(config.GetString("database.database")),
		pgdriver.WithInsecure(config.GetBool("database.insecure")),
	))

	sqldb.SetMaxOpenConns(25)
	sqldb.SetMaxIdleConns(10)
	sqldb.SetConnMaxLifetime(5 * time.Minute)
	sqldb.SetConnMaxIdleTime(5 * time.Minute)

	db := bun.NewDB(sqldb, pgdialect.New())
	return db
}

func (gc *GoobControl) TestDatabase() error {
	if err := gc.DB.Ping(); err != nil {
		//
		return err
	} else {
		return nil
	}
}

func SetupDbMigrator(db *bun.DB) *migrate.Migrator {
	return migrate.NewMigrator(db, migrations.Migrations)
}

func (gc *GoobControl) RunDbMigrations() error {
	migrationCtx := context.TODO()
	var migrationErr error
	if migrationErr = gc.DbMigrator.Init(migrationCtx); migrationErr != nil {
		return migrationErr
	}
	if migrationErr = gc.DbMigrator.Lock(migrationCtx); migrationErr != nil {
		return migrationErr
	}
	defer gc.DbMigrator.Unlock(migrationCtx)

	var migrationGroup *migrate.MigrationGroup
	migrationGroup, migrationErr = gc.DbMigrator.Migrate(migrationCtx)
	if migrationErr != nil {
		// return nil if there have not been any migrations yet, this is fine
		if migrationErr.Error() == "migrate: there are no migrations" {
			return nil
		}
		return migrationErr
	}
	if migrationGroup.IsZero() {
		gc.Logger.Info("No Database migrations to run")
		return nil
	}
	gc.Logger.Info("Database migration done:", "group", migrationGroup.ID, "migrations", migrationGroup.Migrations.String())
	return nil

}
