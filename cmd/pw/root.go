package pw

import (
	"context"
	"fmt"
	"slices"

	"github.com/charmbracelet/huh"
	"github.com/urfave/cli/v3"
)

func Cmd() *cli.Command {
	return &cli.Command{
		Name:    "pw",
		Aliases: []string{"pwd", "pass", "password"},
		Usage:   "Stateless passwords",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "full-name",
				Usage:    "Your full name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "site",
				Usage:    "The site for the password",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "variant",
				Usage: "The kind of credentials to generate",
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
				Name:  "kind",
				Usage: "The kind of the password",
				Value: "maximum",
				Validator: func(v string) error {
					allowed := []string{"maximum", "long", "medium", "basic", "short", "pin", "name", "phrase"}
					if !slices.Contains(allowed, v) {
						return fmt.Errorf("Choose from %v", allowed)
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:  "context",
				Usage: "Optional, useful for variant answer. Empty for a universal site answer or the most significant word(s) of the question",
				Value: "",
			},
			&cli.IntFlag{
				Name:  "counter",
				Usage: "The counter of the current password for the site",
				Value: 1,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			var mainPass string

			if err := huh.NewInput().
				Title("Enter your main password").
				Value(&mainPass).
				EchoMode(huh.EchoModePassword).
				Run(); err != nil {
				return err
			}

			fullName := cmd.String("full-name")

			identicon, err := identiconOf(fullName, mainPass)
			if err != nil {
				return err
			}

			variant := SiteVariant(cmd.String("variant"))

			mainKey, err := mainKey(fullName, mainPass, variant)
			if err != nil {
				return err
			}

			password, err := password(mainKey, cmd.String("site"), cmd.String("context"), cmd.Int("counter"), variant, Kind(cmd.String("kind")))
			if err != nil {
				return err
			}

			fmt.Println(identicon, password)

			return nil
		},
	}
}
