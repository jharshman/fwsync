package cmd

import (
	"fmt"
	"os"

	"github.com/jharshman/fwsync/internal/auth"
	"github.com/jharshman/fwsync/internal/user"
	"github.com/spf13/cobra"
	"google.golang.org/api/compute/v1"
)

// Update will intelligently update the firewall rule if the user's public IP has changed and doesn't exist in the
// current rule. If the IP is to be added and the number of IPs in the rule exceeds 5, the oldest IP is dropped from the list.
func Update() *cobra.Command {

	// Local variable shared between the closures.
	var local *user.Config
	var skipSync bool

	return &cobra.Command{
		Use:   "update",
		Short: "Allow a new IP on the firewall.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// get local configuration
			f, err := os.OpenFile(cfgFilePath, os.O_RDWR, 0666)
			if err != nil {
				return err
			}
			defer f.Close()

			cfg, err := user.NewFromFile(f)
			if err != nil {
				return err
			}

			currentIP, _ := user.PublicIP()
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

// Sync initiates a manual synchronization of the local configuration stored in ~/.bitly_firewall to the desired GCP Firewall.
func Sync() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Synchronize local config with firewall",
		RunE: func(cmd *cobra.Command, args []string) error {
			// get local configuration
			f, err := os.OpenFile(cfgFilePath, os.O_RDWR, 0666)
			if err != nil {
				return err
			}
			defer f.Close()

			cfg, err := user.NewFromFile(f)
			if err != nil {
				return err
			}

			return synchronize(cfg)
		},
	}
}

// synchronize will use the local configuration update the desired firewall rule.
func synchronize(cfg *user.Config) error {
	// CIDR notation required by GoogleAPIs.
	for idx := range cfg.SourceIPs {
		cfg.SourceIPs[idx] = cfg.SourceIPs[idx] + "/32"
	}

	fw := &compute.Firewall{
		SourceRanges: cfg.SourceIPs,
	}
	_, err := auth.GoogleCloudAuthorizedClient.Firewalls.Patch(project, cfg.Name, fw).Do()
	if err != nil {
		return err
	}
	return nil
}
