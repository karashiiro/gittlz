package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/karashiiro/gittlz/pkg/api"
	"github.com/spf13/cobra"
)

// createRepoCmd represents the createRepo command
var createRepoCmd = &cobra.Command{
	Use:   "create-repo",
	Short: "Create a Git repository",
	Long:  `Creates a Git repository.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		params := &api.CreateRepositoryParams{
			Name: args[0],
		}

		buf, err := json.Marshal(params)
		if err != nil {
			fmt.Printf("Failed to parse arguments: %v", err)
			os.Exit(1)
		}

		reader := new(bytes.Buffer)
		_, err = reader.Write(buf)
		if err != nil {
			fmt.Printf("Failed to process arguments: %v", err)
			os.Exit(1)
		}

		res, err := http.Post(fmt.Sprintf("http://%s:%d/repo", Hostname, ApiPort), "application/json", reader)
		if err != nil {
			fmt.Printf("Failed to complete operation: %v", err)
			os.Exit(1)
		}
		defer res.Body.Close()

		resBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Failed to read response data: %v", err)
			os.Exit(1)
		}

		fmt.Printf(string(resBytes))
	},
}

func init() {
	rootCmd.AddCommand(createRepoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createRepoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createRepoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
