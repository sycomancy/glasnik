package commands

import (
	"github.com/spf13/cobra"
	"github.com/sycomancy/glasnik/internal/proxy"
)

func NewRegistryCommand() *cobra.Command {
	var port string

	cmd := &cobra.Command{
		Use:   "registry",
		Short: "Start a proxy registry server",
		RunE: func(cmd *cobra.Command, args []string) error {
			registry := proxy.NewRegistry(port)
			return registry.Start()
		},
	}

	cmd.Flags().StringVar(&port, "port", "8080", "Port to run the registry server on")

	return cmd
}
