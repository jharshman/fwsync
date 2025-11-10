package generic

// Firewaller describes the behavior that a provider should implement in order to
// be usable by fwsync.
type Firewaller interface {
	List() ([]string, error)
	Get(name string) (*Firewall, error)
	Update(name string, sourceRanges []string) error
}

// Firewall hyper-simplifies what a firewall is. Fwsync only needs to use a small subset of the full
// definition a provider's firewall for Get and List operations.
type Firewall struct {
	Name       string
	AllowedIPs []string
}
