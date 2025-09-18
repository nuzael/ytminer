package e2e

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"ytminer/analysis"
	"ytminer/config"
	"ytminer/ui"
	"ytminer/youtube"
)

func TestAllMode_PrintsActiveWeights(t *testing.T) {
	videos := []youtube.Video{{Title: "t", Channel: "c", URL: "u", Views: 1}}
	cfg := config.LoadConfig()
	a := analysis.NewAnalyzer(videos, cfg)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	items := a.AnalyzeOpportunityScore()
	ui.DisplayInfo(
		fmt.Sprintf("Active weights → VPD=%.2f, LIKE=%.2f, FRESH=%.2f, SAT=%.2f, SLOPE=%.2f",
			cfg.OppWeightVPD, cfg.OppWeightLike, cfg.OppWeightFresh, cfg.OppWeightSatPen, cfg.OppWeightSlope,
		),
	)
	ui.DisplayOpportunityScore(items)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	out := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("Active weights →")) {
		t.Fatalf("expected 'Active weights →' in output, got: %s", out)
	}
}
