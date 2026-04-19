package ts

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/lispyclouds/dei/pkg"
)

type Parser struct {
	Url      string `json:"url"`
	RepoPath string `json:"repoPath"`
}

type Parsers struct {
	InstallPath string            `json:"installPath"`
	Langs       map[string]Parser `json:"langs"`
}

type Queries struct {
	RepoPrefix string   `json:"repoPrefix"`
	Langs      []string `json:"langs"`
}

type Conf struct {
	Parsers     Parsers `json:"parsers"`
	Queries     Queries `json:"queries"`
	QueriesPath string  `json:"queriesPath"`
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

func gitClone(url, dir string) error {
	log.Info("Cloning", "url", url)
	_, err := pkg.Run(exec.Command("git", "clone", "--depth=1", url, dir))
	return err
}

func checkUptoDate(url, dir string) (bool, error) {
	log.Info("Checking for updates", "url", url)
	out, err := pkg.Run(exec.Command("git", "-C", dir, "pull"))
	if err != nil {
		return false, err
	}

	// TODO: Any better checks?
	if strings.TrimSpace(out) == "Already up to date." {
		log.Info("Upto date", "url", url)
		return true, nil
	}

	return false, nil
}
