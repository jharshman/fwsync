package main

import (
	"fmt"
	"os"
	"strings"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/jharshman/fwsync/cmd"
	"github.com/jharshman/fwsync/internal/auth"
	"github.com/spf13/cobra"
)

var version string

func main() {

	rootCmd := &cobra.Command{
		Use:   "fwsync",
		Short: "A CLI utility to keep your development VM firewall up to date.",
		Long: `fwsync uses a local file to keep track of the latest IP addresses you've been
connecting from and keeps your development VM firewall rule up to date with that list.`,
		PersistentPreRunE: auth.Auth(),
	}

	rootCmd.AddCommand(cmd.Initialize())
	rootCmd.AddCommand(cmd.Update())
	rootCmd.AddCommand(cmd.List())
	rootCmd.AddCommand(cmd.Sync())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Operation complete, no errors.")

	err := notifyIfUpdateAvailable()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func notifyIfUpdateAvailable() error {
	cli, err := ghapi.DefaultRESTClient()
	if err != nil {
		return err
	}
	latestTag := struct {
		Tag string `json:"tag_name"`
	}{}
	err = cli.Get("repos/jharshman/fwsync/releases/latest", &latestTag)
	if err != nil {
		return err
	}
	tag := strings.TrimPrefix(latestTag.Tag, "v")
	if tag != version {
		fmt.Printf("\n\033[0;32mA new version (%q) is available for fwsync.\033[0m\n", latestTag.Tag)
		fmt.Printf("\033[0;32mTo update run:\ncurl https://raw.githubusercontent.com/jharshman/fwsync/master/install.sh | sh\033[0m\n")
	}
	return nil
}
