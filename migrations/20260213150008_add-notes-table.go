package migrations

import (
	"context"

	"github.com/gooberspace/goobcontrol/internal/models"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		if _, error := db.NewCreateTable().
			Model((*models.Note)(nil)).
			IfNotExists().
			Exec(ctx); error != nil {
			return error
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		if _, error := db.NewDropTable().
			Model((*models.Note)(nil)).
			Exec(ctx); error != nil {
			return error
		}
		return nil
	})
}
