package ui

import (
	"fmt"
	"os"
	"strconv"

	"ytminer/analysis"
	"ytminer/utils"
	"ytminer/youtube"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/olekukonko/tablewriter"
)

// Colors and styles
var (
	primaryColor   = lipgloss.Color("#FF6B6B") // Red
	secondaryColor = lipgloss.Color("#4ECDC4") // Cyan
	accentColor    = lipgloss.Color("#45B7D1") // Blue
	textColor      = lipgloss.Color("#FFFFFF") // Orange
	successColor   = lipgloss.Color("#27AE60") // Green
	warningColor   = lipgloss.Color("#F39C12") // Yellow
	errorColor     = lipgloss.Color("#E74C3C") // Red
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Margin(1, 0).
			Align(lipgloss.Center)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true).
			Margin(0, 0, 1, 0).
			Align(lipgloss.Center)

	infoStyle = lipgloss.NewStyle().
			Foreground(textColor)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	headerStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	sectionStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Margin(1, 0, 0, 0)
)

func DisplayWelcome() {
	fmt.Println(titleStyle.Render("🚀 YTMiner"))
	fmt.Println(subtitleStyle.Render("Beautiful YouTube Analytics CLI"))
	fmt.Println()
}

func DisplayVideos(videos []youtube.Video) {
	if len(videos) == 0 {
		fmt.Println(warningStyle.Render("No videos found"))
		return
	}

	fmt.Println(headerStyle.Render("📺 Search Results"))
	fmt.Println()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "Title", "Channel", "Views", "Likes", "Published", "URL"})
	table.SetBorder(true)
	table.SetCenterSeparator("|")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")

	for i, video := range videos {
		// Truncate long titles
		title := video.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}

		// Truncate long channel names
		channel := video.Channel
		if len(channel) > 15 {
			channel = channel[:12] + "..."
		}

		// Format numbers
		views := utils.FormatNumber(video.Views)
		likes := utils.FormatNumber(video.Likes)

		// Format date
		published := video.PublishedAt.Format("2006-01-02")

		// Keep full URL
		url := video.URL

		table.Append([]string{
			strconv.Itoa(i + 1),
			title,
			channel,
			views,
			likes,
			published,
			url,
		})
	}

	table.Render()
	fmt.Println()
}

func DisplayGrowthAnalysis(growth analysis.GrowthPattern) {
	fmt.Println(sectionStyle.Render("📈 Growth Pattern Analysis"))
	fmt.Println()

	// Summary
	fmt.Println(infoStyle.Render(fmt.Sprintf("Total Videos: %d", growth.TotalVideos)))
	fmt.Println(infoStyle.Render(fmt.Sprintf("Average Views: %s", utils.FormatNumber(growth.AvgViews))))
	fmt.Println(infoStyle.Render(fmt.Sprintf("Average Likes: %s", utils.FormatNumber(growth.AvgLikes))))
	fmt.Println(infoStyle.Render(fmt.Sprintf("Growth Rate: %.1f%%", growth.GrowthRate)))
	fmt.Println()

	// Top Performers
	if len(growth.TopPerformers) > 0 {
		fmt.Println(headerStyle.Render("🏆 Top Performing Videos"))
		fmt.Println()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Title", "Channel", "Views", "Likes", "Engagement", "URL"})
		table.SetBorder(true)
		table.SetCenterSeparator("|")
		table.SetColumnSeparator("|")
		table.SetRowSeparator("-")

		for _, video := range growth.TopPerformers {
			title := video.Title
			if len(title) > 35 {
				title = title[:32] + "..."
			}

			channel := video.Channel
			if len(channel) > 12 {
				channel = channel[:9] + "..."
			}

			// Keep full URL
			url := video.URL

			table.Append([]string{
				title,
				channel,
				utils.FormatNumber(video.Views),
				utils.FormatNumber(video.Likes),
				fmt.Sprintf("%.2f%%", video.Engagement),
				url,
			})
		}

		table.Render()
		fmt.Println()
	}

	// Insights
	if len(growth.Insights) > 0 {
		fmt.Println(headerStyle.Render("💡 Insights"))
		for _, insight := range growth.Insights {
			fmt.Println(infoStyle.Render("• " + insight))
		}
		fmt.Println()
	}
}

func DisplayTitleAnalysis(titles analysis.TitleAnalysis) {
	fmt.Println(sectionStyle.Render("📝 Title Pattern Analysis"))
	fmt.Println()

	// Common Words
	if len(titles.CommonWords) > 0 {
		fmt.Println(headerStyle.Render("🔤 Most Common Words"))
		for i, word := range titles.CommonWords {
			if i >= 5 {
				break
			}
			fmt.Println(infoStyle.Render(fmt.Sprintf("%d. %s (%d times)", i+1, word.Word, word.Count)))
		}
		fmt.Println()
	}

	// Common Phrases
	if len(titles.CommonPhrases) > 0 {
		fmt.Println(headerStyle.Render("📄 Most Common Phrases"))
		for i, phrase := range titles.CommonPhrases {
			if i >= 3 {
				break
			}
			fmt.Println(infoStyle.Render(fmt.Sprintf("%d. %s (%d times)", i+1, phrase.Phrase, phrase.Count)))
		}
		fmt.Println()
	}

	// Emojis
	if len(titles.Emojis) > 0 {
		fmt.Println(headerStyle.Render("😀 Most Used Emojis"))
		for i, emoji := range titles.Emojis {
			if i >= 5 {
				break
			}
			fmt.Println(infoStyle.Render(fmt.Sprintf("%d. %s (%d times)", i+1, emoji.Emoji, emoji.Count)))
		}
		fmt.Println()
	}

	// Patterns
	if len(titles.Patterns) > 0 {
		fmt.Println(headerStyle.Render("🎯 Common Patterns"))
		for _, pattern := range titles.Patterns {
			fmt.Println(infoStyle.Render("• " + pattern))
		}
		fmt.Println()
	}

	// Insights
	if len(titles.Insights) > 0 {
		fmt.Println(headerStyle.Render("💡 Insights"))
		for _, insight := range titles.Insights {
			fmt.Println(infoStyle.Render("• " + insight))
		}
		fmt.Println()
	}
}

func DisplayCompetitorAnalysis(competitors analysis.CompetitorAnalysis) {
	fmt.Println(sectionStyle.Render("🏢 Competitor Analysis"))
	fmt.Println()

	// Top Channels
	if len(competitors.TopChannels) > 0 {
		fmt.Println(headerStyle.Render("📊 Top Performing Channels"))
		fmt.Println()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Channel", "Videos", "Total Views", "Avg Views", "Engagement"})
		table.SetBorder(true)
		table.SetCenterSeparator("|")
		table.SetColumnSeparator("|")
		table.SetRowSeparator("-")

		for _, channel := range competitors.TopChannels {
			channelName := channel.Channel
			if len(channelName) > 20 {
				channelName = channelName[:17] + "..."
			}

			table.Append([]string{
				channelName,
				strconv.Itoa(channel.VideoCount),
				utils.FormatNumber(channel.TotalViews),
				utils.FormatNumber(channel.AvgViews),
				fmt.Sprintf("%.2f%%", channel.Engagement),
			})
		}

		table.Render()
		fmt.Println()
	}

	// Market Share
	if len(competitors.MarketShare) > 0 {
		fmt.Println(headerStyle.Render("📈 Market Share"))
		for channel, share := range competitors.MarketShare {
			fmt.Println(infoStyle.Render(fmt.Sprintf("%s: %.1f%%", channel, share)))
		}
		fmt.Println()
	}

	// Opportunities
	if len(competitors.Opportunities) > 0 {
		fmt.Println(headerStyle.Render("🎯 Opportunities"))
		for _, opportunity := range competitors.Opportunities {
			fmt.Println(infoStyle.Render("• " + opportunity))
		}
		fmt.Println()
	}

	// Insights
	if len(competitors.Insights) > 0 {
		fmt.Println(headerStyle.Render("💡 Insights"))
		for _, insight := range competitors.Insights {
			fmt.Println(infoStyle.Render("• " + insight))
		}
		fmt.Println()
	}
}

func DisplayTemporalAnalysis(temporal analysis.TemporalAnalysis) {
	fmt.Println(sectionStyle.Render("⏰ Temporal Analysis"))
	fmt.Println()

	// Best Hours
	if len(temporal.BestHours) > 0 {
		fmt.Println(headerStyle.Render("🕐 Best Posting Hours"))
		fmt.Println()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Hour", "Avg Views", "Avg Likes", "Engagement"})
		table.SetBorder(true)
		table.SetCenterSeparator("|")
		table.SetColumnSeparator("|")
		table.SetRowSeparator("-")

		for _, hour := range temporal.BestHours {
			table.Append([]string{
				fmt.Sprintf("%d:00", hour.Hour),
				utils.FormatNumber(hour.AvgViews),
				utils.FormatNumber(hour.AvgLikes),
				fmt.Sprintf("%.2f%%", hour.Engagement),
			})
		}

		table.Render()
		fmt.Println()
	}

	// Best Days
	if len(temporal.BestDays) > 0 {
		fmt.Println(headerStyle.Render("📅 Best Posting Days"))
		fmt.Println()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Day", "Avg Views", "Avg Likes", "Engagement"})
		table.SetBorder(true)
		table.SetCenterSeparator("|")
		table.SetColumnSeparator("|")
		table.SetRowSeparator("-")

		for _, day := range temporal.BestDays {
			table.Append([]string{
				day.Day,
				utils.FormatNumber(day.AvgViews),
				utils.FormatNumber(day.AvgLikes),
				fmt.Sprintf("%.2f%%", day.Engagement),
			})
		}

		table.Render()
		fmt.Println()
	}

	// Insights
	if len(temporal.Insights) > 0 {
		fmt.Println(headerStyle.Render("💡 Insights"))
		for _, insight := range temporal.Insights {
			fmt.Println(infoStyle.Render("• " + insight))
		}
		fmt.Println()
	}
}

func DisplayKeywordAnalysis(keywords analysis.KeywordAnalysis) {
	fmt.Println(sectionStyle.Render("🔍 Keyword Analysis"))
	fmt.Println()

	// Trending Keywords
	if len(keywords.TrendingKeywords) > 0 {
		fmt.Println(headerStyle.Render("📈 Trending Keywords"))
		fmt.Println()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Keyword", "Frequency", "Avg Views", "Engagement"})
		table.SetBorder(true)
		table.SetCenterSeparator("|")
		table.SetColumnSeparator("|")
		table.SetRowSeparator("-")

		for _, keyword := range keywords.TrendingKeywords {
			table.Append([]string{
				keyword.Keyword,
				strconv.Itoa(keyword.Frequency),
				utils.FormatNumber(keyword.AvgViews),
				fmt.Sprintf("%.2f%%", keyword.Engagement),
			})
		}

		table.Render()
		fmt.Println()
	}

	// Long Tail Keywords
	if len(keywords.LongTailKeywords) > 0 {
		fmt.Println(headerStyle.Render("🎯 Long-tail Keywords"))
		fmt.Println()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Keyword", "Frequency", "Avg Views", "Engagement"})
		table.SetBorder(true)
		table.SetCenterSeparator("|")
		table.SetColumnSeparator("|")
		table.SetRowSeparator("-")

		for _, keyword := range keywords.LongTailKeywords {
			table.Append([]string{
				keyword.Keyword,
				strconv.Itoa(keyword.Frequency),
				utils.FormatNumber(keyword.AvgViews),
				fmt.Sprintf("%.2f%%", keyword.Engagement),
			})
		}

		table.Render()
		fmt.Println()
	}

	// SEO Opportunities
	if len(keywords.SEOOpportunities) > 0 {
		fmt.Println(headerStyle.Render("🚀 SEO Opportunities"))
		for _, opportunity := range keywords.SEOOpportunities {
			fmt.Println(infoStyle.Render("• " + opportunity))
		}
		fmt.Println()
	}

	// Insights
	if len(keywords.Insights) > 0 {
		fmt.Println(headerStyle.Render("💡 Insights"))
		for _, insight := range keywords.Insights {
			fmt.Println(infoStyle.Render("• " + insight))
		}
		fmt.Println()
	}
}

func DisplayExecutiveReport(report analysis.ExecutiveReport) {
	fmt.Println(sectionStyle.Render("💼 Executive Report"))
	fmt.Println()

	// Summary
	if report.Summary != "" {
		fmt.Println(headerStyle.Render("📋 Executive Summary"))
		fmt.Println(infoStyle.Render(report.Summary))
		fmt.Println()
	}

	// Key Insights
	if len(report.KeyInsights) > 0 {
		fmt.Println(headerStyle.Render("💡 Key Insights"))
		for _, insight := range report.KeyInsights {
			fmt.Println(infoStyle.Render("• " + insight))
		}
		fmt.Println()
	}

	// Recommendations
	if len(report.Recommendations) > 0 {
		fmt.Println(headerStyle.Render("🎯 Strategic Recommendations"))
		for _, rec := range report.Recommendations {
			fmt.Println(infoStyle.Render("• " + rec))
		}
		fmt.Println()
	}

	// Content Strategy
	if len(report.ContentStrategy) > 0 {
		fmt.Println(headerStyle.Render("📝 Content Strategy"))
		for _, strategy := range report.ContentStrategy {
			fmt.Println(infoStyle.Render("• " + strategy))
		}
		fmt.Println()
	}

	// Competitive Intelligence
	if len(report.CompetitiveIntel) > 0 {
		fmt.Println(headerStyle.Render("🏢 Competitive Intelligence"))
		for _, intel := range report.CompetitiveIntel {
			fmt.Println(infoStyle.Render("• " + intel))
		}
		fmt.Println()
	}

	// Performance Benchmarks
	if len(report.PerformanceBench) > 0 {
		fmt.Println(headerStyle.Render("📊 Performance Benchmarks"))
		for _, bench := range report.PerformanceBench {
			fmt.Println(infoStyle.Render("• " + bench))
		}
		fmt.Println()
	}

	// Next Steps
	if len(report.NextSteps) > 0 {
		fmt.Println(headerStyle.Render("🚀 Next Steps"))
		for _, step := range report.NextSteps {
			fmt.Println(infoStyle.Render("• " + step))
		}
		fmt.Println()
	}
}

func DisplayError(message string) {
	fmt.Println(errorStyle.Render("❌ Error: " + message))
	fmt.Println()
}

func DisplaySuccess(message string) {
	fmt.Println(successStyle.Render("✅ " + message))
	fmt.Println()
}

func DisplayWarning(message string) {
	fmt.Println(warningStyle.Render("⚠️ " + message))
	fmt.Println()
}

func DisplayInfo(message string) {
	fmt.Println(infoStyle.Render("ℹ️ " + message))
	fmt.Println()
}

func DisplayMarkdown(content string) {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	
	out, _ := r.Render(content)
	fmt.Print(out)
}


