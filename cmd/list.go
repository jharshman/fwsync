package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/jharshman/fwsync/config"
	"github.com/spf13/cobra"
)

// List prints out the current source IPs configured in ~/.fwsync and the source IPs active on the configured firewall.
func List() *cobra.Command {
	return &cobra.Command{
		SilenceErrors: true, // errors are always propogated to main, no need to print again
		Use:           "list",
		Short:         "Display your firewall's allowed IPs.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, provider, err := loadConfig()
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
			defer cancel()

			fw, err := provider.Get(ctx, cfg.Name)
			if err != nil {
				return err
			}

			fmt.Printf("fwsync configurations\n----------------------\nlocal: (%s)\n%s", cfgFilePath, prettyPrint(cfg.SourceIPs))
			fmt.Printf("\nremote: (%s)\n%s", cfg.Name, prettyPrint(fw.AllowedIPv4Addresses))
			return nil
		},
	}
}

// GetCurrentIP will return the current IP and print it to standard out.
// It invokes the same mechanisms to fetch the IP as the update command.
// This command will not update the fwsync configuration.
func GetCurrentIP() *cobra.Command {
	return &cobra.Command{
		SilenceErrors: true, // errors are always propogated to main, no need to print again
		Use:           "get-ip",
		Short:         "Fetches your current public IP.",
		RunE: func(cmd *cobra.Command, args []string) error {
			currentIP, err := config.PublicIP()
			fmt.Printf("current public IP: %s\n", currentIP)
			return err
		},
	}
}

func prettyPrint(in []string) string {
	builder := strings.Builder{}
	for _, v := range in {
		builder.WriteString(v + "\n")
	}
	return builder.String()
}
