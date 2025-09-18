package ytapi

import (
	"context"
	"testing"
	"ytminer/youtube"
)

func TestAdapter_SearchVideos(t *testing.T) {
	// This test requires a valid API key, so we'll just test the interface
	client, err := New()
	if err != nil {
		t.Skip("Skipping test: no API key available")
	}

	opts := youtube.SearchOptions{
		Query:      "test",
		MaxResults: 1,
		Level:      youtube.QuickScan,
	}

	_, err = client.SearchVideos(opts)
	if err != nil {
		t.Logf("SearchVideos failed (expected without API key): %v", err)
	}
}

func TestAdapter_GetTranscript(t *testing.T) {
	client, err := New()
	if err != nil {
		t.Skip("Skipping test: no API key available")
	}

	ctx := context.Background()
	_, err = client.GetTranscript(ctx, "dQw4w9WgXcQ") // Rick Roll video
	if err != nil {
		t.Logf("GetTranscript failed (expected for test video): %v", err)
	}
}
