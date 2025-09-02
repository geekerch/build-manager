package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Config represents the application configuration
type Config struct {
	Server     ServerConfig           `json:"server"`
	GitConfigs map[string]GitConfig   `json:"git_configs"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port         string `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// GitConfig represents Git repository configuration
type GitConfig struct {
	URL         string `json:"url"`
	Token       string `json:"token"`       // PAT token for authentication
	Description string `json:"description"` // Human-readable description
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         "8080",
			ReadTimeout:  15,
			WriteTimeout: 15,
		},
		GitConfigs: map[string]GitConfig{
			"main": {
				URL:         "https://github.com/your-org/build-scripts.git",
				Token:       "",
				Description: "主要構建配置倉庫",
			},
		},
	}
}

// LoadConfig loads configuration from file or returns default
func LoadConfig(filename string) *Config {
	config := DefaultConfig()

	// Try to load from file
	if data, err := ioutil.ReadFile(filename); err == nil {
		if err := json.Unmarshal(data, config); err != nil {
			log.Printf("Error parsing config file %s: %v, using defaults", filename, err)
		} else {
			log.Printf("Loaded configuration from %s", filename)
		}
	} else {
		log.Printf("Config file %s not found, using defaults", filename)
	}

	// Override with environment variables
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = port
	}

	return config
}

// SaveConfig saves configuration to file
func (c *Config) SaveConfig(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}