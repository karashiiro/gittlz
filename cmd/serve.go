package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/karashiiro/ditz/pkg/protocol"
	"github.com/spf13/cobra"
)

var Protocol string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Git server",
	Long:  `Starts the Git server.`,
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working directory: %v", err)
		}

		switch Protocol {
		case "git":
			daemon, err := protocol.StartGit(cwd, cwd)
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
			err := protocol.StartSmartHTTP(cwd)
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

	serveCmd.Flags().StringVarP(&Protocol, "protocol", "p", "git", "Git server protocol")
}
