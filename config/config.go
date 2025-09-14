package config

import (
	"os"
	"strconv"
)

// AppConfig represents the application configuration
type AppConfig struct {
	MaxResults      int
	DefaultRegion   string
	DefaultDuration string
	APIKey          string
}

// LoadConfig loads configuration from environment
func LoadConfig() *AppConfig {
	config := &AppConfig{
		MaxResults:      25,
		DefaultRegion:   "any",
		DefaultDuration: "any",
	}

	// Load configuration from environment
	if maxResults := os.Getenv("YTMINER_MAX_RESULTS"); maxResults != "" {
		if parsed, err := strconv.Atoi(maxResults); err == nil {
			config.MaxResults = parsed
		}
	}

	if region := os.Getenv("YTMINER_DEFAULT_REGION"); region != "" {
		config.DefaultRegion = region
	}

	if duration := os.Getenv("YTMINER_DEFAULT_DURATION"); duration != "" {
		config.DefaultDuration = duration
	}

	config.APIKey = os.Getenv("YOUTUBE_API_KEY")

	return config
}

// SaveConfig saves configuration to .env file
func (c *AppConfig) SaveConfig() error {
	envContent := ""
	
	if c.APIKey != "" {
		envContent += "YOUTUBE_API_KEY=" + c.APIKey + "\n"
	}
	
	if c.MaxResults != 25 {
		envContent += "YTMINER_MAX_RESULTS=" + strconv.Itoa(c.MaxResults) + "\n"
	}
	
	if c.DefaultRegion != "any" {
		envContent += "YTMINER_DEFAULT_REGION=" + c.DefaultRegion + "\n"
	}
	
	if c.DefaultDuration != "any" {
		envContent += "YTMINER_DEFAULT_DURATION=" + c.DefaultDuration + "\n"
	}

	return os.WriteFile(".env", []byte(envContent), 0644)
}
