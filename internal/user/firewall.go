package user

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// IPLimit enforces that there only be 5 IPs at any given time
// configured on the firewall.
const IPLimit = 5

// IPURL use https://icanhazip.com/.
// Reliable and owned and operated by cloudflare.
const IPURL = "https://ipv4.icanhazip.com"

// Config provides the structure of the local configuration file for fwsync.
type Config struct {
	Project   string   `yaml:"project,omitempty"`
	Name      string   `yaml:"name"`
	SourceIPs []string `yaml:"ips"`
}

// NewConfig creates a new fwsync configuration.
func NewConfig(name string, ips ...string) (*Config, error) {
	if len(ips) == 0 {
		return nil, errors.New("must have at least one IP")
	}
	if len(ips) > IPLimit {
		ips = ips[:IPLimit]
	}
	return &Config{
		Name:      name,
		SourceIPs: ips,
	}, nil
}

// NewFromFile loads an existing configuration from file.
func NewFromFile(r io.Reader) (*Config, error) {
	config := &Config{}
	err := yaml.NewDecoder(r).Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
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
// over the limit defined by IPLimit then the oldest IP is removed.
func (c *Config) Add(ip string) {
	if ip == "" {
		return
	}
	if len(c.SourceIPs) >= IPLimit {
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

// PublicIP fetches the current public IP from the URL defined by IPURL.
// On error, it will return an empty string and error.
func PublicIP() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", IPURL, nil)
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
