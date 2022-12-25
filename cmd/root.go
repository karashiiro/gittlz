package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var BasePath string
var Path string
var Hostname string
var ApiPort int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gittlz",
	Short: "A zero-maintenance Git server for high-maintenance people.",
	Long: `gittlz is a thin wrapper around the Git server that enables quick
setup and teardown of a managed Git environment. Use gittlz when
you want to test tools that interact with managed Git providers
such as GitHub and GitLab.

To spin up a no-auth Git server, just run "gittlz serve".`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gittlz.yaml)")
	rootCmd.PersistentFlags().StringVar(&Hostname, "host", "0.0.0.0", "The hostname to run the servers on.")
	rootCmd.PersistentFlags().IntVar(&ApiPort, "api-port", 6177, "The port to run the control API server on.")
	rootCmd.PersistentFlags().StringVar(&BasePath, "base-path", "/srv/git", "Base path for the Git repositories directory. This acts a prefix removed from --path.")
	rootCmd.PersistentFlags().StringVar(&Path, "path", "/srv/git", "Full path for the Git repositories directory.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
