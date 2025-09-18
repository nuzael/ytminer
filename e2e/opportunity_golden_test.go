package e2e

import (
	"bytes"
	"os"
	"testing"
	"time"

	"ytminer/analysis"
	"ytminer/config"
	"ytminer/ui"
	"ytminer/youtube"
)

func TestOpportunityScore_Golden_VPDOnly(t *testing.T) {
	// Force weights via env: VPD=1, others=0
	os.Setenv("YTMINER_OPP_W_VPD", "1")
	os.Setenv("YTMINER_OPP_W_LIKE", "0")
	os.Setenv("YTMINER_OPP_W_FRESH", "0")
	os.Setenv("YTMINER_OPP_W_SAT", "0")
	t.Cleanup(func() {
		os.Unsetenv("YTMINER_OPP_W_VPD")
		os.Unsetenv("YTMINER_OPP_W_LIKE")
		os.Unsetenv("YTMINER_OPP_W_FRESH")
		os.Unsetenv("YTMINER_OPP_W_SAT")
	})

	now := time.Now()
	pub := now.AddDate(0, 0, -10)
	videos := []youtube.Video{
		{ID: "v1", Title: "Alpha guide", Channel: "C1", ChannelID: "c1", PublishedAt: pub, Views: 100, Likes: 1, URL: "http://x/1", VPD: 5000},
		{ID: "v2", Title: "Beta review", Channel: "C2", ChannelID: "c2", PublishedAt: pub, Views: 100, Likes: 1, URL: "http://x/2", VPD: 1000},
		{ID: "v3", Title: "Gamma tips", Channel: "C3", ChannelID: "c3", PublishedAt: pub, Views: 100, Likes: 1, URL: "http://x/3", VPD: 3000},
	}
	cfg := config.LoadConfig()
	a := analysis.NewAnalyzer(videos, cfg)
	items := a.AnalyzeOpportunityScore()
	if len(items) == 0 {
		t.Fatalf("expected items")
	}
	if items[0].Title != "Alpha guide" {
		t.Fatalf("expected top by VPD to be 'Alpha guide', got %q", items[0].Title)
	}

	// Capture stdout of UI to ensure stable header and rank formatting
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	ui.DisplayOpportunityScore(items)
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	out := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("Opportunity Score (Top Candidates)")) {
		t.Fatalf("expected header in output; got: %s", out)
	}
	if !bytes.Contains(buf.Bytes(), []byte("#1")) {
		t.Fatalf("expected rank #1 in output")
	}
}
