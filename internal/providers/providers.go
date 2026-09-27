// Package providers maps provider names, as written in the fwsync configuration file,
// to their generic.Provider implementations.
//
// Adding a provider means implementing generic.Provider in its own package under
// internal/providers and adding a single entry to the registry below.
package providers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jharshman/fwsync/internal/providers/gcp"
	"github.com/jharshman/fwsync/internal/providers/generic"
	"github.com/jharshman/fwsync/internal/providers/linode"
)

// Provider names as they appear in the configuration file and the --provider flag.
const (
	Google = "google"
	Linode = "linode"
)

// Settings holds the provider specific values stored in the configuration file.
// Providers ignore the fields they don't use.
type Settings struct {
	Project string
}

type constructor func(s Settings) (generic.Provider, error)

var registry = map[string]constructor{
	Google: func(s Settings) (generic.Provider, error) {
		c, err := gcp.New(s.Project)
		if err != nil {
			return nil, err
		}
		return c, nil
	},
	Linode: func(Settings) (generic.Provider, error) {
		c, err := linode.New()
		if err != nil {
			return nil, err
		}
		return c, nil
	},
}

// New authenticates with the named provider and returns its generic.Provider implementation.
func New(name string, s Settings) (generic.Provider, error) {
	newProvider, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("invalid provider: %q (supported: %s)", name, strings.Join(Names(), ", "))
	}
	return newProvider(s)
}

// Names returns the names of all supported providers in sorted order.
func Names() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
