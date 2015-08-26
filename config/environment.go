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

var JsonConfig = Config{"", ""}

func (cfg Config) String() string {
	return cfg.ListenAddress + "; " + cfg.DatabaseUri
}
