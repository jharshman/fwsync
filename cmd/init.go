package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jharshman/fwsync/config"
	"github.com/jharshman/fwsync/internal/providers"
	"github.com/spf13/cobra"
)

// Initialize performs the first sync of the firewall rule. It will prompt the user to select
// the firewall rule to manage and will then update that firewall rule with their current
// public IP. Any existing source IPs on the firewall rule will be overwritten.
func Initialize() *cobra.Command {
	var cloudProvider string
	var cloudProject string
	var ipLimit int

	initCmd := &cobra.Command{
		SilenceErrors: true, // errors are always propogated to main, no need to print again
		Use:           "init",
		Short:         "Initialize fwsync configuration.",
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
		RunE: func(cmd *cobra.Command, args []string) error {
			if ipLimit < 1 {
				return errors.New("--ip-limit must be at least 1")
			}

			cfg := config.New(
				config.WithProvider(cloudProvider),
				config.WithProject(cloudProject),
				config.WithIPLimit(ipLimit))

			provider, err := providers.New(cfg.Provider, cfg.ProviderSettings())
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
			defer cancel()

			firewalls, err := provider.List(ctx)
			if err != nil {
				return err
			}
			if len(firewalls) == 0 {
				return fmt.Errorf("no firewalls found for provider %s, create one and run init again", cfg.Provider)
			}

			for idx, fw := range firewalls {
				fmt.Printf("%d:\t%s\n", idx, fw.Name)
			}

		ASK:
			selection, ok := ask(fmt.Sprintf("Select Firewall to use 0-%d: ", len(firewalls)-1), false, func(val string) bool {
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
						return true
					default:
						return false
					}
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

			ip, err := config.PublicIP()
			if err != nil {
				return err
			}
			fmt.Printf("IP determined to be: %s\n", ip)
			cfg.Name = firewalls[fwSelection].Name
			cfg.SourceIPs = []string{ip}

			if err := synchronize(provider, cfg); err != nil {
				return err
			}
			return saveConfig(cfg)
		},
	}
	initCmd.Flags().StringVar(&cloudProvider, "provider", providers.Google,
		fmt.Sprintf("Cloud Provider (%s)", strings.Join(providers.Names(), ", ")))
	initCmd.Flags().StringVar(&cloudProject, "project", "", "Cloud Project (required for google)")
	initCmd.Flags().IntVar(&ipLimit, "ip-limit", 5, "Maximum number of IPs to allow")
	return initCmd
}

func ask(prompt string, skipRetry bool, check func(val string) bool) (string, bool) {
	var answer string
PROMPT:
	answer = "" // Scanln leaves answer untouched on empty input
	fmt.Print(prompt)
	fmt.Scanln(&answer)

	b := check(answer)
	if !b && !skipRetry {
		goto PROMPT
	}
	return answer, b
}
