package metrics

import (
	"math"
	"time"
)

// VPD computes Views Per Day given total views and publication time.
func VPD(views int64, publishedAt time.Time, now time.Time) float64 {
	days := AgeDays(now, publishedAt)
	if days < 1 {
		days = 1
	}
	return float64(views) / float64(days)
}

// VPDWindow computes views per day over a capped window (e.g., 7, 30 days).
// It approximates recent velocity by dividing total views by min(ageDays, windowDays),
// guarding against division by zero.
func VPDWindow(views int64, publishedAt time.Time, now time.Time, windowDays int) float64 {
	if windowDays <= 0 {
		windowDays = 1
	}
	age := AgeDays(now, publishedAt)
	den := age
	if den > windowDays {
		den = windowDays
	}
	if den < 1 {
		den = 1
	}
	return float64(views) / float64(den)
}

// SlopeVPD estimates acceleration comparing short vs long window VPD.
// Returns a ratio-like measure: (v7 - v30) / max(v30, 1). Positive means accelerating.
// The value is clamped to [-5, 5] to limit outliers.
func SlopeVPD(v7 float64, v30 float64) float64 {
	base := v30
	if base < 1 {
		base = 1
	}
	val := (v7 - v30) / base
	if val > 5 {
		return 5
	}
	if val < -5 {
		return -5
	}
	return val
}

// LikeRatePerThousand computes likes per 1k views.
func LikeRatePerThousand(views int64, likes int64) float64 {
	if views <= 0 {
		return 0
	}
	return (float64(likes) / float64(views)) * 1000.0
}

// AgeDays returns the integer number of days between now and publishedAt.
func AgeDays(now time.Time, publishedAt time.Time) int {
	age := int(now.Sub(publishedAt).Hours() / 24)
	if age < 0 {
		return 0
	}
	return age
}

// FreshnessFromAges computes a [0,1] freshness score from an array of ages and the current age value.
// Newer items (smaller ages) receive higher freshness via inverted min-max scaling.
func FreshnessFromAges(allAges []float64, ageDays int) float64 {
	if len(allAges) == 0 {
		return 0
	}
	minV, maxV := allAges[0], allAges[0]
	for _, v := range allAges {
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
	}
	if maxV == minV {
		return 0
	}
	norm := (float64(ageDays) - minV) / (maxV - minV)
	return 1.0 - norm
}

// NormalizeSaturation returns a [0,1] saturation given frequency in sample and sample size.
func NormalizeSaturation(freq int, sampleSize int) float64 {
	if sampleSize <= 0 {
		return 0
	}
	v := float64(freq) / float64(sampleSize)
	return math.Max(0, math.Min(1, v))
}
