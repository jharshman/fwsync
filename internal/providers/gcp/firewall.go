package gcp

import (
	"context"

	"github.com/jharshman/fwsync/internal/providers/generic"
	"google.golang.org/api/compute/v1"
)

// Client is a simple type containing a GCP client connection to compute.Service. This
// type implements the generic.Firewaller interface.
type Client struct {
	conn    *compute.Service
	project string
}

// New creates a new instance of the Client.
func New(project string) (*Client, error) {
	conn, err := compute.NewService(context.Background())
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, project: project}, nil
}

// List returns all the available Firewall Policies in the Project.
// It distills that information into a simpler generic.Firewall type and
// returns it to the caller.
func (c *Client) List(ctx context.Context) ([]string, error) {
	fw, err := c.conn.Firewalls.List(c.project).Do()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(fw.Items))
	for _, item := range fw.Items {
		names = append(names, item.Name)
	}

	return names, nil
}

// Get returns a generic.Firewall if one exists by the given name parameter.
func (c *Client) Get(ctx context.Context, name string) (*generic.Firewall, error) {
	fw, err := c.conn.Firewalls.Get(c.project, name).Do()
	if err != nil {
		return nil, err
	}

	userFirewall := &generic.Firewall{
		Name:       fw.Name,
		AllowedIPs: fw.SourceRanges,
	}

	return userFirewall, nil
}

// Update performs a Patch operation on an existing Firewall and sets the SourceRanges of allowed IPs
// to the provided parameter sourceRanges.
func (c *Client) Update(ctx context.Context, name string, sourceRanges []string) error {
	_, err := c.conn.Firewalls.Patch(c.project, name, &compute.Firewall{SourceRanges: sourceRanges}).Do()
	return err
}
