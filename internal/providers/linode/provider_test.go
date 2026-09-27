package linode

import (
	"testing"

	"github.com/linode/linodego"
	"github.com/matryer/is"
)

func addrs(ips ...string) *[]string { return &ips }

func TestWithManagedRule(t *testing.T) {
	ssh := linodego.FirewallRule{Action: "ACCEPT", Label: "ssh-office", Protocol: linodego.TCP, Ports: "22", Addresses: linodego.NetworkAddresses{IPv4: addrs("10.0.0.0/8")}}
	outbound := []linodego.FirewallRule{{Action: "ACCEPT", Label: "all-out", Protocol: linodego.TCP}}

	tests := []struct {
		description string
		rules       linodego.FirewallRuleSet
		expect      []linodego.FirewallRule
	}{
		{
			description: "no managed rule appends one",
			rules:       linodego.FirewallRuleSet{Inbound: []linodego.FirewallRule{ssh}},
			expect: []linodego.FirewallRule{ssh, {
				Action:      "ACCEPT",
				Label:       ruleLabel,
				Description: "Managed by fwsync",
				Protocol:    linodego.TCP,
				Addresses:   linodego.NetworkAddresses{IPv4: addrs("1.1.1.1/32")},
			}},
		},
		{
			description: "managed rule only has its IPv4 addresses replaced",
			rules: linodego.FirewallRuleSet{Inbound: []linodego.FirewallRule{
				ssh,
				{Action: "ACCEPT", Label: ruleLabel, Protocol: linodego.TCP, Ports: "22,443", Addresses: linodego.NetworkAddresses{IPv4: addrs("9.9.9.9/32"), IPv6: addrs("::1/128")}},
			}},
			expect: []linodego.FirewallRule{
				ssh,
				{Action: "ACCEPT", Label: ruleLabel, Protocol: linodego.TCP, Ports: "22,443", Addresses: linodego.NetworkAddresses{IPv4: addrs("1.1.1.1/32"), IPv6: addrs("::1/128")}},
			},
		},
		{
			description: "legacy rule labeled with the firewall label is managed",
			rules: linodego.FirewallRuleSet{Inbound: []linodego.FirewallRule{
				{Action: "ACCEPT", Label: "my-firewall", Protocol: linodego.TCP, Addresses: linodego.NetworkAddresses{IPv4: addrs("9.9.9.9/32")}},
			}},
			expect: []linodego.FirewallRule{
				{Action: "ACCEPT", Label: "my-firewall", Protocol: linodego.TCP, Addresses: linodego.NetworkAddresses{IPv4: addrs("1.1.1.1/32")}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			is := is.New(t)
			tc.rules.InboundPolicy = "DROP"
			tc.rules.OutboundPolicy = "ACCEPT"
			tc.rules.Outbound = outbound

			got := withManagedRule(tc.rules, "my-firewall", []string{"1.1.1.1/32"})

			is.Equal(got.Inbound, tc.expect)
			is.Equal(got.InboundPolicy, "DROP")    // policies are untouched
			is.Equal(got.OutboundPolicy, "ACCEPT") // policies are untouched
			is.Equal(got.Outbound, outbound)       // outbound rules are untouched
		})
	}
}

func TestWithManagedRuleDoesNotMutateInput(t *testing.T) {
	is := is.New(t)
	in := linodego.FirewallRuleSet{Inbound: []linodego.FirewallRule{
		{Label: ruleLabel, Addresses: linodego.NetworkAddresses{IPv4: addrs("9.9.9.9/32")}},
	}}
	withManagedRule(in, "my-firewall", []string{"1.1.1.1/32"})
	is.Equal(*in.Inbound[0].Addresses.IPv4, []string{"9.9.9.9/32"})
}

func TestToGeneric(t *testing.T) {
	is := is.New(t)

	// firewalls without a managed rule, or without any rules, must not panic.
	is.Equal(toGeneric(linodego.Firewall{Label: "empty"}).AllowedIPv4Addresses, nil)
	is.Equal(toGeneric(linodego.Firewall{Label: "no-ipv4", Rules: linodego.FirewallRuleSet{
		Inbound: []linodego.FirewallRule{{Label: ruleLabel}},
	}}).AllowedIPv4Addresses, nil)

	got := toGeneric(linodego.Firewall{Label: "fw", Rules: linodego.FirewallRuleSet{Inbound: []linodego.FirewallRule{
		{Label: "other", Addresses: linodego.NetworkAddresses{IPv4: addrs("10.0.0.0/8")}},
		{Label: ruleLabel, Addresses: linodego.NetworkAddresses{IPv4: addrs("1.1.1.1/32")}},
	}}})
	is.Equal(got.Name, "fw")
	is.Equal(got.AllowedIPv4Addresses, []string{"1.1.1.1/32"})
}
