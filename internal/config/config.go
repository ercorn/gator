package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const config_file_name = ".gatorconfig.json"

func Read() *Config {
	path := getConfigFilePath() + "/"
	data, err := os.ReadFile(path + config_file_name)
	if err != nil {
		log.Fatal(err)
	}
	cfg := Config{}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}

func getConfigFilePath() string {
	path, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func (cfg *Config) SetUser(name string) {
	cfg.CurrentUserName = name
	data, err := json.Marshal(cfg)
	if err != nil {
		log.Fatal(err)
	}
	path := getConfigFilePath()
	err = os.WriteFile(path+"/"+config_file_name, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
