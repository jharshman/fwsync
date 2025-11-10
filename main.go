package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v53/github"
	"github.com/jharshman/fwsync/cmd"
	"github.com/spf13/cobra"
)

var version string

func main() {

	rootCmd := &cobra.Command{
		Use:   "fwsync",
		Short: "A CLI utility to keep your development VM firewall up to date.",
		Long: `fwsync uses a local file to keep track of the latest IP addresses you've been
connecting from and keeps your development VM firewall rule up to date with that list.`,
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	rootCmd.AddCommand(cmd.Initialize())
	rootCmd.AddCommand(cmd.Update())
	rootCmd.AddCommand(cmd.List())
	rootCmd.AddCommand(cmd.Sync())
	rootCmd.AddCommand(cmd.GetCurrentIP())
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Operation complete, no errors.")

	err := notifyIfUpdateAvailable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for updates to fwsync: %q", err)
	}
}

func notifyIfUpdateAvailable() error {
	cli := github.NewClient(nil)
	repoRelease, _, err := cli.Repositories.GetLatestRelease(context.Background(), "jharshman", "fwsync")
	if err != nil {
		return err
	}
	latest := strings.TrimPrefix(repoRelease.GetTagName(), "v")
	if latest != version {
		fmt.Printf("\n\033[0;32mA new version (%q) is available for fwsync.\033[0m\n", latest)
		fmt.Printf("\033[0;32mTo update run:\ncurl https://raw.githubusercontent.com/jharshman/fwsync/master/install.sh | sh\033[0m\n")
	}
	return nil
}
