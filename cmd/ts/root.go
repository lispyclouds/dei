package ts

import (
	"github.com/lispyclouds/dei/pkg"
	"github.com/urfave/cli/v3"
)

func TsCmd(cache *pkg.Cache) *cli.Command {
	return &cli.Command{
		Name:            "treesitter",
		Usage:           "Babysit treesitter",
		Aliases:         []string{"ts"},
		CommandNotFound: pkg.CommandNotFound,
		Commands: []*cli.Command{
			{
				Name:            "sync",
				Usage:           "Manage treesitter artifacts",
				CommandNotFound: pkg.CommandNotFound,
				Action:          syncCmd,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "conf",
						Aliases:  []string{"c"},
						Usage:    "The config file",
						Required: true,
					},
					&cli.StringFlag{
						Name:  "cli",
						Usage: "The path to tree-sitter cli",
						Value: "tree-sitter",
					},
				},
			},
		},
	}
}
