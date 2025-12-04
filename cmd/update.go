package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jharshman/fwsync/config"
	"github.com/spf13/cobra"
)

// Update will intelligently update the firewall rule if the user's public IP has changed and doesn't exist in the
// current rule. If the IP is to be added and the number of IPs in the rule exceeds 5, the oldest IP is dropped from the list.
func Update() *cobra.Command {

	// Local variable shared between the closures.
	var local *config.Config
	var skipSync bool

	return &cobra.Command{
		SilenceErrors: true, // errors are always propogated to main, no need to print again
		Use:           "update",
		Short:         "Allow a new IP on the firewall.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// get local configuration
			f, err := os.OpenFile(cfgFilePath, os.O_RDWR, 0666)
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

			currentIP, err := config.PublicIP()
			if err != nil {
				return err
			}
			_, ipExists := cfg.HasIP(currentIP)
			if ipExists {
				skipSync = true
				return nil
			}

			// Remove oldest in list.
			// Update appends at end so oldest will be front of list.
			cfg.Add(currentIP)
			local = cfg

			// truncate file for writing
			// cannot use Create or os.O_TRUNC
			// The file must read/writable and not truncated before the read
			f.Truncate(0)
			f.Seek(0, 0)

			return cfg.Write(f)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			if skipSync {
				fmt.Println("IPs are up-to-date, skipping sync.")
				return nil
			}
			fmt.Println("syncing firewall rule")
			return synchronize(local)
		},
	}
}

// Sync initiates a manual synchronization of the local configuration stored in ~/.fwsync to the desired GCP Firewall.
func Sync() *cobra.Command {
	return &cobra.Command{
		SilenceErrors: true, // errors are always propogated to main, no need to print again
		Use:           "sync",
		Short:         "Synchronize local config with firewall",
		RunE: func(cmd *cobra.Command, args []string) error {
			// get local configuration
			f, err := os.OpenFile(cfgFilePath, os.O_RDWR, 0666)
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

			return synchronize(cfg)
		},
	}
}

// synchronize will use the local configuration update the desired firewall rule.
func synchronize(config *config.Config) error {
	// CIDR notation required by GoogleAPIs.
	for idx := range config.SourceIPs {
		config.SourceIPs[idx] = config.SourceIPs[idx] + "/32"
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return FirewallClient.Update(ctx, config.Name, config.SourceIPs)
}
