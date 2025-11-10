package cmd

import (
	"github.com/jharshman/fwsync/internal/providers/generic"
)

var (
	// FirewallClient holds the initialized provider client implementing the Firewaller interface.
	FirewallClient generic.Firewaller
)
