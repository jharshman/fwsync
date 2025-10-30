package user

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"
)

const IPLimit = 5

type Config struct {
	Project   string   `yaml:"project,omitempty"`
	Name      string   `yaml:"name"`
	SourceIPs []string `yaml:"ips"`
}

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

func NewFromFile(r io.Reader) (*Config, error) {
	config := &Config{}
	err := yaml.NewDecoder(r).Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) Write(w io.Writer) error {
	enc := yaml.NewEncoder(w)
	err := enc.Encode(c)
	defer enc.Close()
	return err
}

func (c *Config) HasIP(ip string) (int, bool) {
	for idx, sip := range c.SourceIPs {
		if sip == ip {
			return idx, true
		}
	}
	return -1, false
}

func (c *Config) Add(ip string) {
	if ip == "" {
		return
	}
	if len(c.SourceIPs) >= IPLimit {
		c.SourceIPs = c.SourceIPs[1:]
	}
	c.SourceIPs = append(c.SourceIPs, ip)
}

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

func PublicIP() (string, error) {
	// use https://icanhazip.com/
	// reliable and owned and operated by cloudflare.
	res, err := http.Get("https://ipv4.icanhazip.com")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	ip := strings.TrimSpace(string(body))
	return ip, nil
}
