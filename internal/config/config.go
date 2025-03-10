package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting home directory: %w", err)
	}

	return filepath.Join(homeDir, configFileName), nil
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("error getting config file path: %w", err)
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshaling json data: %w", err)
	}

	return cfg, nil
}

func write(data []byte) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error getting config file path: %w", err)
	}

	err = os.WriteFile(configFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing to config file: %w", err)
	}

	return nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	err = write(data)
	if err != nil {
		return fmt.Errorf("error writing user %s to config file: %w", username, err)
	}

	return nil
}
