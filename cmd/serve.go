package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/karashiiro/gittlz/pkg/api"
	"github.com/karashiiro/gittlz/pkg/protocol"
	"github.com/spf13/cobra"
)

var Username string
var Password string
var Protocol string
var GitPort int

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Git server",
	Long: `Starts the Git server. By default, this uses the Git
server protocol, which has no authentication scheme.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Start the control API
		controller, err := api.NewServer(Hostname, ApiPort, Path)
		if err != nil {
			log.Fatalf("Failed to start API server: %v", err)
		}

		log.Println("API server started on port", ApiPort)

		// Start the Git server
		switch Protocol {
		case "git":
			port := getGitPort(9418)
			daemon, err := protocol.StartGit(Hostname, port, BasePath, Path)
			if err != nil {
				log.Fatalf("Failed to start Git server: %v", err)
			}

			log.Println("Git server started on port", port)

			err = daemon.Wait()
			if err != nil {
				log.Fatalf("Git server failed unexpectedly: %v", err)
			}
		case "http":
			port := getGitPort(80)
			log.Println("Git server started on port", port)
			err := protocol.StartSmartHTTP(Hostname, port, Path, Username, Password)
			if err != nil {
				log.Fatalf("Git server failed unexpectedly: %v", err)
			}
		case "ssh":
			port := getGitPort(22)
			wh, err := protocol.StartSSH(Hostname, port, Path, Password)
			if err != nil {
				log.Fatalf("Failed to start Git server: %v", err)
			}

			log.Println("Git server started on port", port)

			err = wh.Wait()
			if err != nil {
				log.Fatalf("Git server failed unexpectedly: %v", err)
			}
		default:
			fmt.Printf("Unknown protocol: %s\n", Protocol)
			os.Exit(1)
		}

		// Stop the control API
		err = controller.Shutdown()
		if err != nil {
			log.Fatalf("API server failed unexpectedly: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&Protocol, "protocol", "P", "git", "Git server protocol. Valid options: \"git\", \"http\", \"ssh\".")
	serveCmd.Flags().IntVar(&GitPort, "git-port", -1, "The port to run the Git server on. Set this to -1 to use the default ports for each protocol.")
	serveCmd.Flags().StringVarP(&Username, "username", "u", "", "Git username for authentication. Leave this empty to disable the username.")
	serveCmd.Flags().StringVarP(&Password, "password", "p", "", "Git password for authentication. Leave this empty to disable the password.")
}

func getGitPort(fallback int) int {
	if GitPort == -1 {
		return fallback
	}

	return GitPort
}
