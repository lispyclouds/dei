package commit

import (
	"context"

	"github.com/lispyclouds/dei/pkg"
	"github.com/urfave/cli/v3"
)

func CommitCmd(cache *pkg.Cache) *cli.Command {
	return &cli.Command{
		Name:  "commit",
		Usage: "Committed companion",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "short",
				Usage: "Set this to omit the extended commit message",
				Value: false,
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			return commit(cache, cmd)
		},
		Commands: []*cli.Command{
			{
				Name:  "co-authors",
				Usage: "Manage co-authors",
				Commands: []*cli.Command{
					{
						Name:  "add",
						Usage: "Add a co-author",
						Action: func(_ context.Context, cmd *cli.Command) error {
							return manageCoAuthor(cache, cmd.String("name"), cmd.String("email"), "add")
						},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Usage:    "The name of the co-author",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "email",
								Usage:    "The email of the co-author",
								Required: true,
							},
						},
					},
					{
						Name:  "remove",
						Usage: "Remove a co-author",
						Action: func(_ context.Context, cmd *cli.Command) error {
							return manageCoAuthor(cache, cmd.String("name"), cmd.String("email"), "remove")
						},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "email",
								Usage:    "The email of the co-author",
								Required: true,
							},
						},
					},
					{
						Name:  "list",
						Usage: "List all co-authors",
						Action: func(_ context.Context, _ *cli.Command) error {
							return listCoAuthors(cache)
						},
					},
				},
			},
		},
	}
}
