package main

import (
	"context"
	"log"
	"os"

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
			pw.Cmd(cache),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		cache.Close()
		log.Fatal(err)
	}
}
