# Configuration Guide

Complete guide to configuring YTMiner for your specific needs.

## Environment Variables

### Required Configuration

#### YouTube API Key
```bash
YOUTUBE_API_KEY=your_api_key_here
```
- **Required**: Yes
- **Description**: Your YouTube Data API v3 key
- **How to get**: [YouTube Data API v3 Setup Guide](https://developers.google.com/youtube/v3/getting-started)

### Optional Configuration

#### Default Search Parameters
```bash
YTMINER_DEFAULT_REGION=any
YTMINER_DEFAULT_DURATION=any
YTMINER_DEFAULT_TIME_RANGE=any
YTMINER_DEFAULT_ORDER=relevance
```
- **Description**: Default values for search parameters
- **Values**: See [Usage Guide](USAGE.md) for available options

#### Transcript Configuration
```bash
YTMINER_WITH_TRANSCRIPTS=false
YTMINER_TRANSCRIPT_LANGS=en,pt,es
YTMINER_CACHE_DIR=.cache
```
- **YTMINER_WITH_TRANSCRIPTS**: Enable transcript fetching by default
- **YTMINER_TRANSCRIPT_LANGS**: Preferred languages (comma-separated)
- **YTMINER_CACHE_DIR**: Cache directory for transcripts

#### Opportunity Score Weights
```bash
YTMINER_OPP_W_VPD=0.45
YTMINER_OPP_W_LIKE=0.25
YTMINER_OPP_W_FRESH=0.20
YTMINER_OPP_W_SAT=0.30
YTMINER_OPP_W_SLOPE=0.15
```
- **Description**: Default weights for Opportunity Score calculation
- **Range**: 0.0 to 1.0 (or higher for emphasis)
- **Note**: Weights are not automatically normalized

## Configuration Files

### .env File
Create a `.env` file in the project root:
```bash
# Copy from env.example
cp env.example .env

# Edit with your values
nano .env
```

### Example .env
```bash
# Required
YOUTUBE_API_KEY=AIzaSyBxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# Optional - Search defaults
YTMINER_DEFAULT_REGION=BR
YTMINER_DEFAULT_DURATION=medium
YTMINER_DEFAULT_TIME_RANGE=30d
YTMINER_DEFAULT_ORDER=viewCount

# Optional - Transcripts
YTMINER_WITH_TRANSCRIPTS=true
YTMINER_TRANSCRIPT_LANGS=en,pt,es
YTMINER_CACHE_DIR=.cache

# Optional - Opportunity Score weights
YTMINER_OPP_W_VPD=0.50
YTMINER_OPP_W_LIKE=0.20
YTMINER_OPP_W_FRESH=0.25
YTMINER_OPP_W_SAT=0.20
YTMINER_OPP_W_SLOPE=0.25
```

## Weight Profiles

### Predefined Profiles
YTMiner comes with four predefined weight profiles:

#### Exploration Profile
```bash
ytminer --profile exploration
```
- **VPD**: 0.30 (moderate velocity focus)
- **Like**: 0.20 (low engagement focus)
- **Fresh**: 0.35 (high freshness focus)
- **Sat**: 0.10 (low saturation penalty)
- **Slope**: 0.25 (moderate slope focus)
- **Best for**: Discovering new niches and emerging trends

#### Evergreen Profile
```bash
ytminer --profile evergreen
```
- **VPD**: 0.25 (low velocity focus)
- **Like**: 0.40 (high engagement focus)
- **Fresh**: 0.10 (low freshness focus)
- **Sat**: 0.20 (moderate saturation penalty)
- **Slope**: 0.05 (low slope focus)
- **Best for**: Long-term content strategy, quality focus

#### Trending Profile
```bash
ytminer --profile trending
```
- **VPD**: 0.50 (high velocity focus)
- **Like**: 0.15 (low engagement focus)
- **Fresh**: 0.20 (moderate freshness focus)
- **Sat**: 0.05 (low saturation penalty)
- **Slope**: 0.30 (high slope focus)
- **Best for**: Viral content, trending topics

#### Balanced Profile
```bash
ytminer --profile balanced
```
- **VPD**: 0.45 (moderate-high velocity focus)
- **Like**: 0.25 (moderate engagement focus)
- **Fresh**: 0.20 (moderate freshness focus)
- **Sat**: 0.30 (moderate saturation penalty)
- **Slope**: 0.15 (moderate slope focus)
- **Best for**: General use, balanced analysis

### Custom Weights
You can set custom weights via CLI flags or environment variables:

#### CLI Flags
```bash
ytminer -k "keyword" --opp-w-vpd 0.6 --opp-w-like 0.1 -a opportunity
```

#### Environment Variables
```bash
export YTMINER_OPP_W_VPD=0.6
export YTMINER_OPP_W_LIKE=0.1
ytminer -k "keyword" -a opportunity
```

## Cache Configuration

### Transcript Cache
YTMiner caches transcripts locally to improve performance and reduce API calls.

#### Cache Directory
```bash
YTMINER_CACHE_DIR=.cache
```
- **Default**: `.cache/transcripts/`
- **Structure**: `{videoID}.{language}.txt`
- **Benefits**: Faster repeated analysis, reduced API usage

#### Cache Management
```bash
# Clear cache
rm -rf .cache/transcripts/

# Check cache size
du -sh .cache/transcripts/

# List cached videos
ls .cache/transcripts/
```

## API Quota Management

### YouTube Data API v3 Quotas
- **Search**: 100 units per request
- **Video details**: 1 unit per request
- **Transcripts**: No quota (but rate limited)

### Quota Optimization
- **Use appropriate analysis levels**: Quick vs Deep
- **Cache transcripts**: Avoid repeated fetching
- **Filter by time range**: Reduce dataset size
- **Use specific keywords**: More targeted results

### Quota Monitoring
```bash
# Check your quota usage
# Visit: https://console.developers.google.com/apis/api/youtube.googleapis.com/quotas
```

## Performance Tuning

### Analysis Level Selection
- **Quick**: ~300 API units, 30-60s
- **Balanced**: ~800 API units, 1-2min
- **Deep**: ~1500 API units, 2-3min

### Memory Usage
- **Videos in memory**: ~1MB per 1000 videos
- **Transcripts**: ~100KB per video
- **Cache**: ~1MB per 1000 transcripts

### Network Optimization
- **Concurrent requests**: Limited by YouTube API
- **Timeout settings**: 60s for transcripts
- **Retry logic**: Built-in for failed requests

## Troubleshooting

### Common Issues

#### API Key Issues
```bash
# Check API key
echo $YOUTUBE_API_KEY

# Test API key
ytminer -k "test" -l quick
```

#### Transcript Issues
```bash
# Check transcript languages
echo $YTMINER_TRANSCRIPT_LANGS

# Clear transcript cache
rm -rf .cache/transcripts/

# Test transcript fetching
ytminer -k "test" --with-transcripts
```

#### Performance Issues
```bash
# Use quick analysis level
ytminer -k "keyword" -l quick

# Filter by time range
ytminer -k "keyword" -t 7d

# Use specific region
ytminer -k "keyword" -r US
```

### Debug Mode
```bash
# Enable verbose logging
export YTMINER_DEBUG=true
ytminer -k "keyword" -l quick
```

### Log Files
```bash
# Check logs
tail -f ytminer.log

# Clear logs
rm ytminer.log
```

## Advanced Configuration

### Custom Analysis Weights
Create custom weight profiles by modifying the source code:

```go
// In config/profiles.go
"custom": {
    Name:        "custom",
    Description: "My custom strategy",
    Weights: ProfileWeights{
        VPD:   0.40,
        Like:  0.30,
        Fresh: 0.20,
        Sat:   0.25,
        Slope: 0.20,
    },
},
```

### Custom Metrics
Add custom metrics by extending the domain layer:

```go
// In domain/metrics/metrics.go
func CustomMetric(video Video) float64 {
    // Your custom calculation
    return value
}
```

### Custom Analysis Types
Add new analysis types by extending the analysis layer:

```go
// In analysis/analyzer.go
func (a *Analyzer) AnalyzeCustom() CustomResult {
    // Your custom analysis
    return result
}
```
