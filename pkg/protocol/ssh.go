package protocol

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/git"
	"github.com/charmbracelet/wish/logging"
)

type gitBackend struct {
	access git.AccessLevel
}

func (b gitBackend) AuthRepo(repo string, pk ssh.PublicKey) git.AccessLevel {
	return b.access
}

func (b gitBackend) Push(repo string, pk ssh.PublicKey) {
	log.Printf("Pushed %s", repo)
}

func (b gitBackend) Fetch(repo string, pk ssh.PublicKey) {
	log.Printf("Fetched %s", repo)
}

var _ git.Hooks = &gitBackend{}

type Waiter struct {
	server *ssh.Server
	done   chan os.Signal
}

func (wh *Waiter) Wait() error {
	<-wh.done

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := wh.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func StartSSH(port int, path, password string) (*Waiter, error) {
	b := &gitBackend{git.ReadWriteAccess}

	// Create the SSH server
	s, err := wish.NewServer(
		ssh.PublicKeyAuth(pkAuth),
		ssh.PasswordAuth(passwordAuth(password != "", password)),
		wish.WithAddress(fmt.Sprintf(":%d", port)),
		wish.WithMiddleware(
			git.Middleware(path, b),
			logging.Middleware(),
		),
	)
	if err != nil {
		return nil, err
	}

	// Start the server
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err = s.ListenAndServe(); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	}()

	// Return a wait handle
	wh := &Waiter{s, done}
	return wh, nil
}

func pkAuth(ctx ssh.Context, key ssh.PublicKey) bool {
	return true
}

func passwordAuth(usePassword bool, password string) ssh.PasswordHandler {
	return func(ctx ssh.Context, p string) bool {
		// No authentication
		if !usePassword {
			return false
		}

		return password == p
	}
}
