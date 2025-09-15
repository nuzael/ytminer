package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

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
		maxResults  = flag.Int("n", appConfig.MaxResults, "Maximum number of results")
		region      = flag.String("r", appConfig.DefaultRegion, "Search region (any, BR, US, GB)")
		duration    = flag.String("d", appConfig.DefaultDuration, "Video duration (any, short, medium, long)")
		analysis    = flag.String("a", "", "Analysis type (growth, titles, competitors, temporal, keywords, executive, all)")
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
		runCLIMode(*keyword, *maxResults, *region, *duration, *analysis)
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
	fmt.Println("  -n int       Maximum number of results (default: 25)")
	fmt.Println("  -r string    Search region: any, BR, US, GB (default: any)")
	fmt.Println("  -d string    Video duration: any, short, medium, long (default: any)")
	fmt.Println("  -a string    Analysis type: growth, titles, competitors, temporal, keywords, executive, all")
	fmt.Println("  -help        Show this help message")
	fmt.Println("  -version     Show version information")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  ytminer -k \"Python tutorial\" -n 50")
	fmt.Println("  ytminer -k \"Pokemon\" -r BR -d short -a growth")
	fmt.Println("  ytminer -k \"Machine Learning\" -a all")
	fmt.Println()
	fmt.Println("INTERACTIVE MODE:")
	fmt.Println("  ytminer")
	fmt.Println("    Run without parameters for interactive mode")
	fmt.Println()
}

func runCLIMode(keyword string, maxResults int, region, duration, analysis string) {
	// Create YouTube client
	client, err := utils.CreateYouTubeClient()
	if err != nil {
		return
	}

	// Search videos with loading
	fmt.Printf("Searching for: %s\n", keyword)
	
	searchOpts := youtube.SearchOptions{
		Query:      keyword,
		MaxResults: int64(maxResults),
		Region:     region,
		Duration:   duration,
		Order:      "relevance",
	}

	// Show loading while searching
	stopLoading := utils.ShowLoading("🔍 Searching YouTube videos...")
	
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
	var maxResults string
	var region string
	var duration string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("🔍 Search Keyword").
				Description("What would you like to search for?").
				Placeholder("e.g., Python tutorial").
				Value(&keyword),

			huh.NewInput().
				Title("📊 Max Results").
				Description("How many videos to analyze?").
				Placeholder("25").
				Value(&maxResults),

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
		),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	// Convert maxResults to int
	maxResultsInt, err := strconv.ParseInt(maxResults, 10, 64)
	if err != nil {
		maxResultsInt = 25
	}

	// Create YouTube client
	client, err := utils.CreateYouTubeClient()
	if err != nil {
		showMainMenu()
		return
	}

	// Search videos with loading
	ui.DisplayInfo("🔍 Searching for: " + keyword)
	
	searchOpts := youtube.SearchOptions{
		Query:      keyword,
		MaxResults: maxResultsInt,
		Region:     region,
		Duration:   duration,
		Order:      "relevance",
	}

	// Show loading while searching
	stopLoading := utils.ShowLoading("🔍 Searching YouTube videos...")
	
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
		),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	// Create YouTube client
	client, err := utils.CreateYouTubeClient()
	if err != nil {
		showMainMenu()
		return
	}

	// Search videos with loading
	ui.DisplayInfo("🔍 Searching for: " + keyword)
	
	searchOpts := youtube.SearchOptions{
		Query:      keyword,
		MaxResults: 25,
		Region:     "any",
		Duration:   "any",
		Order:      "relevance",
	}

	// Show loading while searching
	stopLoading := utils.ShowLoading("🔍 Searching YouTube videos...")
	
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
	var maxResults string
	var defaultRegion string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("🔑 YouTube API Key").
				Description("Enter your YouTube Data API v3 key").
				Placeholder("AIzaSy...").
				Value(&apiKey),

			huh.NewInput().
				Title("📊 Default Max Results").
				Description("Default number of videos to analyze").
				Placeholder("25").
				Value(&maxResults),

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
		
		if maxResults != "" {
			if parsed, err := strconv.Atoi(maxResults); err == nil {
				config.MaxResults = parsed
			}
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

