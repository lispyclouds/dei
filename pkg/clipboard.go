package pkg

import (
	"fmt"
	"maps"
	"os/exec"
	"runtime"
	"slices"
)

func getCopyCmd(tools map[string][]string) (*exec.Cmd, error) {
	for tool, args := range tools {
		if _, err := exec.LookPath(tool); err == nil {
			return exec.Command(tool, args...), nil
		}
	}

	return nil, fmt.Errorf("Cannot find a clipboard tool, tried %v", slices.Collect(maps.Keys(tools)))
}

func CopyToClipboard(data string) error {
	conf := map[string]map[string][]string{
		"linux":   {"wl-copy": {}, "xsel": {"-ib"}, "xclip": {}},
		"darwin":  {"pbcopy": {}},
		"windows": {"clip": {}},
	}

	tools, ok := conf[runtime.GOOS]
	if !ok {
		return fmt.Errorf("Unsupported OS: %s", runtime.GOOS)
	}

	cmd, err := getCopyCmd(tools)
	if err != nil {
		return err
	}

	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	if _, err = in.Write([]byte(data)); err != nil {
		return err
	}

	if err = in.Close(); err != nil {
		return err
	}

	return cmd.Wait()
}
