package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sycomancy/glasnik/cmd/cli/commands"
	"github.com/sycomancy/glasnik/internal/infra"
)

var rootCmd = &cobra.Command{
	Use:   "glasnik",
	Short: "Glasnik is a web scraping tool",
	Long:  `A web scraping tool for collecting and analyzing data from various sources`,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(commands.GetCleanupCmd())
	rootCmd.AddCommand(commands.GetFetchCmd())
}

func initConfig() {
	infra.LoadConfig()
	infra.MongoConnect("mongodb://root:example@localhost:27017/?authSource=admin")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
