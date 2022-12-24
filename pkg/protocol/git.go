package protocol

import (
	"fmt"
	"os"
	"os/exec"
)

func StartGit(basePath, path string) (*exec.Cmd, error) {
	cmd := exec.Command("git", "daemon", "--reuseaddr", fmt.Sprintf("--base-path=%s", basePath), path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	return cmd, err
}
