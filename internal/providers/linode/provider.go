package linode

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/jharshman/fwsync/internal/providers/generic"
	"github.com/linode/linodego"
)

const (
	// ruleLabel is the label of the inbound rule that fwsync creates and manages.
	ruleLabel = "fwsync"

	policyAccept = "ACCEPT"
)

// Client is an implementation of generic.Provider for Akamai Linode.
//
// fwsync manages a single inbound rule on the firewall, identified by its label. Every other
// rule and the firewall's inbound/outbound policies are left as they are.
type Client struct {
	conn *linodego.Client
}

// New returns a new Client authenticated with the token in the LINODE_TOKEN environment variable.
func New() (*Client, error) {
	conn, err := linodego.NewClientFromEnv(http.DefaultClient)
	if err != nil {
		return nil, fmt.Errorf("linode: %w (set the LINODE_TOKEN environment variable)", err)
	}

	return &Client{conn: conn}, nil
}

// List will list all firewalls present in the account.
func (c *Client) List(ctx context.Context) ([]generic.Firewall, error) {
	fws, err := c.conn.ListFirewalls(ctx, nil)
	if err != nil {
		return nil, err
	}

	out := make([]generic.Firewall, 0, len(fws))
	for _, fw := range fws {
		out = append(out, toGeneric(fw))
	}
	return out, nil
}

// Get returns the firewall with the given label.
func (c *Client) Get(ctx context.Context, name string) (*generic.Firewall, error) {
	fw, err := c.find(ctx, name)
	if err != nil {
		return nil, err
	}
	g := toGeneric(*fw)
	return &g, nil
}

// Update sets the IPv4 addresses of the fwsync managed inbound rule to cidrs, creating the rule
// if it does not exist yet.
func (c *Client) Update(ctx context.Context, name string, cidrs []string) error {
	fw, err := c.find(ctx, name)
	if err != nil {
		return err
	}

	rules, err := c.conn.GetFirewallRules(ctx, fw.ID)
	if err != nil {
		return err
	}

	if _, err := c.conn.UpdateFirewallRules(ctx, fw.ID, withManagedRule(*rules, fw.Label, cidrs)); err != nil {
		return err
	}

	if rules.InboundPolicy == policyAccept {
		fmt.Fprintf(os.Stderr, "warning: firewall %q has an inbound policy of ACCEPT, so all inbound traffic is allowed "+
			"regardless of the IPs fwsync manages. Set the inbound policy to DROP to restrict access.\n", fw.Label)
	}
	return nil
}

// find looks up a firewall by label. Linode addresses firewalls by ID, but fwsync stores the label.
func (c *Client) find(ctx context.Context, label string) (*linodego.Firewall, error) {
	filter, err := json.Marshal(map[string]string{"label": label})
	if err != nil {
		return nil, err
	}

	fws, err := c.conn.ListFirewalls(ctx, &linodego.ListOptions{Filter: string(filter)})
	if err != nil {
		return nil, err
	}

	switch len(fws) {
	case 0:
		return nil, fmt.Errorf("no firewall found with label: %s", label)
	case 1:
		return &fws[0], nil
	default:
		return nil, fmt.Errorf("more than one firewall found with label: %s", label)
	}
}

func toGeneric(fw linodego.Firewall) generic.Firewall {
	g := generic.Firewall{Name: fw.Label}
	if i := managedRuleIndex(fw.Rules.Inbound, fw.Label); i >= 0 && fw.Rules.Inbound[i].Addresses.IPv4 != nil {
		g.AllowedIPv4Addresses = *fw.Rules.Inbound[i].Addresses.IPv4
	}
	return g
}

// managedRuleIndex returns the index of the fwsync managed rule in rules, or -1 if there isn't one.
// Versions of fwsync up to v0.0.3 labeled the rule with the firewall's own label, so that is
// recognized as well.
func managedRuleIndex(rules []linodego.FirewallRule, firewallLabel string) int {
	for i, r := range rules {
		if r.Label == ruleLabel || r.Label == firewallLabel {
			return i
		}
	}
	return -1
}

// withManagedRule returns a copy of rules with the managed inbound rule's IPv4 addresses set to
// cidrs. If there is no managed rule, one allowing TCP on all ports is appended.
func withManagedRule(rules linodego.FirewallRuleSet, firewallLabel string, cidrs []string) linodego.FirewallRuleSet {
	addrs := append([]string(nil), cidrs...)
	inbound := append([]linodego.FirewallRule(nil), rules.Inbound...)

	if i := managedRuleIndex(inbound, firewallLabel); i >= 0 {
		inbound[i].Addresses.IPv4 = &addrs
	} else {
		inbound = append(inbound, linodego.FirewallRule{
			Action:      policyAccept,
			Label:       ruleLabel,
			Description: "Managed by fwsync",
			Protocol:    linodego.TCP,
			Addresses:   linodego.NetworkAddresses{IPv4: &addrs},
		})
	}

	rules.Inbound = inbound
	return rules
}
