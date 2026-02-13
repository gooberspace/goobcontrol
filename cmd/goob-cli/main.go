package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gooberspace/goobcontrol/internal/goobcontrol"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "goob-cli",
		Usage: "CLI for managing Goob Control",
		Commands: []*cli.Command{
			{
				Name:  "db",
				Usage: "Commands for managing the database",
				Commands: []*cli.Command{
					{
						Name:    "create-migration",
						Aliases: []string{"cm"},
						Usage:   "Create a new database migration",
						Action:  createMigration,
					},
				},
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func createMigration(ctx context.Context, c *cli.Command) error {
	name := strings.Join(c.Args().Slice(), "_")
	config := goobcontrol.CreateConfig()
	db := goobcontrol.SetupDatabase(config)
	migrator := goobcontrol.SetupDbMigrator(db)
	mf, err := migrator.CreateGoMigration(ctx, name)
	if err != nil {
		return err
	}
	fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
	return nil
}
