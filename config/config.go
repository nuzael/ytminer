package config

import (
	"os"
	"strconv"
)

// AppConfig represents the application configuration
type AppConfig struct {
	DefaultRegion          string
	DefaultDuration        string
	DefaultTimeRange       string
	DefaultOrder           string
	APIKey                 string
	RisingStarMultiplier   float64
	LongTailMinEngagement  float64
	LongTailMaxFreq        int
}

// LoadConfig loads configuration from environment
func LoadConfig() *AppConfig {
	config := &AppConfig{
		DefaultRegion:         "any",
		DefaultDuration:       "any",
		DefaultTimeRange:      "any",
		DefaultOrder:          "relevance",
		RisingStarMultiplier:  1.5,
		LongTailMinEngagement: 5.0,
		LongTailMaxFreq:       2,
	}

	// Load configuration from environment
	if region := os.Getenv("YTMINER_DEFAULT_REGION"); region != "" {
		config.DefaultRegion = region
	}

	if duration := os.Getenv("YTMINER_DEFAULT_DURATION"); duration != "" {
		config.DefaultDuration = duration
	}

	if tr := os.Getenv("YTMINER_DEFAULT_TIME_RANGE"); tr != "" {
		config.DefaultTimeRange = tr
	}

	if ord := os.Getenv("YTMINER_DEFAULT_ORDER"); ord != "" {
		config.DefaultOrder = ord
	}

	if v := os.Getenv("YTMINER_RISING_STAR_MULTIPLIER"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
			config.RisingStarMultiplier = f
		}
	}

	if v := os.Getenv("YTMINER_LONG_TAIL_MIN_ENGAGEMENT"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 0 {
			config.LongTailMinEngagement = f
		}
	}

	if v := os.Getenv("YTMINER_LONG_TAIL_MAX_FREQ"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 1 {
			config.LongTailMaxFreq = n
		}
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
	
	envContent += "YTMINER_DEFAULT_REGION=" + c.DefaultRegion + "\n"
	envContent += "YTMINER_DEFAULT_DURATION=" + c.DefaultDuration + "\n"
	envContent += "YTMINER_DEFAULT_TIME_RANGE=" + c.DefaultTimeRange + "\n"
	envContent += "YTMINER_DEFAULT_ORDER=" + c.DefaultOrder + "\n"
	envContent += "YTMINER_RISING_STAR_MULTIPLIER=" + strconv.FormatFloat(c.RisingStarMultiplier, 'f', -1, 64) + "\n"
	envContent += "YTMINER_LONG_TAIL_MIN_ENGAGEMENT=" + strconv.FormatFloat(c.LongTailMinEngagement, 'f', -1, 64) + "\n"
	envContent += "YTMINER_LONG_TAIL_MAX_FREQ=" + strconv.Itoa(c.LongTailMaxFreq) + "\n"

	return os.WriteFile(".env", []byte(envContent), 0644)
}
