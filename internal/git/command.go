package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func gitCommand(dir string, errMsg string, args ...string) (string, error) {
	// #nosec G204 -- args are intentional git CLI flags/subcommands
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	var stdOut bytes.Buffer

	var stdErr bytes.Buffer

	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %s: %s", err, errMsg, stdErr.String())
	}

	return strings.Trim(stdOut.String(), "\n"), nil
}
