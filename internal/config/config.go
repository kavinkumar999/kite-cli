package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	APIKey      string `mapstructure:"api_key"`
	APISecret   string `mapstructure:"api_secret"`
	AccessToken string `mapstructure:"access_token"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Environment variables override config file (more secure)
	if envKey := os.Getenv("KITE_API_KEY"); envKey != "" {
		cfg.APIKey = envKey
	}
	if envSecret := os.Getenv("KITE_API_SECRET"); envSecret != "" {
		cfg.APISecret = envSecret
	}
	if envToken := os.Getenv("KITE_ACCESS_TOKEN"); envToken != "" {
		cfg.AccessToken = envToken
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api_key not configured. Run 'kite auth' or set KITE_API_KEY env var")
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	viper.Set("api_key", cfg.APIKey)
	viper.Set("api_secret", cfg.APISecret)
	viper.Set("access_token", cfg.AccessToken)

	configPath := filepath.Join(home, ".kite.yaml")
	if err := viper.WriteConfigAs(configPath); err != nil {
		return err
	}

	// Set secure file permissions (owner read/write only)
	return os.Chmod(configPath, 0600)
}

func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kite.yaml")
}
