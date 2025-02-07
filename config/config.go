package config

import (
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

type Config struct {
	LinkUpUsername     string  `yaml:"link_up_username"`
	LinkUpPassword     string  `yaml:"link_up_password"`
	LinkUpRegion       int32   `yaml:"link_up_region"`
	LinkUpTimeInterval int64   `yaml:"link_up_time_interval"`
	LinkUpConnection   string  `yaml:"link_up_connection"`
	BotFatherToken     string  `yaml:"bot_father_token"`
	ChatIDsToNotify    []int64 `yaml:"chat_ids_to_notify"`
}

func Read(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %v", err)
	}

	// Validate required fields
	requiredEnvs := []string{
		config.LinkUpUsername,
		config.LinkUpPassword,
		config.BotFatherToken,
	}

	for _, env := range requiredEnvs {
		if env == "" {
			return nil, fmt.Errorf("required environment variable is not set")
		}
	}

	if config.LinkUpTimeInterval < 0 {
		return nil, fmt.Errorf("LINK_UP_TIME_INTERVAL expected to be a positive integer, but got '%d'", config.LinkUpTimeInterval)
	}

	return &config, nil
}
