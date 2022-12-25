package api

import "github.com/labstack/echo/v4"

type serverContext struct {
	echo.Context
	repoDir string
}

func (sc *serverContext) RepoDirectory() string {
	return sc.repoDir
}
