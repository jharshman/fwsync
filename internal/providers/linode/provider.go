package linode

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jharshman/fwsync/internal/providers/generic"
	"github.com/linode/linodego"
)

// Client is an implementation of generic.Provider for Akamai Linode.
type Client struct {
	conn *linodego.Client
}

// New returns a new Client.
func New() (*Client, error) {
	conn, err := linodego.NewClientFromEnv(http.DefaultClient)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

// List will list all firewalls present in the account. It returns an unfiltered list of firewall names.
func (c Client) List(ctx context.Context) ([]generic.Firewall, error) {
	fw, err := c.conn.ListFirewalls(ctx, nil)
	if err != nil {
		return nil, err
	}

	// process firewalls
	fws := make([]generic.Firewall, 0, len(fw))
	for _, v := range fw {
		fws = append(fws, generic.Firewall{
			Name:                 v.Label,
			AllowedIPv4Addresses: *v.Rules.Inbound[0].Addresses.IPv4, // this is sketch since it may or may not exist at that index.
			Misc:                 map[string]any{"id": v.ID},
		})
	}

	return fws, nil
}

// Get searches for a specific firewall in the account by name. It then returns a pointer to a *generic.Firewall which
// is a aggressively paired down representation of a Cloud Provider Firewall.
func (c Client) Get(ctx context.Context, name string) (*generic.Firewall, error) {
	// linode's GetFirewall uses the firewall's ID instead of the firewall's name. Instead of doing that, I filter
	// here by the name (label) of the firewall using the Filter field in ListOptions.
	fw, err := c.conn.ListFirewalls(ctx, &linodego.ListOptions{Filter: fmt.Sprintf(`{ "label": "%s" }`, name)})
	if err != nil {
		return nil, err
	}

	if len(fw) == 0 {
		return nil, fmt.Errorf("no firewall found matching filter: label:%s AND is:firewall", name)
	}

	if len(fw) > 1 {
		return nil, fmt.Errorf("more than one firewall matching filter: label:%s AND is:firewall", name)
	}

	return &generic.Firewall{
		Name:                 fw[0].Label,
		Misc:                 map[string]any{"id": fw[0].ID},
		AllowedIPv4Addresses: *fw[0].Rules.Inbound[0].Addresses.IPv4,
	}, nil
}

// Update will update the given firewall rule with the provided IPs in sourceRanges.
func (c Client) Update(ctx context.Context, name string, sourceRanges []string) error {
	// get the firewall by name
	fw, err := c.Get(ctx, name)
	if err != nil {
		return err
	}

	id, ok := fw.Misc["id"].(int)
	if !ok {
		return fmt.Errorf("no id found")
	}

	_, err = c.conn.UpdateFirewallRules(ctx, id, linodego.FirewallRuleSet{
		InboundPolicy:  "ACCEPT",
		OutboundPolicy: "ACCEPT",
		Inbound: []linodego.FirewallRule{
			{
				Action:   "ACCEPT",
				Label:    fw.Name,
				Protocol: "TCP",
				Addresses: linodego.NetworkAddresses{
					IPv4: &sourceRanges,
				},
			},
		}})

	return err
}
