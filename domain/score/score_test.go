package score

import (
	"testing"
	"time"

	"ytminer/youtube"
)

func sampleVideos() []youtube.Video {
	now := time.Now().Add(-24 * time.Hour) // ensure age >= 1 day for all
	return []youtube.Video{
		{ID: "a", Title: "alpha", Views: 10000, Likes: 200, PublishedAt: now.Add(-5 * 24 * time.Hour), VPD: 2000},
		{ID: "b", Title: "beta", Views: 5000, Likes: 150, PublishedAt: now.Add(-10 * 24 * time.Hour), VPD: 500},
		{ID: "c", Title: "alpha guide", Views: 8000, Likes: 80, PublishedAt: now.Add(-2 * 24 * time.Hour), VPD: 4000},
	}
}

func TestCompute_MonotonicityVPD(t *testing.T) {
	videos := sampleVideos()
	now := time.Now()
	w := Weights{VPD: 1.0, Like: 0.0, Fresh: 0.0, Sat: 0.0}
	items := Compute(videos, w, now)
	if len(items) != len(videos) {
		t.Fatalf("expected %d items, got %d", len(videos), len(items))
	}
	// Expect item with highest VPD (video c) to rank first
	if items[0].Title != "alpha guide" {
		t.Fatalf("expected top by VPD to be 'alpha guide', got %q", items[0].Title)
	}
}

func TestCompute_MonotonicityLikeRate(t *testing.T) {
	videos := sampleVideos()
	now := time.Now()
	w := Weights{VPD: 0.0, Like: 1.0, Fresh: 0.0, Sat: 0.0}
	items := Compute(videos, w, now)
	if len(items) != len(videos) {
		t.Fatalf("expected %d items, got %d", len(videos), len(items))
	}
	// Compute like rates manually
	// a: 200/10000=0.02 -> 20/1k, b: 150/5000=0.03 -> 30/1k, c: 80/8000=0.01 -> 10/1k
	if items[0].Title != "beta" {
		t.Fatalf("expected top by like rate to be 'beta', got %q", items[0].Title)
	}
}

func TestCompute_FreshnessPreference(t *testing.T) {
	videos := sampleVideos()
	now := time.Now()
	w := Weights{VPD: 0.0, Like: 0.0, Fresh: 1.0, Sat: 0.0}
	items := Compute(videos, w, now)
	// newest is c (2 days) vs a (5) vs b (10)
	if items[0].Title != "alpha guide" {
		t.Fatalf("expected freshest to be first, got %q", items[0].Title)
	}
}

func TestCompute_SaturationPenalty(t *testing.T) {
	videos := sampleVideos()
	now := time.Now()
	w := Weights{VPD: 0.0, Like: 0.0, Fresh: 0.0, Sat: 1.0}
	items := Compute(videos, w, now)
	// primary token frequency: "alpha" appears in a and c; "beta" appears once -> lower penalty ranks higher
	if items[0].Title != "beta" {
		t.Fatalf("expected least saturated token to rank first, got %q", items[0].Title)
	}
}

func TestCompute_EmptyInput(t *testing.T) {
	items := Compute(nil, Weights{VPD: 1}, time.Now())
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}
