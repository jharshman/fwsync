package cmd

import (
	"fmt"
	"os"

	"github.com/jharshman/fwsync/internal/user"
	"github.com/spf13/cobra"
)

func Add() *cobra.Command {
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

			// truncate file for writing
			// cannot use Create or os.O_TRUNC
			// The file must read/writable and not truncated before the read
			f.Truncate(0)
			f.Seek(0, 0)

			return cfg.Write(f)
		},
	}
}
