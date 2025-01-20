package commands

import (
	"github.com/spf13/cobra"
	"github.com/sycomancy/glasnik/internal/proxy"
)

func NewRegistryCommand() *cobra.Command {
	var port, username, password string

	cmd := &cobra.Command{
		Use:   "registry",
		Short: "Start a proxy registry server",
		RunE: func(cmd *cobra.Command, args []string) error {
			registry := proxy.NewRegistry(port, username, password)
			return registry.Start()
		},
	}

	cmd.Flags().StringVar(&port, "port", "8082", "Port to run the registry server on")
	cmd.Flags().StringVar(&username, "username", "", "Username for registry authentication (optional)")
	cmd.Flags().StringVar(&password, "password", "", "Password for registry authentication (optional)")

	return cmd
}
