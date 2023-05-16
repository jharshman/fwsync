package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jharshman/fwsync/internal/auth"
	"github.com/jharshman/fwsync/internal/user"
	"github.com/spf13/cobra"
)

// List prints out the current source IPs configured in ~/.bitly_firewall and the source IPs active on the GCP Firewall.
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

			cfg, err := user.NewFromFile(f)
			if err != nil {
				return err
			}

			localIPs := cfg.SourceIPs

			// get configured fw ips
			fw, err := auth.GoogleCloudAuthorizedClient.Firewalls.Get(project, cfg.Name).Do()
			if err != nil {
				return err
			}
			remoteIPs := fw.SourceRanges

			// pretty print
			fmt.Printf("fwsync configurations\n----------------------\nlocal: (%s)\n%s", f.Name(), prettyPrint(localIPs))
			fmt.Printf("\nremote: (%s)\n%s", cfg.Name, prettyPrint(remoteIPs))
			return nil
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
