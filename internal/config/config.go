package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const config_file_name = ".gatorconfig.json"

func Read() (Config, error) {
	path := getConfigFilePath() + "/"
	data, err := os.ReadFile(path + config_file_name)
	if err != nil {
		return Config{}, err
	}
	cfg := Config{}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func getConfigFilePath() string {
	path, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func (cfg *Config) SetUser(name string) error {
	cfg.CurrentUserName = name
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal user config: %w", err)
	}
	path := getConfigFilePath()
	err = os.WriteFile(path+"/"+config_file_name, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
