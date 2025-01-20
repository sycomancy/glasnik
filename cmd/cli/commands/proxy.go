package commands

import (
	"github.com/spf13/cobra"
	"github.com/sycomancy/glasnik/internal/proxy"
)

func NewProxyCommand() *cobra.Command {
	var host, port, username, password, registryURL string

	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "Start a proxy server",
		RunE: func(cmd *cobra.Command, args []string) error {
			server := proxy.NewProxyServer(host, port, username, password)

			if registryURL != "" {
				if err := server.RegisterWithRegistry(registryURL); err != nil {
					return err
				}
			}

			return server.Start()
		},
	}

	cmd.Flags().StringVar(&host, "host", "localhost", "Host to bind the proxy server to")
	cmd.Flags().StringVar(&port, "port", "8083", "Port to run the proxy server on")
	cmd.Flags().StringVar(&username, "username", "", "Username for proxy authentication")
	cmd.Flags().StringVar(&password, "password", "", "Password for proxy authentication")
	cmd.Flags().StringVar(&registryURL, "registry", "", "URL of the registry server to register with")

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")

	return cmd
}
