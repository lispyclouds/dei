package main

import (
	"context"
	"log"
	"os"

	"github.com/lispyclouds/dei/cmd/pw"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := cli.Command{
		Name:  "dei",
		Usage: "me in the CLI",
		Commands: []*cli.Command{
			pw.Cmd(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
