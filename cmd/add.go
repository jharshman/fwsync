package cmd

import (
	"fmt"
	"os"

	"github.com/jharshman/fwsync/internal/auth"
	"github.com/jharshman/fwsync/internal/user"
	"github.com/spf13/cobra"
	"google.golang.org/api/compute/v1"
)

func Add() *cobra.Command {

	// Local variable shared between the closures.
	var local *user.Config

	return &cobra.Command{
		Use:   "add",
		Short: "Allow a new IP on the firewall.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// get local configuration
			home, _ := os.UserHomeDir()
			f, err := os.OpenFile(fmt.Sprintf("%s/%s", home, transactionFile), os.O_RDWR, 0666)
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
				return nil
			}

			// Remove oldest in list.
			// Add appends at end so oldest will be front of list.
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
			return synchronize(local)
		},
	}
}

func synchronize(cfg *user.Config) error {
	// CIDR notation required by GoogleAPIs.
	for idx, _ := range cfg.SourceIPs {
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
