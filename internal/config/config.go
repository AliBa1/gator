package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("Couldn't find the home directory of this machine")
	}

	return fmt.Sprintf("%s/%s", homeDirectory, configFileName), nil
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, errors.New("Problem opening config file")
	}

	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, errors.New("Problem decoding config into struct")
	}

	return config, nil
}

func write(c Config) error {
	configJSON, err := json.Marshal(c)
	if err != nil {
		return errors.New("Problem marshaling config to JSON")
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, configJSON, os.ModeDevice)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) SetUser(username string) error {
	config, err := Read()
	if err != nil {
		return err
	}

	config.CurrentUserName = username

	err = write(config)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) Print() {
	fmt.Println("Database URL:", c.DbURL)
	fmt.Println("Current Username:", c.CurrentUserName)
}
