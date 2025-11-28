package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jharshman/fwsync/config"
	"github.com/spf13/cobra"
)

const (
	transactionFile = ".fwsync"
)

var (
	home, _     = os.UserHomeDir()
	cfgFilePath = fmt.Sprintf("%s/%s", home, transactionFile)
)

// Initialize performs the first sync of the firewall rule. It will prompt the user to select
// the firewall rule with his or her name and then will update that firewall rule with their current
// public IP. Any existing source IPs on the firewall rule will be overwritten.
func Initialize() *cobra.Command {
	var local *config.Config
	var cloudProvider string
	var cloudProject string
	var ipLimit int

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize fwsync configuration.",
		RunE: func(cmd *cobra.Command, args []string) error {

			if cloudProvider == config.ProviderGoogle {
				return fmt.Errorf("the provider: %s requires the --project argument", config.ProviderGoogle)
			}

			cfg := config.New(
				config.WithProvider(cloudProvider),
				config.WithProject(cloudProject),
				config.WithIPLimit(ipLimit))

			var err error
			FirewallClient, err = cfg.AuthForProvider()
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			firewalls, err := FirewallClient.List(ctx)
			if err != nil {
				return err
			}

			for idx, fw := range firewalls {
				fmt.Printf("%d:\t%s\n", idx, fw.Name)
			}

		ASK:
			selection, ok := ask(fmt.Sprintf("Select Firewall to use 0-%d: ", len(firewalls)), false, func(val string) bool {
				i, err := strconv.Atoi(val)
				if err != nil {
					return false
				}
				if i >= len(firewalls) || i < 0 {
					return false
				}

				_, ok := ask(fmt.Sprintf("You've selected %s, is that correct? [Y/n]: ", firewalls[i].Name), true, func(val string) bool {
					switch val {
					case "Y", "y", "yes", "":
					case "N", "n", "no":
						return false
					default:
						return false
					}
					return true
				})
				return ok
			})
			if !ok {
				goto ASK
			}

			fwSelection, err := strconv.Atoi(selection)
			if err != nil {
				return err
			}

			ip, _ := config.PublicIP()
			fmt.Printf("IP determined to be: %s\n", ip)
			cfg.Name = firewalls[fwSelection].Name
			cfg.SourceIPs = []string{ip}

			local = cfg

			// write file
			f, err := os.Create(cfgFilePath)
			if err != nil && !os.IsExist(err) {
				return err
			}
			defer f.Close()
			return cfg.Write(f)
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(cfgFilePath); err != nil {
				// config file doesn't exist, continue to RunE to go through creation.
				return nil
			}
			// prompt to nuke existing configuration file.
			ask("Existing configuration file detected. Continue anyway? [Y/n]: ", false, func(val string) bool {
				switch val {
				case "Y", "y", "yes", "":
				case "N", "n", "no":
					os.Exit(0)
				default:
					return false
				}
				return true
			})
			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("syncing firewall rule")
			return synchronize(local)
		},
	}
	initCmd.Flags().StringVar(&cloudProvider, "provider", "", "Cloud Provider")
	initCmd.Flags().StringVar(&cloudProject, "project", "", "Cloud Project")
	initCmd.Flags().IntVar(&ipLimit, "ip-limit", 5, "IP Limit")
	initCmd.MarkFlagRequired("provider")
	return initCmd
}

func ask(prompt string, skipRetry bool, check func(val string) bool) (string, bool) {
	var answer string
PROMPT:
	fmt.Print(prompt)
	fmt.Scanln(&answer)

	b := check(answer)
	if !b && !skipRetry {
		goto PROMPT
	}
	return answer, b
}
