package ts

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/lispyclouds/dei/pkg"
)

func buildParserIfChanged(cli, name, dir, url, repoPath, artifact, queriesDest, parserDest string, flushCache bool) error {
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

		log.Info("Linking queries", "url", url)
		destPath := filepath.Join(queriesDest, name)

		if _, err := os.Lstat(destPath); err != nil {
			if !os.IsNotExist(err) {
				return err
			}

			if err = os.Symlink(filepath.Join(dir, repoPath, "queries"), filepath.Join(queriesDest, name)); err != nil {
				return err
			}
		}

		cloned = true
	}

	if !cloned {
		yes, err := checkUptoDate(url, dir)
		if err != nil {
			return err
		}

		if yes {
			return nil
		}
	}

	log.Info("Generating parser", "dir", dir)
	c := exec.Command(cli, "generate")
	c.Dir = filepath.Join(dir, repoPath)
	_, err := pkg.Run(c)
	if err != nil {
		return err
	}

	log.Info("Building parser", "dir", dir)
	c = exec.Command(cli, "build", "-o", filepath.Join(parserDest, artifact))
	c.Dir = filepath.Join(dir, repoPath)
	_, err = pkg.Run(c)
	if err != nil {
		return err
	}

	log.Info("Cleaning up", "dir", dir)
	_, err = pkg.Run(exec.Command("git", "-C", dir, "clean", "-fd"))
	if err != nil {
		return err
	}

	_, err = pkg.Run(exec.Command("git", "-C", dir, "checkout", "."))
	return err
}

func syncParsers(conf Conf, cacheDir, cli string, flushCache bool, langs []string) error {
	queriesPath, err := expandHome(conf.QueriesPath)
	if err != nil {
		return err
	}

	installPath, err := expandHome(conf.Parsers.InstallPath)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for name, info := range conf.Parsers.Langs {
		if len(langs) > 0 && !slices.Contains(langs, name) {
			continue
		}

		wg.Go(func() {
			if err = buildParserIfChanged(
				cli,
				name,
				filepath.Join(cacheDir, "parsers", name),
				info.Url,
				info.RepoPath,
				name+sharedLibExt(),
				queriesPath,
				installPath,
				flushCache,
			); err != nil {
				log.Error("Error in parser", "name", name, "err", err)
			}
		})
	}

	wg.Wait()
	return nil
}
