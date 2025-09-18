# YTMiner - YouTube Analytics CLI

A command-line tool for YouTube content creators, marketers, and researchers to analyze video performance, discover trends, and optimize content strategy.

[![Go](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![YouTube API](https://img.shields.io/badge/YouTube-API%20v3-red.svg)](https://developers.google.com/youtube/v3)

## Quick Start

### Option 1: Download Pre-built Binary (Recommended)

#### Automatic Installation

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/nuzael/ytminer/main/install.sh | bash
```

**Windows (PowerShell):**
```powershell
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/nuzael/ytminer/main/install.ps1" -OutFile "install.ps1"
.\install.ps1
```

#### Manual Installation
1. Go to [Releases](https://github.com/nuzael/ytminer/releases)
2. Download the binary for your platform
3. Make it executable and run:

```bash
# Windows
ytminer.exe

# Linux/Mac
chmod +x ytminer
./ytminer
```

### Option 2: Build from Source
```bash
git clone https://github.com/nuzael/ytminer.git
cd ytminer
go build -o ytminer.exe
```

## Usage

### Interactive Mode (Recommended)
```bash
# Windows
.\ytminer.exe

# Linux/Mac
./ytminer
```

### CLI Mode
```bash
# Quick analysis
ytminer -k "ai tools" -l quick

# Balanced analysis with filters
ytminer -k "meditation" -r BR -t 30d -l balanced

# Deep dive with opportunity scoring
ytminer -k "content marketing" -l deep -a opportunity --profile trending
```

## Prerequisites

- YouTube Data API v3 key
- Go 1.24+ (only if building from source)

## Configuration

1. Copy `env.example` to `.env`
2. Add your YouTube API key:
```bash
YOUTUBE_API_KEY=your_api_key_here
```

## Core Features

- **Growth Analysis** - Identify trending videos and momentum
- **Title Analysis** - Discover high-performing title patterns
- **Competitor Analysis** - Find rising stars and market opportunities
- **Opportunity Score** - AI-powered ranking of content opportunities
- **Transcript Analysis** - Content insights from video transcripts
- **Temporal Analysis** - Best posting times and patterns

## Analysis Levels

- **Quick** (~150 videos, 30-60s) - Fast exploration
- **Balanced** (~400 videos, 1-2min) - Regular analysis
- **Deep** (~750 videos, 2-3min) - Comprehensive research

## Weight Profiles

Pre-configured analysis strategies:
- `exploration` - Discover new niches and emerging trends
- `evergreen` - Focus on timeless, high-quality content
- `trending` - Catch viral content and momentum
- `balanced` - Default balanced approach

## Documentation

- [Complete Usage Guide](docs/USAGE.md) - Detailed CLI reference
- [Analysis Types](docs/ANALYSIS.md) - Deep dive into each analysis
- [Configuration](docs/CONFIGURATION.md) - Advanced settings
- [Architecture](docs/ARCHITECTURE.md) - Technical details
- [Metrics](docs/METRICS.md) - Methodology and formulas

## Examples

```bash
# Find trending opportunities
ytminer -k "python tutorial" --profile trending -a opportunity

# Analyze competitor landscape
ytminer -k "machine learning" -l deep -a competitors

# Executive summary
ytminer -k "content marketing" -l deep -a executive

# Show available profiles
ytminer --profiles
```

## Contributing

See [CONTRIBUTING.md](docs/CONTRIBUTING.md) for development setup and guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.
