package utils

import (
	"fmt"
)

// FormatNumber formats numbers for friendly display (K, M)
func FormatNumber(num interface{}) string {
	var f float64
	switch v := num.(type) {
	case int64:
		f = float64(v)
	case float64:
		f = v
	default:
		return fmt.Sprintf("%v", num)
	}

	if f >= 1000000 {
		return fmt.Sprintf("%.1fM", f/1000000)
	} else if f >= 1000 {
		return fmt.Sprintf("%.1fK", f/1000)
	}
	return fmt.Sprintf("%.0f", f)
}

// FormatEngagement formats engagement rate to descriptive text
func FormatEngagement(rate float64) string {
	if rate > 10 {
		return "excellent"
	} else if rate > 5 {
		return "very good"
	} else if rate > 2 {
		return "good"
	} else if rate > 1 {
		return "average"
	}
	return "low"
}

// FormatEngagementRate formats engagement percentage with proper decimals
func FormatEngagementRate(rate float64) string {
	if rate >= 10 {
		return fmt.Sprintf("%.1f%%", rate)
	} else if rate >= 1 {
		return fmt.Sprintf("%.2f%%", rate)
	} else {
		return fmt.Sprintf("%.3f%%", rate)
	}
}

// FormatVPD formats VPD (Views Per Day) for better readability
func FormatVPD(vpd float64) string {
	if vpd >= 1000000 {
		return fmt.Sprintf("%.1fM", vpd/1000000)
	} else if vpd >= 1000 {
		return fmt.Sprintf("%.1fK", vpd/1000)
	} else if vpd >= 1 {
		return fmt.Sprintf("%.0f", vpd)
	} else {
		return fmt.Sprintf("%.2f", vpd)
	}
}
