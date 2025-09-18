# YTMiner Usage Guide

Complete reference for using YTMiner CLI and interactive modes.

## Command Line Interface

### Basic Syntax

**Note**: If you downloaded a pre-built binary, use the full path or add it to your PATH. If you built from source, the binary will be in the project directory.
```bash
ytminer [OPTIONS]
```

### Required Options
- `-k string` - Search keyword (required for CLI mode)

### Optional Options

#### Search Filters
- `-r string` - Search region: any, BR, US, GB (default: any)
- `-d string` - Video duration: any, short, medium, long (default: any)
- `-t string` - Time range: any, 1h, 24h, 7d, 30d, 90d, 180d, 1y (default: any)
- `-o string` - Search order: relevance, date, viewCount, rating, title (default: relevance)

#### Analysis Configuration
- `-a string` - Analysis type: growth, titles, competitors, temporal, keywords, executive, all, opportunity
- `-l string` - Analysis level: quick, balanced, deep (default: balanced)
- `--no-preview` - Skip preview table and run analysis directly

#### Opportunity Score Configuration
- `--opp-w-vpd float` - VPD weight (default: 0.45)
- `--opp-w-like float` - Like rate weight (default: 0.25)
- `--opp-w-fresh float` - Freshness weight (default: 0.20)
- `--opp-w-sat float` - Saturation penalty weight (default: 0.30)
- `--opp-w-slope float` - Slope weight (default: 0.15)
- `--profile string` - Apply predefined profile: exploration, evergreen, trending, balanced
- `--profiles` - Show available profiles and exit

#### General Options
- `--help` - Show help
- `--version` - Show version

## Analysis Types

### Growth Analysis (`-a growth`)
Identifies trending videos and momentum patterns.
- **Output**: Top videos by velocity, momentum insights
- **Best for**: Finding viral content, trending topics
- **Use case**: "What's trending in this niche?"

### Title Analysis (`-a titles`)
Analyzes high-performing title patterns.
- **Output**: Title patterns, keywords, length insights
- **Best for**: Content optimization, A/B testing
- **Use case**: "What titles work best in this niche?"

### Competitor Analysis (`-a competitors`)
Finds rising stars and market opportunities.
- **Output**: Rising channels, market share analysis
- **Best for**: Competitive intelligence, partnership opportunities
- **Use case**: "Who are the rising stars in this space?"

### Temporal Analysis (`-a temporal`)
Analyzes best posting times and patterns.
- **Output**: Optimal posting windows, engagement patterns
- **Best for**: Content scheduling, audience behavior
- **Use case**: "When should I post in this niche?"

### Keyword Analysis (`-a keywords`)
Identifies breakout keywords and opportunities.
- **Output**: High-velocity keywords, SEO opportunities
- **Best for**: SEO strategy, content planning
- **Use case**: "What keywords should I target?"

### Opportunity Score (`-a opportunity`)
AI-powered ranking of content opportunities.
- **Output**: Ranked opportunities with detailed reasoning
- **Best for**: Content strategy, opportunity discovery
- **Use case**: "What content should I create next?"

### Executive Summary (`-a executive`)
Comprehensive business intelligence report.
- **Output**: Market overview, strategic recommendations
- **Best for**: Business planning, investor presentations
- **Use case**: "What's the market landscape?"

### All Analyses (`-a all`)
Runs all analysis types in sequence.
- **Output**: Complete analysis suite
- **Best for**: Comprehensive research, deep dives
- **Use case**: "I need to understand this market completely"

## Analysis Levels

### Quick Scan (`-l quick`)
- **Videos**: ~150 videos
- **Time**: 30-60 seconds
- **API Units**: ~300
- **Best for**: Quick exploration, demos, trend validation
- **Strategy**: Prioritizes your chosen order, adapts automatically

### Balanced (`-l balanced`)
- **Videos**: ~400 videos
- **Time**: 1-2 minutes
- **API Units**: ~800
- **Best for**: Regular content creators, marketers
- **Strategy**: Deep dive into preferences, intelligent supplementation

### Deep Dive (`-l deep`)
- **Videos**: ~750 videos
- **Time**: 2-3 minutes
- **API Units**: ~1500
- **Best for**: Research, competitive analysis, content strategy
- **Strategy**: Exhaustive search across multiple orders

## Interactive Mode

Run without parameters to enter interactive mode:
```bash
ytminer
```

### Interactive Features
- **Step-by-step configuration**: Guided setup
- **Preview before analysis**: See what will be analyzed
- **Transcript options**: Choose whether to fetch transcripts
- **Analysis selection**: Pick specific analyses to run
- **Real-time feedback**: Progress indicators and status updates

## Weight Profiles

### Exploration (`--profile exploration`)
Discover new niches and emerging trends.
- **VPD**: 0.30, **Like**: 0.20, **Fresh**: 0.35, **Sat**: 0.10, **Slope**: 0.25
- **Best for**: Finding new opportunities, emerging trends

### Evergreen (`--profile evergreen`)
Focus on timeless, high-quality content.
- **VPD**: 0.25, **Like**: 0.40, **Fresh**: 0.10, **Sat**: 0.20, **Slope**: 0.05
- **Best for**: Long-term content strategy, quality focus

### Trending (`--profile trending`)
Catch viral content and momentum.
- **VPD**: 0.50, **Like**: 0.15, **Fresh**: 0.20, **Sat**: 0.05, **Slope**: 0.30
- **Best for**: Viral content, trending topics

### Balanced (`--profile balanced`)
Default balanced approach.
- **VPD**: 0.45, **Like**: 0.25, **Fresh**: 0.20, **Sat**: 0.30, **Slope**: 0.15
- **Best for**: General use, balanced analysis

## Environment Variables

### Required
- `YOUTUBE_API_KEY` - Your YouTube Data API v3 key

### Optional
- `YTMINER_DEFAULT_REGION` - Default search region
- `YTMINER_DEFAULT_DURATION` - Default video duration
- `YTMINER_DEFAULT_TIME_RANGE` - Default time range
- `YTMINER_DEFAULT_ORDER` - Default search order
- `YTMINER_WITH_TRANSCRIPTS` - Enable transcript fetching by default
- `YTMINER_TRANSCRIPT_LANGS` - Preferred transcript languages (comma-separated)
- `YTMINER_CACHE_DIR` - Cache directory for transcripts

### Opportunity Score Weights
- `YTMINER_OPP_W_VPD` - VPD weight
- `YTMINER_OPP_W_LIKE` - Like rate weight
- `YTMINER_OPP_W_FRESH` - Freshness weight
- `YTMINER_OPP_W_SAT` - Saturation penalty weight
- `YTMINER_OPP_W_SLOPE` - Slope weight

## Examples

### Quick Exploration
```bash
# Quick trend check
ytminer -k "ai tools" -l quick

# Recent videos only
ytminer -k "python tutorial" -t 7d -l quick
```

### Balanced Analysis
```bash
# Standard analysis
ytminer -k "meditation" -l balanced

# With region filter
ytminer -k "cooking" -r BR -l balanced

# Specific analysis type
ytminer -k "fitness" -a competitors -l balanced
```

### Deep Research
```bash
# Comprehensive analysis
ytminer -k "machine learning" -l deep -a all

# Executive summary
ytminer -k "content marketing" -l deep -a executive

# Opportunity discovery
ytminer -k "programming" -l deep -a opportunity --profile trending
```

### Advanced Configuration
```bash
# Custom weights
ytminer -k "gaming" --opp-w-vpd 0.6 --opp-w-like 0.1 -a opportunity

# Show profiles
ytminer --profiles

# Apply profile
ytminer -k "tech" --profile exploration -a opportunity
```

## Tips and Best Practices

### Keyword Selection
- Use specific, targeted keywords
- Avoid overly broad terms
- Consider niche-specific terminology
- Test different keyword variations

### Analysis Level Selection
- **Quick**: For exploration and demos
- **Balanced**: For regular content planning
- **Deep**: For comprehensive research

### Region and Time Filters
- Use region filters for local markets
- Time filters help find recent trends
- Combine filters for targeted analysis

### Opportunity Score
- Experiment with different profiles
- Customize weights for your strategy
- Focus on high-scoring opportunities
- Consider the reasoning provided

### Transcript Analysis
- Enable for content-rich analysis
- Be aware of YouTube restrictions
- Use language preferences
- Cache helps with repeated analysis
