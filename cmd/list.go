package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jharshman/fwsync/config"
	"github.com/spf13/cobra"
)

// List prints out the current source IPs configured in ~/.fwsync and the source IPs active on the GCP Firewall.
func List() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Display your firewall's allowed IPs.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// get local configured ips
			home, _ := os.UserHomeDir()
			f, err := os.Open(fmt.Sprintf("%s/%s", home, transactionFile))
			if err != nil {
				return err
			}
			defer f.Close()

			cfg, err := config.LoadFromFile(f)
			if err != nil {
				return err
			}

			FirewallClient, err = cfg.AuthForProvider()
			if err != nil {
				return err
			}

			localIPs := cfg.SourceIPs

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			// get configured fw ips
			fw, err := FirewallClient.Get(ctx, cfg.Name)
			if err != nil {
				return err
			}
			remoteIPs := fw.AllowedIPs

			// pretty print
			fmt.Printf("fwsync configurations\n----------------------\nlocal: (%s)\n%s", f.Name(), prettyPrint(localIPs))
			fmt.Printf("\nremote: (%s)\n%s", cfg.Name, prettyPrint(remoteIPs))
			return nil
		},
	}
}

// GetCurrentIP will return the current IP and print it to standard out.
// It invokes the same mechanisms to fetch the IP as the update command.
// This command will not update the fwsync configuration.
func GetCurrentIP() *cobra.Command {
	return &cobra.Command{
		Use:   "get-ip",
		Short: "Fetches your current public IP.",
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
