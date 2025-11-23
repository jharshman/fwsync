package generic

import "context"

// Firewaller describes the behavior that a provider should implement in order to
// be usable by fwsync.
type Firewaller interface {
	List(ctx context.Context) ([]string, error)
	Get(ctx context.Context, name string) (*Firewall, error)
	Update(ctx context.Context, name string, sourceRanges []string) error
}

// Firewall hyper-simplifies what a firewall is. Fwsync only needs to use a small subset of the full
// definition a provider's firewall for Get and List operations.
type Firewall struct {
	Name       string
	UUID       int // used when the firewall has a unique identifier associated with it
	AllowedIPs []string
}
