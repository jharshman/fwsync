package cmd

import (
	"fmt"

	"github.com/jharshman/fwsync/config"
	"github.com/spf13/cobra"
)

// Update will intelligently update the firewall rule if the user's public IP has changed and doesn't exist in the
// current rule. If the IP is to be added and the number of IPs in the rule exceeds the IP limit, the oldest IP is dropped from the list.
func Update() *cobra.Command {
	return &cobra.Command{
		SilenceErrors: true, // errors are always propogated to main, no need to print again
		Use:           "update",
		Short:         "Allow a new IP on the firewall.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, provider, err := loadConfig()
			if err != nil {
				return err
			}

			currentIP, err := config.PublicIP()
			if err != nil {
				return err
			}
			if _, ok := cfg.HasIP(currentIP); ok {
				fmt.Println("IPs are up-to-date, skipping sync.")
				return nil
			}

			// Update appends at end so oldest will be front of list.
			cfg.Add(currentIP)

			// Sync before saving so a failed sync doesn't leave the new IP recorded locally,
			// which would cause the next update to skip the sync.
			if err := synchronize(provider, cfg); err != nil {
				return err
			}
			return saveConfig(cfg)
		},
	}
}

// Sync initiates a manual synchronization of the local configuration stored in ~/.fwsync to the configured firewall.
func Sync() *cobra.Command {
	return &cobra.Command{
		SilenceErrors: true, // errors are always propogated to main, no need to print again
		Use:           "sync",
		Short:         "Synchronize local config with firewall",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, provider, err := loadConfig()
			if err != nil {
				return err
			}
			return synchronize(provider, cfg)
		},
	}
}
