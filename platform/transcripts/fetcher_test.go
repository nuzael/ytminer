package transcripts

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultFetcher_Get_EmptyInput(t *testing.T) {
	f := DefaultFetcher{}
	ctx := context.Background()

	// This will fail in real usage but should not panic
	_, err := f.Get(ctx, "")
	if err == nil {
		t.Fatal("expected error for empty video ID")
	}
}

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"en,pt", []string{"en", "pt"}},
		{"en, pt, es", []string{"en", "pt", "es"}},
		{"", []string{}},
		{"en", []string{"en"}},
		{"en,,pt", []string{"en", "pt"}},
	}

	for _, tt := range tests {
		result := splitAndTrim(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("splitAndTrim(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestHtmlEntityDecode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"&amp;", "&"},
		{"&lt;", "<"},
		{"&gt;", ">"},
		{"&#39;", "'"},
		{"&quot;", "\""},
		{"hello &amp; world", "hello & world"},
		{"no entities", "no entities"},
	}

	for _, tt := range tests {
		result := htmlEntityDecode(tt.input)
		if result != tt.expected {
			t.Errorf("htmlEntityDecode(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestCachePath_RespectsEnvDir(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("YTMINER_CACHE_DIR", tmp)
	t.Cleanup(func() { os.Unsetenv("YTMINER_CACHE_DIR") })

	p := cachePath("abc", "en")
	if filepath.Dir(p) != filepath.Join(tmp, "transcripts") {
		t.Fatalf("expected cache dir %s, got %s", filepath.Join(tmp, "transcripts"), filepath.Dir(p))
	}
}

func TestReadWriteCache_WithLangsFallback(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("YTMINER_CACHE_DIR", tmp)
	os.Setenv("YTMINER_TRANSCRIPT_LANGS", "es,pt,en")
	t.Cleanup(func() {
		os.Unsetenv("YTMINER_CACHE_DIR")
		os.Unsetenv("YTMINER_TRANSCRIPT_LANGS")
	})

	// write english cache only
	w := &Transcript{VideoID: "vid1", Language: "en", Text: "hello world"}
	writeToCache(w)

	// read should find en when es and pt miss
	tr, ok := readFromCache("vid1")
	if !ok || tr == nil {
		t.Fatalf("expected to read from cache")
	}
	if tr.Text != "hello world" || tr.Language != "en" || tr.Source != "cache" {
		t.Fatalf("unexpected transcript: %+v", tr)
	}
}
