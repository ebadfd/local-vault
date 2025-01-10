package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/pelletier/go-toml/v2"
)

const (
	APP_NAME    = "local-vault"
	DATA_FOLDER = "data"
)

type Config struct {
	StorePath string
	DbName    string
	Recipient *string
}

func (c *Config) GetFullDbPath() string {
	return filepath.Join(c.StorePath, c.DbName)
}

func (c *Config) GetFullDataPath() string {
	return filepath.Join(c.StorePath, DATA_FOLDER)
}

func LoadConfig() (*Config, error) {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", APP_NAME)
	configFile := filepath.Join(configDir, "config.toml")

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create config directory: %w", err)
		}
	}

	configDirData := filepath.Join(os.Getenv("HOME"), ".config", APP_NAME, DATA_FOLDER)

	if _, err := os.Stat(configDirData); os.IsNotExist(err) {
		err := os.MkdirAll(configDirData, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create config directory for data: %w", err)
		}
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		defaultConfig := &Config{
			StorePath: configDir,
			DbName:    "core.db",
		}
		err := saveConfig(configFile, defaultConfig)
		if err != nil {
			return nil, err
		}
		return defaultConfig, nil
	}

	return loadConfig(configFile)
}

func loadConfig(configFile string) (*Config, error) {
	var config Config
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if config.Recipient == nil {
		email, err := gpguser()

		if err != nil {
			return nil, fmt.Errorf("gpg user generation failed")
		}

		config.Recipient = &email

		err = saveConfig(configFile, &config)
		if err != nil {
			return nil, err
		}
	}

	return &config, nil
}

func gpguser() (string, error) {
	validate := func(input string) error {
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		re := regexp.MustCompile(emailRegex)

		if !re.MatchString(input) {
			return errors.New("invalid email address")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "GPG recipient",
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return result, err
}

func saveConfig(configFile string, config *Config) error {
	data, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
