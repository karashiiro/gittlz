package protocol

import (
	"fmt"
	"os"
	"os/exec"
)

func StartGit(basePath, path string) (*exec.Cmd, error) {
	cmd := exec.Command("git", "daemon", "--verbose", "--reuseaddr", fmt.Sprintf("--base-path=%s", basePath), path)
	// TODO: Pass these through a custom logger that prefixes output appropriately
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	return cmd, err
}
