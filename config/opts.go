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
		if len(sourceIPs) > cfg.IPLimit {
			sourceIPs = sourceIPs[:cfg.IPLimit]
		}
		cfg.SourceIPs = sourceIPs
	}
}

// WithIPLimit sets the number of allowed IPs.
func WithIPLimit(limit int) configOpts {
	return func(cfg *Config) {
		cfg.IPLimit = limit
	}
}
