package git

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func gitCommand(dir string, errMsg string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Run(); err != nil {
		return "", errors.Wrapf(err, "%s: %s", errMsg, stdErr.String())
	}

	return strings.Trim(stdOut.String(), "\n"), nil
}
