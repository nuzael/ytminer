package analysis

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"math"
	"ytminer/config"
	"ytminer/utils"
	"ytminer/youtube"
)

type Analyzer struct {
	videos []youtube.Video
	cfg    *config.AppConfig
}

type GrowthPattern struct {
	TotalVideos        int                `json:"total_videos"`
	AvgViews           float64            `json:"avg_views"`
	AvgLikes           float64            `json:"avg_likes"`
	AvgComments        float64            `json:"avg_comments"`
	NicheVelocityScore float64            `json:"niche_velocity_score"`
	TopPerformers      []VideoPerformance `json:"top_performers"`
	Insights           []string           `json:"insights"`
}

type VideoPerformance struct {
	Title      string  `json:"title"`
	Channel    string  `json:"channel"`
	Views      int64   `json:"views"`
	Likes      int64   `json:"likes"`
	Engagement float64 `json:"engagement"`
	VPD        float64 `json:"vpd"`
	URL        string  `json:"url"`
}

type TitleAnalysis struct {
	CommonWords   []WordCount   `json:"common_words"`
	CommonPhrases []PhraseCount `json:"common_phrases"`
	Emojis        []EmojiCount  `json:"emojis"`
	Patterns      []string      `json:"patterns"`
	Insights      []string      `json:"insights"`
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
	TopChannels   []ChannelStats     `json:"top_channels"`
	RisingStars   []ChannelStats     `json:"rising_stars"`
	MarketShare   map[string]float64 `json:"market_share"`
	Opportunities []string           `json:"opportunities"`
	Insights      []string           `json:"insights"`
}

type ChannelStats struct {
	Channel      string  `json:"channel"`
	ChannelID    string  `json:"channel_id"`
	ChannelURL   string  `json:"channel_url"`
	VideoCount   int     `json:"video_count"`
	TotalViews   int64   `json:"total_views"`
	AvgViews     float64 `json:"avg_views"`
	AvgVPD       float64 `json:"avg_vpd"`
	Engagement   float64 `json:"engagement"`
	IsRisingStar bool    `json:"is_rising_star"`
}

type TemporalAnalysis struct {
	BestHours []HourStats `json:"best_hours"`
	BestDays  []DayStats  `json:"best_days"`
	Insights  []string    `json:"insights"`
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

type KeywordAnalysis struct {
	TrendingKeywords []KeywordStats `json:"trending_keywords"`
	CoreKeywords     []KeywordStats `json:"core_keywords"`
	LongTailKeywords []KeywordStats `json:"long_tail_keywords"`
	SEOOpportunities []string       `json:"seo_opportunities"`
	Insights         []string       `json:"insights"`
}

type KeywordStats struct {
	Keyword    string  `json:"keyword"`
	Frequency  int     `json:"frequency"`
	AvgViews   float64 `json:"avg_views"`
	AvgVPD     float64 `json:"avg_vpd"`
	Engagement float64 `json:"engagement"`
}

type ExecutiveReport struct {
	Summary          string   `json:"summary"`
	KeyInsights      []string `json:"key_insights"`
	Recommendations  []string `json:"recommendations"`
	ContentStrategy  []string `json:"content_strategy"`
	CompetitiveIntel []string `json:"competitive_intel"`
	PerformanceBench []string `json:"performance_bench"`
	NextSteps        []string `json:"next_steps"`
}

type OpportunityItem struct {
	Title      string   `json:"title"`
	Channel    string   `json:"channel"`
	URL        string   `json:"url"`
	Score      float64  `json:"score"`
	VPD        float64  `json:"vpd"`
	LikeRate   float64  `json:"like_rate"`
	AgeDays    int      `json:"age_days"`
	Saturation float64  `json:"saturation"`
	Reasons    []string `json:"reasons"`
}

func NewAnalyzer(videos []youtube.Video, cfg *config.AppConfig) *Analyzer {
	return &Analyzer{videos: videos, cfg: cfg}
}

func (a *Analyzer) AnalyzeGrowthPatterns() GrowthPattern {
	if len(a.videos) == 0 {
		return GrowthPattern{}
	}

	var totalViews, totalLikes, totalComments int64
	var totalVPD float64
	var topPerformers []VideoPerformance

	for _, video := range a.videos {
		totalViews += video.Views
		totalLikes += video.Likes
		totalComments += video.Comments
		totalVPD += video.VPD

		engagement := float64(video.Likes+video.Comments) / max1(float64(video.Views)) * 100
		topPerformers = append(topPerformers, VideoPerformance{
			Title:      video.Title,
			Channel:    video.Channel,
			Views:      video.Views,
			Likes:      video.Likes,
			Engagement: engagement,
			VPD:        video.VPD,
			URL:        video.URL,
		})
	}

	// Sort by VPD (Views Per Day) in descending order - this shows "Highest Velocity Videos"
	sort.Slice(topPerformers, func(i, j int) bool {
		return topPerformers[i].VPD > topPerformers[j].VPD
	})
	if len(topPerformers) > 5 {
		topPerformers = topPerformers[:5]
	}

	avgViews := float64(totalViews) / float64(len(a.videos))
	avgLikes := float64(totalLikes) / float64(len(a.videos))
	avgComments := float64(totalComments) / float64(len(a.videos))

	// Calculate Niche Velocity Score (Avg VPD)
	nicheVelocityScore := totalVPD / float64(len(a.videos))

	insights := a.generateGrowthInsights(avgViews, avgLikes, nicheVelocityScore)

	return GrowthPattern{
		TotalVideos:        len(a.videos),
		AvgViews:           avgViews,
		AvgLikes:           avgLikes,
		AvgComments:        avgComments,
		NicheVelocityScore: nicheVelocityScore,
		TopPerformers:      topPerformers,
		Insights:           insights,
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
		if strings.Contains(titleLower, "tutorial") {
			patterns = append(patterns, "Tutorial Pattern")
		}
		if strings.Contains(titleLower, "how to") {
			patterns = append(patterns, "How-to Pattern")
		}
		if strings.Contains(titleLower, "2024") || strings.Contains(titleLower, "2023") {
			patterns = append(patterns, "Year Pattern")
		}
	}

	var commonWords []WordCount
	for word, count := range wordCounts {
		commonWords = append(commonWords, WordCount{Word: word, Count: count})
	}
	sort.Slice(commonWords, func(i, j int) bool { return commonWords[i].Count > commonWords[j].Count })
	if len(commonWords) > 10 {
		commonWords = commonWords[:10]
	}

	var commonPhrases []PhraseCount
	for phrase, count := range phraseCounts {
		commonPhrases = append(commonPhrases, PhraseCount{Phrase: phrase, Count: count})
	}
	sort.Slice(commonPhrases, func(i, j int) bool { return commonPhrases[i].Count > commonPhrases[j].Count })
	if len(commonPhrases) > 5 {
		commonPhrases = commonPhrases[:5]
	}

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
	channelVPDs := make(map[string][]float64) // To calculate average VPD per channel
	var totalViews int64
	var totalVPD float64

	for _, video := range a.videos {
		totalVPD += video.VPD

		if stats, exists := channelStats[video.Channel]; exists {
			stats.VideoCount++
			stats.TotalViews += video.Views
			stats.AvgViews = float64(stats.TotalViews) / float64(stats.VideoCount)
			engagement := float64(video.Likes+video.Comments) / max1(float64(video.Views)) * 100
			stats.Engagement += engagement // Accumulate for average calculation later
			channelVPDs[video.Channel] = append(channelVPDs[video.Channel], video.VPD)
		} else {
			engagement := float64(video.Likes+video.Comments) / max1(float64(video.Views)) * 100
			channelURL := fmt.Sprintf("https://www.youtube.com/channel/%s", video.ChannelID)
			channelStats[video.Channel] = &ChannelStats{
				Channel:    video.Channel,
				ChannelID:  video.ChannelID,
				ChannelURL: channelURL,
				VideoCount: 1,
				TotalViews: video.Views,
				AvgViews:   float64(video.Views),
				Engagement: engagement,
			}
			channelVPDs[video.Channel] = []float64{video.VPD}
		}
		totalViews += video.Views
	}

	// Calculate niche average VPD
	nicheAvgVPD := totalVPD / float64(len(a.videos))

	// Calculate average VPD for each channel and identify Rising Stars
	var allChannels []ChannelStats
	for channelName, stats := range channelStats {
		// Calculate average VPD for this channel
		var sumVPD float64
		for _, vpd := range channelVPDs[channelName] {
			sumVPD += vpd
		}
		stats.AvgVPD = sumVPD / float64(len(channelVPDs[channelName]))

		// Calculate average engagement for this channel
		stats.Engagement = stats.Engagement / float64(stats.VideoCount)

		// Check if it's a Rising Star based on configurable multiplier
		mult := a.cfg.RisingStarMultiplier
		if mult <= 0 {
			mult = 1.5
		}
		stats.IsRisingStar = stats.AvgVPD > nicheAvgVPD*mult

		allChannels = append(allChannels, *stats)
	}

	// Sort by total views for top channels
	sort.Slice(allChannels, func(i, j int) bool {
		return allChannels[i].TotalViews > allChannels[j].TotalViews
	})

	var topChannels []ChannelStats
	if len(allChannels) > 5 {
		topChannels = allChannels[:5]
	} else {
		topChannels = allChannels
	}

	// Extract Rising Stars
	var risingStars []ChannelStats
	for _, channel := range allChannels {
		if channel.IsRisingStar {
			risingStars = append(risingStars, channel)
		}
	}
	// Sort Rising Stars by VPD
	sort.Slice(risingStars, func(i, j int) bool {
		return risingStars[i].AvgVPD > risingStars[j].AvgVPD
	})

	// Calculate market share
	marketShare := make(map[string]float64)
	for _, channel := range topChannels {
		den := max1(float64(totalViews))
		marketShare[channel.Channel] = float64(channel.TotalViews) / den * 100
	}

	opportunities := a.generateCompetitorOpportunities(topChannels, marketShare)
	insights := a.generateCompetitorInsights(topChannels, marketShare)

	return CompetitorAnalysis{
		TopChannels:   topChannels,
		RisingStars:   risingStars,
		MarketShare:   marketShare,
		Opportunities: opportunities,
		Insights:      insights,
	}
}

func (a *Analyzer) AnalyzeTemporal() TemporalAnalysis {
	if len(a.videos) == 0 {
		return TemporalAnalysis{}
	}

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
		if n < minN {
			continue
		}
		s.AvgViews /= float64(n)
		s.AvgLikes /= float64(n)
		s.Engagement /= float64(n)
		bestHours = append(bestHours, *s)
	}
	sort.Slice(bestHours, func(i, j int) bool { return bestHours[i].Engagement > bestHours[j].Engagement })

	var bestDays []DayStats
	for d, s := range dayStats {
		n := dayN[d]
		if n < minN {
			continue
		}
		s.AvgViews /= float64(n)
		s.AvgLikes /= float64(n)
		s.Engagement /= float64(n)
		bestDays = append(bestDays, *s)
	}
	sort.Slice(bestDays, func(i, j int) bool { return bestDays[i].Engagement > bestDays[j].Engagement })

	insights := a.generateTemporalInsights(bestHours, bestDays)
	return TemporalAnalysis{BestHours: bestHours, BestDays: bestDays, Insights: insights}
}

func (a *Analyzer) AnalyzeKeywords() KeywordAnalysis {
	if len(a.videos) == 0 {
		return KeywordAnalysis{}
	}

	keywordStats := make(map[string]*KeywordStats)

	for _, video := range a.videos {
		words := utils.Tokenize(video.Title)
		engagement := float64(video.Likes+video.Comments) / max1(float64(video.Views)) * 100
		for _, w := range words {
			if stats, ok := keywordStats[w]; ok {
				stats.Frequency++
				stats.AvgViews += float64(video.Views)
				stats.AvgVPD += video.VPD
				stats.Engagement += engagement
			} else {
				keywordStats[w] = &KeywordStats{
					Keyword:    w,
					Frequency:  1,
					AvgViews:   float64(video.Views),
					AvgVPD:     video.VPD,
					Engagement: engagement,
				}
			}
		}
	}

	// Core Keywords: ranked by frequency (the old trending logic)
	var coreKeywords []KeywordStats
	for _, s := range keywordStats {
		coreKeywords = append(coreKeywords, *s)
	}
	sort.Slice(coreKeywords, func(i, j int) bool { return coreKeywords[i].Frequency > coreKeywords[j].Frequency })
	if len(coreKeywords) > 10 {
		coreKeywords = coreKeywords[:10]
	}

	// Trending Keywords: now ranked by Avg VPD (velocity-based)
	var trendingKeywords []KeywordStats
	for _, s := range keywordStats {
		trendingKeywords = append(trendingKeywords, *s)
	}

	var longTailKeywords []KeywordStats
	for _, s := range keywordStats {
		maxFreq := a.cfg.LongTailMaxFreq
		if maxFreq < 1 {
			maxFreq = 2
		}
		minEng := a.cfg.LongTailMinEngagement
		if minEng < 0 {
			minEng = 5.0
		}
		if s.Frequency <= maxFreq && (s.Engagement/float64(s.Frequency)) > minEng {
			longTailKeywords = append(longTailKeywords, *s)
		}
	}

	// finalize means
	for i := range coreKeywords {
		fk := coreKeywords[i].Frequency
		if fk > 0 {
			coreKeywords[i].AvgViews /= float64(fk)
			coreKeywords[i].AvgVPD /= float64(fk)
			coreKeywords[i].Engagement /= float64(fk)
		}
	}
	for i := range trendingKeywords {
		fk := trendingKeywords[i].Frequency
		if fk > 0 {
			trendingKeywords[i].AvgViews /= float64(fk)
			trendingKeywords[i].AvgVPD /= float64(fk)
			trendingKeywords[i].Engagement /= float64(fk)
		}
	}
	for i := range longTailKeywords {
		fk := longTailKeywords[i].Frequency
		if fk > 0 {
			longTailKeywords[i].AvgViews /= float64(fk)
			longTailKeywords[i].AvgVPD /= float64(fk)
			longTailKeywords[i].Engagement /= float64(fk)
		}
	}

	// Sort trending keywords by Average VPD (velocity)
	sort.Slice(trendingKeywords, func(i, j int) bool { return trendingKeywords[i].AvgVPD > trendingKeywords[j].AvgVPD })
	if len(trendingKeywords) > 10 {
		trendingKeywords = trendingKeywords[:10]
	}

	// Sort long tail by engagement as before
	sort.Slice(longTailKeywords, func(i, j int) bool { return longTailKeywords[i].Engagement > longTailKeywords[j].Engagement })

	seoOpportunities := a.generateSEOOpportunities(trendingKeywords, longTailKeywords)
	insights := a.generateKeywordInsights(trendingKeywords, longTailKeywords)

	return KeywordAnalysis{
		TrendingKeywords: trendingKeywords,
		CoreKeywords:     coreKeywords,
		LongTailKeywords: longTailKeywords,
		SEOOpportunities: seoOpportunities,
		Insights:         insights,
	}
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

// AnalyzeOpportunityScore computes a lightweight in-memory ranking combining
// velocity (VPD), engagement (like_rate), freshness, and a simple saturation penalty.
func (a *Analyzer) AnalyzeOpportunityScore() []OpportunityItem {
	if len(a.videos) == 0 {
		return []OpportunityItem{}
	}

	// Precompute features
	now := time.Now()
	vpdVals := make([]float64, 0, len(a.videos))
	likeRateVals := make([]float64, 0, len(a.videos))
	ageVals := make([]float64, 0, len(a.videos))

	// Simple saturation: frequency of primary tokens (from title)
	tokenFreq := make(map[string]int)
	videoPrimaryToken := make([]string, len(a.videos))
	for i, v := range a.videos {
		// VPD
		vpdVals = append(vpdVals, v.VPD)
		// Like rate (likes per 1k views)
		likeRate := 0.0
		if v.Views > 0 {
			likeRate = (float64(v.Likes) / float64(v.Views)) * 1000.0
		}
		likeRateVals = append(likeRateVals, likeRate)
		// Age in days
		ageDays := int(now.Sub(v.PublishedAt).Hours() / 24)
		if ageDays < 0 {
			ageDays = 0
		}
		ageVals = append(ageVals, float64(ageDays))

		// Primary token: first non-stopword token of title
		toks := utils.Tokenize(v.Title)
		primary := ""
		for _, t := range toks {
			if t == "" {
				continue
			}
			primary = t
			break
		}
		if primary == "" {
			primary = "_na_"
		}
		videoPrimaryToken[i] = primary
		tokenFreq[primary]++
	}

	// Helpers: z-score and min-max freshness
	z := func(vals []float64, x float64) float64 {
		if len(vals) == 0 {
			return 0
		}
		m := mean(vals)
		sd := stddev(vals, m)
		if sd == 0 {
			return 0
		}
		return (x - m) / sd
	}
	minMax := func(vals []float64, x float64, invert bool) float64 {
		if len(vals) == 0 {
			return 0
		}
		minV, maxV := vals[0], vals[0]
		for _, v := range vals {
			if v < minV {
				minV = v
			}
			if v > maxV {
				maxV = v
			}
		}
		if maxV == minV {
			return 0
		}
		norm := (x - minV) / (maxV - minV)
		if invert {
			return 1.0 - norm
		}
		return norm
	}

	items := make([]OpportunityItem, 0, len(a.videos))
	for i, v := range a.videos {
		likeRate := 0.0
		if v.Views > 0 {
			likeRate = (float64(v.Likes) / float64(v.Views)) * 1000.0
		}
		ageDays := int(now.Sub(v.PublishedAt).Hours() / 24)
		if ageDays < 0 {
			ageDays = 0
		}
		freshness := minMax(ageVals, float64(ageDays), true) // younger -> closer to 1
		// Saturation: normalized frequency of primary token
		freq := float64(tokenFreq[videoPrimaryToken[i]])
		// normalize by sample size
		saturation := freq / float64(len(a.videos))

		// Weights (could become config)
		wVPD := a.cfg.OppWeightVPD
		wLike := a.cfg.OppWeightLike
		wFresh := a.cfg.OppWeightFresh
		wSat := a.cfg.OppWeightSatPen // penalty weight

		score := wVPD*z(vpdVals, v.VPD) + wLike*z(likeRateVals, likeRate) + wFresh*freshness - wSat*saturation

		reasons := []string{
			fmt.Sprintf("VPD=%s", utils.FormatVPD(v.VPD)),
			fmt.Sprintf("LikeRate=%.2f/1k", likeRate),
			fmt.Sprintf("Age=%dd", ageDays),
			fmt.Sprintf("Saturation=%.2f", saturation),
		}

		items = append(items, OpportunityItem{
			Title:      v.Title,
			Channel:    v.Channel,
			URL:        v.URL,
			Score:      score,
			VPD:        v.VPD,
			LikeRate:   likeRate,
			AgeDays:    ageDays,
			Saturation: saturation,
			Reasons:    reasons,
		})
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Score > items[j].Score })
	if len(items) > 20 {
		items = items[:20]
	}
	return items
}

func mean(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

func stddev(vals []float64, m float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		dv := v - m
		sum += dv * dv
	}
	return math.Sqrt(sum / float64(len(vals)))
}

// Helper methods

func max1(v float64) float64 {
	if v <= 0 {
		return 1
	}
	return v
}

func (a *Analyzer) generateGrowthInsights(avgViews, avgLikes, nicheVelocityScore float64) []string {
	var insights []string

	if avgViews > 1000000 {
		insights = append(insights, "High-performing content with over 1M average views")
	} else if avgViews > 100000 {
		insights = append(insights, "Good performance with 100K+ average views")
	} else {
		insights = append(insights, "Room for improvement in view counts")
	}

	// Insights based on VPD (velocity)
	if nicheVelocityScore > 10000 {
		insights = append(insights, "Extremely high niche velocity - viral content potential")
	} else if nicheVelocityScore > 1000 {
		insights = append(insights, "Good niche velocity - positive momentum detected")
	} else if nicheVelocityScore > 100 {
		insights = append(insights, "Moderate velocity - growth opportunity")
	} else {
		insights = append(insights, "Low velocity - needs more push to gain traction")
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
	risingStarsInfo := ""
	if len(competitors.RisingStars) > 0 {
		risingStarsInfo = fmt.Sprintf(" %d rising star channel(s) detected with high velocity.", len(competitors.RisingStars))
	}

	if len(competitors.TopChannels) == 0 {
		return fmt.Sprintf("Analysis of %d videos shows Niche Velocity Score of %s VPD with %s average views.%s",
			growth.TotalVideos,
			utils.FormatVPD(growth.NicheVelocityScore),
			utils.FormatNumber(growth.AvgViews),
			risingStarsInfo)
	}
	return fmt.Sprintf("Analysis of %d videos shows Niche Velocity Score of %s VPD with %s average views. Top channel '%s' leads with %.1f%% market share.%s",
		growth.TotalVideos,
		utils.FormatVPD(growth.NicheVelocityScore),
		utils.FormatNumber(growth.AvgViews),
		competitors.TopChannels[0].Channel,
		competitors.MarketShare[competitors.TopChannels[0].Channel],
		risingStarsInfo)
}

func (a *Analyzer) generateKeyInsights(growth GrowthPattern, titles TitleAnalysis, competitors CompetitorAnalysis, temporal TemporalAnalysis, keywords KeywordAnalysis) []string {
	var insights []string

	insights = append(insights, fmt.Sprintf("Average views: %s", utils.FormatNumber(growth.AvgViews)))
	insights = append(insights, fmt.Sprintf("Niche Velocity Score: %s VPD", utils.FormatVPD(growth.NicheVelocityScore)))

	if len(keywords.TrendingKeywords) > 0 {
		insights = append(insights, fmt.Sprintf("Top trending keyword: '%s' (%s VPD)", keywords.TrendingKeywords[0].Keyword, utils.FormatVPD(keywords.TrendingKeywords[0].AvgVPD)))
	}

	if len(competitors.RisingStars) > 0 {
		insights = append(insights, fmt.Sprintf("Rising stars detected: %d channels", len(competitors.RisingStars)))
	}

	if len(temporal.BestHours) > 0 {
		insights = append(insights, fmt.Sprintf("Best posting time: %d:00", temporal.BestHours[0].Hour))
	}

	return insights
}

func (a *Analyzer) generateRecommendations(growth GrowthPattern, titles TitleAnalysis, competitors CompetitorAnalysis, temporal TemporalAnalysis, keywords KeywordAnalysis) []string {
	var recommendations []string

	if growth.NicheVelocityScore < 100 {
		recommendations = append(recommendations, "Focus on trending topics and viral content to improve velocity")
	}

	if len(keywords.TrendingKeywords) > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Target breakout keyword: '%s' (%s VPD)", keywords.TrendingKeywords[0].Keyword, utils.FormatVPD(keywords.TrendingKeywords[0].AvgVPD)))
	}

	if len(competitors.RisingStars) > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Study rising star channel '%s' for momentum strategies", competitors.RisingStars[0].Channel))
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
		intel = append(intel, fmt.Sprintf("Top competitor: %s (%s)", competitors.TopChannels[0].Channel, competitors.TopChannels[0].ChannelURL))
	}

	if len(competitors.RisingStars) > 0 {
		intel = append(intel, fmt.Sprintf("Rising stars detected: %d channels", len(competitors.RisingStars)))
		for i, star := range competitors.RisingStars {
			intel = append(intel, fmt.Sprintf("⭐ Rising Star #%d: %s (VPD: %s) - %s",
				i+1, star.Channel, utils.FormatVPD(star.AvgVPD), star.ChannelURL))
		}
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
