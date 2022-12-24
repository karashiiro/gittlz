package protocol

import (
	"fmt"
	"net/http"
	"net/http/cgi"
	"os/exec"
)

func StartSmartHTTP(projectRoot string) error {
	git, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("exec.LookPath: %w", err)
	}

	handler := &cgi.Handler{
		Path: git,
		Args: []string{"http-backend"},
		Env:  []string{fmt.Sprintf("GIT_PROJECT_ROOT=%s", projectRoot), "GIT_HTTP_EXPORT_ALL="},
	}

	return http.ListenAndServe(":80", handler)
}
