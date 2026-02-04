package goobcontrol

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
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

func (gc *GoobControl) TestDatabase() {
	if err := gc.DB.Ping(); err != nil {
		gc.Logger.Error("Failed to connect to database", slog.Any("err", err))
	} else {
		gc.Logger.Info("Database connection successful")
	}
}
