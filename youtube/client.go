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

type SearchOptions struct {
	Query      string
	MaxResults int64
	Region     string
	Duration   string
	Order      string
	MinViews   int64
	MinLikes   int64
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

	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search videos: %v", err)
	}

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
