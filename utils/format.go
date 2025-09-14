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
	if rate > 5 {
		return "excellent"
	} else if rate > 2 {
		return "good"
	}
	return "low"
}
