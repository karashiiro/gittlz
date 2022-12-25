package protocol

import (
	"fmt"
	"os"
	"os/exec"
)

func StartGit(port int, basePath, path string) (*exec.Cmd, error) {
	cmd := exec.Command("git", "daemon", "--verbose", "--reuseaddr", "--export-all", fmt.Sprintf("--port=%d", port), fmt.Sprintf("--base-path=%s", basePath), path)
	// TODO: Pass these through a custom logger that prefixes output appropriately
	// TODO: Automatically create git-daemon-export-ok?
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	return cmd, err
}
