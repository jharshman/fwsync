package main

import (
	"fmt"
	"os"

	"github.com/jharshman/fwsync/cmd"
	"github.com/jharshman/fwsync/internal/auth"
	"github.com/spf13/cobra"
)

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
}
