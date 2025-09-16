package analysis

import (
	"fmt"
	"sort"
	"strings"

	"ytminer/utils"
	"ytminer/youtube"
)

type Analyzer struct {
	videos []youtube.Video
}

type GrowthPattern struct {
	TotalVideos    int     `json:"total_videos"`
	AvgViews       float64 `json:"avg_views"`
	AvgLikes       float64 `json:"avg_likes"`
	AvgComments     float64 `json:"avg_comments"`
	GrowthRate     float64 `json:"growth_rate"`
	TopPerformers  []VideoPerformance `json:"top_performers"`
	Insights       []string `json:"insights"`
}

type VideoPerformance struct {
	Title      string  `json:"title"`
	Channel    string  `json:"channel"`
	Views      int64   `json:"views"`
	Likes      int64   `json:"likes"`
	Engagement float64 `json:"engagement"`
	URL        string  `json:"url"`
}

type TitleAnalysis struct {
	CommonWords    []WordCount `json:"common_words"`
	CommonPhrases  []PhraseCount `json:"common_phrases"`
	Emojis         []EmojiCount `json:"emojis"`
	Patterns       []string `json:"patterns"`
	Insights       []string `json:"insights"`
}

type WordCount struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type PhraseCount struct {
	Phrase string `json:"phrase"`
	Count  int    `json:"count"`
}

type EmojiCount struct {
	Emoji string `json:"emoji"`
	Count int    `json:"count"`
}

type CompetitorAnalysis struct {
	TopChannels    []ChannelStats `json:"top_channels"`
	MarketShare    map[string]float64 `json:"market_share"`
	Opportunities  []string `json:"opportunities"`
	Insights       []string `json:"insights"`
}

type ChannelStats struct {
	Channel     string  `json:"channel"`
	VideoCount  int     `json:"video_count"`
	TotalViews  int64   `json:"total_views"`
	AvgViews    float64 `json:"avg_views"`
	Engagement  float64 `json:"engagement"`
}

type TemporalAnalysis struct {
	BestHours     []HourStats `json:"best_hours"`
	BestDays      []DayStats `json:"best_days"`
	GrowthPattern []TimeStats `json:"growth_pattern"`
	Insights      []string `json:"insights"`
}

type HourStats struct {
	Hour       int     `json:"hour"`
	AvgViews   float64 `json:"avg_views"`
	AvgLikes   float64 `json:"avg_likes"`
	Engagement float64 `json:"engagement"`
}

type DayStats struct {
	Day        string  `json:"day"`
	AvgViews   float64 `json:"avg_views"`
	AvgLikes   float64 `json:"avg_likes"`
	Engagement float64 `json:"engagement"`
}

type TimeStats struct {
	Period     string  `json:"period"`
	AvgViews   float64 `json:"avg_views"`
	GrowthRate float64 `json:"growth_rate"`
}

type KeywordAnalysis struct {
	TrendingKeywords []KeywordStats `json:"trending_keywords"`
	LongTailKeywords []KeywordStats `json:"long_tail_keywords"`
	SEOOpportunities []string `json:"seo_opportunities"`
	Insights         []string `json:"insights"`
}

type KeywordStats struct {
	Keyword     string  `json:"keyword"`
	Frequency   int     `json:"frequency"`
	AvgViews    float64 `json:"avg_views"`
	Engagement  float64 `json:"engagement"`
}

type ExecutiveReport struct {
	Summary           string `json:"summary"`
	KeyInsights       []string `json:"key_insights"`
	Recommendations   []string `json:"recommendations"`
	ContentStrategy   []string `json:"content_strategy"`
	CompetitiveIntel  []string `json:"competitive_intel"`
	PerformanceBench  []string `json:"performance_bench"`
	NextSteps         []string `json:"next_steps"`
}

func NewAnalyzer(videos []youtube.Video) *Analyzer {
	return &Analyzer{videos: videos}
}

func (a *Analyzer) AnalyzeGrowthPatterns() GrowthPattern {
	if len(a.videos) == 0 {
		return GrowthPattern{}
	}

	var totalViews, totalLikes, totalComments int64
	var topPerformers []VideoPerformance

	for _, video := range a.videos {
		totalViews += video.Views
		totalLikes += video.Likes
		totalComments += video.Comments

		engagement := float64(video.Likes+video.Comments) / max1(float64(video.Views)) * 100
		topPerformers = append(topPerformers, VideoPerformance{
			Title:      video.Title,
			Channel:    video.Channel,
			Views:      video.Views,
			Likes:      video.Likes,
			Engagement: engagement,
			URL:        video.URL,
		})
	}

	// Sort by engagement
	sort.Slice(topPerformers, func(i, j int) bool {
		return topPerformers[i].Engagement > topPerformers[j].Engagement
	})
	if len(topPerformers) > 5 {
		topPerformers = topPerformers[:5]
	}

	avgViews := float64(totalViews) / float64(len(a.videos))
	avgLikes := float64(totalLikes) / float64(len(a.videos))
	avgComments := float64(totalComments) / float64(len(a.videos))

	// Calculate growth slope (simple linear regression on views vs index)
	growthRate := a.calculateGrowthSlope()

	insights := a.generateGrowthInsights(avgViews, avgLikes, growthRate)

	return GrowthPattern{
		TotalVideos:   len(a.videos),
		AvgViews:      avgViews,
		AvgLikes:      avgLikes,
		AvgComments:   avgComments,
		GrowthRate:    growthRate,
		TopPerformers: topPerformers,
		Insights:      insights,
	}
}

func (a *Analyzer) AnalyzeTitles() TitleAnalysis {
	if len(a.videos) == 0 {
		return TitleAnalysis{}
	}

	wordCounts := make(map[string]int)
	phraseCounts := make(map[string]int)
	emojiCounts := make(map[string]int)
	var patterns []string

	for _, video := range a.videos {
		titleLower := strings.ToLower(video.Title)

		// Words (normalized, stopwords removed)
		words := utils.Tokenize(video.Title)
		for _, w := range words {
			wordCounts[w]++
		}

		// Phrases (bigrams) from tokenized words
		for i := 0; i < len(words)-1; i++ {
			phrase := words[i] + " " + words[i+1]
			phraseCounts[phrase]++
		}

		// Emojis
		for _, e := range utils.ExtractEmojis(video.Title) {
			emojiCounts[e]++
		}

		// Patterns
		if strings.Contains(titleLower, "tutorial") { patterns = append(patterns, "Tutorial Pattern") }
		if strings.Contains(titleLower, "how to") { patterns = append(patterns, "How-to Pattern") }
		if strings.Contains(titleLower, "2024") || strings.Contains(titleLower, "2023") { patterns = append(patterns, "Year Pattern") }
	}

	var commonWords []WordCount
	for word, count := range wordCounts {
		commonWords = append(commonWords, WordCount{Word: word, Count: count})
	}
	sort.Slice(commonWords, func(i, j int) bool { return commonWords[i].Count > commonWords[j].Count })
	if len(commonWords) > 10 { commonWords = commonWords[:10] }

	var commonPhrases []PhraseCount
	for phrase, count := range phraseCounts {
		commonPhrases = append(commonPhrases, PhraseCount{Phrase: phrase, Count: count})
	}
	sort.Slice(commonPhrases, func(i, j int) bool { return commonPhrases[i].Count > commonPhrases[j].Count })
	if len(commonPhrases) > 5 { commonPhrases = commonPhrases[:5] }

	var emojis []EmojiCount
	for emoji, count := range emojiCounts {
		emojis = append(emojis, EmojiCount{Emoji: emoji, Count: count})
	}
	sort.Slice(emojis, func(i, j int) bool { return emojis[i].Count > emojis[j].Count })

	insights := a.generateTitleInsights(commonWords, commonPhrases, patterns)

	return TitleAnalysis{
		CommonWords:   commonWords,
		CommonPhrases: commonPhrases,
		Emojis:        emojis,
		Patterns:      patterns,
		Insights:      insights,
	}
}

func (a *Analyzer) AnalyzeCompetitors() CompetitorAnalysis {
	if len(a.videos) == 0 {
		return CompetitorAnalysis{}
	}

	channelStats := make(map[string]*ChannelStats)
	var totalViews int64

	for _, video := range a.videos {
		if stats, exists := channelStats[video.Channel]; exists {
			stats.VideoCount++
			stats.TotalViews += video.Views
			stats.AvgViews = float64(stats.TotalViews) / float64(stats.VideoCount)
			stats.Engagement = float64(video.Likes+video.Comments) / float64(video.Views) * 100
		} else {
			engagement := float64(video.Likes+video.Comments) / float64(video.Views) * 100
			channelStats[video.Channel] = &ChannelStats{
				Channel:    video.Channel,
				VideoCount: 1,
				TotalViews: video.Views,
				AvgViews:   float64(video.Views),
				Engagement: engagement,
			}
		}
		totalViews += video.Views
	}

	// Convert to slice and sort
	var topChannels []ChannelStats
	for _, stats := range channelStats {
		topChannels = append(topChannels, *stats)
	}
	sort.Slice(topChannels, func(i, j int) bool {
		return topChannels[i].TotalViews > topChannels[j].TotalViews
	})
	if len(topChannels) > 5 {
		topChannels = topChannels[:5]
	}

	// Calculate market share
	marketShare := make(map[string]float64)
	for _, channel := range topChannels {
		marketShare[channel.Channel] = float64(channel.TotalViews) / float64(totalViews) * 100
	}

	opportunities := a.generateCompetitorOpportunities(topChannels, marketShare)
	insights := a.generateCompetitorInsights(topChannels, marketShare)

	return CompetitorAnalysis{
		TopChannels:   topChannels,
		MarketShare:   marketShare,
		Opportunities: opportunities,
		Insights:      insights,
	}
}

func (a *Analyzer) AnalyzeTemporal() TemporalAnalysis {
	if len(a.videos) == 0 { return TemporalAnalysis{} }

	hourStats := make(map[int]*HourStats)
	hourN := make(map[int]int)
	dayStats := make(map[string]*DayStats)
	dayN := make(map[string]int)

	for _, video := range a.videos {
		hour := video.PublishedAt.Hour()
		day := video.PublishedAt.Weekday().String()

		eng := float64(video.Likes+video.Comments) / max1(float64(video.Views)) * 100

		if stats, ok := hourStats[hour]; ok {
			stats.AvgViews += float64(video.Views)
			stats.AvgLikes += float64(video.Likes)
			stats.Engagement += eng
			hourN[hour]++
		} else {
			hourStats[hour] = &HourStats{Hour: hour, AvgViews: float64(video.Views), AvgLikes: float64(video.Likes), Engagement: eng}
			hourN[hour] = 1
		}

		if stats, ok := dayStats[day]; ok {
			stats.AvgViews += float64(video.Views)
			stats.AvgLikes += float64(video.Likes)
			stats.Engagement += eng
			dayN[day]++
		} else {
			dayStats[day] = &DayStats{Day: day, AvgViews: float64(video.Views), AvgLikes: float64(video.Likes), Engagement: eng}
			dayN[day] = 1
		}
	}

	// finalize means
	const minN = 5
	var bestHours []HourStats
	for h, s := range hourStats {
		n := hourN[h]
		if n < minN { continue }
		s.AvgViews /= float64(n)
		s.AvgLikes /= float64(n)
		s.Engagement /= float64(n)
		bestHours = append(bestHours, *s)
	}
	sort.Slice(bestHours, func(i, j int) bool { return bestHours[i].Engagement > bestHours[j].Engagement })

	var bestDays []DayStats
	for d, s := range dayStats {
		n := dayN[d]
		if n < minN { continue }
		s.AvgViews /= float64(n)
		s.AvgLikes /= float64(n)
		s.Engagement /= float64(n)
		bestDays = append(bestDays, *s)
	}
	sort.Slice(bestDays, func(i, j int) bool { return bestDays[i].Engagement > bestDays[j].Engagement })

	insights := a.generateTemporalInsights(bestHours, bestDays)
	return TemporalAnalysis{ BestHours: bestHours, BestDays: bestDays, Insights: insights }
}

func (a *Analyzer) AnalyzeKeywords() KeywordAnalysis {
	if len(a.videos) == 0 { return KeywordAnalysis{} }

	keywordStats := make(map[string]*KeywordStats)

	for _, video := range a.videos {
		words := utils.Tokenize(video.Title)
		engagement := float64(video.Likes+video.Comments) / max1(float64(video.Views)) * 100
		for _, w := range words {
			if stats, ok := keywordStats[w]; ok {
				stats.Frequency++
				stats.AvgViews += float64(video.Views)
				stats.Engagement += engagement
			} else {
				keywordStats[w] = &KeywordStats{ Keyword: w, Frequency: 1, AvgViews: float64(video.Views), Engagement: engagement }
			}
		}
	}

	var trendingKeywords []KeywordStats
	for _, s := range keywordStats { trendingKeywords = append(trendingKeywords, *s) }
	sort.Slice(trendingKeywords, func(i, j int) bool { return trendingKeywords[i].Frequency > trendingKeywords[j].Frequency })
	if len(trendingKeywords) > 10 { trendingKeywords = trendingKeywords[:10] }

	var longTailKeywords []KeywordStats
	for _, s := range keywordStats {
		if s.Frequency < 3 && s.Engagement/float64(s.Frequency) > 5.0 {
			longTailKeywords = append(longTailKeywords, *s)
		}
	}
	sort.Slice(longTailKeywords, func(i, j int) bool { return longTailKeywords[i].Engagement/float64(longTailKeywords[i].Frequency) > longTailKeywords[j].Engagement/float64(longTailKeywords[j].Frequency) })

	// finalize means
	for i := range trendingKeywords {
		fk := trendingKeywords[i].Frequency
		if fk > 0 {
			trendingKeywords[i].AvgViews /= float64(fk)
			trendingKeywords[i].Engagement /= float64(fk)
		}
	}
	for i := range longTailKeywords {
		fk := longTailKeywords[i].Frequency
		if fk > 0 {
			longTailKeywords[i].AvgViews /= float64(fk)
			longTailKeywords[i].Engagement /= float64(fk)
		}
	}

	seoOpportunities := a.generateSEOOpportunities(trendingKeywords, longTailKeywords)
	insights := a.generateKeywordInsights(trendingKeywords, longTailKeywords)

	return KeywordAnalysis{ TrendingKeywords: trendingKeywords, LongTailKeywords: longTailKeywords, SEOOpportunities: seoOpportunities, Insights: insights }
}

func (a *Analyzer) GenerateExecutiveReport() ExecutiveReport {
	growth := a.AnalyzeGrowthPatterns()
	titles := a.AnalyzeTitles()
	competitors := a.AnalyzeCompetitors()
	temporal := a.AnalyzeTemporal()
	keywords := a.AnalyzeKeywords()

	summary := a.generateExecutiveSummary(growth, titles, competitors)
	keyInsights := a.generateKeyInsights(growth, titles, competitors, temporal, keywords)
	recommendations := a.generateRecommendations(growth, titles, competitors, temporal, keywords)
	contentStrategy := a.generateContentStrategy(titles, keywords, temporal)
	competitiveIntel := a.generateCompetitiveIntel(competitors)
	performanceBench := a.generatePerformanceBenchmarks(growth, competitors)
	nextSteps := a.generateNextSteps(growth, titles, competitors)

	return ExecutiveReport{
		Summary:          summary,
		KeyInsights:      keyInsights,
		Recommendations:  recommendations,
		ContentStrategy:  contentStrategy,
		CompetitiveIntel: competitiveIntel,
		PerformanceBench: performanceBench,
		NextSteps:        nextSteps,
	}
}

// Helper methods
func (a *Analyzer) calculateGrowthSlope() float64 {
	if len(a.videos) < 2 { return 0 }
	// Sort by published date
	sorted := make([]youtube.Video, len(a.videos))
	copy(sorted, a.videos)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].PublishedAt.Before(sorted[j].PublishedAt) })
	// Simple linear regression on index (time proxy) vs views
	n := float64(len(sorted))
	var sumX, sumY, sumXY, sumXX float64
	for i, v := range sorted {
		x := float64(i)
		y := float64(v.Views)
		sumX += x; sumY += y; sumXY += x*y; sumXX += x*x
	}
	den := (n*sumXX - sumX*sumX)
	if den == 0 { return 0 }
	slope := (n*sumXY - sumX*sumY) / den
	// Express as percent of first views if possível
	base := float64(sorted[0].Views)
	if base <= 0 { return 0 }
	return slope / base * 100
}

func max1(v float64) float64 { if v <= 0 { return 1 } ; return v }

func (a *Analyzer) generateGrowthInsights(avgViews, avgLikes, growthRate float64) []string {
	var insights []string

	if avgViews > 1000000 {
		insights = append(insights, "High-performing content with over 1M average views")
	} else if avgViews > 100000 {
		insights = append(insights, "Good performance with 100K+ average views")
	} else {
		insights = append(insights, "Room for improvement in view counts")
	}

	if growthRate > 50 {
		insights = append(insights, "Strong growth trend detected")
	} else if growthRate > 0 {
		insights = append(insights, "Moderate growth trend")
	} else {
		insights = append(insights, "Declining trend - needs attention")
	}

	engagement := avgLikes / avgViews * 100
	if engagement > 5 {
		insights = append(insights, "Excellent engagement rate")
	} else if engagement > 2 {
		insights = append(insights, "Good engagement rate")
	} else {
		insights = append(insights, "Low engagement - consider content improvements")
	}

	return insights
}

func (a *Analyzer) generateTitleInsights(words []WordCount, phrases []PhraseCount, patterns []string) []string {
	var insights []string

	if len(words) > 0 {
		insights = append(insights, fmt.Sprintf("Most common word: '%s' (%d times)", words[0].Word, words[0].Count))
	}

	if len(phrases) > 0 {
		insights = append(insights, fmt.Sprintf("Most common phrase: '%s' (%d times)", phrases[0].Phrase, phrases[0].Count))
	}

	if len(patterns) > 0 {
		insights = append(insights, fmt.Sprintf("Common patterns: %s", strings.Join(patterns, ", ")))
	}

	return insights
}

func (a *Analyzer) generateCompetitorOpportunities(channels []ChannelStats, marketShare map[string]float64) []string {
	var opportunities []string

	if len(channels) > 0 {
		topChannel := channels[0]
		opportunities = append(opportunities, fmt.Sprintf("Top channel '%s' has %.1f%% market share - opportunity to compete", topChannel.Channel, marketShare[topChannel.Channel]))
	}

	// Find gaps
	var totalShare float64
	for _, share := range marketShare {
		totalShare += share
	}
	if totalShare < 80 {
		opportunities = append(opportunities, "Significant market opportunity - top channels don't dominate completely")
	}

	return opportunities
}

func (a *Analyzer) generateCompetitorInsights(channels []ChannelStats, marketShare map[string]float64) []string {
	var insights []string

	if len(channels) > 0 {
		insights = append(insights, fmt.Sprintf("Top performing channel: %s with %.0f average views", channels[0].Channel, channels[0].AvgViews))
	}

	// Analyze engagement
	var totalEngagement float64
	for _, channel := range channels {
		totalEngagement += channel.Engagement
	}
	avgEngagement := totalEngagement / float64(len(channels))
	insights = append(insights, fmt.Sprintf("Average engagement rate: %.2f%%", avgEngagement))

	return insights
}

func (a *Analyzer) generateTemporalInsights(hours []HourStats, days []DayStats) []string {
	var insights []string

	if len(hours) > 0 {
		bestHour := hours[0]
		insights = append(insights, fmt.Sprintf("Best posting hour: %d:00 with %.2f%% engagement", bestHour.Hour, bestHour.Engagement))
	}

	if len(days) > 0 {
		bestDay := days[0]
		insights = append(insights, fmt.Sprintf("Best posting day: %s with %.2f%% engagement", bestDay.Day, bestDay.Engagement))
	}

	return insights
}

func (a *Analyzer) generateSEOOpportunities(trending, longTail []KeywordStats) []string {
	var opportunities []string

	if len(trending) > 0 {
		opportunities = append(opportunities, fmt.Sprintf("High-volume keyword: '%s' (%d mentions)", trending[0].Keyword, trending[0].Frequency))
	}

	if len(longTail) > 0 {
		opportunities = append(opportunities, fmt.Sprintf("High-engagement long-tail keyword: '%s' (%.2f%% engagement)", longTail[0].Keyword, longTail[0].Engagement))
	}

	return opportunities
}

func (a *Analyzer) generateKeywordInsights(trending, longTail []KeywordStats) []string {
	var insights []string

	insights = append(insights, fmt.Sprintf("Analyzed %d unique keywords", len(trending)+len(longTail)))
	insights = append(insights, fmt.Sprintf("Found %d long-tail opportunities", len(longTail)))

	return insights
}

func (a *Analyzer) generateExecutiveSummary(growth GrowthPattern, titles TitleAnalysis, competitors CompetitorAnalysis) string {
	return fmt.Sprintf("Analysis of %d videos shows %s with %s engagement. Top channel '%s' leads with %.1f%% market share.", 
		growth.TotalVideos, 
		utils.FormatNumber(growth.AvgViews), 
		utils.FormatEngagement(growth.AvgLikes/growth.AvgViews*100),
		competitors.TopChannels[0].Channel,
		competitors.MarketShare[competitors.TopChannels[0].Channel])
}

func (a *Analyzer) generateKeyInsights(growth GrowthPattern, titles TitleAnalysis, competitors CompetitorAnalysis, temporal TemporalAnalysis, keywords KeywordAnalysis) []string {
	var insights []string

	insights = append(insights, fmt.Sprintf("Average views: %s", utils.FormatNumber(growth.AvgViews)))
	insights = append(insights, fmt.Sprintf("Growth rate: %.1f%%", growth.GrowthRate))
	insights = append(insights, fmt.Sprintf("Top keyword: '%s'", keywords.TrendingKeywords[0].Keyword))
	
	if len(temporal.BestHours) > 0 {
		insights = append(insights, fmt.Sprintf("Best posting time: %d:00", temporal.BestHours[0].Hour))
	}

	return insights
}

func (a *Analyzer) generateRecommendations(growth GrowthPattern, titles TitleAnalysis, competitors CompetitorAnalysis, temporal TemporalAnalysis, keywords KeywordAnalysis) []string {
	var recommendations []string

	if growth.GrowthRate < 0 {
		recommendations = append(recommendations, "Focus on content quality to reverse declining trend")
	}

	if len(temporal.BestHours) > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Post at %d:00 for maximum engagement", temporal.BestHours[0].Hour))
	}

	if len(keywords.LongTailKeywords) > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Target long-tail keyword: '%s'", keywords.LongTailKeywords[0].Keyword))
	}

	return recommendations
}

func (a *Analyzer) generateContentStrategy(titles TitleAnalysis, keywords KeywordAnalysis, temporal TemporalAnalysis) []string {
	var strategy []string

	if len(titles.CommonWords) > 0 {
		strategy = append(strategy, fmt.Sprintf("Use common words: %s", titles.CommonWords[0].Word))
	}

	if len(keywords.TrendingKeywords) > 0 {
		strategy = append(strategy, fmt.Sprintf("Focus on trending keyword: '%s'", keywords.TrendingKeywords[0].Keyword))
	}

	return strategy
}

func (a *Analyzer) generateCompetitiveIntel(competitors CompetitorAnalysis) []string {
	var intel []string

	if len(competitors.TopChannels) > 0 {
		intel = append(intel, fmt.Sprintf("Top competitor: %s", competitors.TopChannels[0].Channel))
	}

	return intel
}

func (a *Analyzer) generatePerformanceBenchmarks(growth GrowthPattern, competitors CompetitorAnalysis) []string {
	var benchmarks []string

	benchmarks = append(benchmarks, fmt.Sprintf("Your average: %s views", utils.FormatNumber(growth.AvgViews)))
	
	if len(competitors.TopChannels) > 0 {
		benchmarks = append(benchmarks, fmt.Sprintf("Top competitor: %s views", utils.FormatNumber(competitors.TopChannels[0].AvgViews)))
	}

	return benchmarks
}

func (a *Analyzer) generateNextSteps(growth GrowthPattern, titles TitleAnalysis, competitors CompetitorAnalysis) []string {
	var steps []string

	steps = append(steps, "Analyze top-performing content patterns")
	steps = append(steps, "Implement keyword strategy")
	steps = append(steps, "Optimize posting schedule")
	steps = append(steps, "Monitor competitor activity")

	return steps
}

