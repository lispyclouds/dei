package pw

import (
	"bytes"
	"encoding/json/jsontext"
	json "encoding/json/v2"
	"log/slog"
	"os"

	"github.com/lispyclouds/dei/pkg"
	"github.com/urfave/cli/v3"
)

const sitesCacheKey = "dei.password.sites"

func loadSites(cache *pkg.Cache) (Sites, error) {
	cachedSites, err := cache.Get(sitesCacheKey)
	if err != nil {
		return nil, err
	}

	var sites Sites
	if err = json.UnmarshalRead(bytes.NewReader(cachedSites), &sites); err != nil {
		return nil, err
	}

	return sites, nil
}

func saveSites(cache *pkg.Cache, sites Sites) error {
	buffer := bytes.NewBuffer([]byte{})
	if err := json.MarshalWrite(buffer, &sites); err != nil {
		return err
	}

	return cache.Put(sitesCacheKey, buffer.Bytes())
}

func cachePut(cache *pkg.Cache, cmd *cli.Command) error {
	sites, err := loadSites(cache)
	if err != nil {
		return err
	}

	site := cmd.String("site")
	newInfo := SiteInfo{
		Counter: getCounter(false, site, sites, cmd),
		Class:   getClass(false, site, sites, cmd),
		Variant: getVariant(false, site, sites, cmd),
	}
	toAdd := false
	info, ok := sites[site]

	if !ok {
		slog.Info("Added", "site", site, "info", newInfo)
		toAdd = true
	}

	if info != newInfo {
		slog.Info("Updated", "site", site, "prev", info, "new", newInfo)
		toAdd = true
	}

	if toAdd {
		sites[site] = newInfo
		return saveSites(cache, sites)
	}

	slog.Info("Unchanged")
	return nil
}

func cacheRemove(cache *pkg.Cache, cmd *cli.Command) error {
	sites, err := loadSites(cache)
	if err != nil {
		return err
	}

	site := cmd.String("site")
	_, ok := sites[site]
	if ok {
		delete(sites, site)
		return saveSites(cache, sites)
	}

	return nil
}

func cacheShow(cache *pkg.Cache) error {
	sites, err := loadSites(cache)
	if err != nil {
		return err
	}

	return json.MarshalWrite(os.Stdout, &sites, jsontext.WithIndent("  "))
}
