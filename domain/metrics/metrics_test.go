package metrics

import (
	"testing"
	"time"
)

func TestVPD_Basic(t *testing.T) {
	now := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	pub := now.AddDate(0, 0, -10)
	v := VPD(1000, pub, now)
	if v <= 0 {
		t.Fatalf("expected positive vpd")
	}
}

func TestVPD_DivByZeroGuard(t *testing.T) {
	now := time.Now()
	pub := now
	v := VPD(1000, pub, now)
	if v != 1000 {
		t.Fatalf("expected vpd 1000 for age<1d, got %v", v)
	}
}

func TestLikeRatePerThousand(t *testing.T) {
	if LikeRatePerThousand(0, 10) != 0 {
		t.Fatal("guard expected")
	}
	if LikeRatePerThousand(1000, 100) != 100 {
		t.Fatal("100 likes over 1k views => 100/1k")
	}
}

func TestAgeDays(t *testing.T) {
	now := time.Now()
	pub := now.Add(-25 * time.Hour)
	if AgeDays(now, pub) < 1 {
		t.Fatal("expected at least 1 day")
	}
}

func TestFreshnessFromAges(t *testing.T) {
	ages := []float64{1, 5, 10}
	f := FreshnessFromAges(ages, 1)
	if f <= 0.5 {
		t.Fatal("newest should have higher freshness")
	}
}

func TestNormalizeSaturation(t *testing.T) {
	if NormalizeSaturation(5, 0) != 0 {
		t.Fatal("guard expected")
	}
	if NormalizeSaturation(5, 10) <= 0.4 {
		t.Fatal("5/10 should be 0.5")
	}
}

func TestVPDWindow(t *testing.T) {
	now := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	pub := now.AddDate(0, 0, -20)
	v7 := VPDWindow(7000, pub, now, 7)
	v30 := VPDWindow(7000, pub, now, 30)
	if v7 <= v30 {
		t.Fatalf("expected v7 > v30 (window smaller) got %v <= %v", v7, v30)
	}
}

func TestSlopeVPD(t *testing.T) {
	if SlopeVPD(100, 0) <= 0 {
		t.Fatalf("v30 guard should avoid division by zero")
	}
	if SlopeVPD(200, 100) <= 0 {
		t.Fatalf("accelerating should be positive")
	}
	if SlopeVPD(50, 100) >= 0 {
		t.Fatalf("decelerating should be negative")
	}
	if SlopeVPD(1e9, 1) > 5 {
		t.Fatalf("clamp upper bound to 5")
	}
	if SlopeVPD(-1e9, 1) < -5 {
		t.Fatalf("clamp lower bound to -5")
	}
}
