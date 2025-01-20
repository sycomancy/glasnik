package commands

import (
	"github.com/spf13/cobra"
	"github.com/sycomancy/glasnik/internal/infra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Cleanup database collections",
	Long:  `Removes all documents from specified collections or all collections if none specified`,
	Run: func(cmd *cobra.Command, args []string) {
		collections, _ := cmd.Flags().GetStringSlice("collections")
		if len(collections) == 0 {
			infra.CleanupAllCollections()
		} else {
			infra.CleanupCollections(collections)
		}
	},
}

func init() {
	cleanupCmd.Flags().StringSliceP("collections", "c", []string{}, "Collections to cleanup (comma-separated)")
}

func GetCleanupCmd() *cobra.Command {
	return cleanupCmd
}
