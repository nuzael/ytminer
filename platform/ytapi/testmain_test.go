package ytapi

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
)

// TestMain loads .env before running tests in this package so tests that depend on
// YOUTUBE_API_KEY can run without manual environment export.
func TestMain(m *testing.M) {
	// Try multiple locations: current dir and parents
	candidates := []string{
		".env",
		filepath.Join("..", ".env"),
		filepath.Join("..", "..", ".env"),
		filepath.Join("..", "..", "..", ".env"),
	}
	loaded := false
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			if err := godotenv.Overload(p); err == nil {
				loaded = true
				break
			}
		}
	}
	if !loaded {
		_ = godotenv.Load() // last attempt in CWD without erroring
	}

	if os.Getenv("YOUTUBE_API_KEY") == "" {
		log.Println("ytapi tests: YOUTUBE_API_KEY not set; tests that require API key may be skipped")
	}
	os.Exit(m.Run())
}
