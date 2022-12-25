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

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Git server",
	Long:  `Starts the Git server.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch Protocol {
		case "git":
			daemon, err := protocol.StartGit(BasePath, Path)
			if err != nil {
				log.Fatalf("Failed to start Git server: %v", err)
			}

			log.Println("Git server started on port 9418")

			err = daemon.Wait()
			if err != nil {
				log.Fatalf("Git server failed unexpectedly: %v", err)
			}
		case "http":
			log.Println("Git server started on port 80")
			err := protocol.StartSmartHTTP(80, Path, Username, Password)
			if err != nil {
				log.Fatalf("Git server failed unexpectedly: %v", err)
			}
		case "ssh":
			wh, err := protocol.StartSSH(22, Path, Password)
			if err != nil {
				log.Fatalf("Failed to start Git server: %v", err)
			}

			log.Println("Git server started on port 22")

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

	serveCmd.Flags().StringVarP(&Protocol, "protocol", "P", "git", "Git server protocol")
	serveCmd.Flags().StringVarP(&Username, "username", "u", "", "Git username for authentication")
	serveCmd.Flags().StringVarP(&Password, "password", "p", "", "Git password for authentication")
	serveCmd.Flags().StringVar(&BasePath, "base-path", "/srv/git", "Base path for Git repos directory (removes this prefix from URLs)")
	serveCmd.Flags().StringVar(&Path, "path", "/srv/git", "Full path for Git repos directory")
}
