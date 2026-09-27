package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jharshman/fwsync/config"
	"github.com/jharshman/fwsync/internal/providers"
	"github.com/jharshman/fwsync/internal/providers/generic"
)

const (
	transactionFile = ".fwsync"
	apiTimeout      = 10 * time.Second
)

var (
	home, _     = os.UserHomeDir()
	cfgFilePath = filepath.Join(home, transactionFile)
)

// loadConfig reads the fwsync configuration from disk and authenticates with its provider.
func loadConfig() (*config.Config, generic.Provider, error) {
	f, err := os.Open(cfgFilePath)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	cfg, err := config.LoadFromFile(f)
	if err != nil {
		return nil, nil, fmt.Errorf("reading %s: %w", cfgFilePath, err)
	}

	provider, err := providers.New(cfg.Provider, cfg.ProviderSettings())
	if err != nil {
		return nil, nil, err
	}
	return cfg, provider, nil
}

// saveConfig writes the fwsync configuration to disk, replacing any existing file.
func saveConfig(cfg *config.Config) error {
	f, err := os.Create(cfgFilePath)
	if err != nil {
		return err
	}
	if err := cfg.Write(f); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// synchronize sets the allowed IPs on the configured firewall to the IPs in the local configuration.
func synchronize(provider generic.Provider, cfg *config.Config) error {
	fmt.Println("syncing firewall rule")
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()
	return provider.Update(ctx, cfg.Name, toCIDRs(cfg.SourceIPs))
}

// toCIDRs converts bare IPv4 addresses to /32 CIDR notation. Entries already in CIDR notation are kept as is.
func toCIDRs(ips []string) []string {
	cidrs := make([]string, 0, len(ips))
	for _, ip := range ips {
		if !strings.Contains(ip, "/") {
			ip += "/32"
		}
		cidrs = append(cidrs, ip)
	}
	return cidrs
}
