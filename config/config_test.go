package config

import (
	"os"
	"testing"
)

func TestLoadConfig_WeightsFromEnv(t *testing.T) {
	os.Setenv("YTMINER_OPP_W_VPD", "1.5")
	os.Setenv("YTMINER_OPP_W_LIKE", "0.3")
	os.Setenv("YTMINER_OPP_W_FRESH", "0.2")
	os.Setenv("YTMINER_OPP_W_SAT", "0.7")
	t.Cleanup(func() {
		os.Unsetenv("YTMINER_OPP_W_VPD")
		os.Unsetenv("YTMINER_OPP_W_LIKE")
		os.Unsetenv("YTMINER_OPP_W_FRESH")
		os.Unsetenv("YTMINER_OPP_W_SAT")
	})

	cfg := LoadConfig()
	if cfg.OppWeightVPD != 1.5 {
		t.Fatalf("VPD weight = %v", cfg.OppWeightVPD)
	}
	if cfg.OppWeightLike != 0.3 {
		t.Fatalf("Like weight = %v", cfg.OppWeightLike)
	}
	if cfg.OppWeightFresh != 0.2 {
		t.Fatalf("Fresh weight = %v", cfg.OppWeightFresh)
	}
	if cfg.OppWeightSatPen != 0.7 {
		t.Fatalf("Sat weight = %v", cfg.OppWeightSatPen)
	}
}
