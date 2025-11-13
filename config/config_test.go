package config

import (
	"bytes"
	"testing"

	"github.com/matryer/is"
)

func TestNewConfig(t *testing.T) {

	tests := []struct {
		description string
		argFwID     string
		argIpLimit  int
		argIPs      []string
		expected    *Config
	}{
		{
			description: "new fwsync config max ips",
			argFwID:     "firstname-lastname-firewall-rule",
			argIPs: []string{
				"1.1.1.1",
				"2.2.2.2",
				"3.3.3.3",
				"4.4.4.4",
				"5.5.5.5",
			},
			expected: &Config{
				Name:    "firstname-lastname-firewall-rule",
				IPLimit: defaultIPLimit,
				SourceIPs: []string{
					"1.1.1.1",
					"2.2.2.2",
					"3.3.3.3",
					"4.4.4.4",
					"5.5.5.5",
				},
			},
		},
		{
			description: "new fwsync config min ips",
			argFwID:     "firstname-lastname-firewall-rule",
			argIPs: []string{
				"1.1.1.1",
			},
			expected: &Config{
				Name:    "firstname-lastname-firewall-rule",
				IPLimit: defaultIPLimit,
				SourceIPs: []string{
					"1.1.1.1",
				},
			},
		},
		{
			description: "new fwsync config no ips",
			argFwID:     "firstname-lastname-firewall-rule",
			argIPs:      nil,
			expected: &Config{
				Name:    "firstname-lastname-firewall-rule",
				IPLimit: defaultIPLimit,
			},
		},
		{
			description: "new fwsync config too many ips",
			argFwID:     "firstname-lastname-firewall-rule",
			argIPs: []string{
				"1.1.1.1",
				"2.2.2.2",
				"3.3.3.3",
				"4.4.4.4",
				"5.5.5.5",
				"6.6.6.6.",
			},
			expected: &Config{
				Name:    "firstname-lastname-firewall-rule",
				IPLimit: defaultIPLimit,
				SourceIPs: []string{
					"1.1.1.1",
					"2.2.2.2",
					"3.3.3.3",
					"4.4.4.4",
					"5.5.5.5",
				},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			is := is.New(t)
			got := New(
				WithFirewall(tc.argFwID),
				WithSourceIPs(tc.argIPs...))

			is.Equal(got, tc.expected)
		})
	}
}

func TestNewFromFile(t *testing.T) {
	in := []byte(`
name: firstname-lastname-firewall-rule
ips:
  - 1.1.1.1
  - 2.2.2.2
  - 3.3.3.3
`)
	is := is.New(t)
	buf := bytes.NewBuffer(in)
	got, err := LoadFromFile(buf)
	is.NoErr(err)
	is.Equal(got, &Config{
		Name: "firstname-lastname-firewall-rule",
		SourceIPs: []string{
			"1.1.1.1",
			"2.2.2.2",
			"3.3.3.3",
		},
	})
}

func TestConfig_Write(t *testing.T) {
	expected := []byte(`provider: google
project: myproject
ip_limit: 5
name: firstname-lastname-firewall-rule
ips:
- 1.1.1.1
- 2.2.2.2
- 3.3.3.3
`)
	cfg := New(
		WithProvider("google"),
		WithProject("myproject"),
		WithFirewall("firstname-lastname-firewall-rule"),
		WithSourceIPs("1.1.1.1", "2.2.2.2", "3.3.3.3"))
	var b []byte
	got := bytes.NewBuffer(b)
	cfg.Write(got)

	is := is.New(t)
	is.Equal(string(expected), got.String())
}

func TestConfig_HasIP(t *testing.T) {

	tests := []struct {
		description string
		ip          string
		cfg         *Config
		expect      bool
	}{
		{
			description: "Config does not contain IP",
			ip:          "1.1.1.1",
			cfg: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"2.2.2.2", "3.3.3.3"},
			},
			expect: false,
		},
		{
			description: "Config contains IP",
			ip:          "1.1.1.1",
			cfg: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			},
			expect: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			is := is.New(t)
			_, got := tc.cfg.HasIP(tc.ip)
			is.Equal(got, tc.expect)
		})
	}

}

func TestConfig_Add(t *testing.T) {

	tests := []struct {
		description string
		ip          string
		cfg         *Config
		expectedIPs []string
	}{
		{
			description: "Empty parameter",
			ip:          "",
			cfg: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1"},
			},
			expectedIPs: []string{"1.1.1.1"},
		},
		{
			description: "Add IP",
			ip:          "2.2.2.2",
			cfg: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1"},
			},
			expectedIPs: []string{"1.1.1.1", "2.2.2.2"},
		},
		{
			description: "Add IP on full list",
			ip:          "6.6.6.6",
			cfg: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3", "4.4.4.4", "5.5.5.5"},
			},
			expectedIPs: []string{"2.2.2.2", "3.3.3.3", "4.4.4.4", "5.5.5.5", "6.6.6.6"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			is := is.New(t)
			tc.cfg.Add(tc.ip)
			is.Equal(tc.cfg.SourceIPs, tc.expectedIPs)
		})
	}

}

func TestConfig_Remove(t *testing.T) {
	tests := []struct {
		description string
		ip          string
		cfg         *Config
		expect      *Config
	}{
		{
			description: "empty param",
			ip:          "",
			cfg: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1", "2.2.2.2"},
			},
			expect: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1", "2.2.2.2"},
			},
		},
		{
			description: "does not contain ip",
			ip:          "3.3.3.3",
			cfg: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1", "2.2.2.2"},
			},
			expect: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1", "2.2.2.2"},
			},
		},
		{
			description: "contains ip at end index",
			ip:          "3.3.3.3",
			cfg: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			},
			expect: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1", "2.2.2.2"},
			},
		},
		{
			description: "contains ip at mid index",
			ip:          "3.3.3.3",
			cfg: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1", "3.3.3.3", "2.2.2.2"},
			},
			expect: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1", "2.2.2.2"},
			},
		},
		{
			description: "contains ip at first index",
			ip:          "3.3.3.3",
			cfg: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"3.3.3.3", "1.1.1.1", "2.2.2.2"},
			},
			expect: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1", "2.2.2.2"},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			is := is.New(t)
			tc.cfg.Remove(tc.ip)
			is.Equal(tc.cfg, tc.expect)
		})
	}
}
