package config

// Constants for environment.
const (
	// DEVELOPMENT | TEST | PRODUCTION
	Environment = "DEVELOPMENT"
	// Environment            = "PRODUCTION"
)

type Config struct {
	ListenAddress string
	DatabaseUri   string
}

func (cfg Config) String() string {
	return cfg.ListenAddress + "; " + cfg.DatabaseUri
}
