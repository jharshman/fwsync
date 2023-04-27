package user

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/matryer/is"
)

func TestNewConfig(t *testing.T) {

	tests := []struct {
		description string
		argFwID     string
		argIPs      []string
		expected    *Config
		hasError    bool
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
				Name: "firstname-lastname-firewall-rule",
				SourceIPs: []string{
					"1.1.1.1",
					"2.2.2.2",
					"3.3.3.3",
					"4.4.4.4",
					"5.5.5.5",
				},
			},
			hasError: false,
		},
		{
			description: "new fwsync config min ips",
			argFwID:     "firstname-lastname-firewall-rule",
			argIPs: []string{
				"1.1.1.1",
			},
			expected: &Config{
				Name: "firstname-lastname-firewall-rule",
				SourceIPs: []string{
					"1.1.1.1",
				},
			},
			hasError: false,
		},
		{
			description: "new fwsync config no ips",
			argFwID:     "firstname-lastname-firewall-rule",
			argIPs:      []string{},
			expected:    nil,
			hasError:    true,
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
				Name: "firstname-lastname-firewall-rule",
				SourceIPs: []string{
					"1.1.1.1",
					"2.2.2.2",
					"3.3.3.3",
					"4.4.4.4",
					"5.5.5.5",
				},
			},
			hasError: true,
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Log(tc.description)
			is := is.New(t)
			got, err := NewConfig(tc.argFwID, tc.argIPs...)
			if !tc.hasError {
				is.NoErr(err)
			}
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
	got, err := NewFromFile(buf)
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
	expected := []byte(`name: firstname-lastname-firewall-rule
ips:
- 1.1.1.1
- 2.2.2.2
- 3.3.3.3
`)
	cfg, _ := NewConfig("firstname-lastname-firewall-rule", "1.1.1.1", "2.2.2.2", "3.3.3.3")
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

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Log(tc.description)
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
		expect      *Config
	}{
		{
			description: "Empty parameter",
			ip:          "",
			cfg: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1"},
			},
			expect: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1"},
			},
		},
		{
			description: "Add IP",
			ip:          "2.2.2.2",
			cfg: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1"},
			},
			expect: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1", "2.2.2.2"},
			},
		},
		{
			description: "Add IP on full list",
			ip:          "6.6.6.6",
			cfg: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3", "4.4.4.4", "5.5.5.5"},
			},
			expect: &Config{
				Name:      "firstname-lastname-firewall-rule",
				SourceIPs: []string{"2.2.2.2", "3.3.3.3", "4.4.4.4", "5.5.5.5", "6.6.6.6"},
			},
		},
	}

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Log(tc.description)
			is := is.New(t)
			tc.cfg.Add(tc.ip)
			is.Equal(tc.cfg, tc.expect)
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

	for i, tc := range tests {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Log(tc.description)
			is := is.New(t)
			tc.cfg.Remove(tc.ip)
			is.Equal(tc.cfg, tc.expect)
		})
	}
}
