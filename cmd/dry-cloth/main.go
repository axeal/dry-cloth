package main

import (
	"context"
	"log"
	"os"

	"github.com/axeal/dry-cloth/pkg/drycloth"
	"github.com/urfave/cli/v2"
)

func main() {

	var accessToken string
	var preserveTag string
	var maxAgeDays int
	var dryRun bool

	app := &cli.App{
		Name:  "dry-cloth",
		Usage: "clean up old droplets",
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "access-token",
			Usage:       "Digital Ocean access token",
			EnvVars:     []string{"DIGITALOCEAN_ACCESS_TOKEN"},
			Destination: &accessToken,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "preserve-tag",
			Usage:       "Tag to prevent droplet deletion",
			Destination: &preserveTag,
		},
		&cli.IntFlag{
			Name:        "max-age-days",
			Usage:       "Maximum age of droplets to keep",
			Destination: &maxAgeDays,
			Value:       14,
		},
		&cli.BoolFlag{
			Name:        "dry-run",
			Usage:       "Dry run without deleting droplets",
			Destination: &dryRun,
		},
	}
	app.Action = func(c *cli.Context) error {
		ctx := context.TODO()
		err := drycloth.Run(ctx, accessToken, preserveTag, maxAgeDays, dryRun)
		return err
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
