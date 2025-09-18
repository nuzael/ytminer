package utils

import (
	"strings"

	"github.com/rivo/uniseg"
)

// emojiRanges covers main emoji ranges (pictographic, symbols, dingbats, etc.)
var emojiRanges = []struct{ start, end rune }{
	{0x1F600, 0x1F64F}, // Emoticons
	{0x1F300, 0x1F5FF}, // Misc Symbols and Pictographs
	{0x1F680, 0x1F6FF}, // Transport & Map
	{0x1F900, 0x1F9FF}, // Supplemental Symbols & Pictographs
	{0x1FA70, 0x1FAFF}, // Symbols & Pictographs Extended-A
	{0x2700, 0x27BF},   // Dingbats
	{0x2600, 0x26FF},   // Misc Symbols
	{0x1F780, 0x1F7FF}, // Geometric Shapes Extended (subset)
}

// isRegionalIndicator tests if a rune is in the range of regional indicators (for flags)
func isRegionalIndicator(r rune) bool {
	return r >= 0x1F1E6 && r <= 0x1F1FF
}

func isEmojiRune(r rune) bool {
	for _, rg := range emojiRanges {
		if r >= rg.start && r <= rg.end {
			return true
		}
	}
	return false
}

// ExtractEmojis returns a list of emojis (as complete clusters), including flags and ZWJ sequences.
func ExtractEmojis(s string) []string {
	var result []string
	g := uniseg.NewGraphemes(s)
	prevRI := rune(0)
	for g.Next() {
		cluster := g.Str()
		// If cluster has ZWJ, consider emoji (families, combined)
		if strings.Contains(cluster, "\u200d") { // ZWJ
			result = append(result, cluster)
			prevRI = 0
			continue
		}
		// Flags: pair of regional indicators
		runes := []rune(cluster)
		if len(runes) == 1 && isRegionalIndicator(runes[0]) {
			if prevRI != 0 {
				result = append(result, string([]rune{prevRI, runes[0]}))
				prevRI = 0
				continue
			}
			prevRI = runes[0]
			continue
		}
		prevRI = 0

		// Variation Selector-16 (0xFE0F) indicates emoji presentation
		if strings.Contains(cluster, "\uFE0F") {
			result = append(result, cluster)
			continue
		}
		// Simple pictographic
		for _, r := range runes {
			if isEmojiRune(r) {
				result = append(result, cluster)
				break
			}
		}
	}
	return result
}
