package config

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/jharshman/fwsync/internal/providers/gcp"
	"github.com/jharshman/fwsync/internal/providers/generic"
	"github.com/jharshman/fwsync/internal/providers/linode"
	"gopkg.in/yaml.v2"
)

const (
	defaultIPLimit = 5
	ipURL          = "https://ipv4.icanhazip.com"
)

var (
	// providers
	ProviderGoogle = "google"
	ProviderLinode = "linode"

	// todo: implement the following providers
	//providerAWS          = "amazon"
	//providerAzure        = "azure"
	//providerDigitalOcean = "digitalocean"
)

// Config describes the fwsync configuration. It is used to hold basic information about the
// firewall and the desired IPs that are to be allowed.
type Config struct {
	Provider  string   `yaml:"provider"`
	Project   string   `yaml:"project,omitempty"`
	IPLimit   int      `yaml:"ip_limit,omitempty"`
	Name      string   `yaml:"name"`
	SourceIPs []string `yaml:"ips"`
}

// New creates a new Config and returns a pointer to it.
func New(opts ...configOpts) *Config {
	cfg := &Config{}
	cfg.IPLimit = defaultIPLimit // always set the IPLimit equal to the defaultIPLimit.

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// LoadFromFile creates a new Config from the .fwsync configuration file.
func LoadFromFile(r io.Reader) (*Config, error) {
	config := &Config{}
	err := yaml.NewDecoder(r).Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// AuthForProvider authenticates for a given supported Cloud Provider and returns the
// provider's implementation of generic.Firewaller.
func (c *Config) AuthForProvider() (generic.Firewaller, error) {
	var client generic.Firewaller
	var err error
	switch c.Provider {
	case ProviderGoogle:
		client, err = gcp.New(c.Project)
	case ProviderLinode:
		client, err = linode.New()
	default:
		err = fmt.Errorf("invalid provider: %s", c.Provider)
	}
	return client, err
}

// Write will write the fwsync configuration from memory to disk.
func (c *Config) Write(w io.Writer) error {
	enc := yaml.NewEncoder(w)
	err := enc.Encode(c)
	defer enc.Close()
	return err
}

// HasIP checks if the current configuration has a given IP. Returns the index of the IP and true if found.
// Returns -1 and false if not found.
func (c *Config) HasIP(ip string) (int, bool) {
	for idx, sip := range c.SourceIPs {
		if sip == ip {
			return idx, true
		}
	}
	return -1, false
}

// Add will add the given IP to the configuration file.
// If the new IP puts the number of IPs held in the configuration file
// over the limit defined by ipLimit then the oldest IP is removed.
func (c *Config) Add(ip string) {
	if ip == "" {
		return
	}
	if c.IPLimit == 0 {
		// ip limit cannot be 0
		c.IPLimit = defaultIPLimit
	}
	if len(c.SourceIPs) >= c.IPLimit {
		c.SourceIPs = c.SourceIPs[1:]
	}
	c.SourceIPs = append(c.SourceIPs, ip)
}

// Remove will remove an IP from the configuration.
func (c *Config) Remove(ip string) {
	if ip == "" {
		return
	}
	idx, ipExists := c.HasIP(ip)
	if !ipExists {
		return
	}

	chunkOne := c.SourceIPs[:idx]
	chunkTwo := c.SourceIPs[idx+1:]
	var newIPs []string
	newIPs = append(newIPs, chunkOne...)
	newIPs = append(newIPs, chunkTwo...)

	c.SourceIPs = newIPs
}

// PublicIP fetches the current public IP from the URL defined by ipURL.
// On error, it will return an empty string and error.
func PublicIP() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", ipURL, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}
