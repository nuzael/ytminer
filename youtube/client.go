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
	ChannelID   string    `json:"channel_id"`
	PublishedAt time.Time `json:"published_at"`
	Views       int64     `json:"views"`
	Likes       int64     `json:"likes"`
	Comments    int64     `json:"comments"`
	Duration    string    `json:"duration"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	VPD         float64   `json:"vpd"`
}

type AnalysisLevel int

const (
	QuickScan AnalysisLevel = iota
	Balanced
	DeepDive
)

type SearchOptions struct {
	Query           string
	MaxResults      int64
	Region          string
	Duration        string
	Order           string
	MinViews        int64
	MinLikes        int64
	Level           AnalysisLevel
	PublishedAfter  string
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
	// Use level-based search for all levels now
	return c.SearchVideosWithLevel(opts)
}

func (c *Client) SearchVideosWithLevel(opts SearchOptions) ([]Video, error) {
	var allVideos []Video
	var seenIDs = make(map[string]bool)

	switch opts.Level {
	case QuickScan:
		allVideos = c.searchQuickScan(opts, seenIDs)
	case Balanced:
		allVideos = c.searchBalanced(opts, seenIDs)
	case DeepDive:
		allVideos = c.searchDeepDive(opts, seenIDs)
	}

	return allVideos, nil
}

func (c *Client) searchQuickScan(opts SearchOptions, seenIDs map[string]bool) []Video {
	var allVideos []Video
	totalSearches := 0
	totalItemsIn := 0
	totalItemsOut := 0
	targetVideos := 150 // Target for Quick Scan

	userOrder := opts.Order
	if userOrder == "" {
		userOrder = "relevance"
	}

	log.Printf("searchQuickScan: Starting adaptive strategy targeting %d videos", targetVideos)

	// Phase 1: Try up to 3 pages of user's chosen order
	pageToken := ""
	for page := 1; page <= 3; page++ {
		videos, nextToken := c.searchWithPagination(opts.Query, opts.Region, opts.Duration, userOrder, 50, pageToken, opts)
		totalSearches++
		totalItemsIn += len(videos)
		unique := c.filterUniqueVideos(videos, seenIDs)
		totalItemsOut += len(unique)
		allVideos = append(allVideos, unique...)

		pageToken = nextToken
		if pageToken == "" {
			log.Printf("searchQuickScan: No more pages available for order=%s after %d pages", userOrder, page)
			break
		}
	}

	// Phase 2: If we have less than 50% of target, complement with other orders
	currentCount := len(allVideos)
	if currentCount < targetVideos/2 {
		log.Printf("searchQuickScan: Only %d videos from primary order, complementing with backup orders", currentCount)

		// Try relevance as backup (if not already used)
		if userOrder != "relevance" {
			backupVideos := c.searchWithBackupOrder(opts, seenIDs, "relevance", 2) // 2 pages max
			totalSearches += 2
			totalItemsIn += len(backupVideos)
			unique := c.filterUniqueVideos(backupVideos, seenIDs)
			totalItemsOut += len(unique)
			allVideos = append(allVideos, unique...)
		}

		// If still need more, try viewCount (if not already used)
		if len(allVideos) < targetVideos/2 && userOrder != "viewCount" {
			backupVideos := c.searchWithBackupOrder(opts, seenIDs, "viewCount", 1) // 1 page
			totalSearches += 1
			totalItemsIn += len(backupVideos)
			unique := c.filterUniqueVideos(backupVideos, seenIDs)
			totalItemsOut += len(unique)
			allVideos = append(allVideos, unique...)
		}
	}

	duplicatesRemoved := totalItemsIn - totalItemsOut
	log.Printf("searchQuickScan SUMMARY: searches=%d items_in=%d duplicates_removed=%d final_out=%d target=%d", totalSearches, totalItemsIn, duplicatesRemoved, len(allVideos), targetVideos)
	return allVideos
}

func (c *Client) searchBalanced(opts SearchOptions, seenIDs map[string]bool) []Video {
	var allVideos []Video
	totalSearches := 0
	totalItemsIn := 0
	totalItemsOut := 0
	targetVideos := 400 // Target for Balanced

	userOrder := opts.Order
	if userOrder == "" {
		userOrder = "relevance"
	}

	log.Printf("searchBalanced: Starting adaptive strategy targeting %d videos", targetVideos)

	// Phase 1: Try up to 5 pages of user's chosen order
	pageToken := ""
	for page := 1; page <= 5; page++ {
		videos, nextToken := c.searchWithPagination(opts.Query, opts.Region, opts.Duration, userOrder, 50, pageToken, opts)
		totalSearches++
		totalItemsIn += len(videos)
		unique := c.filterUniqueVideos(videos, seenIDs)
		totalItemsOut += len(unique)
		allVideos = append(allVideos, unique...)

		pageToken = nextToken
		if pageToken == "" {
			log.Printf("searchBalanced: No more pages available for order=%s after %d pages", userOrder, page)
			break
		}
	}

	// Phase 2: If we have less than 50% of target, complement with other orders
	currentCount := len(allVideos)
	if currentCount < targetVideos/2 {
		log.Printf("searchBalanced: Only %d videos from primary order, complementing with backup orders", currentCount)

		// Backup sequence: relevance, viewCount, date
		backupOrders := []string{"relevance", "viewCount", "date"}
		pagesPerBackup := []int{3, 2, 1} // 3+2+1 = 6 additional searches max

		for i, backupOrder := range backupOrders {
			if backupOrder != userOrder && len(allVideos) < targetVideos*3/4 {
				backupVideos := c.searchWithBackupOrder(opts, seenIDs, backupOrder, pagesPerBackup[i])
				totalSearches += pagesPerBackup[i]
				totalItemsIn += len(backupVideos)
				unique := c.filterUniqueVideos(backupVideos, seenIDs)
				totalItemsOut += len(unique)
				allVideos = append(allVideos, unique...)
			}
		}
	}

	duplicatesRemoved := totalItemsIn - totalItemsOut
	log.Printf("searchBalanced SUMMARY: searches=%d items_in=%d duplicates_removed=%d final_out=%d target=%d", totalSearches, totalItemsIn, duplicatesRemoved, len(allVideos), targetVideos)
	return allVideos
}

func (c *Client) searchDeepDive(opts SearchOptions, seenIDs map[string]bool) []Video {
	var allVideos []Video
	totalSearches := 0
	totalItemsIn := 0
	totalItemsOut := 0
	targetVideos := 1000 // Target for Deep Dive

	userOrder := opts.Order
	if userOrder == "" {
		userOrder = "relevance"
	}

	log.Printf("searchDeepDive: Starting adaptive strategy targeting %d videos", targetVideos)

	// Phase 1: Try up to 10 pages of user's chosen order
	pageToken := ""
	for page := 1; page <= 10; page++ {
		videos, nextToken := c.searchWithPagination(opts.Query, opts.Region, opts.Duration, userOrder, 50, pageToken, opts)
		totalSearches++
		totalItemsIn += len(videos)
		unique := c.filterUniqueVideos(videos, seenIDs)
		totalItemsOut += len(unique)
		allVideos = append(allVideos, unique...)

		pageToken = nextToken
		if pageToken == "" {
			log.Printf("searchDeepDive: No more pages available for order=%s after %d pages", userOrder, page)
			break
		}
	}

	// Phase 2: If we have less than 50% of target, complement with other orders
	currentCount := len(allVideos)
	if currentCount < targetVideos/2 {
		log.Printf("searchDeepDive: Only %d videos from primary order, complementing with backup orders", currentCount)

		// Comprehensive backup sequence: relevance, viewCount, date, rating
		backupOrders := []string{"relevance", "viewCount", "date", "rating"}
		pagesPerBackup := []int{5, 3, 2, 2} // 5+3+2+2 = 12 additional searches max

		for i, backupOrder := range backupOrders {
			if backupOrder != userOrder && len(allVideos) < targetVideos*3/4 {
				backupVideos := c.searchWithBackupOrder(opts, seenIDs, backupOrder, pagesPerBackup[i])
				totalSearches += pagesPerBackup[i]
				totalItemsIn += len(backupVideos)
				unique := c.filterUniqueVideos(backupVideos, seenIDs)
				totalItemsOut += len(unique)
				allVideos = append(allVideos, unique...)
			}
		}
	}

	duplicatesRemoved := totalItemsIn - totalItemsOut
	log.Printf("searchDeepDive SUMMARY: searches=%d items_in=%d duplicates_removed=%d final_out=%d target=%d", totalSearches, totalItemsIn, duplicatesRemoved, len(allVideos), targetVideos)
	return allVideos
}

func (c *Client) searchWithPagination(query, region, duration, order string, maxResults int64, pageToken string, originalOpts SearchOptions) ([]Video, string) {
	call := c.service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(50). // Always request max 50 per API call
		Type("video")

	if pageToken != "" {
		call = call.PageToken(pageToken)
	}
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
		return []Video{}, ""
	}

	videos, tele := c.processSearchResultsBatched(response.Items, originalOpts)
	log.Printf("searchWithPagination: q=%q r=%s d=%s o=%s items_in=%d stats_failed=%d filtered=%d items_out=%d", query, region, duration, order, len(response.Items), tele.statsFailed, tele.filtered, len(videos))

	return videos, response.NextPageToken
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
			ChannelID:   item.Snippet.ChannelId,
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

		// Calculate VPD (Views Per Day)
		video.VPD = calculateVPD(video.Views, video.PublishedAt)

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

func (c *Client) processSearchResultsBatched(collected []*youtube.SearchResult, opts SearchOptions) ([]Video, Telemetry) {
	var videos []Video
	var tele Telemetry

	for _, item := range collected {
		// Parse published date
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			publishedAt = time.Now()
		}

		video := Video{
			ID:          item.Id.VideoId,
			Title:       item.Snippet.Title,
			Channel:     item.Snippet.ChannelTitle,
			ChannelID:   item.Snippet.ChannelId,
			PublishedAt: publishedAt,
			URL:         fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.Id.VideoId),
			Description: item.Snippet.Description,
		}

		// Get video statistics
		stats, err := c.getVideoStats(item.Id.VideoId)
		if err != nil {
			log.Printf("Warning: failed to get stats for video %s: %v", item.Id.VideoId, err)
			tele.statsFailed++
			continue
		}

		video.Views = stats.Views
		video.Likes = stats.Likes
		video.Comments = stats.Comments
		video.Duration = stats.Duration

		// Calculate VPD (Views Per Day)
		video.VPD = calculateVPD(video.Views, video.PublishedAt)

		// Apply filters
		if opts.MinViews > 0 && video.Views < opts.MinViews {
			tele.filtered++
			continue
		}
		if opts.MinLikes > 0 && video.Likes < opts.MinLikes {
			tele.filtered++
			continue
		}

		videos = append(videos, video)
	}

	return videos, tele
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

// calculateVPD calculates Views Per Day using VPD = video.viewCount / max(1, days_since(video.publishedAt))
func calculateVPD(views int64, publishedAt time.Time) float64 {
	now := time.Now()
	daysSince := int(now.Sub(publishedAt).Hours() / 24)

	// Protection against division by zero using max(1, daysSince)
	if daysSince < 1 {
		daysSince = 1
	}

	return float64(views) / float64(daysSince)
}

type Telemetry struct {
	statsFailed int
	filtered    int
}

// Helper function to search with backup orders
func (c *Client) searchWithBackupOrder(opts SearchOptions, seenIDs map[string]bool, order string, maxPages int) []Video {
	var backupVideos []Video
	pageToken := ""

	log.Printf("searchWithBackupOrder: Trying %d pages of order=%s", maxPages, order)

	for page := 1; page <= maxPages; page++ {
		videos, nextToken := c.searchWithPagination(opts.Query, opts.Region, opts.Duration, order, 50, pageToken, opts)
		backupVideos = append(backupVideos, videos...)

		pageToken = nextToken
		if pageToken == "" {
			log.Printf("searchWithBackupOrder: No more pages available for order=%s after %d pages", order, page)
			break
		}
	}

	return backupVideos
}
