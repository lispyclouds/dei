package cmd

import (
	"bytes"
	"context"
	json "encoding/json/v2"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/lispyclouds/dei/pkg"
	"github.com/urfave/cli/v3"
)

const coAuthorsKey = "dei.commit.coAuthors"

type CoAuthors = map[string]map[string]string

func loadCoAuthors(cache *pkg.Cache) (CoAuthors, error) {
	var coAuthors CoAuthors

	data, err := cache.Get(coAuthorsKey)
	if err != nil {
		return nil, err
	}

	if data == nil {
		coAuthors = make(CoAuthors)
	} else {
		if err = json.UnmarshalRead(bytes.NewReader(data), &coAuthors); err != nil {
			return nil, err
		}
	}

	return coAuthors, nil
}

func manageCoAuthor(cache *pkg.Cache, name, email, op string) error {
	coAuthors, err := loadCoAuthors(cache)
	if err != nil {
		return err
	}

	switch op {
	case "add":
		info, ok := coAuthors[email]
		if !ok {
			info = make(map[string]string)
		}

		info["name"] = name
		coAuthors[email] = info
	case "remove":
		delete(coAuthors, email)
	}

	buffer := bytes.NewBuffer([]byte{})
	if err = json.MarshalWrite(buffer, &coAuthors); err != nil {
		return err
	}

	return cache.Put(coAuthorsKey, buffer.Bytes())
}

func CommitCmd(cache *pkg.Cache) *cli.Command {
	return &cli.Command{
		Name:  "commit",
		Usage: "Committed companion",
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
						Action: func(_ context.Context, cmd *cli.Command) error {
							coAuthors, err := loadCoAuthors(cache)
							if err != nil {
								return err
							}

							if len(coAuthors) == 0 {
								return nil
							}

							w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)

							fmt.Fprintln(w, "Name\tEmail")
							for email, info := range coAuthors {
								fmt.Fprintf(w, "%s\t%s\n", info["name"], email)
							}
							w.Flush()

							return nil
						},
					},
				},
			},
		},
	}
}
