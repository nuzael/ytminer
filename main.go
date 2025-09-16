package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"ytminer/analysis"
	"ytminer/config"
	"ytminer/ui"
	"ytminer/utils"
	"ytminer/youtube"

	"github.com/charmbracelet/huh"
	"github.com/joho/godotenv"
)


func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	appConfig := config.LoadConfig()

	// Parse command line flags
	var (
		keyword     = flag.String("k", "", "Search keyword")
		region      = flag.String("r", appConfig.DefaultRegion, "Search region (any, BR, US, GB)")
		duration    = flag.String("d", appConfig.DefaultDuration, "Video duration (any, short, medium, long)")
		analysis    = flag.String("a", "", "Analysis type (growth, titles, competitors, temporal, keywords, executive, all)")
		level       = flag.String("l", "balanced", "Analysis level (quick, balanced, deep)")
		timeRange   = flag.String("t", "any", "Time range (any, 7d, 30d, 90d, 1y)")
		help        = flag.Bool("help", false, "Show help")
		version     = flag.Bool("version", false, "Show version")
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

	// Show welcome message
	ui.DisplayWelcome()

	// If keyword is provided, run in CLI mode
	if *keyword != "" {
		runCLIMode(*keyword, *region, *duration, *analysis, *level, *timeRange)
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
	fmt.Println("  -t string    Time range: any, 7d, 30d, 90d, 1y (default: any)")
	fmt.Println("  -help        Show this help message")
	fmt.Println("  -version     Show version information")
	fmt.Println()
	fmt.Println("ANALYSIS LEVELS:")
	fmt.Println("  quick        Fast analysis (~200 units, 50 videos, 30-60s)")
	fmt.Println("  balanced     Balanced analysis (~1000 units, 200 videos, 2-3min)")
	fmt.Println("  deep         Deep analysis (~3000 units, 600 videos, 5-8min)")
	fmt.Println()
	fmt.Println("TIME RANGES:")
	fmt.Println("  any          No time filter (all videos)")
	fmt.Println("  7d           Last 7 days")
	fmt.Println("  30d          Last 30 days")
	fmt.Println("  90d          Last 90 days")
	fmt.Println("  1y           Last year")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  ytminer -k \"Python tutorial\" -l quick -t 7d")
	fmt.Println("  ytminer -k \"Pokemon\" -r BR -d short -a growth -l balanced -t 30d")
	fmt.Println("  ytminer -k \"Machine Learning\" -a all -l deep -t 90d")
	fmt.Println()
	fmt.Println("INTERACTIVE MODE:")
	fmt.Println("  ytminer")
	fmt.Println("    Run without parameters for interactive mode")
	fmt.Println()
}

func runCLIMode(keyword string, region, duration, analysis, level, timeRange string) {
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
	fmt.Printf("Searching for: %s (Level: %s, Time: %s)\n", keyword, level, timeRange)
	
	searchOpts := youtube.SearchOptions{
		Query:           keyword,
		MaxResults:      50, // Fixed at 50 per search (controlled by level)
		Region:          region,
		Duration:        duration,
		Order:           "relevance",
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
		return
	}

	fmt.Printf("Found %d videos!\n\n", len(videos))

	// Display results
	ui.DisplayVideos(videos)

	// If analysis is specified, run it
	if analysis != "" {
		runAnalysis(videos, analysis)
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
					huh.NewOption("📊 Run Analysis", "analysis"),
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
	case "analysis":
		showAnalysisForm()
	case "settings":
		showSettingsForm()
	case "exit":
		fmt.Println("👋 Thanks for using YTMiner!")
		os.Exit(0)
	}
}

func showSearchForm() {
	var keyword string
	var region string
	var duration string
	var level string
	var timeRange string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("🔍 Search Keyword").
				Description("What would you like to search for?").
				Placeholder("e.g., Python tutorial").
				Value(&keyword),

			huh.NewSelect[string]().
				Title("🌍 Region").
				Description("Search region").
				Options(
					huh.NewOption("Any", "any"),
					huh.NewOption("Brazil", "BR"),
					huh.NewOption("United States", "US"),
					huh.NewOption("United Kingdom", "GB"),
				).
				Value(&region),

			huh.NewSelect[string]().
				Title("⏱️ Duration").
				Description("Video duration filter").
				Options(
					huh.NewOption("Any", "any"),
					huh.NewOption("Short (< 4min)", "short"),
					huh.NewOption("Medium (4-20min)", "medium"),
					huh.NewOption("Long (> 20min)", "long"),
				).
				Value(&duration),

			huh.NewSelect[string]().
				Title("📅 Time Range").
				Description("When were the videos published?").
				Options(
					huh.NewOption("Any time", "any"),
					huh.NewOption("Last 7 days", "7d"),
					huh.NewOption("Last 30 days", "30d"),
					huh.NewOption("Last 90 days", "90d"),
					huh.NewOption("Last year", "1y"),
				).
				Value(&timeRange),

			huh.NewSelect[string]().
				Title("📊 Analysis Level").
				Description("Choose analysis depth").
				Options(
					huh.NewOption("🔍 Quick Scan (~200 units, 50 videos)", "quick"),
					huh.NewOption("⚖️ Balanced (~1000 units, 200 videos)", "balanced"),
					huh.NewOption("🚀 Deep Dive (~3000 units, 600 videos)", "deep"),
				).
				Value(&level),
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
	ui.DisplayInfo("🔍 Searching for: " + keyword + " (Level: " + level + ", Time: " + timeRange + ")")
	
	searchOpts := youtube.SearchOptions{
		Query:           keyword,
		MaxResults:      50, // Fixed at 50 per search (controlled by level)
		Region:          region,
		Duration:        duration,
		Order:           "relevance",
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

	// Display results
	ui.DisplayVideos(videos)

	// Ask if user wants to analyze
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
}

func showAnalysisForm() {
	var analysisType string
	var keyword string
	var level string
	var timeRange string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("🔍 Keyword").
				Description("What topic to analyze?").
				Placeholder("e.g., Python tutorial").
				Value(&keyword),

			huh.NewSelect[string]().
				Title("📊 Analysis Type").
				Description("Choose the type of analysis").
				Options(
					huh.NewOption("🎯 Creator Analysis", "creator"),
					huh.NewOption("📈 Marketing Analysis", "marketing"),
					huh.NewOption("🔬 Research Analysis", "research"),
					huh.NewOption("💼 Executive Report", "executive"),
					huh.NewOption("📋 Basic Analysis", "basic"),
				).
				Value(&analysisType),

			huh.NewSelect[string]().
				Title("📅 Time Range").
				Description("When were the videos published?").
				Options(
					huh.NewOption("Any time", "any"),
					huh.NewOption("Last 7 days", "7d"),
					huh.NewOption("Last 30 days", "30d"),
					huh.NewOption("Last 90 days", "90d"),
					huh.NewOption("Last year", "1y"),
				).
				Value(&timeRange),

			huh.NewSelect[string]().
				Title("📊 Analysis Level").
				Description("Choose analysis depth").
				Options(
					huh.NewOption("🔍 Quick Scan (~200 units, 50 videos)", "quick"),
					huh.NewOption("⚖️ Balanced (~1000 units, 200 videos)", "balanced"),
					huh.NewOption("🚀 Deep Dive (~3000 units, 600 videos)", "deep"),
				).
				Value(&level),
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
	ui.DisplayInfo("🔍 Searching for: " + keyword + " (Level: " + level + ", Time: " + timeRange + ")")
	
	searchOpts := youtube.SearchOptions{
		Query:           keyword,
		MaxResults:      50, // Fixed at 50 per search (controlled by level)
		Region:          "any",
		Duration:        "any",
		Order:           "relevance",
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

	// Run analysis
	runAnalysis(videos, "")
}

func runAnalysis(videos []youtube.Video, analysisType string) {
	if len(videos) == 0 {
		ui.DisplayWarning("No videos to analyze")
		showMainMenu()
		return
	}

	// Create analyzer
	analyzer := analysis.NewAnalyzer(videos)

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
	var apiKey string
	var defaultRegion string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("🔑 YouTube API Key").
				Description("Enter your YouTube Data API v3 key").
				Placeholder("AIzaSy...").
				Value(&apiKey),

			huh.NewSelect[string]().
				Title("🌍 Default Region").
				Description("Default search region").
				Options(
					huh.NewOption("Any", "any"),
					huh.NewOption("Brazil", "BR"),
					huh.NewOption("United States", "US"),
					huh.NewOption("United Kingdom", "GB"),
				).
				Value(&defaultRegion),
		),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	// Save settings to .env file
	if apiKey != "" {
		config := &config.AppConfig{
			APIKey: apiKey,
		}
		
		if defaultRegion != "" {
			config.DefaultRegion = defaultRegion
		}

		err := config.SaveConfig()
		if err != nil {
			utils.HandleError(err, "Failed to save settings")
		} else {
			ui.DisplaySuccess("Settings saved to .env file!")
		}
	} else {
		ui.DisplayWarning("No API key provided. Settings not saved.")
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
	case youtube.QuickScan:
		return "🔍 Quick scan - searching YouTube videos..."
	case youtube.Balanced:
		return "⚖️ Balanced analysis - searching multiple sources..."
	case youtube.DeepDive:
		return "🚀 Deep dive analysis - comprehensive search..."
	default:
		return "🔍 Searching YouTube videos..."
	}
}

func parseTimeRange(timeRange string) (string, string) {
	now := time.Now()
	
	switch timeRange {
	case "7d":
		sevenDaysAgo := now.AddDate(0, 0, -7)
		return sevenDaysAgo.Format(time.RFC3339), ""
	case "30d":
		thirtyDaysAgo := now.AddDate(0, 0, -30)
		return thirtyDaysAgo.Format(time.RFC3339), ""
	case "90d":
		ninetyDaysAgo := now.AddDate(0, 0, -90)
		return ninetyDaysAgo.Format(time.RFC3339), ""
	case "1y":
		oneYearAgo := now.AddDate(-1, 0, 0)
		return oneYearAgo.Format(time.RFC3339), ""
	default:
		return "", ""
	}
}

