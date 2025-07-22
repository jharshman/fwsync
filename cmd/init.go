package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jharshman/fwsync/internal/auth"
	"github.com/jharshman/fwsync/internal/user"
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
	var local *user.Config
	var gcpProject string
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize fwsync configuration.",
		RunE: func(cmd *cobra.Command, args []string) error {
			firewalls, err := auth.GoogleCloudAuthorizedClient.Firewalls.List(gcpProject).Do()
			if err != nil {
				return err
			}

			for idx, fw := range firewalls.Items {
				fmt.Printf("%d:\t%s\n", idx, fw.Name)
			}
		ASK:
			selection, ok := ask(fmt.Sprintf("Select Firewall with your name 0-%d: ", len(firewalls.Items)), false, func(val string) bool {
				i, err := strconv.Atoi(val)
				if err != nil {
					return false
				}
				if i > len(firewalls.Items) || i < 0 {
					return false
				}

				_, ok := ask(fmt.Sprintf("You've selected %s, is that correct? [Y/n]: ", firewalls.Items[i].Name), true, func(val string) bool {
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

			ip, _ := user.PublicIP()
			fmt.Printf("IP determined to be: %s\n", ip)
			cfg, err := user.NewConfig(firewalls.Items[fwSelection].Name, ip)
			if err != nil {
				return err
			}
			cfg.Project = gcpProject
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
	initCmd.Flags().StringVar(&gcpProject, "project", "", "GCP Project")
	initCmd.MarkFlagRequired("project")
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
