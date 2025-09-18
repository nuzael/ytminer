package analysis

import (
	"testing"
	"time"
	"ytminer/config"
	"ytminer/youtube"
)

func sampleVideos() []youtube.Video {
	now := time.Now()
	return []youtube.Video{
		{ID: "a", Title: "Alpha tutorial", Channel: "C1", ChannelID: "c1", Views: 10000, Likes: 200, Comments: 20, PublishedAt: now.AddDate(0, 0, -5), URL: "http://x/a", VPD: 2000},
		{ID: "b", Title: "Beta review", Channel: "C2", ChannelID: "c2", Views: 5000, Likes: 150, Comments: 50, PublishedAt: now.AddDate(0, 0, -10), URL: "http://x/b", VPD: 500},
		{ID: "c", Title: "Alpha 2024 guide", Channel: "C1", ChannelID: "c1", Views: 8000, Likes: 80, Comments: 10, PublishedAt: now.AddDate(0, 0, -2), URL: "http://x/c", VPD: 4000},
	}
}

func newAnalyzerForTest() *Analyzer {
	cfg := config.LoadConfig()
	return NewAnalyzer(sampleVideos(), cfg)
}

func TestAnalyzeGrowthPatterns(t *testing.T) {
	a := newAnalyzerForTest()
	g := a.AnalyzeGrowthPatterns()
	if g.TotalVideos != 3 {
		t.Fatalf("expected 3 videos, got %d", g.TotalVideos)
	}
	if g.NicheVelocityScore <= 0 {
		t.Fatalf("expected positive niche velocity score")
	}
	if len(g.TopPerformers) == 0 {
		t.Fatalf("expected at least one top performer")
	}
}

func TestAnalyzeTitles(t *testing.T) {
	a := newAnalyzerForTest()
	titles := a.AnalyzeTitles()
	if len(titles.CommonWords) == 0 {
		t.Fatalf("expected common words")
	}
}

func TestAnalyzeCompetitors(t *testing.T) {
	a := newAnalyzerForTest()
	comp := a.AnalyzeCompetitors()
	if len(comp.TopChannels) == 0 {
		t.Fatalf("expected top channels")
	}
}

func TestAnalyzeTemporal(t *testing.T) {
	a := newAnalyzerForTest()
	temp := a.AnalyzeTemporal()
	// May be empty if timestamps don't cover enough buckets, but should not panic
	_ = temp
}

func TestAnalyzeKeywords(t *testing.T) {
	a := newAnalyzerForTest()
	kw := a.AnalyzeKeywords()
	// Should produce some stats from titles
	if len(kw.TrendingKeywords) == 0 && len(kw.CoreKeywords) == 0 {
		t.Fatalf("expected some keyword stats")
	}
}

func TestAnalyzeOpportunityScore(t *testing.T) {
	a := newAnalyzerForTest()
	items := a.AnalyzeOpportunityScore()
	if len(items) == 0 {
		t.Fatalf("expected ranked items")
	}
	if items[0].Score < items[len(items)-1].Score {
		t.Fatalf("expected descending scores")
	}
}
