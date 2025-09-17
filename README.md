# YTMiner - YouTube Analytics CLI

A command-line tool for YouTube content creators, marketers, and researchers to analyze video performance, discover trends, and optimize content strategy.

[![Go](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![YouTube API](https://img.shields.io/badge/YouTube-API%20v3-red.svg)](https://developers.google.com/youtube/v3)

## Core Analysis Features

### Growth Pattern Analysis (`-a growth`)
- **Niche Velocity Score (Avg. VPD)**: Measures overall niche momentum based on Views Per Day.
- **Highest Velocity Videos**: Identifies videos that are trending now, ranked by VPD.
- Performance metrics based on velocity, not just total views.
- Momentum analysis and viral potential detection.

### Title Pattern Analysis (`-a titles`) 
- Winning title formulas and patterns.
- Most common words in successful videos.
- Emoji usage patterns and effectiveness.
- Title length optimization insights.

### Competitor Analysis (`-a competitors`)
- **Rising Stars Detection**: Identifies emerging channels with high velocity (VPD > 1.5x niche average).
- **Direct Channel Links**: Provides clickable YouTube channel URLs for immediate analysis.
- Analysis of established competitors vs rising channels.
- Velocity metrics (Avg. VPD) per channel to identify momentum.
- Market share combined with velocity analysis.

### Temporal Analysis (`-a temporal`)
- Optimal posting times (best hours and days).
- Performance by hour and day of week.
- Peak engagement periods.
- Posting schedule optimization.

### Keyword Analysis (`-a keywords`)
- **Trending Keywords (Breakout Topics)**: Ranked by Avg. VPD to identify "hot" topics with momentum.
- **Core Keywords**: Most common words based on frequency (traditional analysis).
- Long-tail keywords with high engagement for SEO opportunities.
- Clear separation between viral vs popular topics.

### Executive Reports (`-a executive`)
- Comprehensive reports focusing on market velocity and momentum.
- Strategic recommendations based on VPD metrics.
- Identifies "Rising Stars" and breakout topics for monitoring.
- Performance benchmarks including Niche Velocity Score.

---

## Quick Start

### Interactive Mode (Recommended for beginners)
```bash
# Windows (PowerShell/CMD)
.\ytminer.exe

# Linux/Mac
./ytminer.exe
```
The interactive mode guides you through the entire process with helpful prompts and recommendations. Perfect for users who want a guided experience without memorizing command-line flags.

Defaults used in interactive search (configurable in Settings or .env):
- Region (`YTMINER_DEFAULT_REGION`)
- Duration (`YTMINER_DEFAULT_DURATION`)
- Time Range (`YTMINER_DEFAULT_TIME_RANGE`)
- Order (`YTMINER_DEFAULT_ORDER`)

### Analysis Levels
YTMiner offers three analysis levels to optimize API usage and analysis depth:

#### Quick Scan (~200 units, 50 videos, 30-60s)
- **Best for**: Quick exploration, demos, trend validation
- **Data**: Single search with 50 videos
- **Insights**: Basic patterns and velocity indicators
- **Use case**: "Is this topic hot right now?"

#### Balanced (~1000 units, 200 videos, 2-3min)
- **Best for**: Regular content creators, marketers
- **Data**: 4 searches across different parameters
- **Insights**: Reliable patterns and recommendations
- **Use case**: "What is the current growth strategy in this niche?"

#### Deep Dive (~3000 units, 600 videos, 5-8min)
- **Best for**: Strategic analysis, research, business decisions
- **Data**: 12 searches across regions, durations, and related topics
- **Insights**: Comprehensive market and competitor velocity intelligence
- **Use case**: "What's the complete competitive landscape, including emerging players?"

### 1. Get YouTube API Key
1. Visit [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project
3. Enable **YouTube Data API v3**
4. Create credentials (API Key)
5. Copy your API key

### 2. Download & Setup
```bash
# Download the pre-built executable
# Or build from source: go build -o ytminer.exe main.go

# Make executable (Linux/Mac)
chmod +x ytminer.exe

# Test the tool
# Windows
.\ytminer.exe --help
# Linux/Mac
./ytminer.exe --help
```

### 3. Configure API Key
```bash
# Create .env file from template
cp env.example .env

# Edit .env file and add your API key:
YOUTUBE_API_KEY=your_youtube_api_key_here
```

### 4. Test the Tool
```bash
# Test help
# Windows
.\ytminer.exe --help
# Linux/Mac
./ytminer.exe --help

# Test different analysis levels with time filters
# Quick scan - last 7 days
# Windows
.\ytminer.exe -k "programming tutorial" -l quick -t 7d
# Linux/Mac
./ytminer.exe -k "programming tutorial" -l quick -t 7d

# Balanced analysis - last 30 days
# Windows
.\ytminer.exe -k "programming tutorial" -l balanced -a all -t 30d
# Linux/Mac
./ytminer.exe -k "programming tutorial" -l balanced -a all -t 30d

# Deep dive analysis - last 90 days
# Windows
.\ytminer.exe -k "programming tutorial" -l deep -a all -t 90d
# Linux/Mac
./ytminer.exe -k "programming tutorial" -l deep -a all -t 90d
```

---

## Interactive Mode

The easiest way to use YTMiner is through interactive mode. Simply run:

```bash
# Windows (PowerShell/CMD)
.\ytminer.exe

# Linux/Mac
./ytminer.exe
```

### What happens in interactive mode:

1. **Topic Selection**: Enter the keyword you want to analyze
2. **Filters**: Choose Region, Duration, Time Range, **Order**, and Analysis Level
3. **Preview Choice**: Decide if you want to preview the results before analysis
4. **Results/Analysis**:
   - If you chose to preview, a results table is shown and you can confirm running the analysis
   - If you skip preview, the selected analysis runs immediately after the search

### Example Interactive Session:
```
Welcome to YTMiner!

1. Enter keyword → "Pokemon"
2. Choose filters → Region: BR, Duration: any, Time Range: 30d, Order: viewCount, Level: balanced
3. Preview results before analysis? → Yes
4. Table is shown → Confirm to run analysis
```

---

## Usage Examples

### Go straight to analysis (skip preview)
```bash
# Skip preview; you'll be prompted for the analysis type interactively
.\ytminer.exe -k "python tutorial" -l quick --no-preview

# Skip preview and specify analysis type, order by most viewed
.\ytminer.exe -k "python tutorial" -l balanced -a growth -o viewCount --no-preview
```

### Digital Marketers
```bash
# Market research with full analysis
# Windows
.\ytminer.exe -k "marketing strategies" -l balanced -a executive -o date
# Linux/Mac
./ytminer.exe -k "marketing strategies" -l balanced -a executive -o date

# Competitor intelligence (run analyses separately)
# Windows
.\ytminer.exe -k "your brand" -l balanced -a competitors -o relevance
.\ytminer.exe -k "your brand" -l balanced -a keywords -o relevance
# Linux/Mac
./ytminer.exe -k "your brand" -l balanced -a competitors -o relevance
./ytminer.exe -k "your brand" -l balanced -a keywords -o relevance

# Trend analysis
# Windows
.\ytminer.exe -k "trending topic" -l quick -a growth -o viewCount
# Linux/Mac
./ytminer.exe -k "trending topic" -l quick -a growth -o viewCount
```

### Region and defaults
```bash
# Use region from defaults (.env/Settings) — no -r flag provided
.\ytminer.exe -k "ai tools" -a growth

# Override defaults: set region explicitly to BR
.\ytminer.exe -k "ai tools" -a growth -r BR
```

### Quick demo: all analyses with quick level
```bash
.\ytminer.exe -k "ai tools" -l quick -a all -t 7d
```

---

## Command Reference

### Basic Parameters
| Parameter | Description | Example |
|-----------|-------------|---------|
| `-k, --keyword` | Search keyword | `-k "Python tutorial"` |
| `-r, --region` | Country code (BR, US, GB, any) | `-r BR` |
| `-d, --duration` | Video length (short, medium, long, any) | `-d short` |
| `-a, --analysis` | Analysis type | `-a growth` |
| `-l, --level` | Analysis level (quick, balanced, deep) | `-l balanced` |
| `-t, --time` | Time range (any, 7d, 30d, 90d, 1y) | `-t 30d` |
| `-o, --order` | Search order (`relevance`, `date`, `viewCount`, `rating`, `title`) | `-o viewCount` |
| `--no-preview` | Skip preview and run analysis directly | `--no-preview` |
| `--help` | Show help message | `--help` |
| `--version` | Show version information | `--version` |

Note:
- Flags override defaults loaded from `.env`/Settings for that run.
- For `relevance`/`rating`, ordering follows YouTube API; `viewCount`/`date`/`title` are explicitly sorted in the preview table.

For methodology details, see `docs/METRICS.md`.

---

## Sample Output

### Growth Pattern Analysis
```
📈 Growth Pattern Analysis

Total Videos (N=50)
Average Views: 250.4K
Average Likes: 12.1K
🚀 Niche Velocity Score (Avg. VPD): 1.9K

⚡ Highest Velocity Videos (Trending Now)
┌─────────────────────────────────┬─────────────┬─────────┬─────────┬────────────┐
│ Title                           │ Channel     │ Views   │ VPD     │ Engagement │
├─────────────────────────────────┼─────────────┼─────────┼─────────┼────────────┤
│ New AI Tool Changes Everything  │ TechGuru    │ 55K     │ 27.5K   │ 4.20%      │
│ My 2025 Productivity System     │ ProductPro  │ 150K    │ 10.0K   │ 3.85%      │
│ Is This The End of Photoshop?   │ DesignMaster│ 90K     │ 7.5K    │ 5.12%      │
└─────────────────────────────────┴─────────────┴─────────┴─────────┴────────────┘

💡 Insights
• Good niche velocity - positive momentum detected
• Excellent engagement rate
```

### Executive Report Preview
```
💼 Executive Report

📋 Executive Summary
Analysis of 200 videos shows Niche Velocity Score of 2.1K VPD with 450K average views. 
Top channel 'ProductivityGuru' leads with 12.5% market share. 3 rising star channel(s) detected with high velocity.

💡 Key Insights
• Average views: 450K
• Niche Velocity Score: 2.1K VPD
• Top trending keyword: 'ai-integration' (3.2K VPD)
• Rising stars detected: 3 channels

🎯 Strategic Recommendations
• Target breakout keyword: 'ai-integration' (3.2K VPD)
• Study rising star channel 'ModernWorker' for momentum strategies
• Post at 14:00 for maximum engagement

🏢 Competitive Intelligence
• Top competitor: ProductivityGuru (https://www.youtube.com/channel/UCxxxxxxx)
• Rising stars detected: 3 channels
• ⭐ Rising Star #1: ModernWorker (VPD: 3.2K) - https://www.youtube.com/channel/UCyyyyyyy
• ⭐ Rising Star #2: TechFlow (VPD: 2.8K) - https://www.youtube.com/channel/UCzzzzzzz
• ⭐ Rising Star #3: ProductivityHacks (VPD: 2.5K) - https://www.youtube.com/channel/UCaaaaaaa
```

---

## Configuration

Create a `.env` file with your YouTube API key:

```bash
# Required
YOUTUBE_API_KEY=your_youtube_api_key_here

# Optional settings
YTMINER_DEFAULT_REGION=any
YTMINER_DEFAULT_DURATION=any
YTMINER_DEFAULT_TIME_RANGE=any
YTMINER_DEFAULT_ORDER=relevance

# Velocity/Keyword thresholds (optional)
# Rising Star: channel AvgVPD > (multiplier * niche AvgVPD)
YTMINER_RISING_STAR_MULTIPLIER=1.5
# Long Tail: keyword Frequency <= max_freq AND Avg Engagement > min_engagement (%)
YTMINER_LONG_TAIL_MIN_ENGAGEMENT=5.0
YTMINER_LONG_TAIL_MAX_FREQ=2
```

- Changes made in the in-app **Settings** are saved back to `.env` and become the new defaults for future runs.

---

## Limitations

- **API Quota**: 10,000 units per day (free tier)
- **Search Results**: 50 videos per individual search (controlled by analysis level)
- **Data**: Some videos may have limited public stats

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

## Support

- **Issues**: Report bugs and feature requests on GitHub
- **Documentation**: Check this README for detailed usage
