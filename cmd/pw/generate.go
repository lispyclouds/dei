package pw

import (
	"bytes"
	json "encoding/json/v2"
	"fmt"
	"log/slog"

	"github.com/lispyclouds/dei/pkg"
	"github.com/urfave/cli/v3"
)

type SiteInfo struct {
	Counter int           `json:"counter"`
	Class   TemplateClass `json:"class"`
	Variant SiteVariant   `json:"variant"`
}

type Sites = map[string]SiteInfo

const sitesCacheKey = "dei.password.sites"

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
	buffer := bytes.NewBuffer([]byte{})
	if err := json.MarshalWrite(buffer, &sites); err != nil {
		return err
	}

	return cache.Put(sitesCacheKey, buffer.Bytes())
}

func generate(cache *pkg.Cache, cmd *cli.Command) error {
	var (
		key       []byte = nil
		identicon string
	)
	keyCacheKey := "dei.password.mainKey"
	identiconCacheKey := "dei.password.identicon"
	sites := Sites{}
	noCache := cmd.Bool("no-cache")

	if !noCache {
		cachedKey, err := cache.Get(keyCacheKey)
		if err != nil {
			return err
		}

		if cachedKey != nil {
			key = cachedKey
		}

		cachedIdenticon, err := cache.Get(identiconCacheKey)
		if err != nil {
			return err
		}

		if cachedIdenticon != nil {
			identicon = string(cachedIdenticon)
		}

		cachedSites, err := cache.Get(sitesCacheKey)
		if err != nil {
			return err
		}

		if cachedSites != nil {
			if err = json.UnmarshalRead(bytes.NewReader(cachedSites), &sites); err != nil {
				return err
			}
		}
	}

	fullName := cmd.String("full-name")
	site := cmd.String("site")
	variant := getVariant(noCache, site, sites, cmd)
	class := getClass(noCache, site, sites, cmd)
	counter := getCounter(noCache, site, sites, cmd)

	if cmd.Bool("flush-cache") || key == nil || len(identicon) == 0 {
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

		if !noCache {
			if err = cache.
				WithWriteTxn().
				Put(keyCacheKey, key).
				Put(identiconCacheKey, []byte(identicon)).
				Run(); err != nil {
				return err
			}
		}
	}

	pass, err := password(key, site, counter, variant, class)
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

	fmt.Println(identicon)

	if !cmd.Bool("to-clipboard") {
		fmt.Println(pass)
		return nil
	}

	return pkg.CopyToClipboard(pass)
}
