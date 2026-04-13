package pkg

import (
	"context"
	"fmt"
	"os/exec"

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

func Which(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func Run(cmd *exec.Cmd) (string, error) {
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Error running command", "cmd", cmd)
		fmt.Println(string(out))

		return "", err
	}

	return string(out), nil
}
