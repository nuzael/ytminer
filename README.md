# YTMiner - YouTube Analytics CLI

A command-line tool for YouTube content creators, marketers, and researchers to analyze video performance, discover trends, and optimize content strategy.

[![Go](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![YouTube API](https://img.shields.io/badge/YouTube-API%20v3-red.svg)](https://developers.google.com/youtube/v3)
[![CI](https://github.com/nuzael/ytminer/actions/workflows/ci.yml/badge.svg)](https://github.com/nuzael/ytminer/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-80%25+-brightgreen.svg)](#testing)

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

### Transcript-Aware Topic Insights (Limited)
- ⚠️ **Limited availability**: Transcript fetching may be blocked by YouTube restrictions
- Optional transcript fetching for public videos (timedtext); language preference via `YTMINER_TRANSCRIPT_LANGS`
- Enriches keyword/topic analysis and future clustering (title + description + transcript)
- Auto-generated captions may be used when manual captions are not available
- Tip: set `YTMINER_CACHE_DIR` to control where transcript cache files are stored (useful for CI/dev isolation)

### Opportunity Score (`-a opportunity`) (Enhanced)
Note: weights are not auto-normalized. You can sum > 1; the score is linear and components are standardized (z-scores) where applicable.
- Ranks videos/themes by a combined signal: velocity (VPD), engagement (likes per 1k views), freshness (younger is better), minus saturation penalty
- Runs entirely in-memory; configurable weights via env (`YTMINER_OPP_W_VPD`, `YTMINER_OPP_W_LIKE`, `YTMINER_OPP_W_FRESH`, `YTMINER_OPP_W_SAT`, `YTMINER_OPP_W_SLOPE`); transcripts (if enabled) enrich future topic grouping, not this score

### Weight Profiles
Predefined strategies to steer the Opportunity Score:
- exploration: Discover new niches and emerging trends
- evergreen: Focus on timeless, high-quality content
- trending: Catch viral content and momentum
- balanced: Default balanced approach

Usage:
```bash
# List profiles
./ytminer.exe --profiles

# Apply profile in CLI
./ytminer.exe -k "ai tools" -a opportunity --profile trending

# Apply profile then run all
./ytminer.exe -k "python tutorial" -a all --profile exploration
```

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
- **Insights**: Comprehensive patterns, velocity trends, competitor analysis
- **Use case**: "What's working in this niche?"

#### Deep Dive (~1500 units, up to ~750 videos, 2-3min)
- **Best for**: Research, competitive analysis, content strategy
- **Strategy**: Exhaustive search across multiple orders and criteria
- **Data**: 10+ pages across relevance, viewCount, date, rating
- **Insights**: Complete market picture, emerging trends, saturation analysis
- **Use case**: "I need to understand this market completely"

### Command Line Interface

#### Basic Usage
```bash
ytminer -k "your keyword" -l [quick|balanced|deep] -a [analysis_type]
```

#### Advanced Options
```bash
# Region and time filtering
ytminer -k "python tutorial" -r US -t 30d -l balanced

# Specific analysis with custom order
ytminer -k "machine learning" -o date -l deep -a competitors

# Executive report with custom filters
ytminer -k "content marketing" -o viewCount -l deep -a executive
```

#### Analysis Types (`-a`)
- `growth` - Growth pattern analysis
- `titles` - Title pattern analysis  
- `competitors` - Competitor analysis
- `temporal` - Temporal analysis
- `keywords` - Keyword analysis
- `opportunity` - Opportunity Score ranking
- `executive` - Executive summary report
- `all` - Run all analyses

#### Order Options (`-o`)
- `relevance` - Most relevant (default)
- `date` - Most recent
- `viewCount` - Most viewed
- `rating` - Highest rated

#### Region Codes (`-r`)
- `US`, `BR`, `GB`, `CA`, `AU`, `DE`, `FR`, `ES`, `IT`, `JP`, `KR`, `IN`, `MX`, `RU`, `CN`

#### Time Ranges (`-t`)
- `1d`, `7d`, `30d`, `90d`, `1y` (days, weeks, months, years)

---

## Installation

### Prerequisites
- Go 1.24+ installed
- YouTube Data API v3 key

### Build from Source
```bash
git clone https://github.com/nuzael/ytminer.git
cd ytminer
go build -o ytminer.exe
```

### Configuration
1. Copy `env.example` to `.env`
2. Add your YouTube API key:
```bash
YOUTUBE_API_KEY=your_api_key_here
```

### Transcript Configuration (Optional)
```bash
# Enable transcript fetching by default
YTMINER_WITH_TRANSCRIPTS=true

# Preferred languages (comma-separated)
YTMINER_TRANSCRIPT_LANGS=en,pt,es

# Cache directory (useful for CI/dev isolation)
YTMINER_CACHE_DIR=.cache
```

---

## Architecture Overview

YTMiner follows Clean Architecture principles with clear separation of concerns:

```
ytminer/
├── domain/           # Business logic (pure functions)
│   ├── metrics/      # Core metrics (VPD, Slope, etc.)
│   └── score/        # Opportunity Score computation
├── platform/         # External integrations
│   ├── ytapi/        # YouTube API adapter
│   └── transcripts/  # Transcript fetching & caching
├── analysis/         # Analysis orchestration
├── config/           # Configuration management
├── ui/               # User interface
└── utils/            # Utilities
```

### Key Design Principles
- **Functional Core, Imperative Shell**: Business logic is pure and testable
- **Ports & Adapters**: External dependencies are abstracted
- **Domain-Driven Design**: Code organized around business concepts
- **Testability**: High test coverage with unit and integration tests

---

## Testing

YTMiner includes comprehensive test coverage with beautiful output using `gotestsum`.

### Running Tests

#### Quick Test Run
```bash
# Git Bash / Linux / Mac
./test.sh

# PowerShell
.\test.ps1
```

#### Manual Test Commands
```bash
# All tests with coverage
gotestsum --format testdox -- -coverprofile=coverage.out ./...

# Specific package
gotestsum --format testdox ./analysis

# Different output formats
gotestsum --format dots ./...           # Minimalist
gotestsum --format short-verbose ./...  # Detailed
```

### Test Coverage

Current coverage by package:
- **`domain/metrics`**: 86.7% ✅ (excellent)
- **`platform/ytapi`**: 83.3% ✅ (very good)
- **`domain/score`**: 73.3% ✅ (good)
- **`analysis`**: 64.4% ⚠️ (can improve)
- **`config`**: 48.7% ⚠️ (can improve)
- **`platform/transcripts`**: 31.6% ⚠️ (can improve)
- **Total**: 26.0% (low due to packages without tests)

### Viewing Coverage
```bash
# Summary in terminal
go tool cover -func=coverage.out

# Detailed HTML report
go tool cover -html=coverage.out
```

### Test Types
- **Unit Tests**: Pure functions in `domain/`
- **Integration Tests**: API adapters in `platform/`
- **Golden Tests**: End-to-end validation in `e2e/`
- **Property-Based Tests**: Monotonicity and edge cases

---

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
