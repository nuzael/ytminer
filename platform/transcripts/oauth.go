package transcripts

import (
	"context"
	"fmt"
)

// OAuth2Fetcher implements Fetcher using YouTube Data API v3 with OAuth 2.0
// This is prepared for future implementation when OAuth 2.0 is needed
type OAuth2Fetcher struct {
	// Future: OAuth 2.0 client, access token, refresh token, etc.
}

// Get fetches transcript using YouTube Data API v3 captions.download
// TODO: Implement OAuth 2.0 flow and captions.download API call
func (f *OAuth2Fetcher) Get(ctx context.Context, videoID string) (*Transcript, error) {
	return nil, fmt.Errorf("OAuth2Fetcher not implemented yet - requires OAuth 2.0 setup")
}

// NewOAuth2Fetcher creates a new OAuth2Fetcher instance
// TODO: Implement OAuth 2.0 initialization
func NewOAuth2Fetcher() (*OAuth2Fetcher, error) {
	return nil, fmt.Errorf("OAuth2Fetcher not implemented yet - requires OAuth 2.0 setup")
}
