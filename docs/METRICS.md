# YTMiner Metrics and Methodology

This document explains how YTMiner computes each analysis. Our methodology is centered around **View Velocity** to provide insights into current momentum, not just historical popularity.

## Global Inputs and Filters
- Keyword, Region, Duration, Time Range, Order, Analysis Level.
- Defaults for interactive mode are configurable via env or Settings UI:
  - `YTMINER_DEFAULT_REGION`, `YTMINER_DEFAULT_DURATION`, `YTMINER_DEFAULT_TIME_RANGE`, `YTMINER_DEFAULT_ORDER`.
- All analyses operate on the same filtered set of videos. The report header shows the parameters and the total sample size (N videos).

## Data Cleaning and Normalization
- Unicode normalization, tokenization, stopword removal, and emoji extraction are performed on titles.
- When transcripts are available, the same cleaning steps are applied to transcript text (title + description + transcript).

## Transcript Usage
- Transcripts are used to better understand topical context and improve keyword/topic analyses.
- They do not change the core velocity metric (VPD), but enable richer cluster/topic grouping and keyword ranking.
- Transcript availability varies by video and language; auto-generated captions may be used with lower reliability.

## Core Analytical Metric: View Velocity (VPD)
The cornerstone of our analysis is **Views Per Day (VPD)**. It normalizes total views by the video's age, allowing us to compare the momentum of new and older videos on a level playing field.

- **VPD** = `video.viewCount / max(1, days_since_publication)`
- `days_since_publication` is the integer number of days from the video's publication date to today.
- This metric is calculated for every video in the sample and used as a basis for most analyses.

## Advanced Velocity Metrics: VPD7, VPD30, and Slope

### VPD7 and VPD30 (Windowed Velocity)
To better understand recent momentum vs. overall performance, we calculate velocity over specific time windows:

- **VPD7** = `video.viewCount / min(7, days_since_publication)`
- **VPD30** = `video.viewCount / min(30, days_since_publication)`

These metrics approximate "recent velocity" by capping the denominator at 7 or 30 days respectively. This helps identify:
- Videos that gained momentum recently (high VPD7 vs. VPD30)
- Content that peaked early but has slowed down
- Steady performers vs. viral spikes

### Slope (Acceleration Indicator)
The **Slope** metric compares short-term vs. long-term velocity to detect acceleration or deceleration:

- **Slope** = `(VPD7 - VPD30) / max(VPD30, 1)`
- Values are clamped to [-5, 5] to limit outliers
- **Positive Slope**: Content is accelerating (recent momentum > overall momentum)
- **Negative Slope**: Content is decelerating (recent momentum < overall momentum)
- **Zero Slope**: Steady performance

**Examples:**
- Slope = 2.0: VPD7 is 3x higher than VPD30 (strong acceleration)
- Slope = -0.5: VPD7 is 50% lower than VPD30 (deceleration)
- Slope = 0.0: VPD7 equals VPD30 (steady state)

## Growth Pattern Analysis
- **Niche Velocity Score (Avg. VPD)**: The primary health indicator of a niche.
  - `Avg. VPD = sum(all_video_VPDs) / N`
- **Highest Velocity Videos (Trending Now)**: Ranked by **VPD** descending, not by total views or engagement. Highlights what's trending *now*.
- Standard aggregations (Avg Views, Likes, Comments) are still provided for context.
- The previous "Growth Trend" (linear regression) has been deprecated in favor of the more robust Avg. VPD.

## Title Pattern Analysis
- Methodology remains the same: word/phrase frequency, emoji counts, and pattern matching. This analysis is independent of velocity.

## Competitor Analysis
- For each channel in the sample, we calculate both traditional and velocity-based metrics:
  - Traditional: VideoCount, TotalViews, AvgViews, Market Share %.
  - **Velocity Metrics**:
    - **Channel Average VPD**: The mean VPD of all videos from a specific channel in the sample. Measures the channel's current "hotness".
- **Rising Stars**: New category to identify emerging competitors. A channel is flagged as "Rising Star" if its `Channel Average VPD` is significantly higher than the overall `Niche Average VPD`.
  - Default threshold: `YTMINER_RISING_STAR_MULTIPLIER=1.5` (configurable via env)
- Ranking: The main competitor list is ranked by `TotalViews` to show established leaders, but the report highlights channels with highest `Channel Average VPD` and identified "Rising Stars".

## Temporal Analysis
- Methodology remains the same: aggregates views and engagement by hour and day of the week. This analysis is independent of velocity.

## Keyword Analysis
- **Trending Keywords (redefined)**: This analysis now identifies keywords driving the most momentum.
  - For each keyword (token), we calculate the **Average VPD** of all videos in the sample containing that keyword.
  - Keywords are ranked by this `Average VPD` to reveal which specific topics are currently "hot" and generating high-velocity views.
- **Core Keywords**: The previous frequency-based ranking is available as "Core Keywords", identifying the most common terms in the niche.
- **Long-tail Candidates**: Low frequency but high `Engagement%`.
  - Default: `Frequency <= YTMINER_LONG_TAIL_MAX_FREQ` (default 2) AND `Avg Engagement > YTMINER_LONG_TAIL_MIN_ENGAGEMENT` (default 5%).
  - Both are configurable via env.

## Opportunity Score (Enhanced)
The Opportunity Score is a comprehensive ranking computed in-memory over the current result set. It combines multiple signals using z-score normalization:

**Formula:**
```
score = w_vpd*z(VPD) + w_like*z(LikeRate) + w_fresh*Freshness - w_sat*Saturation + w_slope*z(Slope)
```

**Components:**
- **z(VPD)**: Standardized velocity (views per day)
- **z(LikeRate)**: Standardized engagement (likes per 1k views)
- **Freshness**: Min-max normalized age (inverted: newer = higher score)
- **Saturation**: Penalty based on primary keyword frequency in sample
- **z(Slope)**: Standardized acceleration (VPD7 vs VPD30 comparison)

**Default Weights (configurable via env):**
- `YTMINER_OPP_W_VPD=0.45` (velocity)
- `YTMINER_OPP_W_LIKE=0.25` (engagement)
- `YTMINER_OPP_W_FRESH=0.20` (freshness)
- `YTMINER_OPP_W_SAT=0.30` (saturation penalty)
- `YTMINER_OPP_W_SLOPE=0.15` (acceleration)

**Interpretation:**
- Higher scores indicate better opportunities
- Slope component rewards accelerating content
- Saturation penalty reduces overused themes
- All components are z-score normalized for fair comparison

**Notes:**
- No persistence required; operates on the sample collected for the current search.
- Transcripts are not required; future versions may incorporate topic clustering.
- Weights can be adjusted via environment variables or CLI flags.

## Engagement Analysis
- **Engagement Rate**: `(Likes + Comments) / max(1, Views) * 100`
- Channel engagement is calculated as the average engagement across all videos from that channel in the sample.
- Engagement formatting adjusts precision based on value range for better readability.
- Engagement categories: Excellent (>10%), Very Good (5-10%), Good (2-5%), Average (1-2%), Low (<1%).

## Executive Report
- Synthesizes insights from all analyses, with strong emphasis on velocity and momentum.
- Highlights the `Niche Velocity Score`, identifies "Breakout" topics from Keyword Analysis, and flags "Rising Star" competitors.
- Includes direct links to Rising Star channels for immediate competitive analysis.
