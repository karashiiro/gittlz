package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/karashiiro/gittlz/pkg/protocol"
	"github.com/spf13/cobra"
)

var Username string
var Password string
var BasePath string
var Path string
var Protocol string
var GitPort int

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Git server",
	Long: `Starts the Git server. By default, this uses the Git
server protocol, which has no authentication scheme.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch Protocol {
		case "git":
			port := getGitPort(9418)
			daemon, err := protocol.StartGit(port, BasePath, Path)
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
			err := protocol.StartSmartHTTP(port, Path, Username, Password)
			if err != nil {
				log.Fatalf("Git server failed unexpectedly: %v", err)
			}
		case "ssh":
			port := getGitPort(22)
			wh, err := protocol.StartSSH(port, Path, Password)
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
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&Protocol, "protocol", "P", "git", "Git server protocol. Valid options: \"git\", \"http\", \"ssh\".")
	serveCmd.Flags().IntVar(&GitPort, "git-port", -1, "The port to run the Git server on. Set this to -1 to use the default ports for each protocol.")
	serveCmd.Flags().StringVarP(&Username, "username", "u", "", "Git username for authentication. Leave this empty to disable the username.")
	serveCmd.Flags().StringVarP(&Password, "password", "p", "", "Git password for authentication. Leave this empty to disable the password.")
	serveCmd.Flags().StringVar(&BasePath, "base-path", "/srv/git", "Base path for the Git repositories directory. This acts a prefix removed from --path.")
	serveCmd.Flags().StringVar(&Path, "path", "/srv/git", "Full path for the Git repositories directory.")
}

func getGitPort(fallback int) int {
	if GitPort == -1 {
		return fallback
	}

	return GitPort
}
