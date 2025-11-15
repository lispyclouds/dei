package pkg

import "github.com/charmbracelet/huh"

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
