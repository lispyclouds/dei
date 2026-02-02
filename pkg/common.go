package pkg

import (
	"context"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v3"
)

func Input(prompt, placeholder string, noEcho bool) (string, error) {
	var response string
	echoMode := huh.EchoModeNormal

	if noEcho {
		echoMode = huh.EchoModePassword
	}

	if err := huh.NewInput().
		Title(prompt).
		Placeholder(placeholder).
		Value(&response).
		EchoMode(echoMode).
		Run(); err != nil {
		return "", err
	}

	return response, nil
}

func CommandNotFound(_ context.Context, cmd *cli.Command, command string) {
	log.Error("Unknown command", "command", command)
	cli.ShowSubcommandHelpAndExit(cmd, 1)
}
