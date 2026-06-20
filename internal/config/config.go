package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return home + "/" + configFileName, nil
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config

	err = json.Unmarshal(file, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c *Config) SetUser(user string) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	c.CurrentUserName = user

	newConfig, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, newConfig, 0644)
	if err != nil {
		return err
	}

	return nil
}
