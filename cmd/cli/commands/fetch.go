package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/sycomancy/glasnik/internal/job"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch data from a source",
	Long:  `Fetch data from a source and store it in the database`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			cmd.Println("URL is required")
			return
		}

		job, err := job.NewJob(url)
		if err != nil {
			fmt.Println("unable to start job", err)
			os.Exit(-1)
		}

		entries := job.FetchEntries()
		fmt.Println(entries)
	},
}

func init() {
	fetchCmd.Flags().StringP("url", "u", "", "URL to fetch data from")
	fetchCmd.Flags().BoolP("details", "d", false, "Fetch detailed data for each entry")
}

func GetFetchCmd() *cobra.Command {
	return fetchCmd
}
