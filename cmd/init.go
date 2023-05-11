package cmd

import (
	"fmt"
	"os"

	"github.com/jharshman/fwsync/internal/auth"
	"github.com/jharshman/fwsync/internal/user"
	"github.com/spf13/cobra"
)

const (
	transactionFile = ".bitly_firewall"
	project         = "bitly-devvm"
)

func Initialize() *cobra.Command {

	var local *user.Config

	return &cobra.Command{
		Use:   "init",
		Short: "Initialize fwsync configuration.",
		RunE: func(cmd *cobra.Command, args []string) error {
			firewalls, err := auth.GoogleCloudAuthorizedClient.Firewalls.List(project).Do()
			if err != nil {
				return err
			}

			for idx, fw := range firewalls.Items {
				fmt.Printf("%d:\t%s\n", idx, fw.Name)
			}

			var fwSelection int
		SELECTVM:
			fmt.Printf("Select Firewall with your name 0-%d: ", len(firewalls.Items))
			fmt.Scanln(&fwSelection)

			if fwSelection > len(firewalls.Items) || fwSelection < 0 {
				fmt.Fprintf(os.Stderr, "%d is an invalid selection.\n", fwSelection)
				goto SELECTVM
			}

			ip, _ := user.PublicIP()
			cfg, err := user.NewConfig(firewalls.Items[fwSelection].Name, ip)
			if err != nil {
				return err
			}
			local = cfg

			// write file
			home, _ := os.UserHomeDir()
			f, err := os.Create(fmt.Sprintf("%s/%s", home, transactionFile))
			if err != nil && !os.IsExist(err) {
				return err
			}
			defer f.Close()
			return cfg.Write(f)
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return synchronize(local)
		},
	}
}
