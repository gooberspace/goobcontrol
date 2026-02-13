package migrationfiles

import (
	"context"
	"fmt"

	"github.com/gooberspace/goobcontrol/internal/migrations"
	"github.com/uptrace/bun"
)

func init() {
	migrations.Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
