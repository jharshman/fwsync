package generic

import (
	"context"
)

// Provider describes the behavior that a provider should implement in order to
// be usable by fwsync.
type Provider interface {
	// List returns every firewall visible to the authenticated account or project.
	List(ctx context.Context) ([]Firewall, error)
	// Get returns the firewall identified by name. AllowedIPv4Addresses should be
	// empty, not an error, if fwsync has not managed the firewall yet.
	Get(ctx context.Context, name string) (*Firewall, error)
	// Update sets the source addresses allowed by the named firewall to cidrs.
	// Implementations must only replace what fwsync manages and leave any other
	// rules or settings on the firewall untouched.
	Update(ctx context.Context, name string, cidrs []string) error
}

// Firewall is a general type to represent a unique firewall from any provider implementing the Provider interface.
type Firewall struct {
	// Name of the firewall
	Name string
	// Allowed IPv4 Addresses in CIDR notation
	AllowedIPv4Addresses []string
}
