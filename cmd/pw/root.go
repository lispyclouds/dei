package pw

import (
	"context"
	"fmt"
	"slices"

	"github.com/lispyclouds/dei/pkg"
	"github.com/urfave/cli/v3"
)

func PwdCmd(cache *pkg.Cache) *cli.Command {
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
			&cli.BoolFlag{
				Name:  "to-clipboard",
				Usage: "Copy the password to clipboard instead of displaying",
				Value: false,
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			var (
				key       []byte
				identicon string
			)
			keyCacheKey, identiconCacheKey := "dei.password.mainKey", "dei.password.identicon"

			cachedKey, err := cache.Get(keyCacheKey)
			if err != nil {
				return err
			}

			cachedIdenticon, err := cache.Get(identiconCacheKey)
			if err != nil {
				return err
			}

			fullName := cmd.String("full-name")
			variant := SiteVariant(cmd.String("variant"))

			if cmd.Bool("flush-cache") || cachedKey == nil || cachedIdenticon == nil {
				mainPass, err := pkg.Input("Enter your main password", "", true)
				if err != nil {
					return err
				}

				identicon, err = identiconOf(fullName, mainPass)
				if err != nil {
					return err
				}

				key, err = mainKey(fullName, mainPass, variant)
				if err != nil {
					return err
				}

				if err = cache.Put(keyCacheKey, key); err != nil {
					return err
				}

				if err = cache.Put(identiconCacheKey, []byte(identicon)); err != nil {
					return err
				}
			} else {
				key = cachedKey
				identicon = string(cachedIdenticon)
			}

			pass, err := password(key, cmd.String("site"), cmd.Int("counter"), variant, TemplateClass(cmd.String("class")))
			if err != nil {
				return err
			}

			fmt.Println(identicon)

			if !cmd.Bool("to-clipboard") {
				fmt.Println(pass)
				return nil
			}

			return pkg.CopyToClipboard(pass)
		},
	}
}
