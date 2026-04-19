package ts

import (
	"context"
	json "encoding/json/v2"
	"errors"
	"os"
	"path/filepath"

	"github.com/lispyclouds/dei/pkg"
	"github.com/urfave/cli/v3"
)

func action(_ context.Context, cmd *cli.Command) error {
	if !pkg.Which("git") {
		return errors.New("Cannot find git on the PATH")
	}

	if !pkg.Which("tree-sitter") {
		return errors.New("Cannot find tree-sitter on the PATH")
	}

	f, err := os.Open(cmd.String("conf"))
	if err != nil {
		return err
	}
	defer f.Close()

	var conf Conf
	if err := json.UnmarshalRead(f, &conf); err != nil {
		return err
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	cacheDir = filepath.Join(cacheDir, "dei", "ts")

	if err = syncQueries(conf, filepath.Join(cacheDir, "queries")); err != nil {
		return err
	}

	return syncParsers(conf, cacheDir, cmd.String("cli"))
}

func TsCmd() *cli.Command {
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
				Action:          action,
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
