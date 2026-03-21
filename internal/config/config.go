package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// Account represents a single Kite account
type Account struct {
	APIKey      string `yaml:"api_key"`
	APISecret   string `yaml:"api_secret"`
	AccessToken string `yaml:"access_token"`
	UserName    string `yaml:"user_name,omitempty"`
}

// MultiConfig represents the multi-account configuration
type MultiConfig struct {
	CurrentAccount string              `yaml:"current_account"`
	Accounts       map[string]*Account `yaml:"accounts"`
}

// Config is kept for backward compatibility with existing code
type Config struct {
	APIKey      string
	APISecret   string
	AccessToken string
	UserName    string
	Alias       string
}

// LoadMulti loads the full multi-account configuration
func LoadMulti() (*MultiConfig, error) {
	configPath := GetConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &MultiConfig{
				Accounts: make(map[string]*Account),
			}, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Try to parse as multi-account config
	var multiCfg MultiConfig
	if err := yaml.Unmarshal(data, &multiCfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Check if this is old format (no accounts key, has flat api_key)
	if multiCfg.Accounts == nil {
		multiCfg, err = migrateOldConfig(data)
		if err != nil {
			return nil, err
		}
	}

	return &multiCfg, nil
}

// migrateOldConfig converts old single-account config to new multi-account format
func migrateOldConfig(data []byte) (MultiConfig, error) {
	var oldCfg struct {
		APIKey      string `yaml:"api_key"`
		APISecret   string `yaml:"api_secret"`
		AccessToken string `yaml:"access_token"`
	}

	if err := yaml.Unmarshal(data, &oldCfg); err != nil {
		return MultiConfig{}, fmt.Errorf("failed to parse old config: %w", err)
	}

	// If old config has credentials, migrate them
	if oldCfg.APIKey != "" {
		multiCfg := MultiConfig{
			CurrentAccount: "default",
			Accounts: map[string]*Account{
				"default": {
					APIKey:      oldCfg.APIKey,
					APISecret:   oldCfg.APISecret,
					AccessToken: oldCfg.AccessToken,
				},
			},
		}

		// Save migrated config
		if err := SaveMulti(&multiCfg); err != nil {
			return MultiConfig{}, fmt.Errorf("failed to save migrated config: %w", err)
		}

		return multiCfg, nil
	}

	return MultiConfig{
		Accounts: make(map[string]*Account),
	}, nil
}

// Load loads the current account configuration (backward compatible)
func Load() (*Config, error) {
	// Environment variables override config file (more secure)
	envKey := os.Getenv("KITE_API_KEY")
	envSecret := os.Getenv("KITE_API_SECRET")
	envToken := os.Getenv("KITE_ACCESS_TOKEN")

	if envKey != "" {
		return &Config{
			APIKey:      envKey,
			APISecret:   envSecret,
			AccessToken: envToken,
		}, nil
	}

	multiCfg, err := LoadMulti()
	if err != nil {
		return nil, err
	}

	if multiCfg.CurrentAccount == "" {
		return nil, fmt.Errorf("no account configured. Run 'kite auth' to add an account")
	}

	account, exists := multiCfg.Accounts[multiCfg.CurrentAccount]
	if !exists {
		return nil, fmt.Errorf("current account '%s' not found. Run 'kite alias' to see available accounts", multiCfg.CurrentAccount)
	}

	return &Config{
		APIKey:      account.APIKey,
		APISecret:   account.APISecret,
		AccessToken: account.AccessToken,
		UserName:    account.UserName,
		Alias:       multiCfg.CurrentAccount,
	}, nil
}

// SaveMulti saves the full multi-account configuration
func SaveMulti(cfg *MultiConfig) error {
	configPath := GetConfigPath()

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// SaveAccount saves or updates an account
func SaveAccount(alias string, account *Account) error {
	multiCfg, err := LoadMulti()
	if err != nil {
		return err
	}

	if multiCfg.Accounts == nil {
		multiCfg.Accounts = make(map[string]*Account)
	}

	multiCfg.Accounts[alias] = account

	// If this is the first account, set it as current
	if multiCfg.CurrentAccount == "" {
		multiCfg.CurrentAccount = alias
	}

	return SaveMulti(multiCfg)
}

// Save saves the current account (backward compatible)
func Save(cfg *Config) error {
	alias := cfg.Alias
	if alias == "" {
		alias = "default"
	}

	return SaveAccount(alias, &Account{
		APIKey:      cfg.APIKey,
		APISecret:   cfg.APISecret,
		AccessToken: cfg.AccessToken,
		UserName:    cfg.UserName,
	})
}

// UpdateAccessToken updates the access token for a specific account
func UpdateAccessToken(alias, accessToken, userName string) error {
	multiCfg, err := LoadMulti()
	if err != nil {
		return err
	}

	account, exists := multiCfg.Accounts[alias]
	if !exists {
		return fmt.Errorf("account '%s' not found", alias)
	}

	account.AccessToken = accessToken
	if userName != "" {
		account.UserName = userName
	}

	return SaveMulti(multiCfg)
}

// ListAccounts returns a sorted list of account aliases
func ListAccounts() ([]string, string, error) {
	multiCfg, err := LoadMulti()
	if err != nil {
		return nil, "", err
	}

	aliases := make([]string, 0, len(multiCfg.Accounts))
	for alias := range multiCfg.Accounts {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)

	return aliases, multiCfg.CurrentAccount, nil
}

// GetAccount returns a specific account by alias
func GetAccount(alias string) (*Account, error) {
	multiCfg, err := LoadMulti()
	if err != nil {
		return nil, err
	}

	account, exists := multiCfg.Accounts[alias]
	if !exists {
		return nil, fmt.Errorf("account '%s' not found", alias)
	}

	return account, nil
}

// SetCurrentAccount switches the current account
func SetCurrentAccount(alias string) error {
	multiCfg, err := LoadMulti()
	if err != nil {
		return err
	}

	if _, exists := multiCfg.Accounts[alias]; !exists {
		return fmt.Errorf("account '%s' not found. Run 'kite alias' to see available accounts", alias)
	}

	multiCfg.CurrentAccount = alias
	return SaveMulti(multiCfg)
}

// RemoveAccount removes an account by alias
func RemoveAccount(alias string) error {
	multiCfg, err := LoadMulti()
	if err != nil {
		return err
	}

	if _, exists := multiCfg.Accounts[alias]; !exists {
		return fmt.Errorf("account '%s' not found", alias)
	}

	delete(multiCfg.Accounts, alias)

	// If we removed the current account, switch to another one
	if multiCfg.CurrentAccount == alias {
		multiCfg.CurrentAccount = ""
		for a := range multiCfg.Accounts {
			multiCfg.CurrentAccount = a
			break
		}
	}

	return SaveMulti(multiCfg)
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kite.yaml")
}

// AccountExists checks if an account with the given alias exists
func AccountExists(alias string) bool {
	multiCfg, err := LoadMulti()
	if err != nil {
		return false
	}
	_, exists := multiCfg.Accounts[alias]
	return exists
}
