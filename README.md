# 🚀 YTMiner - Advanced YouTube Analytics CLI

**A powerful command-line tool for YouTube content creators, marketers, and researchers to analyze video performance, discover trends, and optimize content strategy.**

[![Python](https://img.shields.io/badge/Python-3.8+-blue.svg)](https://python.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![YouTube API](https://img.shields.io/badge/YouTube-API%20v3-red.svg)](https://developers.google.com/youtube/v3)

## ✨ **Core Analysis Features**

### 🎯 **Growth Pattern Analysis** (`-a`)
- Performance metrics and engagement rates
- Top performers identification
- Recent trends analysis (last 7 days)
- Best engagement content discovery

### 📝 **Title Pattern Analysis** (`-t`) 
- Winning title formulas and patterns
- Most common words in successful videos
- Emoji usage patterns and effectiveness
- Title length optimization insights
- Creator-specific recommendations

### 🏆 **Competitor Analysis** (`-c`)
- Market leaders and their strategies
- Market concentration analysis
- Content patterns by top channels
- Market share distribution
- Strategic positioning recommendations

### ⏰ **Temporal Analysis** (`-time`)
- Optimal posting times (best hours and days)
- Performance by hour and day of week
- Peak engagement periods
- Recent vs evergreen content analysis
- Posting schedule optimization

### 🔍 **Keyword Analysis** (`-kwd`)
- High-performing keywords identification
- Long-tail keyword opportunities
- SEO suggestions based on data
- Trending keywords in your niche
- Search visibility optimization

### 📊 **Executive Reports** (`-rpt`)
- Comprehensive analysis with actionable insights
- Strategic recommendations based on data
- Content strategy planning
- Performance benchmarks for comparison
- Prioritized next steps with timelines

---

## 🚀 **Quick Start**

### 🎯 **Interactive Mode (Recommended for beginners)**
```bash
ytminer --interactive
```
The interactive mode guides you through the entire process with helpful prompts and recommendations based on your role. Perfect for users who want a guided experience without memorizing command-line flags.

### 1. **Get YouTube API Key**
1. Visit [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project
3. Enable **YouTube Data API v3**
4. Create credentials (API Key)
5. Copy your API key

### 2. **Setup Project**
```bash
# Clone/download the project
cd ytminer

# Create virtual environment
python -m venv .venv

# Activate virtual environment
# Windows:
.venv\Scripts\activate
# Linux/Mac:
source .venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Install as command-line tool
pip install -e .
```

### 3. **Configure API Key**
```bash
# Create .env file from template
copy env.example .env

# Edit .env file and add your API key:
YOUTUBE_API_KEY=your_youtube_api_key_here
```

### 4. **Test the Tool**
```bash
# Test help
ytminer --help

# Test comprehensive analysis
ytminer -k "programming tutorial" -a -t -c -time -kwd -rpt
```

---

## 🎯 **Interactive Mode**

The easiest way to use YTMiner is through interactive mode. Simply run:

```bash
ytminer --interactive
```

### **What happens in interactive mode:**

1. **Topic Selection**: Enter the keyword you want to analyze
2. **Analysis Scope**: Choose how many videos to analyze (10-50)
3. **Role Selection**: Pick your role for personalized recommendations:
   - **Content Creator**: Title analysis, temporal patterns, competitor insights
   - **Digital Marketer**: Market research, keyword analysis, executive reports
   - **Researcher/Analyst**: Comprehensive data analysis, temporal patterns
   - **Business Owner**: Strategic insights, competitor analysis, market intelligence
4. **Analysis Selection**: Choose recommended analyses or select specific ones
5. **Filters**: Optionally set time, region, and duration filters
6. **Results**: Get a comprehensive analysis with actionable insights

### **Example Interactive Session:**
```
🚀 Welcome to YTMiner Interactive Mode!

1. What topic would you like to analyze? → "Pokemon"
2. How many videos should I analyze? → 25
3. What best describes you? → Content Creator
4. Recommended analyses: Title, Temporal, Competitor, Keyword
5. Any specific filters? → Last 7 days, Brazil region
6. Analysis runs automatically with results!
```

---

## 💡 **Usage Examples**

### **Content Creators**
```bash
# Complete creator analysis
ytminer -k "Python tutorial" -n 25 -a -t -c -time -kwd -rpt

# Find winning title patterns
ytminer -k "cooking" -n 20 -t

# Discover optimal posting times
ytminer -k "fitness" -n 30 -time

# Analyze competitors
ytminer -k "gaming" -n 20 -c
```

### **Digital Marketers**
```bash
# Market research with full analysis
ytminer -k "marketing strategies" -n 50 -rpt

# Competitor intelligence
ytminer -k "your brand" -n 30 -c -kwd

# Trend analysis
ytminer -k "trending topic" -d 7 -a
```

### **Researchers & Analysts**
```bash
# Academic research
ytminer -k "machine learning" -n 100 -a -f research_data.json

# Engagement studies
ytminer -k "education" -n 50 -a -time

# Content pattern analysis
ytminer -k "science" -n 30 -t -c
```

---

## 📋 **Command Reference**

### **Interactive Mode**
| Command | Description |
|---------|-------------|
| `ytminer --interactive` | **Recommended**: Guided experience with prompts |
| `ytminer -i` | Short form of interactive mode |

### **Basic Parameters**
| Parameter | Description | Example |
|-----------|-------------|---------|
| `-k, --keyword` | Search keyword | `-k "Python tutorial"` |
| `-n, --max-results` | Maximum results (1-50) | `-n 25` |
| `-o, --order` | Sort order | `-o viewCount` |
| `-d, --days-back` | Last X days | `-d 7` |
| `-r, --region` | Country code | `-r BR` |
| `--duration` | Video length | `--duration short` |
| `--min-views` | Minimum views | `--min-views 100000` |
| `--min-likes` | Minimum likes | `--min-likes 1000` |
| `-f, --output` | Export to JSON | `-f data.json` |

### **Analysis Features**
| Feature | Flag | Description |
|---------|------|-------------|
| **Growth Analysis** | `-a` | Performance patterns and trends |
| **Title Analysis** | `-t` | Winning title formulas and patterns |
| **Competitor Analysis** | `-c` | Market leaders and positioning |
| **Temporal Analysis** | `-time` | Optimal posting times |
| **Keyword Analysis** | `-kwd` | SEO and keyword opportunities |
| **Executive Report** | `-rpt` | Comprehensive strategic analysis |

---

## 🎯 **Use Cases**

### **Content Creators**
- Discover what works in your niche
- Optimize titles for maximum reach
- Find best posting times for your audience
- Analyze competitors and find gaps
- Plan content strategy based on data

### **Digital Marketers**
- Market research and trend analysis
- Competitor intelligence gathering
- Influencer identification and analysis
- Campaign optimization insights
- ROI measurement and benchmarking

### **Researchers & Analysts**
- Academic research on video content
- Engagement pattern studies
- Social media behavior analysis
- Content performance research
- Trend prediction and forecasting

### **Business Intelligence**
- Market opportunity assessment
- Competitive landscape mapping
- Content strategy development
- Performance benchmarking
- Strategic planning support

---

## 📊 **Sample Output**

### **Executive Report Preview**
```
📊 EXECUTIVE REPORT - Comprehensive Analysis:

🎯 EXECUTIVE SUMMARY
   Keyword: Python tutorial
   Market Size: Large
   Competition Level: High
   Opportunity Score: 5/10
   Total Views Analyzed: 153,396,259
   Average Views: 6,135,850

🚀 ACTIONABLE INSIGHTS:
   1. 🔴 SEO Strategy: 'python' appears 37 times in successful videos
      Action: Create content series around "python"
      Expected Impact: 20-30% increase in search visibility

   2. 🟡 Posting Strategy: Best performance: Monday at 18:00
      Action: Schedule your next 3 videos for this time slot
      Expected Impact: 10-20% increase in initial views

✅ NEXT STEPS (Prioritized):
   1. 🔴 Create content series around "python"
      Timeline: This week | Category: SEO Strategy
```

---

## ⚙️ **Configuration**

### **Required Setup**
- API Key: Configured in `.env` file
- Python: 3.8 or higher
- Dependencies: Auto-installed via `requirements.txt`

### **Optional Settings**
- Output Format: Table (default) or JSON export
- Analysis Depth: Choose specific features with flags
- Filters: Use command-line parameters for customization

---

## 🚫 **Limitations**

- YouTube API Quota: 10,000 units per day (free tier)
- Maximum Results: 50 videos per search
- Data Availability: Some videos may have limited public stats
- Rate Limits: Respect YouTube's API rate limits

---

## 🤝 **Contributing**

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 **License**

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🆘 **Support**

- Issues: Report bugs and feature requests on GitHub
- Documentation: Check this README for detailed usage
- Examples: See the usage examples above

---

**Built for the YouTube creator community**