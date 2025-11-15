package main

import (
	"context"
	"log"
	"os"

	"github.com/lispyclouds/dei/cmd"
	"github.com/lispyclouds/dei/cmd/pw"
	"github.com/lispyclouds/dei/pkg"
	"github.com/urfave/cli/v3"
)

func main() {
	cache, err := pkg.NewCache()
	if err != nil {
		log.Fatal(err)
	}
	defer cache.Close()

	cmd := cli.Command{
		Name:  "dei",
		Usage: "me in the CLI",
		Commands: []*cli.Command{
			pw.PwdCmd(cache),
			cmd.CommitCmd(cache),
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "flush-cache",
				Usage: "Ignore current cache and refresh values",
				Value: false,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		cache.Close()
		log.Fatal(err)
	}
}
