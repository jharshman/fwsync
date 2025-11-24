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

type Firewall struct {
	Name                 string
	AllowedIPv4Addresses []string
	Misc                 map[string]any
}
