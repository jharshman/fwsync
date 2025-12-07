package generic

import (
	"context"
)

// Provider describes the behavior that a provider should implement in order to
// be usable by fwsync.
type Provider interface {
	List(ctx context.Context) ([]Firewall, error)
	Get(ctx context.Context, name string) (*Firewall, error)
	Update(ctx context.Context, name string, sourceRanges []string) error
}

// Firewall is a general type to represent a unique firewall from any provider implementing the Provider interface.
type Firewall struct {
	// Name of the firewall
	Name string
	// Allowed IPv4 Addresses
	AllowedIPv4Addresses []string
	// Misc key/value pair field. Any extra information needed by the Provider implementation to
	// perform the basic firewall operations can be stored here.
	Misc map[string]any
}
