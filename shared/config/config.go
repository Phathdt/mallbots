package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token   TokenConfig   `yaml:"token"`
	Solver  SolverConfig  `yaml:"solver"`
	Watcher WatcherConfig `yaml:"watcher"`
}

type TokenConfig struct {
	TokenURL      string `yaml:"token_url"`
	CacheDuration int    `yaml:"cache_duration"`
}

type SolverConfig struct {
	SolverURL string `yaml:"solver_url"`
}

type WatcherConfig struct {
	RpcURL          string `yaml:"rpc_url"`
	ContractAddress string `yaml:"contract_address"`
	RedisURL        string `yaml:"redis_url"`
	PollingDelay    int    `yaml:"polling_delay"`
}

// LoadConfig reads and parses the YAML configuration file
func LoadConfig(configPath string) (*Config, error) {
	// If configPath is empty, use default path
	if configPath == "" {
		configPath = "config/config.yml"
	}

	// Ensure the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	// Read the config file
	configData, err := os.ReadFile(filepath.Clean(configPath))
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse the YAML
	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("error parsing config YAML: %w", err)
	}

	return &config, nil
}
