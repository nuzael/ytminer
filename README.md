# YTMiner - YouTube Analytics CLI

A command-line tool for YouTube content creators, marketers, and researchers to analyze video performance, discover trends, and optimize content strategy.

[![Go](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![YouTube API](https://img.shields.io/badge/YouTube-API%20v3-red.svg)](https://developers.google.com/youtube/v3)

## Quick Reference

- Interactive (recommended):
```bash
# Windows
.\ytminer.exe
# Linux/Mac
./ytminer.exe
```

- CLI basics:
```bash
# Quick scan
ytminer -k "ai tools" -l quick
# Balanced with region and time range
ytminer -k "meditation" -r BR -t 30d -l balanced
# Deep dive, most recent, executive report
ytminer -k "content marketing" -o date -l deep -a executive
```

## Core Analysis Features

### Growth Pattern Analysis (`-a growth`)
- Velocity-first metrics (Avg. VPD, momentum)
- Trending videos ranked by velocity
- Early signals for breakout topics

### Title Pattern Analysis (`-a titles`)
- High-performing title patterns
- Keywords, emojis, and length insights

### Competitor Analysis (`-a competitors`)
- Rising Stars detection (VPD > niche baseline)
- Channel velocity and market share

### Temporal Analysis (`-a temporal`)
- Best posting windows (day/hour)
- Engagement over time-of-day and weekday

### Keyword Analysis (`-a keywords`)
- Breakout keywords by Avg. VPD
- Core keywords by frequency

### Transcript-Aware Topic Insights (new)
- Optional transcript fetching for public videos (timedtext); language preference via `YTMINER_TRANSCRIPT_LANGS`
- Enriches keyword/topic analysis and future clustering (title + description + transcript)
- Auto-generated captions may be used when manual captions are not available

### Opportunity Score (`-a opportunity`)
- Ranks videos/themes by a combined signal: velocity (VPD), engagement (likes per 1k views), freshness (younger is better), minus saturation penalty
- Runs entirely in-memory; configurable weights via env (`YTMINER_OPP_W_VPD`, `YTMINER_OPP_W_LIKE`, `YTMINER_OPP_W_FRESH`, `YTMINER_OPP_W_SAT`); transcripts (if enabled) enrich future topic grouping, not this score

### Executive Reports (`-a executive`)
- Summary of niche momentum and leaders
- Strategic recommendations

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
YTMiner offers three analysis levels with adaptive search strategy to optimize API usage and maximize data collection:

#### Quick Scan (~300 units, up to ~150 videos, 30-60s)
- **Best for**: Quick exploration, demos, trend validation
- **Strategy**: Prioritizes your chosen order, adapts automatically if topic has limited content
- **Data**: 3 pages of your chosen criteria, complemented with backup orders if needed
- **Insights**: Basic patterns and velocity indicators
- **Use case**: "Is this topic hot right now?"

#### Balanced (~800 units, up to ~400 videos, 1-2min)
- **Best for**: Regular content creators, marketers
- **Strategy**: Deep dive into your preferences, intelligently supplements for comprehensive coverage
- **Data**: 5 pages of your chosen criteria + backup orders (relevance, viewCount, date) as needed
- **Insights**: Comprehensive patterns and recommendations
- **Use case**: "What is the current growth strategy in this niche?"

#### Deep Dive (~2000 units, up to ~1000 videos, 3-5min)
- **Best for**: Strategic analysis, research, business decisions
- **Strategy**: Maximum depth with intelligent adaptation for small niches
- **Data**: 10 pages of your exact configuration + comprehensive backup orders (relevance, viewCount, date, rating) as needed
- **Insights**: Exhaustive market and competitor velocity intelligence
- **Use case**: "What's the complete competitive landscape, including emerging players?"

**Why this approach?** If you want trending content (`order=date`) but the niche only has 50 recent videos, our system automatically supplements with the most relevant content available. This ensures you get maximum value from your API quota while respecting your primary analysis intent. The strategy creates **gold-standard data** for identifying trends, growth opportunities, and niche analysis.

### How Search Order Works

The `order` parameter serves two important functions:

1. **Data Collection Priority**: Your chosen order determines what type of content we prioritize during data collection
2. **Final Display Sorting**: The final results table is sorted by your chosen order for consistent analysis

**Available Orders:**
- `relevance`: YouTube's relevance algorithm (best overall match)
- `date`: Most recent uploads first (trending/fresh content)
- `viewCount`: Highest view counts first (viral/popular content)
- `rating`: Highest rated content first (quality-focused)
- `title`: Alphabetical order (A-Z)

**Adaptive Behavior:**
- **Small niches**: If your chosen order exhausts available content, we automatically supplement with other orders
- **Large niches**: We focus entirely on your chosen order for maximum analytical depth
- **Quality guarantee**: All searches respect your region, duration, and time range filters

**Example**: If you choose `order=viewCount` for "meditation techniques":
1. We first collect the highest-viewed meditation videos
2. If there aren't enough, we supplement with most relevant ones
3. Final results are sorted by view count for your analysis
4. This gives you both viral content AND comprehensive coverage

---

## Architecture Overview

- Domain (pure):
  - `domain/metrics`: core metrics functions (VPD, like_rate, freshness, saturation)
  - `domain/score`: Opportunity Score computation (pure, deterministic)
- Platform (adapters):
  - `platform/ytapi`: port/adapter wrapping `youtube.Client` (search, transcripts)
- App/Analysis:
  - `analysis`: orchestrates analyses using domain and platform
- UI:
  - `ui`: rendering only, no business logic
- Config:
  - `config`: defaults and env overrides (including Opportunity weights)

## Installation

### 1. Get a YouTube API Key
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable the **YouTube Data API v3**
4. Create credentials (API Key)
5. Copy your API key

### 2. Download YTMiner
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

### 3. Configure API Key and Transcripts
```bash
# Create .env file from template
cp env.example .env

# Edit .env file and add your API key and optional transcript languages
YOUTUBE_API_KEY=your_youtube_api_key_here
YTMINER_TRANSCRIPT_LANGS=en,pt
```

### 4. Test the Tool
```bash
# Interactive mode (recommended for first-time users)
.\ytminer.exe

# Quick CLI test
.\ytminer.exe -k "productivity tips" -l quick

# See all options
.\ytminer.exe --help
```

---

## Usage Examples

### Basic Usage
```bash
# Interactive mode (guided experience)
.\ytminer.exe

# Quick keyword analysis
.\ytminer.exe -k "ai tools" -l quick

# Deep competitive analysis
.\ytminer.exe -k "content marketing" -l deep -a growth

# Regional analysis with time constraints
.\ytminer.exe -k "meditation" -r BR -t 7d -l balanced
```

### Advanced Filtering
```bash
# Analyze short-form content trends
.\ytminer.exe -k "cooking hacks" -d short -o date -l deep

# Find high-engagement long-form content
.\ytminer.exe -k "productivity" -d long -o viewCount -l balanced

# Recent trends in specific region
.\ytminer.exe -k "tech reviews" -r US -t 30d -o date -l quick
```

### Market Analysis Workflows
```bash
# Trend validation workflow
.\ytminer.exe -k "sustainable living" -l quick -a velocity
.\ytminer.exe -k "sustainable living" -l deep -a growth

# Competitive intelligence workflow
.\ytminer.exe -k "crypto education" -l balanced -a intel
.\ytminer.exe -t 90d -l deep -a executive

# Content strategy research
.\ytminer.exe -k "fitness motivation" -d any -o relevance -l deep -a growth
```

---

## Command Reference

### Basic Parameters
| Parameter | Description | Example |
|-----------|-------------|---------|
| `-k, --keyword` | Search keyword | `-k "Python tutorial"` |
| `-r, --region` | Country code (BR, US, GB, any) | `-r BR` |
| `-d, --duration` | Video length (short, medium, long, any) | `-d short` |
| `-a, --analysis` | Analysis type | `-a opportunity` |
| `-l, --level` | Analysis level (quick, balanced, deep) | `-l balanced` |
| `-t, --time` | Time range (any, 1h, 24h, 7d, 30d, 90d, 180d, 1y) | `-t 30d` |
| `-o, --order` | Search order (`relevance`, `date`, `viewCount`, `rating`, `title`) | `-o viewCount` |
| `--no-preview` | Skip preview and run analysis directly | `--no-preview` |
| `--help` | Show help message | `--help` |
| `--version` | Show version information | `--version` |

Note:
- Flags override defaults loaded from `.env`/Settings for that run.
- For `relevance`/`rating`, ordering follows YouTube API; `viewCount`/`date`/`title` are explicitly sorted in the preview table.

For methodology details, see `docs/METRICS.md`.

---

## Limitations

- **API Quotas**: Each analysis level uses different amounts of YouTube API quota
- **Rate Limits**: YouTube API has rate limiting; tool includes automatic retry logic
- **Search Scope**: Limited to publicly available YouTube videos
- **Data Freshness**: Statistics reflect YouTube's caching and may have slight delays
- **Adaptive Results**: Video count depends on content availability for your specific criteria

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

## Support

- **Issues**: Report bugs and feature requests on GitHub
- **Documentation**: Check this README for detailed usage

## Continuous Integration (CI)

This repository includes a GitHub Actions workflow that runs on pushes and PRs to `main`:
- `go vet`, formatting check (`gofmt -s`)
- `go test ./... -race` with coverage report (`coverage.out`)
- Enforces a minimum coverage of 80%
- Uploads `coverage.out` as an artifact

### Run tests with coverage locally

```bash
# From repository root
go test ./... -coverprofile=coverage.out -covermode=atomic

go tool cover -func=coverage.out   # summary on console

go tool cover -html=coverage.out -o coverage.html  # open in the browser
```
