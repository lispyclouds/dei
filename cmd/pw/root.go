package pw

import (
	"context"
	"fmt"
	"slices"

	"github.com/lispyclouds/dei/pkg"
	"github.com/urfave/cli/v3"
)

var siteFlag *cli.StringFlag = &cli.StringFlag{
	Name:     "site",
	Usage:    "The site",
	Required: true,
}

func PwdCmd(cache *pkg.Cache) *cli.Command {
	return &cli.Command{
		Name:    "pw",
		Aliases: []string{"pwd", "pass", "password"},
		Usage:   "Stateless passwords",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "counter",
				Usage: "The current counter/generation of the credential",
				Value: 1,
				Validator: func(v int) error {
					if v <= 0 {
						return fmt.Errorf("counter must be 1 or greater")
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:  "variant",
				Usage: "The type of credentials to generate",
				Value: "password",
				Validator: func(v string) error {
					allowed := []string{"password", "login", "answer"}
					if !slices.Contains(allowed, v) {
						return fmt.Errorf("Choose from %v", allowed)
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:  "class",
				Usage: "The class of the credential generated",
				Value: "maximum",
				Validator: func(v string) error {
					allowed := []string{"maximum", "long", "medium", "basic", "short", "pin", "name", "phrase"}
					if !slices.Contains(allowed, v) {
						return fmt.Errorf("Choose from %v", allowed)
					}

					return nil
				},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "gen",
				Aliases: []string{"generate"},
				Usage:   "Generate a password",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "full-name",
						Usage:    "Your full name",
						Required: true,
					},
					siteFlag,
					&cli.BoolFlag{
						Name:  "to-clipboard",
						Usage: "Copy the password to clipboard instead of displaying",
						Value: false,
					},
					&cli.BoolFlag{
						Name:  "no-cache",
						Usage: "Ignore all cache reads and writes for this session",
						Value: false,
					},
					&cli.StringFlag{
						Name:  "cache-security-scheme",
						Usage: "The method used to secure the cached main password",
						Value: "pin",
						Validator: func(v string) error {
							allowed := []string{"pin"}
							if !slices.Contains(allowed, v) {
								return fmt.Errorf("Choose from %v", allowed)
							}

							return nil
						},
					},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					return generate(cache, cmd)
				},
			},
			{
				Name:  "site-cache",
				Usage: "Manage the site metadata cache",
				Commands: []*cli.Command{
					{
						Name:  "put",
						Usage: "Put/Edit values of a site",
						Flags: []cli.Flag{siteFlag},
						Action: func(_ context.Context, cmd *cli.Command) error {
							return cachePut(cache, cmd)
						},
					},
					{
						Name:  "remove",
						Usage: "Remove a site",
						Flags: []cli.Flag{siteFlag},
						Action: func(_ context.Context, cmd *cli.Command) error {
							return cacheRemove(cache, cmd)
						},
					},
					{
						Name:  "show",
						Usage: "Show all cached site metadata",
						Action: func(_ context.Context, cmd *cli.Command) error {
							return cacheShow(cache)
						},
					},
				},
			},
		},
	}
}
