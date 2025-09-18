package config

import (
	"os"
	"strconv"
	"strings"
)

// AppConfig represents the application configuration
type AppConfig struct {
	DefaultRegion         string
	DefaultDuration       string
	DefaultTimeRange      string
	DefaultOrder          string
	APIKey                string
	RisingStarMultiplier  float64
	LongTailMinEngagement float64
	LongTailMaxFreq       int

	// Opportunity Score Weights
	OppWeightVPD    float64
	OppWeightLike   float64
	OppWeightFresh  float64
	OppWeightSatPen float64
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
		OppWeightVPD:          0.45,
		OppWeightLike:         0.25,
		OppWeightFresh:        0.20,
		OppWeightSatPen:       0.30,
	}

	// Load configuration from environment
	if region := strings.TrimSpace(os.Getenv("YTMINER_DEFAULT_REGION")); region != "" {
		config.DefaultRegion = region
	}

	if duration := strings.TrimSpace(os.Getenv("YTMINER_DEFAULT_DURATION")); duration != "" {
		config.DefaultDuration = duration
	}

	if tr := strings.TrimSpace(os.Getenv("YTMINER_DEFAULT_TIME_RANGE")); tr != "" {
		config.DefaultTimeRange = tr
	}

	if ord := strings.TrimSpace(os.Getenv("YTMINER_DEFAULT_ORDER")); ord != "" {
		config.DefaultOrder = ord
	}

	if v := strings.TrimSpace(os.Getenv("YTMINER_RISING_STAR_MULTIPLIER")); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
			config.RisingStarMultiplier = f
		}
	}

	if v := strings.TrimSpace(os.Getenv("YTMINER_LONG_TAIL_MIN_ENGAGEMENT")); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 0 {
			config.LongTailMinEngagement = f
		}
	}

	if v := strings.TrimSpace(os.Getenv("YTMINER_LONG_TAIL_MAX_FREQ")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 1 {
			config.LongTailMaxFreq = n
		}
	}

	// Opportunity Score Weights from env (optional overrides)
	if v := strings.TrimSpace(os.Getenv("YTMINER_OPP_W_VPD")); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 0 {
			config.OppWeightVPD = f
		}
	}
	if v := strings.TrimSpace(os.Getenv("YTMINER_OPP_W_LIKE")); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 0 {
			config.OppWeightLike = f
		}
	}
	if v := strings.TrimSpace(os.Getenv("YTMINER_OPP_W_FRESH")); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 0 {
			config.OppWeightFresh = f
		}
	}
	if v := strings.TrimSpace(os.Getenv("YTMINER_OPP_W_SAT")); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 0 {
			config.OppWeightSatPen = f
		}
	}

	config.APIKey = strings.TrimSpace(os.Getenv("YOUTUBE_API_KEY"))

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

	// Persist Opportunity Score Weights
	envContent += "YTMINER_OPP_W_VPD=" + strconv.FormatFloat(c.OppWeightVPD, 'f', -1, 64) + "\n"
	envContent += "YTMINER_OPP_W_LIKE=" + strconv.FormatFloat(c.OppWeightLike, 'f', -1, 64) + "\n"
	envContent += "YTMINER_OPP_W_FRESH=" + strconv.FormatFloat(c.OppWeightFresh, 'f', -1, 64) + "\n"
	envContent += "YTMINER_OPP_W_SAT=" + strconv.FormatFloat(c.OppWeightSatPen, 'f', -1, 64) + "\n"

	return os.WriteFile(".env", []byte(envContent), 0644)
}
