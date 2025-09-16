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

### Analysis Levels
YTMiner offers three analysis levels to optimize API usage and analysis depth:

#### Quick Scan (~200 units, 50 videos, 30-60s)
- **Best for**: Quick exploration, demos, trend validation
- **Data**: Single search with 50 videos
- **Insights**: Basic patterns and trends
- **Use case**: "Is this topic worth analyzing further?"

#### Balanced (~1000 units, 200 videos, 2-3min)
- **Best for**: Regular content creators, marketers
- **Data**: 4 searches across different parameters
- **Insights**: Reliable patterns and recommendations
- **Use case**: "What should my content strategy be?"

#### Deep Dive (~3000 units, 600 videos, 5-8min)
- **Best for**: Strategic analysis, research, business decisions
- **Data**: 12 searches across regions, durations, and related topics
- **Insights**: Comprehensive market intelligence
- **Use case**: "What's the complete competitive landscape?"

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
2. **Filters**: Choose Region, Duration, Time Range, and Analysis Level
3. **Preview Choice**: Decide if you want to preview the results before analysis
4. **Results/Analysis**:
   - If you chose to preview, a results table is shown and you can confirm running the analysis
   - If you skip preview, the selected analysis runs immediately after the search

### Example Interactive Session:
```
Welcome to YTMiner!

1. Enter keyword → "Pokemon"
2. Choose filters → Region: BR, Duration: any, Time Range: 30d, Level: balanced
3. Preview results before analysis? → Yes
4. Table is shown → Confirm to run analysis
```

---

## Usage Examples

### Content Creators
```bash
# Quick trend check
# Windows
.\ytminer.exe -k "Python tutorial" -l quick
# Linux/Mac
./ytminer.exe -k "Python tutorial" -l quick

# Balanced content strategy
# Windows
.\ytminer.exe -k "cooking" -l balanced -a titles
# Linux/Mac
./ytminer.exe -k "cooking" -l balanced -a titles

# Deep competitive analysis
# Windows
.\ytminer.exe -k "fitness" -l deep -a all
# Linux/Mac
./ytminer.exe -k "fitness" -l deep -a all
```

### Go straight to analysis (skip preview)
```bash
# Skip preview; you'll be prompted for the analysis type interactively
.\ytminer.exe -k "python tutorial" -l quick --no-preview

# Skip preview and specify analysis type
.\ytminer.exe -k "python tutorial" -l balanced -a growth --no-preview
```

### Digital Marketers
```bash
# Market research with full analysis
# Windows
.\ytminer.exe -k "marketing strategies" -l balanced -a executive
# Linux/Mac
./ytminer.exe -k "marketing strategies" -l balanced -a executive

# Competitor intelligence
# Windows
.\ytminer.exe -k "your brand" -l balanced -a competitors -a keywords
# Linux/Mac
./ytminer.exe -k "your brand" -l balanced -a competitors -a keywords

# Trend analysis
# Windows
.\ytminer.exe -k "trending topic" -l quick -a growth
# Linux/Mac
./ytminer.exe -k "trending topic" -l quick -a growth
```

### Researchers & Analysts
```bash
# Academic research
# Windows
.\ytminer.exe -k "machine learning" -l deep -a all
# Linux/Mac
./ytminer.exe -k "machine learning" -l deep -a all

# Engagement studies
# Windows
.\ytminer.exe -k "education" -l balanced -a growth -a temporal
# Linux/Mac
./ytminer.exe -k "education" -l balanced -a growth -a temporal

# Content pattern analysis
# Windows
.\ytminer.exe -k "science" -l quick -a titles -a competitors
# Linux/Mac
./ytminer.exe -k "science" -l quick -a titles -a competitors
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
| `--no-preview` | Skip preview and run analysis directly | `--no-preview` |
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

### **Analysis Levels**
| Level | Flag | Units | Videos | Time | Best For |
|-------|------|-------|--------|------|----------|
| **Quick Scan** | `-l quick` | ~200 | 50 | 30-60s | Exploration, demos |
| **Balanced** | `-l balanced` | ~1000 | 200 | 2-3min | Regular analysis |
| **Deep Dive** | `-l deep` | ~3000 | 600 | 5-8min | Strategic analysis |

### **Time Ranges**
| Range | Flag | Description | Use Case |
|-------|------|-------------|----------|
| **Any time** | `-t any` | No time filter | Historical analysis |
| **Last 7 days** | `-t 7d` | Recent videos only | Trending content |
| **Last 30 days** | `-t 30d` | Past month | Current trends |
| **Last 90 days** | `-t 90d` | Past quarter | Seasonal analysis |
| **Last year** | `-t 1y` | Past year | Annual trends |

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
YTMINER_DEFAULT_REGION=any
YTMINER_DEFAULT_DURATION=any
```

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
