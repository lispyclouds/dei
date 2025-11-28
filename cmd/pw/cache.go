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

	if cachedSites == nil {
		return Sites{}, nil
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

	slog.Warn("No such site", "site", site)
	return nil
}

func cacheDump(cache *pkg.Cache, cmd *cli.Command) error {
	writer := os.Stdout
	if cmd.IsSet("file") {
		f, err := os.OpenFile(cmd.String("file"), os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		writer = f
		defer f.Close()
	}

	sites, err := loadSites(cache)
	if err != nil {
		return err
	}

	return json.MarshalWrite(writer, &sites, jsontext.WithIndent("  "))
}

func cacheImport(cache *pkg.Cache, cmd *cli.Command) error {
	f, err := os.OpenFile(cmd.String("file"), os.O_RDONLY, 0400)
	if err != nil {
		return err
	}
	defer f.Close()

	var sites Sites
	if err = json.UnmarshalRead(f, &sites); err != nil {
		return err
	}

	cachedSites, err := loadSites(cache)
	if err != nil {
		return err
	}

	// TODO: maybe merge deeper?
	for site, newValue := range sites {
		// TODO: Any better way to do defaults?
		if newValue.Counter == 0 {
			newValue.Counter = cmd.Int("counter")
		}

		if newValue.Class == "" {
			newValue.Class = TemplateClass(cmd.String("class"))
		}

		if newValue.Variant == "" {
			newValue.Variant = SiteVariant(cmd.String("variant"))
		}

		cachedSites[site] = newValue
	}

	return saveSites(cache, cachedSites)
}
