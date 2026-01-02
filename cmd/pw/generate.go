package pw

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lispyclouds/dei/pkg"
	"github.com/urfave/cli/v3"
)

type SiteInfo struct {
	Counter int           `json:"counter"`
	Class   TemplateClass `json:"class"`
	Variant SiteVariant   `json:"variant"`
}

type Sites = map[string]SiteInfo

func getVariant(noCache bool, site string, sites Sites, cmd *cli.Command) SiteVariant {
	variant := SiteVariant(cmd.String("variant"))

	if cmd.IsSet("variant") {
		return variant
	}

	if !noCache {
		info, ok := sites[site]
		if ok {
			return info.Variant
		}
	}

	return variant
}

func getClass(noCache bool, site string, sites Sites, cmd *cli.Command) TemplateClass {
	class := TemplateClass(cmd.String("class"))

	if cmd.IsSet("class") {
		return class
	}

	if !noCache {
		info, ok := sites[site]
		if ok {
			return info.Class
		}
	}

	return class
}

func getCounter(noCache bool, site string, sites Sites, cmd *cli.Command) int {
	counter := cmd.Int("counter")

	if cmd.IsSet("counter") {
		return counter
	}

	if !noCache {
		info, ok := sites[site]
		if ok {
			return info.Counter
		}
	}

	return counter
}

func cacheSite(cache *pkg.Cache, site string, sites Sites, info SiteInfo) error {
	var updateNeeded bool

	currentInfo, ok := sites[site]
	switch {
	case !ok:
		updateNeeded = true
	case currentInfo != info:
		updateNeeded = true
		slog.Info("Site info changed, updating", "site", site, "prev", currentInfo, "next", info)
	}

	if !updateNeeded {
		return nil
	}

	sites[site] = info

	return saveSites(cache, sites)
}

func onlyHosts(site string) string {
	parsed, err := url.Parse(site)
	if err != nil || len(parsed.Hostname()) == 0 {
		slog.Warn("Cannot parse hostname, using as is", "site", site)
		return site
	}

	return strings.TrimPrefix(parsed.Hostname(), "www.")
}

func generate(cache *pkg.Cache, cmd *cli.Command) error {
	var (
		mainKey   []byte = nil
		identicon string
	)
	mainKeyCacheKey := "dei.password.mainKey"
	identiconCacheKey := "dei.password.identicon"
	sites := Sites{}
	noCache := cmd.Bool("no-cache")
	cacheSecurityScheme := cmd.String("cache-security-scheme")
	flushCache := cmd.Bool("flush-cache")
	var cryptoKey []byte

	site := onlyHosts(strings.TrimSpace(cmd.String("site")))
	slog.Info("Generating for", "site", site)

	if !noCache {
		cachedMainKey, err := cache.Get(mainKeyCacheKey)
		if err != nil {
			return err
		}

		if !flushCache && cachedMainKey != nil {
			if cacheSecurityScheme == "pin" {
				pin, err := pkg.Input("Enter the PIN", "", true)
				if err != nil {
					return err
				}

				cryptoKey = []byte(pin)

				mainKey, err = decrypt(cachedMainKey, cryptoKey)
				if err != nil {
					return fmt.Errorf("%s: either wrong PIN or data is corrupted", err)
				}
			}
		}

		cachedIdenticon, err := cache.Get(identiconCacheKey)
		if err != nil {
			return err
		}

		if cachedIdenticon != nil {
			identicon = string(cachedIdenticon)
		}

		sites, err = loadSites(cache)
		if err != nil {
			return err
		}
	}

	fullName := strings.TrimSpace(cmd.String("full-name"))
	variant := getVariant(noCache, site, sites, cmd)
	class := getClass(noCache, site, sites, cmd)
	counter := getCounter(noCache, site, sites, cmd)

	if flushCache || mainKey == nil || len(identicon) == 0 {
		mainPass, err := pkg.Input("Enter your main password", "", true)
		if err != nil {
			return err
		}

		identicon, err = identiconOf(fullName, mainPass)
		if err != nil {
			return err
		}

		mainKey, err = mainKeyOf(fullName, mainPass, variant)
		if err != nil {
			return err
		}

		if !noCache {
			if cacheSecurityScheme == "pin" {
				if cryptoKey == nil {
					pin, err := pkg.Input("Enter the PIN", "", true)
					if err != nil {
						return err
					}

					pin_again, err := pkg.Input("Re-enter the PIN", "", true)
					if err != nil {
						return err
					}

					if pin != pin_again {
						return errors.New("Both pins don't match")
					}

					cryptoKey = []byte(pin)
				}

				secureKey, err := encrypt(mainKey, cryptoKey)
				if err != nil {
					return err
				}

				if err = cache.
					WithWriteTxn().
					Put(mainKeyCacheKey, secureKey).
					Put(identiconCacheKey, []byte(identicon)).
					Run(); err != nil {
					return err
				}
			}
		}
	}

	pass, err := derivePass(mainKey, site, counter, variant, class)
	if err != nil {
		return err
	}

	if !noCache {
		if err = cacheSite(
			cache,
			site,
			sites,
			SiteInfo{Counter: counter, Class: class, Variant: variant},
		); err != nil {
			return err
		}
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Site: %s\nIdenticon: %s", site, identicon))

	if !cmd.Bool("to-clipboard") {
		sb.WriteString(fmt.Sprintf("\nPassword: %s", pass))
	} else {
		if err = pkg.CopyToClipboard(pass); err != nil {
			return err
		}

		sb.WriteString("\nPassword copied to clipboard")
	}

	fmt.Println(
		lipgloss.NewStyle().
			Width(35).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Render(sb.String()),
	)

	return nil
}
