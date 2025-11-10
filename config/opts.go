package config

type configOpts func(*Config)

// WithProvider sets the Provider for the fwsync configuration.
func WithProvider(provider string) configOpts {
	return func(cfg *Config) {
		cfg.Provider = provider
	}
}

// WithProject sets the Project for the fwsync configuration.
func WithProject(project string) configOpts {
	return func(cfg *Config) {
		cfg.Project = project
	}
}

// WithFirewall sets the Firewall's name in the fwsync configuration.
func WithFirewall(name string) configOpts {
	return func(cfg *Config) {
		cfg.Name = name
	}
}

// WithSourceIPs sets the allowed IPs in the fwsync configuration.
func WithSourceIPs(sourceIPs ...string) configOpts {
	return func(cfg *Config) {
		// Silently remove the extra IPs.
		// This probably isn't the best solution because it doesn't give any
		// feedback to the caller.
		if len(sourceIPs) > ipLimit {
			sourceIPs = sourceIPs[:ipLimit]
		}
		cfg.SourceIPs = sourceIPs
	}
}
