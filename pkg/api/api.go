package api

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e    *echo.Echo
	done chan os.Signal
}

func (s *Server) Shutdown() error {
	<-s.done

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.e.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func NewServer(host string, port int, repoDir string) (*Server, error) {
	e := echo.New()

	// Configure the API
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sc := &serverContext{c, repoDir}
			return next(sc)
		}
	})
	// TODO: Use a custom logger that prefixes output appropriately
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/repo", createRepository)

	// Start the server
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := e.Start(fmt.Sprintf("%s:%d", host, port)); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	}()

	s := &Server{e, done}
	return s, nil
}

type CreateRepositoryParams struct {
	Name string `json:"name"`
}

func createRepository(c echo.Context) error {
	params := &CreateRepositoryParams{}
	err := c.Bind(params)
	if err != nil {
		return c.JSON(500, map[string]string{
			"error": fmt.Sprintf("echo.Context.Bind: %v", err),
		})
	}

	name := params.Name
	if name == "" {
		return c.JSON(400, map[string]string{
			"error": "No repository name provided.",
		})
	}

	if !strings.HasSuffix(name, ".git") {
		name += ".git"
	}

	sc := c.(*serverContext)

	cmd := exec.Command("git", "init", "--bare", name)
	cmd.Dir = sc.RepoDirectory()
	// TODO: Pass these through a custom logger that prefixes output appropriately
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return c.JSON(500, map[string]string{
			"error": fmt.Sprintf("exec.Cmd.Run: %v", err),
		})
	}

	return c.JSON(200, map[string]string{
		"name": name,
	})
}
