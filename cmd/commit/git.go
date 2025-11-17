package commit

import (
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/lispyclouds/dei/pkg"
)

// TODO: Use go-git when

func run(cmd string, args ...string) (string, error) {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		slog.Error("Error running command", "cmd", cmd, "args", args)
		fmt.Println(string(out))

		return "", err
	}

	return string(out), nil
}

func preChecks() error {
	out, err := run("git", "status", "--porcelain")
	if err != nil {
		return err
	}

	if len(out) == 0 {
		return errors.New("No changes to commit")
	}

	return nil
}

func commit(cache *pkg.Cache) error {
	if err := preChecks(); err != nil {
		return err
	}

	var feat string
	featCacheKey := "dei.commit.feat"

	cachedFeat, err := cache.Get(featCacheKey)
	if err != nil {
		return err
	}

	if cachedFeat != nil {
		feat = string(cachedFeat)
	}

	featResp, err := pkg.Input("Feature", feat, false)
	if err != nil {
		return err
	}

	if len(featResp) > 0 {
		feat = featResp
	}

	summary, err := pkg.Input("Summary", "", false)
	if err != nil {
		return err
	}

	var extendedCommitMsg string
	if err = huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("Extended Commit Message").
				Value(&extendedCommitMsg),
		),
	).Run(); err != nil {
		return err
	}

	if len(extendedCommitMsg) > 0 {
		extendedCommitMsg = "\n\n" + extendedCommitMsg
	}

	coAuthors, err := loadCoAuthors(cache)
	if err != nil {
		return err
	}

	var coAuthorsText string
	if len(coAuthors) > 0 {
		options := []huh.Option[string]{}

		for email, info := range coAuthors {
			options = append(options, huh.NewOption(info["name"], email))
		}

		selectedEmails := []string{}
		if err := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Choose your co-author(s)").
					Options(options...).
					Value(&selectedEmails),
			),
		).Run(); err != nil {
			return err
		}

		if len(selectedEmails) > 0 {
			chunks := []string{"\n"}

			for _, email := range selectedEmails {
				chunks = append(chunks, fmt.Sprintf("Co-authored-by: %s <%s>", coAuthors[email]["name"], email))
			}

			coAuthorsText = strings.Join(chunks, "\n")
		}
	}

	if _, err = run(
		"git",
		"commit",
		"--cleanup=verbatim",
		"-m",
		fmt.Sprintf("[%s] %s%s%s", feat, summary, extendedCommitMsg, coAuthorsText),
	); err != nil {
		return err
	}

	return cache.Put(featCacheKey, []byte(feat))
}
