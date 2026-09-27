package gcp

import (
	"context"
	"errors"

	"github.com/jharshman/fwsync/internal/providers/generic"
	"google.golang.org/api/compute/v1"
)

// Client is a simple type containing a GCP client connection to compute.Service. This
// type implements the generic.Provider interface.
type Client struct {
	conn    *compute.Service
	project string
}

// New creates a new instance of the Client using Application Default Credentials.
func New(project string) (*Client, error) {
	if project == "" {
		return nil, errors.New("the google provider requires a project, set it with --project")
	}
	// The service holds on to this context for token refreshes, so it must not be one that gets cancelled.
	conn, err := compute.NewService(context.Background())
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, project: project}, nil
}

// List returns all the available Firewall Policies in the Project.
// It distills that information into a simpler generic.Firewall type and
// returns it to the caller.
func (c *Client) List(ctx context.Context) ([]generic.Firewall, error) {
	var fws []generic.Firewall
	err := c.conn.Firewalls.List(c.project).Pages(ctx, func(page *compute.FirewallList) error {
		for _, item := range page.Items {
			fws = append(fws, generic.Firewall{
				Name:                 item.Name,
				AllowedIPv4Addresses: item.SourceRanges,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fws, nil
}

// Get returns a generic.Firewall if one exists by the given name parameter.
func (c *Client) Get(ctx context.Context, name string) (*generic.Firewall, error) {
	fw, err := c.conn.Firewalls.Get(c.project, name).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return &generic.Firewall{
		Name:                 fw.Name,
		AllowedIPv4Addresses: fw.SourceRanges,
	}, nil
}

// Update performs a Patch operation on an existing Firewall and sets the SourceRanges of allowed IPs
// to the provided cidrs.
func (c *Client) Update(ctx context.Context, name string, cidrs []string) error {
	_, err := c.conn.Firewalls.Patch(c.project, name, &compute.Firewall{SourceRanges: cidrs}).Context(ctx).Do()
	return err
}
