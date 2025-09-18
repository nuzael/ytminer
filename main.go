package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"ytminer/analysis"
	"ytminer/config"
	"ytminer/ui"
	"ytminer/utils"
	"ytminer/youtube"

	"github.com/charmbracelet/huh"
	"github.com/joho/godotenv"
)

var globalAppConfig *config.AppConfig

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	appConfig := config.LoadConfig()
	globalAppConfig = appConfig

	// Parse command line flags
	var (
		keyword   = flag.String("k", "", "Search keyword")
		region    = flag.String("r", appConfig.DefaultRegion, "Search region (any, BR, US, GB)")
		duration  = flag.String("d", appConfig.DefaultDuration, "Video duration (any, short, medium, long)")
		analysis  = flag.String("a", "", "Analysis type (growth, titles, competitors, temporal, keywords, executive, all)")
		level     = flag.String("l", "balanced", "Analysis level (quick, balanced, deep)")
		timeRange = flag.String("t", appConfig.DefaultTimeRange, "Time range (any, 1h, 24h, 7d, 30d, 90d, 180d, 1y)")
		order     = flag.String("o", appConfig.DefaultOrder, "Search order (relevance, date, viewCount, rating, title)")
		noPreview = flag.Bool("no-preview", false, "Skip preview table and run analysis directly")
		help      = flag.Bool("help", false, "Show help")
		version   = flag.Bool("version", false, "Show version")
		profiles  = flag.Bool("profiles", false, "Show available weight profiles and exit")
		profile   = flag.String("profile", "", "Apply predefined weight profile (exploration, evergreen, trending, balanced)")
	)
	flag.Parse()

	// Show version
	if *version {
		fmt.Println("YTMiner v1.0.0")
		return
	}

	// Show help
	if *help {
		showHelp()
		return
	}

	// Show profiles and exit
	if *profiles {
		config.DisplayProfiles()
		return
	}

	// Apply profile if provided
	if *profile != "" {
		if err := globalAppConfig.ApplyProfile(*profile); err != nil {
			ui.DisplayError(fmt.Sprintf("failed to apply profile '%s': %v", *profile, err))
			return
		}
		ui.DisplayInfo(fmt.Sprintf("Applied profile: %s", *profile))
	}

	// Show welcome message
	ui.DisplayWelcome()

	// If keyword is provided, run in CLI mode
	if *keyword != "" {
		runCLIMode(*keyword, *region, *duration, *analysis, *level, *timeRange, *order, *noPreview)
		return
	}

	// Otherwise, run in interactive mode
	showMainMenu()
}

func showHelp() {
	fmt.Println("YTMiner - YouTube Analytics CLI")
	fmt.Println("=================================")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  ytminer [OPTIONS]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -k string    Search keyword (required for CLI mode)")
	fmt.Println("  -r string    Search region: any, BR, US, GB (default: any)")
	fmt.Println("  -d string    Video duration: any, short, medium, long (default: any)")
	fmt.Println("  -a string    Analysis type: growth, titles, competitors, temporal, keywords, executive, all")
	fmt.Println("  -l string    Analysis level: quick, balanced, deep (default: balanced)")
	fmt.Println("  -t string    Time range: any, 1h, 24h, 7d, 30d, 90d, 180d, 1y (default: any)")
	fmt.Println("  -o string    Search order: relevance, date, viewCount, rating, title (default: relevance)")
	fmt.Println("  --no-preview Skip preview table and run analysis directly")
	fmt.Println("  --help       Show help")
	fmt.Println("  --version    Show version")
	fmt.Println("  --profiles   Show available weight profiles and exit")
	fmt.Println("  --profile    Apply a predefined profile: exploration, evergreen, trending, balanced")
	fmt.Println()
	fmt.Println("ANALYSIS LEVELS:")
	fmt.Println("  quick        Fast analysis (~300 units, up to ~150 videos, 30-60s)")
	fmt.Println("  balanced     Balanced analysis (~800 units, up to ~400 videos, 1-2min)")
	fmt.Println("  deep         Deep analysis (~2000 units, up to ~1000 videos, 3-5min)")
	fmt.Println()
	fmt.Println("TIME RANGES:")
	fmt.Println("  any          No time filter (all videos)")
	fmt.Println("  1h           Last 1 hour")
	fmt.Println("  24h          Last 24 hours")
	fmt.Println("  7d           Last 7 days")
	fmt.Println("  30d          Last 30 days")
	fmt.Println("  90d          Last 90 days")
	fmt.Println("  180d         Last 180 days")
	fmt.Println("  1y           Last year")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  ytminer -k \"Python tutorial\" -l quick -t 7d")
	fmt.Println("  ytminer -k \"Pokemon\" -r BR -d short -a growth -l balanced -t 30d")
	fmt.Println("  ytminer -k \"Machine Learning\" -a executive -l deep -t 90d")
	fmt.Println("  ytminer -k \"AI tools\" --profile trending -a opportunity")
	fmt.Println("  ytminer --profiles")
	fmt.Println()
	fmt.Println("INTERACTIVE MODE:")
	fmt.Println("  ytminer")
	fmt.Println("    Run without parameters for interactive mode")
	fmt.Println()
}

func runCLIMode(keyword string, region, duration, analysis, level, timeRange, order string, noPreview bool) {
	// Create YouTube client
	client, err := utils.CreateYouTubeClient()
	if err != nil {
		return
	}

	// Parse analysis level
	analysisLevel := parseAnalysisLevel(level)

	// Parse time range
	publishedAfter, publishedBefore := parseTimeRange(timeRange)

	// Search videos with loading
	fmt.Printf("Searching for: %s (Level: %s, Time: %s, Region: %s, Duration: %s, Order: %s)\n", keyword, level, timeRange, region, duration, order)

	scrollOpts := youtube.SearchOptions{
		Query:           keyword,
		MaxResults:      50, // Fixed at 50 per search (controlled by level)
		Region:          region,
		Duration:        duration,
		Order:           order,
		Level:           analysisLevel,
		PublishedAfter:  publishedAfter,
		PublishedBefore: publishedBefore,
	}

	// Show loading while searching
	loadingMessage := getLoadingMessage(analysisLevel)
	stopLoading := utils.ShowLoading(loadingMessage)

	videos, err := client.SearchVideos(scrollOpts)
	stopLoading()

	if err != nil {
		utils.HandleError(err, "Failed to search videos")
		return
	}

	// CLI behavior: optionally preview, then analyze
	if analysis != "" || noPreview {
		if !noPreview {
			ui.DisplayVideos(videos)
		}
		runAnalysis(videos, analysis)
		return
	}

	// Default: just preview results; user can decide next steps interactivos
	// Ensure preview respects requested order
	sortVideosByOrder(videos, order)
	ui.DisplayVideos(videos)
}

func sortVideosByOrder(videos []youtube.Video, order string) {
	switch order {
	case "viewCount":
		sort.Slice(videos, func(i, j int) bool { return videos[i].Views > videos[j].Views })
	case "date":
		sort.Slice(videos, func(i, j int) bool { return videos[i].PublishedAt.After(videos[j].PublishedAt) })
	case "title":
		sort.Slice(videos, func(i, j int) bool { return videos[i].Title < videos[j].Title })
	default:
		// relevance/rating: keep API order
	}
}

func showMainMenu() {
	var choice string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do?").
				Description("Choose an option to get started").
				Options(
					huh.NewOption("🔍 Search Videos", "search"),
					huh.NewOption("⚙️ Settings", "settings"),
					huh.NewOption("❌ Exit", "exit"),
				).
				Value(&choice),
		),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	switch choice {
	case "search":
		showSearchForm()
	case "settings":
		showSettingsForm()
	case "exit":
		fmt.Println("👋 Thanks for using YTMiner!")
		os.Exit(0)
	}
}

func showSearchForm() {
	var keyword string
	var region string = globalAppConfig.DefaultRegion
	var duration string = globalAppConfig.DefaultDuration
	var level string = "balanced"
	var timeRange string = globalAppConfig.DefaultTimeRange
	var order string = globalAppConfig.DefaultOrder
	var previewBefore bool
	var useTranscripts bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("🔍 Search Keyword").
				Description("What would you like to search for?").
				Placeholder("e.g., Python tutorial").
				Value(&keyword),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🌍 Region").
				Description("Search region").
				Value(&region).
				Options(
					huh.NewOption("Any", "any").Selected(region == "any"),
					huh.NewOption("Brazil", "BR").Selected(region == "BR"),
					huh.NewOption("United States", "US").Selected(region == "US"),
					huh.NewOption("United Kingdom", "GB").Selected(region == "GB"),
				),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("⏱️ Duration").
				Description("Video duration filter").
				Value(&duration).
				Options(
					huh.NewOption("Any", "any").Selected(duration == "any"),
					huh.NewOption("Short (< 4min)", "short").Selected(duration == "short"),
					huh.NewOption("Medium (4-20min)", "medium").Selected(duration == "medium"),
					huh.NewOption("Long (> 20min)", "long").Selected(duration == "long"),
				),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("📅 Time Range").
				Description("When were the videos published?").
				Value(&timeRange).
				Options(
					huh.NewOption("Any time", "any").Selected(timeRange == "any"),
					huh.NewOption("Last 1 hour", "1h").Selected(timeRange == "1h"),
					huh.NewOption("Last 24 hours", "24h").Selected(timeRange == "24h"),
					huh.NewOption("Last 7 days", "7d").Selected(timeRange == "7d"),
					huh.NewOption("Last 30 days", "30d").Selected(timeRange == "30d"),
					huh.NewOption("Last 90 days", "90d").Selected(timeRange == "90d"),
					huh.NewOption("Last 180 days", "180d").Selected(timeRange == "180d"),
					huh.NewOption("Last year", "1y").Selected(timeRange == "1y"),
				),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🗂️ Order").
				Description("Search results order").
				Value(&order).
				Options(
					huh.NewOption("Relevance", "relevance").Selected(order == "relevance"),
					huh.NewOption("Most Recent", "date").Selected(order == "date"),
					huh.NewOption("Most Viewed", "viewCount").Selected(order == "viewCount"),
					huh.NewOption("Rating", "rating").Selected(order == "rating"),
					huh.NewOption("Title (A-Z)", "title").Selected(order == "title"),
				),
		),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("📊 Analysis Level").
				Description("Choose analysis depth").
				Value(&level).
				Options(
					huh.NewOption("🔍 Quick Scan (~300 units, up to ~150 videos)", "quick").Selected(level == "quick"),
					huh.NewOption("⚖️ Balanced (~800 units, up to ~400 videos)", "balanced").Selected(level == "balanced"),
					huh.NewOption("🚀 Deep Dive (~2000 units, up to ~1000 videos)", "deep").Selected(level == "deep"),
				),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("📝 Fetch transcripts for topic insights?").
				Description("No extra YouTube API quota; may increase runtime.").
				Value(&useTranscripts),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("👀 Preview results before analysis?").
				Description("If No, we'll run the analysis directly after the search").
				Value(&previewBefore),
		),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	// Parse analysis level
	analysisLevel := parseAnalysisLevel(level)

	// Parse time range
	publishedAfter, publishedBefore := parseTimeRange(timeRange)

	// Create YouTube client
	client, err := utils.CreateYouTubeClient()
	if err != nil {
		showMainMenu()
		return
	}

	// Search videos with loading
	ui.DisplayInfo("🔍 Searching for: " + keyword + " (Level: " + level + ", Time: " + timeRange + ", Region: " + region + ", Duration: " + duration + ", Order: " + order + ")")

	searchOpts := youtube.SearchOptions{
		Query:           keyword,
		MaxResults:      50, // Fixed at 50 per search (controlled by level)
		Region:          region,
		Duration:        duration,
		Order:           order,
		Level:           analysisLevel,
		PublishedAfter:  publishedAfter,
		PublishedBefore: publishedBefore,
	}

	// Show loading while searching
	loadingMessage := getLoadingMessage(analysisLevel)
	stopLoading := utils.ShowLoading(loadingMessage)

	videos, err := client.SearchVideos(searchOpts)
	stopLoading()

	if err != nil {
		utils.HandleError(err, "Failed to search videos")
		showMainMenu()
		return
	}

	// Optional transcript fetching (interactive choice overrides env)
	if useTranscripts || os.Getenv("YTMINER_WITH_TRANSCRIPTS") == "true" {
		ui.DisplayInfo("Attempting to fetch transcripts...")
		ui.DisplayWarning("Note: Transcript fetching may be limited due to YouTube restrictions")
		ui.DisplayInfo("If transcripts fail, analysis will continue without them")

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		fetched := 0
		for i := range videos {
			tr, err := client.GetTranscript(ctx, videos[i].ID)
			if err != nil {
				log.Printf("transcript: %s: %v", videos[i].ID, err)
				continue
			}
			videos[i].Transcript = tr.Text
			videos[i].TranscriptLang = tr.Language
			fetched++
		}
		log.Printf("transcripts fetched: %d/%d", fetched, len(videos))
	}

	// Sort videos by order
	sortVideosByOrder(videos, order)

	if previewBefore {
		// Show preview then ask to run analysis
		ui.DisplayVideos(videos)
		var analyze bool
		analyzeForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("📊 Run Analysis?").
					Description("Would you like to analyze these videos?").
					Value(&analyze),
			),
		)
		if err := analyzeForm.Run(); err != nil {
			log.Fatal(err)
		}
		if analyze {
			runAnalysis(videos, "")
		} else {
			showMainMenu()
		}
		return
	}

	// No preview: run analysis directly (type will be asked inside runAnalysis)
	runAnalysis(videos, "")
}

func runAnalysis(videos []youtube.Video, analysisType string) {
	if len(videos) == 0 {
		ui.DisplayWarning("No videos to analyze")
		showMainMenu()
		return
	}

	// Create analyzer
	analyzer := analysis.NewAnalyzer(videos, globalAppConfig)

	// If no analysis type specified, ask user
	if analysisType == "" {
		var selectedType string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("📊 Choose Analysis Type").
					Description("What type of analysis would you like to run?").
					Options(
						huh.NewOption("📈 Growth Pattern Analysis", "growth"),
						huh.NewOption("📝 Title Pattern Analysis", "titles"),
						huh.NewOption("🏢 Competitor Analysis", "competitors"),
						huh.NewOption("⏰ Temporal Analysis", "temporal"),
						huh.NewOption("🔍 Keyword Analysis", "keywords"),
						huh.NewOption("🎯 Opportunity Score", "opportunity"),
						huh.NewOption("💼 Executive Report", "executive"),
						huh.NewOption("📋 All Analyses", "all"),
					).
					Value(&selectedType),
			),
		)

		if err := form.Run(); err != nil {
			log.Fatal(err)
		}
		analysisType = selectedType
	}

	// Run selected analysis with loading
	switch analysisType {
	case "growth":
		stopLoading := utils.ShowLoading("📈 Analyzing growth patterns...")
		growth := analyzer.AnalyzeGrowthPatterns()
		stopLoading()
		ui.DisplayGrowthAnalysis(growth)
	case "titles":
		stopLoading := utils.ShowLoading("📝 Analyzing title patterns...")
		titles := analyzer.AnalyzeTitles()
		stopLoading()
		ui.DisplayTitleAnalysis(titles)
	case "competitors":
		stopLoading := utils.ShowLoading("🏢 Analyzing competitors...")
		competitors := analyzer.AnalyzeCompetitors()
		stopLoading()
		ui.DisplayCompetitorAnalysis(competitors)
	case "temporal":
		stopLoading := utils.ShowLoading("⏰ Analyzing temporal patterns...")
		temporal := analyzer.AnalyzeTemporal()
		stopLoading()
		ui.DisplayTemporalAnalysis(temporal)
	case "keywords":
		stopLoading := utils.ShowLoading("🔍 Analyzing keywords...")
		keywords := analyzer.AnalyzeKeywords()
		stopLoading()
		ui.DisplayKeywordAnalysis(keywords)
	case "opportunity":
		stopLoading := utils.ShowLoading("🎯 Computing opportunity score...")
		items := analyzer.AnalyzeOpportunityScore()
		stopLoading()
		ui.DisplayOpportunityScore(items)
	case "executive":
		stopLoading := utils.ShowLoading("💼 Generating executive report...")
		report := analyzer.GenerateExecutiveReport()
		stopLoading()
		ui.DisplayExecutiveReport(report)
	case "all":
		// Run all analyses with loading
		ui.DisplayInfo("Running comprehensive analysis...")

		stopLoading := utils.ShowLoading("📈 Analyzing growth patterns...")
		growth := analyzer.AnalyzeGrowthPatterns()
		stopLoading()
		ui.DisplayGrowthAnalysis(growth)

		stopLoading = utils.ShowLoading("📝 Analyzing title patterns...")
		titles := analyzer.AnalyzeTitles()
		stopLoading()
		ui.DisplayTitleAnalysis(titles)

		stopLoading = utils.ShowLoading("🏢 Analyzing competitors...")
		competitors := analyzer.AnalyzeCompetitors()
		stopLoading()
		ui.DisplayCompetitorAnalysis(competitors)

		stopLoading = utils.ShowLoading("⏰ Analyzing temporal patterns...")
		temporal := analyzer.AnalyzeTemporal()
		stopLoading()
		ui.DisplayTemporalAnalysis(temporal)

		stopLoading = utils.ShowLoading("🔍 Analyzing keywords...")
		keywords := analyzer.AnalyzeKeywords()
		stopLoading()
		ui.DisplayKeywordAnalysis(keywords)

		stopLoading = utils.ShowLoading("🎯 Computing opportunity score...")
		items := analyzer.AnalyzeOpportunityScore()
		stopLoading()
		ui.DisplayOpportunityScore(items)

		stopLoading = utils.ShowLoading("💼 Generating executive report...")
		report := analyzer.GenerateExecutiveReport()
		stopLoading()
		ui.DisplayExecutiveReport(report)
	default:
		ui.DisplayError("Unknown analysis type: " + analysisType)
	}

	// Ask if user wants to run another analysis
	var another bool
	anotherForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("🔄 Run Another Analysis?").
				Description("Would you like to run another analysis?").
				Value(&another),
		),
	)

	if err := anotherForm.Run(); err != nil {
		log.Fatal(err)
	}

	if another {
		runAnalysis(videos, "")
	} else {
		showMainMenu()
	}
}

func showSettingsForm() {
	// Load current config to prefill defaults
	current := config.LoadConfig()

	var apiKey string = current.APIKey
	var defaultRegion string = current.DefaultRegion
	var defaultDuration string = current.DefaultDuration
	var defaultTimeRange string = current.DefaultTimeRange
	var defaultOrder string = current.DefaultOrder
	var risingStarMultiplier string = strconv.FormatFloat(current.RisingStarMultiplier, 'f', -1, 64)
	var longTailMinEngagement string = strconv.FormatFloat(current.LongTailMinEngagement, 'f', -1, 64)
	var longTailMaxFreq string = strconv.Itoa(current.LongTailMaxFreq)

	form := huh.NewForm(
		// Step 1: API Key
		huh.NewGroup(
			huh.NewInput().
				Title("🔑 YouTube API Key").
				Description("Enter your YouTube Data API v3 key").
				Placeholder("AIzaSy...").
				Value(&apiKey),
		),

		// Step 2: Default Region
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🌍 Default Region").
				Description("Default search region").
				Value(&defaultRegion).
				Options(
					huh.NewOption("Any", "any").Selected(defaultRegion == "any"),
					huh.NewOption("Brazil", "BR").Selected(defaultRegion == "BR"),
					huh.NewOption("United States", "US").Selected(defaultRegion == "US"),
					huh.NewOption("United Kingdom", "GB").Selected(defaultRegion == "GB"),
				),
		),

		// Step 3: Default Duration
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("⏱️ Default Duration").
				Description("Default video duration filter").
				Value(&defaultDuration).
				Options(
					huh.NewOption("Any", "any").Selected(defaultDuration == "any"),
					huh.NewOption("Short (< 4min)", "short").Selected(defaultDuration == "short"),
					huh.NewOption("Medium (4-20min)", "medium").Selected(defaultDuration == "medium"),
					huh.NewOption("Long (> 20min)", "long").Selected(defaultDuration == "long"),
				),
		),

		// Step 4: Default Time Range
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("📅 Default Time Range").
				Description("Default published time filter").
				Value(&defaultTimeRange).
				Options(
					huh.NewOption("Any time", "any").Selected(defaultTimeRange == "any"),
					huh.NewOption("Last 1 hour", "1h").Selected(defaultTimeRange == "1h"),
					huh.NewOption("Last 24 hours", "24h").Selected(defaultTimeRange == "24h"),
					huh.NewOption("Last 7 days", "7d").Selected(defaultTimeRange == "7d"),
					huh.NewOption("Last 30 days", "30d").Selected(defaultTimeRange == "30d"),
					huh.NewOption("Last 90 days", "90d").Selected(defaultTimeRange == "90d"),
					huh.NewOption("Last 180 days", "180d").Selected(defaultTimeRange == "180d"),
					huh.NewOption("Last year", "1y").Selected(defaultTimeRange == "1y"),
				),
		),

		// Step 5: Default Order
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🗂️ Default Order").
				Description("Default search results order").
				Value(&defaultOrder).
				Options(
					huh.NewOption("Relevance", "relevance").Selected(defaultOrder == "relevance"),
					huh.NewOption("Most Recent", "date").Selected(defaultOrder == "date"),
					huh.NewOption("Most Viewed", "viewCount").Selected(defaultOrder == "viewCount"),
					huh.NewOption("Rating", "rating").Selected(defaultOrder == "rating"),
					huh.NewOption("Title (A-Z)", "title").Selected(defaultOrder == "title"),
				),
		),

		// Step 6: Velocity/Keyword Thresholds
		huh.NewGroup(
			huh.NewInput().
				Title("🌟 Rising Star Multiplier").
				Description("Channel AvgVPD > (multiplier × niche AvgVPD). Default: 1.5").
				Placeholder("1.5").
				Value(&risingStarMultiplier),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("🧵 Long Tail Min Engagement (%)").
				Description("Frequency <= max_freq AND Avg Engagement > min_engagement. Default: 5.0").
				Placeholder("5.0").
				Value(&longTailMinEngagement),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("🧵 Long Tail Max Frequency").
				Description("Max keyword frequency to consider as long tail. Default: 2").
				Placeholder("2").
				Value(&longTailMaxFreq),
		),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	// Build config with possibly updated values
	cfg := &config.AppConfig{
		APIKey:                apiKey,
		DefaultRegion:         defaultRegion,
		DefaultDuration:       defaultDuration,
		DefaultTimeRange:      defaultTimeRange,
		DefaultOrder:          defaultOrder,
		RisingStarMultiplier:  current.RisingStarMultiplier,
		LongTailMinEngagement: current.LongTailMinEngagement,
		LongTailMaxFreq:       current.LongTailMaxFreq,
	}

	// Parse thresholds (keep current on parse error)
	if f, err := strconv.ParseFloat(risingStarMultiplier, 64); err == nil && f > 0 {
		cfg.RisingStarMultiplier = f
	}
	if f, err := strconv.ParseFloat(longTailMinEngagement, 64); err == nil && f >= 0 {
		cfg.LongTailMinEngagement = f
	}
	if n, err := strconv.Atoi(longTailMaxFreq); err == nil && n >= 1 {
		cfg.LongTailMaxFreq = n
	}

	// Save settings to .env file (do not require API key to be changed)
	if err := cfg.SaveConfig(); err != nil {
		utils.HandleError(err, "Failed to save settings")
	} else {
		// Update in-memory defaults for this session
		globalAppConfig = cfg
		ui.DisplaySuccess("Settings saved to .env file!")
	}

	fmt.Println()

	// Return to main menu
	showMainMenu()
}

// Helper functions for analysis levels
func parseAnalysisLevel(level string) youtube.AnalysisLevel {
	switch level {
	case "quick":
		return youtube.QuickScan
	case "balanced":
		return youtube.Balanced
	case "deep":
		return youtube.DeepDive
	default:
		return youtube.Balanced
	}
}

func getLoadingMessage(level youtube.AnalysisLevel) string {
	switch level {
	default:
		return "Searching..."
	}
}

func parseTimeRange(timeRange string) (string, string) {
	now := time.Now()

	switch timeRange {
	case "1h":
		oneHourAgo := now.Add(time.Hour * -1)
		return oneHourAgo.Format(time.RFC3339), ""
	case "24h":
		oneDayAgo := now.AddDate(0, 0, -1)
		return oneDayAgo.Format(time.RFC3339), ""
	case "7d":
		oneWeekAgo := now.AddDate(0, 0, -7)
		return oneWeekAgo.Format(time.RFC3339), ""
	case "30d":
		oneMonthAgo := now.AddDate(0, -1, 0)
		return oneMonthAgo.Format(time.RFC3339), ""
	case "90d":
		threeMonthsAgo := now.AddDate(0, -3, 0)
		return threeMonthsAgo.Format(time.RFC3339), ""
	case "180d":
		sixMonthsAgo := now.AddDate(0, -6, 0)
		return sixMonthsAgo.Format(time.RFC3339), ""
	case "1y":
		oneYearAgo := now.AddDate(-1, 0, 0)
		return oneYearAgo.Format(time.RFC3339), ""
	default:
		return "", ""
	}
}
