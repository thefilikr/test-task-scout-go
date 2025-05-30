package config

import (
	"fmt"
	"os"
)

type Config struct {
	RepositoryType string
	DatabasePath   string
	Port           string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	cfg.RepositoryType = os.Getenv("REPOSITORY_TYPE")
	if cfg.RepositoryType == "" {
		cfg.RepositoryType = "inmemory" 
	}

	if cfg.RepositoryType != "inmemory" && cfg.RepositoryType != "sqlite" {
		return nil, fmt.Errorf("unknown repository type: %s. Use 'inmemory' or 'sqlite'.", cfg.RepositoryType)
	}

	if cfg.RepositoryType == "sqlite" {
		cfg.DatabasePath = os.Getenv("DATABASE_PATH")
		if cfg.DatabasePath == "" {
			cfg.DatabasePath = "./quotes.db" 
		}
	}

	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = "8000"
	}

	return cfg, nil
} 