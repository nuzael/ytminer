# Analysis Types Deep Dive

Detailed explanation of each analysis type and their methodologies.

## Growth Analysis (`-a growth`)

### Purpose
Identifies trending videos and momentum patterns to discover what's gaining traction.

### Methodology
- **Primary Metric**: Views Per Day (VPD) - normalizes views by video age
- **Momentum Detection**: Identifies videos with accelerating view rates
- **Trending Signals**: Flags content with unusual velocity patterns

### Output
- Top videos ranked by VPD
- Momentum insights and early signals
- Trending topic identification
- Velocity distribution analysis

### Use Cases
- Content creators: "What's trending in my niche?"
- Marketers: "What content is gaining momentum?"
- Researchers: "What topics are emerging?"

### Example Output
```
Ì∫Ä Growth Analysis - Top Trending Videos
=====================================

#1  "AI Tools 2024" - 15.2k VPD, 2.3x momentum
#2  "Python Tutorial" - 12.8k VPD, 1.8x momentum
#3  "Machine Learning" - 10.5k VPD, 1.5x momentum
```

## Title Analysis (`-a titles`)

### Purpose
Analyzes high-performing title patterns to understand what makes titles successful.

### Methodology
- **Pattern Recognition**: Identifies common structures in successful titles
- **Keyword Analysis**: Extracts high-performing keywords and phrases
- **Length Analysis**: Correlates title length with performance
- **Emoji Usage**: Analyzes emoji impact on engagement

### Output
- High-performing title patterns
- Keyword frequency and performance
- Optimal title length recommendations
- Emoji usage insights

### Use Cases
- Content creators: "What titles work best?"
- SEO specialists: "What keywords to target?"
- A/B testing: "What title formats perform?"

### Example Output
```
Ì≥ù Title Analysis - High-Performing Patterns
==========================================

Top Patterns:
- "How to [Action] in [Time]" - 85% above average
- "[Number] [Adjective] [Nouns]" - 72% above average
- "Why [Controversial Topic]" - 68% above average

Top Keywords:
- "tutorial" - 2.3x average VPD
- "2024" - 1.8x average VPD
- "beginner" - 1.6x average VPD
```

## Competitor Analysis (`-a competitors`)

### Purpose
Finds rising stars and market opportunities by analyzing channel performance.

### Methodology
- **Rising Stars Detection**: Channels with VPD > niche baseline
- **Market Share Analysis**: Channel distribution and dominance
- **Growth Trajectory**: Channel velocity and momentum
- **Opportunity Gaps**: Underserved content areas

### Output
- Rising star channels
- Market share distribution
- Growth trajectory analysis
- Opportunity gap identification

### Use Cases
- Content creators: "Who are the rising stars?"
- Marketers: "What partnerships are available?"
- Business: "What market gaps exist?"

### Example Output
```
ÌøÜ Competitor Analysis - Rising Stars
===================================

Rising Stars:
#1  "TechTutorials" - 45k VPD, +180% growth
#2  "CodeMaster" - 32k VPD, +150% growth
#3  "DevTips" - 28k VPD, +120% growth

Market Share:
- Top 5 channels: 35% of total views
- Rising stars: 15% of total views
- Long tail: 50% of total views
```

## Temporal Analysis (`-a temporal`)

### Purpose
Analyzes best posting times and patterns to optimize content scheduling.

### Methodology
- **Time-of-Day Analysis**: Hourly engagement patterns
- **Day-of-Week Analysis**: Weekly engagement patterns
- **Seasonal Patterns**: Monthly and seasonal trends
- **Optimal Windows**: Best posting time recommendations

### Output
- Optimal posting windows
- Engagement patterns by time
- Seasonal trend analysis
- Scheduling recommendations

### Use Cases
- Content creators: "When should I post?"
- Social media managers: "What's the best schedule?"
- Marketers: "When is my audience most active?"

### Example Output
```
‚è∞ Temporal Analysis - Optimal Posting Times
==========================================

Best Posting Windows:
- Tuesday 2-4 PM: 1.8x average engagement
- Thursday 10-12 AM: 1.6x average engagement
- Saturday 6-8 PM: 1.4x average engagement

Weekly Patterns:
- Tuesday: Highest engagement
- Thursday: Second highest
- Sunday: Lowest engagement
```

## Keyword Analysis (`-a keywords`)

### Purpose
Identifies breakout keywords and SEO opportunities for content strategy.

### Methodology
- **Keyword Velocity**: Keywords with high VPD growth
- **Frequency Analysis**: Most common keywords in successful content
- **Opportunity Scoring**: Underserved keyword combinations
- **Trend Detection**: Emerging keyword patterns

### Output
- High-velocity keywords
- SEO opportunity keywords
- Keyword trend analysis
- Content gap identification

### Use Cases
- SEO specialists: "What keywords to target?"
- Content creators: "What topics are trending?"
- Marketers: "What content gaps exist?"

### Example Output
```
Ì¥ç Keyword Analysis - SEO Opportunities
=====================================

High-Velocity Keywords:
- "ai tools" - 2.5x average VPD
- "python tutorial" - 2.1x average VPD
- "machine learning" - 1.8x average VPD

SEO Opportunities:
- "beginner python" - Low competition, high potential
- "ai for beginners" - Emerging trend, high growth
- "python projects" - High search volume, moderate competition
```

## Opportunity Score (`-a opportunity`)

### Purpose
AI-powered ranking of content opportunities using multiple metrics.

### Methodology
- **Multi-Metric Scoring**: Combines VPD, engagement, freshness, saturation
- **Z-Score Normalization**: Standardizes metrics for fair comparison
- **Weighted Combination**: Configurable weights for different strategies
- **Reasoning Generation**: Explains why each opportunity scores high

### Output
- Ranked opportunity list
- Detailed scoring breakdown
- Reasoning for each opportunity
- Strategic recommendations

### Use Cases
- Content creators: "What should I create next?"
- Content strategists: "What opportunities exist?"
- Marketers: "What content will perform?"

### Example Output
```
ÌæØ Opportunity Score - Top Content Opportunities
==============================================

#1  "AI Tools Tutorial" - Score: 8.7
    VPD: 15.2k, Like Rate: 45/1k, Age: 3d, Saturation: 0.2
    Why: High velocity, fresh content, low saturation

#2  "Python for Beginners" - Score: 8.1
    VPD: 12.8k, Like Rate: 52/1k, Age: 5d, Saturation: 0.3
    Why: Strong engagement, growing trend, moderate saturation
```

## Executive Summary (`-a executive`)

### Purpose
Comprehensive business intelligence report for strategic decision-making.

### Methodology
- **Market Overview**: Complete landscape analysis
- **Competitive Intelligence**: Key players and trends
- **Opportunity Assessment**: Strategic recommendations
- **Performance Benchmarks**: Industry standards and metrics

### Output
- Executive summary
- Market landscape overview
- Competitive intelligence
- Strategic recommendations
- Performance benchmarks

### Use Cases
- Executives: "What's the market landscape?"
- Investors: "What opportunities exist?"
- Strategists: "What should we focus on?"

### Example Output
```
Ì≥ä Executive Summary - Market Intelligence
========================================

Market Overview:
- Total addressable market: 2.3M videos
- Average VPD: 1.2k
- Top 10% performers: 8.5k+ VPD
- Growth rate: +15% month-over-month

Key Insights:
- AI/ML content showing 200% growth
- Tutorial content dominates engagement
- Short-form content gaining traction
- Mobile-first content performing best

Strategic Recommendations:
1. Focus on AI/ML tutorial content
2. Optimize for mobile viewing
3. Target beginner-friendly content
4. Leverage trending keywords
```

## All Analyses (`-a all`)

### Purpose
Runs all analysis types in sequence for comprehensive research.

### Methodology
- **Sequential Execution**: Runs each analysis type in order
- **Comprehensive Coverage**: Complete market analysis
- **Integrated Insights**: Cross-analysis pattern recognition
- **Strategic Overview**: Holistic market understanding

### Output
- Complete analysis suite
- Integrated insights
- Cross-analysis patterns
- Comprehensive recommendations

### Use Cases
- Research: "I need to understand this market completely"
- Strategy: "What's the complete picture?"
- Planning: "What should our content strategy be?"

### Example Output
```
Ì¥ç Complete Analysis Suite
========================

Growth Analysis: ‚úÖ Complete
Title Analysis: ‚úÖ Complete
Competitor Analysis: ‚úÖ Complete
Temporal Analysis: ‚úÖ Complete
Keyword Analysis: ‚úÖ Complete
Opportunity Score: ‚úÖ Complete
Executive Summary: ‚úÖ Complete

Integrated Insights:
- Market is growing 15% month-over-month
- AI/ML content is the biggest opportunity
- Tutorial format performs best
- Tuesday 2-4 PM is optimal posting time
- "beginner" keywords have highest potential
```
