# YTMiner

A command-line tool for analyzing YouTube videos with advanced filtering and growth pattern analysis.

## Features

- **Keyword search** with advanced filters
- **Date filtering** (last X days)
- **Regional filtering** (specific country)
- **Duration filtering** (short, medium, long)
- **Performance filtering** (minimum views/likes)
- **Growth pattern analysis**
- **Detailed reports** with engagement metrics
- **JSON export**

## Quick Start

### 1. Get YouTube API Key
1. Visit [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project
3. Enable **YouTube Data API v3**
4. Create credentials (API Key)
5. Copy your API key

### 2. Setup Project
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
```

### 3. Configure API Key

#### **Method 1: .env file (Recommended)**
```bash
# Create .env file from template
copy env.example .env

# Edit .env file and add your API key:
YOUTUBE_API_KEY=your_youtube_api_key_here
```

#### **Method 2: Environment variable (Temporary)**
```bash
# Windows Command Prompt
set YOUTUBE_API_KEY=your_key_here

# Windows PowerShell
$env:YOUTUBE_API_KEY="your_key_here"

# Linux/Mac
export YOUTUBE_API_KEY=your_key_here
```

### 4. Test the Tool
```bash
# Test help
python ytminer.py --help

# Test search
python ytminer.py -k "Pokemon" -a
```

## Usage

### Basic Commands
```bash
# Basic search
python ytminer.py -k "Pokemon"

# Search with analysis
python ytminer.py -k "Pokemon" -a

# Search last 7 days
python ytminer.py -k "Pokemon" -d 7 -a

# Most viewed videos from Brazil
python ytminer.py -k "Pokemon" -o viewCount -r BR -n 30

# Short videos with minimum 100k views
python ytminer.py -k "Pokemon" --duration short --min-views 100000

# Save results to file
python ytminer.py -k "Pokemon" -a -f results.json
```

### Available Parameters

| Parameter | Description | Example |
|-----------|-------------|---------|
| `-k, --keyword` | Search keyword | `-k "Pokemon"` |
| `-n, --max-results` | Maximum number of results | `-n 50` |
| `-o, --order` | Ordering (relevance, date, rating, viewCount, title) | `-o viewCount` |
| `-d, --days-back` | Number of days back | `-d 7` |
| `-r, --region` | Country code | `-r BR` |
| `--duration` | Duration (short, medium, long, any) | `--duration short` |
| `--min-views` | Minimum view count | `--min-views 100000` |
| `--min-likes` | Minimum like count | `--min-likes 1000` |
| `-f, --output` | Output file (JSON) | `-f data.json` |
| `-a, --analysis` | Show pattern analysis | `-a` |

## Pattern Analysis

The tool provides detailed analysis including:

- **General statistics** (average views, likes, engagement)
- **Top performers** (most viewed videos)
- **Recent trends** (growing videos in the last 7 days)
- **Best engagement** (most liked videos)

## Use Cases

### For Content Creators
- Discover rapidly growing videos
- Analyze title and description patterns
- Identify successful channels in niche
- Find content opportunities

### For Researchers
- YouTube trend analysis
- Engagement studies
- Topic monitoring

### For Marketers
- Market research
- Competitor analysis
- Influencer identification

## 🔧 Files You Need to Configure

### ✅ **REQUIRED - Edit This File:**

#### **`.env`** - API Key Configuration
```bash
# Create from template
copy env.example .env

# Edit .env file and add your YouTube API key:
YOUTUBE_API_KEY=your_youtube_api_key_here
```

### ℹ️ **Configuration:**
- **API Key**: Configured in `.env` file
- **Other parameters**: Passed via command line
- **See options**: `python ytminer.py --help`

### 📋 **Available Settings:**
- **Output Format**: Always table (no other formats supported)
- **Analysis**: Use `-a` flag to enable
- **Filters**: Use command line parameters (see examples below)
- **Export**: Use `-f filename.json` to save results

### ❌ **DON'T Edit These Files:**
- `ytminer.py` - Main script (ready to use)
- `README.md` - Documentation  
- `requirements.txt` - Dependencies
- All `.bat` files - Helper scripts
- `env.example` - Template file

## Project Structure

```
ytminer/
├── ytminer.py          # Main script
├── requirements.txt    # Dependencies
├── env.example        # Configuration template
├── run.bat           # Windows runner
├── activate.bat      # Virtual env activator
└── README.md         # This file
```

## Limitations

- YouTube API has quota limits (10,000 units per day)
- Maximum 50 results per search
- Some videos may not have public statistics

## License

This project is open source and available under the MIT license.
