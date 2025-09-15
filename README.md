# YTMiner - YouTube Analytics CLI

A command-line tool for YouTube content creators, marketers, and researchers to analyze video performance, discover trends, and optimize content strategy.

[![Go](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![YouTube API](https://img.shields.io/badge/YouTube-API%20v3-red.svg)](https://developers.google.com/youtube/v3)

## Core Analysis Features

### Growth Pattern Analysis (`-a growth`)
- Performance metrics and engagement rates
- Top performers identification
- Growth rate calculations
- Best engagement content discovery

### Title Pattern Analysis (`-a titles`) 
- Winning title formulas and patterns
- Most common words in successful videos
- Emoji usage patterns and effectiveness
- Title length optimization insights
- Creator-specific recommendations

### Competitor Analysis (`-a competitors`)
- Market leaders and their strategies
- Market concentration analysis
- Content patterns by top channels
- Market share distribution
- Strategic positioning recommendations

### Temporal Analysis (`-a temporal`)
- Optimal posting times (best hours and days)
- Performance by hour and day of week
- Peak engagement periods
- Posting schedule optimization

### Keyword Analysis (`-a keywords`)
- High-performing keywords identification
- Long-tail keyword opportunities
- SEO suggestions based on data
- Trending keywords in your niche
- Search visibility optimization

### Executive Reports (`-a executive`)
- Comprehensive analysis with actionable insights
- Strategic recommendations based on data
- Content strategy planning
- Performance benchmarks for comparison
- Prioritized next steps with timelines

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

# Test comprehensive analysis
# Windows
.\ytminer.exe -k "programming tutorial" -n 25 -a all
# Linux/Mac
./ytminer.exe -k "programming tutorial" -n 25 -a all
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
2. **Analysis Scope**: Choose how many videos to analyze (10-50)
3. **Analysis Selection**: Choose specific analyses or run all
4. **Filters**: Optionally set region and duration filters
5. **Results**: Get a comprehensive analysis with actionable insights

### Example Interactive Session:
```
Welcome to YTMiner!

1. What topic would you like to analyze? → "Pokemon"
2. How many videos should I analyze? → 25
3. Choose analysis type → All Analyses
4. Any specific filters? → Brazil region, any duration
5. Analysis runs automatically with results!
```

---

## Usage Examples

### Content Creators
```bash
# Complete creator analysis
# Windows
.\ytminer.exe -k "Python tutorial" -n 25 -a all
# Linux/Mac
./ytminer.exe -k "Python tutorial" -n 25 -a all

# Find winning title patterns
# Windows
.\ytminer.exe -k "cooking" -n 20 -a titles
# Linux/Mac
./ytminer.exe -k "cooking" -n 20 -a titles

# Discover optimal posting times
# Windows
.\ytminer.exe -k "fitness" -n 30 -a temporal
# Linux/Mac
./ytminer.exe -k "fitness" -n 30 -a temporal

# Analyze competitors
# Windows
.\ytminer.exe -k "gaming" -n 20 -a competitors
# Linux/Mac
./ytminer.exe -k "gaming" -n 20 -a competitors
```

### Digital Marketers
```bash
# Market research with full analysis
# Windows
.\ytminer.exe -k "marketing strategies" -n 50 -a executive
# Linux/Mac
./ytminer.exe -k "marketing strategies" -n 50 -a executive

# Competitor intelligence
# Windows
.\ytminer.exe -k "your brand" -n 30 -a competitors -a keywords
# Linux/Mac
./ytminer.exe -k "your brand" -n 30 -a competitors -a keywords

# Trend analysis
# Windows
.\ytminer.exe -k "trending topic" -n 25 -a growth
# Linux/Mac
./ytminer.exe -k "trending topic" -n 25 -a growth
```

### Researchers & Analysts
```bash
# Academic research
# Windows
.\ytminer.exe -k "machine learning" -n 50 -a all
# Linux/Mac
./ytminer.exe -k "machine learning" -n 50 -a all

# Engagement studies
# Windows
.\ytminer.exe -k "education" -n 30 -a growth -a temporal
# Linux/Mac
./ytminer.exe -k "education" -n 30 -a growth -a temporal

# Content pattern analysis
# Windows
.\ytminer.exe -k "science" -n 25 -a titles -a competitors
# Linux/Mac
./ytminer.exe -k "science" -n 25 -a titles -a competitors
```

---

## Command Reference

### Basic Parameters
| Parameter | Description | Example |
|-----------|-------------|---------|
| `-k, --keyword` | Search keyword | `-k "Python tutorial"` |
| `-n, --max-results` | Maximum results (1-50) | `-n 25` |
| `-r, --region` | Country code (BR, US, GB, any) | `-r BR` |
| `-d, --duration` | Video length (short, medium, long, any) | `-d short` |
| `-a, --analysis` | Analysis type | `-a growth` |
| `--help` | Show help message | `--help` |
| `--version` | Show version information | `--version` |

### Analysis Types
| Analysis | Flag | Description |
|----------|------|-------------|
| **Growth Analysis** | `-a growth` | Performance patterns and trends |
| **Title Analysis** | `-a titles` | Winning title formulas and patterns |
| **Competitor Analysis** | `-a competitors` | Market leaders and positioning |
| **Temporal Analysis** | `-a temporal` | Optimal posting times |
| **Keyword Analysis** | `-a keywords` | SEO and keyword opportunities |
| **Executive Report** | `-a executive` | Comprehensive strategic analysis |
| **All Analyses** | `-a all` | Run all analysis types |

---

## Sample Output

### Growth Pattern Analysis
```
Growth Pattern Analysis

Total Videos: 25
Average Views: 1.2M
Average Likes: 45.2K
Growth Rate: 15.3%

Top Performing Videos:
1. "Advanced Python Tutorial" - 2.5M views (3.2% engagement)
2. "Python for Beginners" - 1.8M views (2.8% engagement)
3. "Python Data Science" - 1.5M views (2.5% engagement)

Insights:
• High-performing content with over 1M average views
• Strong growth trend detected
• Excellent engagement rate
```

### Executive Report Preview
```
Executive Report - Comprehensive Analysis

Executive Summary
Keyword: Python tutorial
Market Size: Large
Competition Level: High
Total Views Analyzed: 30,250,000
Average Views: 1,210,000

Actionable Insights:
1. SEO Strategy: 'python' appears 37 times in successful videos
   Action: Create content series around "python"
   Expected Impact: 20-30% increase in search visibility

2. Posting Strategy: Best performance: Monday at 18:00
   Action: Schedule your next 3 videos for this time slot
   Expected Impact: 10-20% increase in initial views

Next Steps (Prioritized):
1. Create content series around "python"
   Timeline: This week | Category: SEO Strategy
```

---

## Configuration

Create a `.env` file with your YouTube API key:

```bash
# Required
YOUTUBE_API_KEY=your_youtube_api_key_here

# Optional settings
YTMINER_MAX_RESULTS=25
YTMINER_DEFAULT_REGION=any
YTMINER_DEFAULT_DURATION=any
```

---

## Limitations

- **API Quota**: 10,000 units per day (free tier)
- **Max Results**: 50 videos per search
- **Data**: Some videos may have limited public stats

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

## Support

- **Issues**: Report bugs and feature requests on GitHub
- **Documentation**: Check this README for detailed usage

---
