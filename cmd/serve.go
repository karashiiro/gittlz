package cmd

import (
	"log"
	"os"

	"github.com/karashiiro/ditz/pkg/protocol"
	"github.com/spf13/cobra"
)

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

		daemon, err := protocol.StartGit(cwd, cwd)
		if err != nil {
			log.Fatalf("Failed to start Git server: %v", err)
		}

		log.Println("Git server started on port 9418")

		err = daemon.Wait()
		if err != nil {
			log.Fatalf("Git server failed unexpectedly: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
