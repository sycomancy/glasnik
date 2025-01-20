package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sycomancy/glasnik/internal/infra"
	"github.com/sycomancy/glasnik/internal/job"
	"github.com/sycomancy/glasnik/internal/proxy"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch data from a source",
	Long:  `Fetch data from a source and store it in the database`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		proxyPort, _ := cmd.Flags().GetString("proxy-port")
		proxyUsername, _ := cmd.Flags().GetString("proxy-username")
		proxyPassword, _ := cmd.Flags().GetString("proxy-password")

		if url == "" {
			cmd.Println("URL is required")
			return
		}

		// Initialize proxy registry
		registry := proxy.NewRegistry(proxyPort, proxyUsername, proxyPassword)
		client := infra.NewIncognitoClient(registry, nil)

		job, err := job.NewJob(url, client)
		if err != nil {
			fmt.Println("unable to start job", err)
			os.Exit(-1)
		}

		err = job.FetchEntries()
		if err != nil {
			fmt.Println("error fetching entries:", err)
			os.Exit(-1)
		}
	},
}

func init() {
	fetchCmd.Flags().StringP("url", "u", "", "URL to fetch data from")
	fetchCmd.Flags().StringP("proxy-port", "p", "8080", "Proxy server port")
	fetchCmd.Flags().String("proxy-username", "", "Proxy server username")
	fetchCmd.Flags().String("proxy-password", "", "Proxy server password")
	fetchCmd.Flags().BoolP("details", "d", false, "Fetch detailed data for each entry")
}

func GetFetchCmd() *cobra.Command {
	return fetchCmd
}
