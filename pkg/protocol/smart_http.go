package protocol

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cgi"
	"os/exec"
)

func StartSmartHTTP(port int, projectRoot, username, password string) error {
	git, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("exec.LookPath: %w", err)
	}

	gitService := &cgi.Handler{
		Path: git,
		Args: []string{"http-backend"},
		Env:  []string{fmt.Sprintf("GIT_PROJECT_ROOT=%s", projectRoot), "GIT_HTTP_EXPORT_ALL="},
	}

	auth := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.RequestURI)

		// No authentication configured
		if username == "" && password == "" {
			gitService.ServeHTTP(w, r)
			return
		}

		// URL-encoded authentication (deprecated by major browsers)
		authInfo := r.URL.User
		u := authInfo.Username()
		p, _ := authInfo.Password()

		// HTTP Basic auth
		if u == "" && p == "" {
			u, p, _ = r.BasicAuth()
		}

		// Check provided credentials against configured credentials
		if u == username && p == password {
			gitService.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		}
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", port), auth)
}
