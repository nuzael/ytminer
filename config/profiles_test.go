package config

import (
	"testing"
)

func TestGetProfile(t *testing.T) {
	// Test existing profile
	profile := GetProfile("exploration")
	if profile.Name != "Exploration" {
		t.Fatalf("expected Exploration, got %s", profile.Name)
	}
	if profile.Weights.Fresh != 0.35 {
		t.Fatalf("expected Fresh=0.35, got %f", profile.Weights.Fresh)
	}

	// Test non-existing profile (should return balanced)
	profile = GetProfile("nonexistent")
	if profile.Name != "Balanced" {
		t.Fatalf("expected Balanced for nonexistent, got %s", profile.Name)
	}
}

func TestApplyProfile(t *testing.T) {
	cfg := &AppConfig{}
	err := cfg.ApplyProfile("trending")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.OppWeightVPD != 0.50 {
		t.Fatalf("expected VPD=0.50, got %f", cfg.OppWeightVPD)
	}
	if cfg.OppWeightSlope != 0.30 {
		t.Fatalf("expected Slope=0.30, got %f", cfg.OppWeightSlope)
	}
}

func TestGetActiveProfileName(t *testing.T) {
	cfg := &AppConfig{
		OppWeightVPD:    0.30,
		OppWeightLike:   0.20,
		OppWeightFresh:  0.35,
		OppWeightSatPen: 0.10,
		OppWeightSlope:  0.25,
	}

	name := cfg.GetActiveProfileName()
	if name != "exploration" {
		t.Fatalf("expected exploration, got %s", name)
	}
}
