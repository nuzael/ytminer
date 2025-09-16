package utils

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	spaceRe   = regexp.MustCompile(`\s+`)
	punctRe   = regexp.MustCompile(`[\p{P}\p{S}]`)
	stopwords = map[string]struct{}{
		// English
		"the": {}, "and": {}, "for": {}, "with": {}, "you": {}, "your": {}, "from": {}, "that": {}, "this": {}, "are": {}, "was": {}, "were": {}, "have": {}, "has": {}, "how": {}, "what": {}, "why": {}, "when": {}, "to": {}, "in": {}, "on": {}, "of": {}, "a": {}, "an": {}, "is": {}, "it": {}, "by": {}, "or": {}, "as": {},
		// Portuguese (BR) (avoid duplicates with English: omit "a", "as")
		"de": {}, "da": {}, "do": {}, "das": {}, "dos": {}, "e": {}, "o": {}, "os": {}, "para": {}, "com": {}, "sem": {}, "em": {}, "no": {}, "na": {}, "nos": {}, "nas": {}, "um": {}, "uma": {}, "que": {}, "por": {}, "como": {}, "porque": {}, "qual": {},
	}
)

// RemoveAccents strips diacritics from a string
func RemoveAccents(s string) string {
	// NFD decompose, remove marks, NFC compose
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	res, _, _ := transform.String(t, s)
	return res
}

// NormalizeText lowercases, removes accents, punctuation, and collapses spaces
func NormalizeText(s string) string {
	s = strings.ToLower(s)
	s = RemoveAccents(s)
	s = punctRe.ReplaceAllString(s, " ")
	s = spaceRe.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// Tokenize splits text into words removing stopwords and short tokens
func Tokenize(s string) []string {
	n := NormalizeText(s)
	if n == "" { return nil }
	parts := strings.Fields(n)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if _, ok := stopwords[p]; ok { continue }
		if utf8.RuneCountInString(p) < 3 { continue }
		out = append(out, p)
	}
	return out
} 