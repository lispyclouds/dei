package ts

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/charmbracelet/log"
)

func syncQueryIfChanged(name, url, dir, queriesDest string, flushCache bool) error {
	if flushCache {
		os.RemoveAll(dir)
	}

	cloned := false

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}

		if err = gitClone(url, dir); err != nil {
			return err
		}

		log.Info("Linking external queries", "name", name)
		destPath := filepath.Join(queriesDest, name)
		if _, err := os.Lstat(destPath); !os.IsNotExist(err) {
			log.Info("Removing existing query link", "path", destPath)
			os.Remove(destPath)
		}

		err = os.Symlink(filepath.Join(dir, "queries"), destPath)
		if err != nil {
			return err
		}

		cloned = true
	}

	if !cloned {
		_, err := checkUptoDate(url, dir)
		return err
	}

	return nil
}

func syncQueries(conf Conf, cacheDir string, flushCache bool) error {
	repoPrefix, err := expandHome(conf.Queries.RepoPrefix)
	if err != nil {
		return err
	}

	installPath, err := expandHome(conf.QueriesPath)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, name := range conf.Queries.Langs {
		wg.Go(func() {
			if err = syncQueryIfChanged(
				name,
				repoPrefix+name,
				filepath.Join(cacheDir, name),
				installPath,
				flushCache,
			); err != nil {
				log.Error("Error fetching query", "name", name, "err", err)
			}
		})
	}

	wg.Wait()
	return nil
}
