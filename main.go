package main

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Repo-PR-Stat",
		Usage: "show pull request stat",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "owner",
				Aliases:  []string{"o"},
				Usage:    "repository owner",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "repository",
				Aliases:  []string{"r"},
				Usage:    "repository name",
				Required: true,
			},
			&cli.TimestampFlag{
				Name:     "start",
				Aliases:  []string{"s"},
				Usage:    "start time",
				Layout:   time.RFC3339,
				Required: true,
			},
			&cli.TimestampFlag{
				Name:     "end",
				Aliases:  []string{"e"},
				Usage:    "end time",
				Layout:   time.RFC3339,
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     "include-base",
				Usage:    "include base branch name",
				Required: false,
			},
			&cli.StringSliceFlag{
				Name:     "exclude-base",
				Usage:    "exclude base branch name",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "token",
				Aliases:  []string{"t"},
				Usage:    "GitHub Access Token",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			owner := ctx.String("owner")
			repository := ctx.String("repository")
			start := ctx.Timestamp("start")
			end := ctx.Timestamp("end")
			includeBases := ctx.StringSlice("include-base")
			excludeBases := ctx.StringSlice("exclude-base")

			token := ctx.String("token")
			if len(token) == 0 {
				token = os.Getenv("GITHUB_ACCESS_TOKEN")
			}
			if len(token) == 0 {
				return errors.New("provide Github Access Token")
			}

			err := showStatAsJson(owner, repository, *start, *end, includeBases, excludeBases, token)
			return err
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
