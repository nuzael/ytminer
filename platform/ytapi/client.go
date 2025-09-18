package ytapi

import (
	"context"
	"ytminer/platform/transcripts"
	"ytminer/youtube"
)

// Client is the port interface for YouTube data access used by the app layer.
type Client interface {
	SearchVideos(opts youtube.SearchOptions) ([]youtube.Video, error)
	GetTranscript(ctx context.Context, videoID string) (*transcripts.Transcript, error)
}

// adapter wraps the existing concrete youtube.Client.
type adapter struct {
	c       *youtube.Client
	fetcher transcripts.Fetcher
}

// New creates a new YouTube API client adapter using the existing implementation.
func New() (Client, error) {
	c, err := youtube.NewClient()
	if err != nil {
		return nil, err
	}
	return &adapter{c: c, fetcher: transcripts.DefaultFetcher{}}, nil
}

func (a *adapter) SearchVideos(opts youtube.SearchOptions) ([]youtube.Video, error) {
	return a.c.SearchVideos(opts)
}

func (a *adapter) GetTranscript(ctx context.Context, videoID string) (*transcripts.Transcript, error) {
	return a.fetcher.Get(ctx, videoID)
}
