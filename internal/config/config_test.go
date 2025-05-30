package config_test

import (
	"os"
	"test-task-scout-go/internal/config"
	"testing"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	os.Setenv(key, value)
	t.Cleanup(func() {
		os.Unsetenv(key)
	})
}

func TestLoadConfig_Defaults(t *testing.T) {
	os.Unsetenv("REPOSITORY_TYPE")
	os.Unsetenv("DATABASE_PATH")
	os.Unsetenv("PORT")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %v", err)
	}

	if cfg.RepositoryType != "inmemory" {
		t.Errorf("Expected RepositoryType 'inmemory', got '%s'", cfg.RepositoryType)
	}
	if cfg.DatabasePath != "" {
		t.Errorf("Expected empty DatabasePath, got '%s'", cfg.DatabasePath)
	}
	if cfg.Port != "8000" {
		t.Errorf("Expected Port '8000', got '%s'", cfg.Port)
	}
}

func TestLoadConfig_EnvironmentVariables(t *testing.T) {
	setEnv(t, "REPOSITORY_TYPE", "sqlite")
	setEnv(t, "DATABASE_PATH", "/path/to/db.sqlite")
	setEnv(t, "PORT", "9000")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %v", err)
	}

	if cfg.RepositoryType != "sqlite" {
		t.Errorf("Expected RepositoryType 'sqlite', got '%s'", cfg.RepositoryType)
	}
	if cfg.DatabasePath != "/path/to/db.sqlite" {
		t.Errorf("Expected DatabasePath '/path/to/db.sqlite', got '%s'", cfg.DatabasePath)
	}
	if cfg.Port != "9000" {
		t.Errorf("Expected Port '9000', got '%s'", cfg.Port)
	}
}

func TestLoadConfig_SQLiteDefaultDBPath(t *testing.T) {
	setEnv(t, "REPOSITORY_TYPE", "sqlite")
	os.Unsetenv("DATABASE_PATH")
	os.Unsetenv("PORT")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned an error: %v", err)
	}

	if cfg.RepositoryType != "sqlite" {
		t.Errorf("Expected RepositoryType 'sqlite', got '%s'", cfg.RepositoryType)
	}
	if cfg.DatabasePath != "./quotes.db" {
		t.Errorf("Expected DatabasePath './quotes.db', got '%s'", cfg.DatabasePath)
	}
	if cfg.Port != "8000" {
		t.Errorf("Expected Port '8000', got '%s'", cfg.Port)
	}
}

func TestLoadConfig_InvalidRepositoryType(t *testing.T) {
	setEnv(t, "REPOSITORY_TYPE", "postgres")

	cfg, err := config.LoadConfig()
	if err == nil {
		t.Error("LoadConfig did not return an error for invalid repository type")
	}
	if cfg != nil {
		t.Errorf("Expected nil config for invalid type, got %+v", cfg)
	}

	expectedErr := "unknown repository type: postgres. Use 'inmemory' or 'sqlite'."
	if err != nil && err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', got '%s'", expectedErr, err.Error())
	}
} 