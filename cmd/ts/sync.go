package ts

import (
	"context"
	json "encoding/json/v2"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/lispyclouds/dei/pkg"
	"github.com/urfave/cli/v3"
)

type Parser struct {
	Url      string `json:"url"`
	RepoPath string `json:"repoPath"`
}

type Conf struct {
	ParsersPath string            `json:"parsersPath"`
	QueriesPath string            `json:"queriesPath"`
	Langs       map[string]Parser `json:"langs"`
}

func buildIfChanged(cli, name, dir, url, repoPath, artifact, queriesDest, parserDest string) error {
	cloned := false

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}

		log.Info("Cloning", "url", url)
		_, err := pkg.Run(exec.Command("git", "clone", "--depth=1", url, dir))
		if err != nil {
			return err
		}

		log.Info("Linking queries", "url", url)
		err = os.Symlink(filepath.Join(dir, "queries"), filepath.Join(queriesDest, name))
		if err != nil {
			return err
		}

		cloned = true
	}

	if !cloned {
		log.Info("Checking for updates", "url", url)
		out, err := pkg.Run(exec.Command("git", "-C", dir, "pull"))
		if err != nil {
			return err
		}

		// TODO: Any better checks?
		if strings.TrimSpace(out) == "Already up to date." {
			log.Info("Upto date", "url", url)
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

func sharedLibExt() string {
	switch runtime.GOOS {
	case "windows":
		return ".dll"
	case "darwin":
		return ".dylib"
	default:
		return ".so"
	}
}

func expandHome(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, path[1:]), nil
}

func syncCmd(_ context.Context, cmd *cli.Command) error {
	if !pkg.Which("git") {
		return errors.New("Cannot find git on the PATH")
	}

	if !pkg.Which("tree-sitter") {
		return errors.New("Cannot find tree-sitter on the PATH")
	}

	f, err := os.Open(cmd.String("conf"))
	if err != nil {
		return err
	}
	defer f.Close()

	var conf Conf
	if err := json.UnmarshalRead(f, &conf); err != nil {
		return err
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	queriesPath, err := expandHome(conf.QueriesPath)
	if err != nil {
		return err
	}

	parsersPath, err := expandHome(conf.ParsersPath)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for name, info := range conf.Langs {
		wg.Go(func() {
			if err = buildIfChanged(
				cmd.String("cli"),
				name,
				filepath.Join(cacheDir, "dei", "ts", "parsers", name),
				info.Url,
				info.RepoPath,
				name+sharedLibExt(),
				queriesPath,
				parsersPath,
			); err != nil {
				log.Error("Error in parser", "name", name, "err", err)
			}
		})
	}

	wg.Wait()
	return nil
}
