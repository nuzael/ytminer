package youtube

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Client struct {
	service *youtube.Service
}

type Video struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Channel     string    `json:"channel"`
	PublishedAt time.Time `json:"published_at"`
	Views       int64     `json:"views"`
	Likes       int64     `json:"likes"`
	Comments    int64     `json:"comments"`
	Duration    string    `json:"duration"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
}

type AnalysisLevel int

const (
	QuickScan AnalysisLevel = iota
	Balanced
	DeepDive
)

type SearchOptions struct {
	Query         string
	MaxResults    int64
	Region        string
	Duration      string
	Order         string
	MinViews      int64
	MinLikes      int64
	Level         AnalysisLevel
	PublishedAfter string
	PublishedBefore string
}

func NewClient() (*Client, error) {
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("YOUTUBE_API_KEY environment variable not set")
	}

	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create YouTube service: %v", err)
	}

	return &Client{service: service}, nil
}

func (c *Client) SearchVideos(opts SearchOptions) ([]Video, error) {
	// Use level-based search if level is specified
	if opts.Level != QuickScan {
		return c.SearchVideosWithLevel(opts)
	}

	// Original single search for QuickScan
	return c.searchSingle(opts)
}

func (c *Client) SearchVideosWithLevel(opts SearchOptions) ([]Video, error) {
	var allVideos []Video
	var seenIDs = make(map[string]bool)

	switch opts.Level {
	case QuickScan:
		return c.searchSingle(opts)
	case Balanced:
		allVideos = c.searchBalanced(opts, seenIDs)
	case DeepDive:
		allVideos = c.searchDeepDive(opts, seenIDs)
	}

	return allVideos, nil
}

func (c *Client) searchSingle(opts SearchOptions) ([]Video, error) {
	call := c.service.Search.List([]string{"id", "snippet"}).
		Q(opts.Query).
		MaxResults(opts.MaxResults).
		Type("video")

	if opts.Region != "" && opts.Region != "any" {
		call = call.RegionCode(opts.Region)
	}

	if opts.Duration != "" && opts.Duration != "any" {
		call = call.VideoDuration(opts.Duration)
	}

	if opts.Order != "" {
		call = call.Order(opts.Order)
	}

	if opts.PublishedAfter != "" {
		call = call.PublishedAfter(opts.PublishedAfter)
	}

	if opts.PublishedBefore != "" {
		call = call.PublishedBefore(opts.PublishedBefore)
	}

	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search videos: %v", err)
	}

	videos, err := c.processSearchResults(response, opts)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (c *Client) searchBalanced(opts SearchOptions, seenIDs map[string]bool) []Video {
	var allVideos []Video
	
	// Search 1: Original query
	videos1 := c.searchWithParams(opts.Query, "any", "any", "relevance", 50, opts)
	allVideos = append(allVideos, c.filterUniqueVideos(videos1, seenIDs)...)

	// Search 2: Different region
	if opts.Region == "any" || opts.Region == "" {
		videos2 := c.searchWithParams(opts.Query, "BR", "any", "relevance", 50, opts)
		allVideos = append(allVideos, c.filterUniqueVideos(videos2, seenIDs)...)
	}

	// Search 3: Different duration
	if opts.Duration == "any" || opts.Duration == "" {
		videos3 := c.searchWithParams(opts.Query, "any", "medium", "relevance", 50, opts)
		allVideos = append(allVideos, c.filterUniqueVideos(videos3, seenIDs)...)
	}

	// Search 4: Different order
	videos4 := c.searchWithParams(opts.Query, "any", "any", "date", 50, opts)
	allVideos = append(allVideos, c.filterUniqueVideos(videos4, seenIDs)...)

	return allVideos
}

func (c *Client) searchDeepDive(opts SearchOptions, seenIDs map[string]bool) []Video {
	var allVideos []Video
	
	// Search 1-3: Original query with different orders
	allVideos = append(allVideos, c.filterUniqueVideos(c.searchWithParams(opts.Query, "any", "any", "relevance", 50, opts), seenIDs)...)
	allVideos = append(allVideos, c.filterUniqueVideos(c.searchWithParams(opts.Query, "any", "any", "viewCount", 50, opts), seenIDs)...)
	allVideos = append(allVideos, c.filterUniqueVideos(c.searchWithParams(opts.Query, "any", "any", "date", 50, opts), seenIDs)...)

	// Search 4-6: Different regions
	allVideos = append(allVideos, c.filterUniqueVideos(c.searchWithParams(opts.Query, "BR", "any", "relevance", 50, opts), seenIDs)...)
	allVideos = append(allVideos, c.filterUniqueVideos(c.searchWithParams(opts.Query, "US", "any", "relevance", 50, opts), seenIDs)...)
	allVideos = append(allVideos, c.filterUniqueVideos(c.searchWithParams(opts.Query, "GB", "any", "relevance", 50, opts), seenIDs)...)

	// Search 7-9: Different durations
	allVideos = append(allVideos, c.filterUniqueVideos(c.searchWithParams(opts.Query, "any", "short", "relevance", 50, opts), seenIDs)...)
	allVideos = append(allVideos, c.filterUniqueVideos(c.searchWithParams(opts.Query, "any", "medium", "relevance", 50, opts), seenIDs)...)
	allVideos = append(allVideos, c.filterUniqueVideos(c.searchWithParams(opts.Query, "any", "long", "relevance", 50, opts), seenIDs)...)

	// Search 10-12: Related queries with different approaches
	relatedQuery1 := opts.Query + " tutorial"
	relatedQuery2 := opts.Query + " guide"
	relatedQuery3 := opts.Query + " tips"
	allVideos = append(allVideos, c.filterUniqueVideos(c.searchWithParams(relatedQuery1, "any", "any", "relevance", 50, opts), seenIDs)...)
	allVideos = append(allVideos, c.filterUniqueVideos(c.searchWithParams(relatedQuery2, "any", "any", "viewCount", 50, opts), seenIDs)...)
	allVideos = append(allVideos, c.filterUniqueVideos(c.searchWithParams(relatedQuery3, "any", "any", "date", 50, opts), seenIDs)...)

	return allVideos
}

func (c *Client) searchWithParams(query, region, duration, order string, maxResults int64, originalOpts SearchOptions) []Video {
	call := c.service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(maxResults).
		Type("video")

	if region != "" && region != "any" {
		call = call.RegionCode(region)
	}

	if duration != "" && duration != "any" {
		call = call.VideoDuration(duration)
	}

	if order != "" {
		call = call.Order(order)
	}

	if originalOpts.PublishedAfter != "" {
		call = call.PublishedAfter(originalOpts.PublishedAfter)
	}

	if originalOpts.PublishedBefore != "" {
		call = call.PublishedBefore(originalOpts.PublishedBefore)
	}

	response, err := call.Do()
	if err != nil {
		log.Printf("Warning: failed to search with params %s: %v", query, err)
		return []Video{}
	}

	videos, err := c.processSearchResults(response, originalOpts)
	if err != nil {
		log.Printf("Warning: failed to process search results: %v", err)
		return []Video{}
	}
	return videos
}

func (c *Client) processSearchResults(response *youtube.SearchListResponse, opts SearchOptions) ([]Video, error) {
	var videos []Video
	for _, item := range response.Items {
		// Parse published date
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			publishedAt = time.Now()
		}

		video := Video{
			ID:          item.Id.VideoId,
			Title:       item.Snippet.Title,
			Channel:     item.Snippet.ChannelTitle,
			PublishedAt: publishedAt,
			URL:         fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId),
			Description: item.Snippet.Description,
		}

		// Get video statistics
		stats, err := c.getVideoStats(item.Id.VideoId)
		if err != nil {
			log.Printf("Warning: failed to get stats for video %s: %v", item.Id.VideoId, err)
			continue
		}

		video.Views = stats.Views
		video.Likes = stats.Likes
		video.Comments = stats.Comments
		video.Duration = stats.Duration

		// Apply filters
		if opts.MinViews > 0 && video.Views < opts.MinViews {
			continue
		}
		if opts.MinLikes > 0 && video.Likes < opts.MinLikes {
			continue
		}

		videos = append(videos, video)
	}

	return videos, nil
}

func (c *Client) filterUniqueVideos(videos []Video, seenIDs map[string]bool) []Video {
	var uniqueVideos []Video
	for _, video := range videos {
		if !seenIDs[video.ID] {
			seenIDs[video.ID] = true
			uniqueVideos = append(uniqueVideos, video)
		}
	}
	return uniqueVideos
}

type VideoStats struct {
	Views    int64
	Likes    int64
	Comments int64
	Duration string
}

func (c *Client) getVideoStats(videoID string) (VideoStats, error) {
	call := c.service.Videos.List([]string{"statistics", "contentDetails"}).
		Id(videoID)

	response, err := call.Do()
	if err != nil {
		return VideoStats{}, err
	}

	if len(response.Items) == 0 {
		return VideoStats{}, fmt.Errorf("video not found")
	}

	item := response.Items[0]
	stats := VideoStats{}

	stats.Views = int64(item.Statistics.ViewCount)
	stats.Likes = int64(item.Statistics.LikeCount)
	stats.Comments = int64(item.Statistics.CommentCount)
	stats.Duration = item.ContentDetails.Duration

	return stats, nil
}
