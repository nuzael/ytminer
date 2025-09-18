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

## Opportunity Score (Lightweight)
The Opportunity Score is a simple, quota-friendly ranking computed in-memory over the current result set. It combines:
- z(VPD)
- z(Like rate per 1k views)
- Freshness (min–max of age, inverted so newer videos score higher)
- Saturation penalty (normalized frequency of the primary title token within the sample)

Formula (weights subject to change):
- `score = 0.45*z(VPD) + 0.25*z(like_rate) + 0.20*freshness - 0.30*saturation`

Notes:
- No persistence required; operates on the sample collected for the current search.
- Transcripts are not required; future versions may incorporate topic clustering.

## Engagement Analysis
- **Engagement Rate**: `(Likes + Comments) / max(1, Views) * 100`
- Channel engagement is calculated as the average engagement across all videos from that channel in the sample.
- Engagement formatting adjusts precision based on value range for better readability.
- Engagement categories: Excellent (>10%), Very Good (5-10%), Good (2-5%), Average (1-2%), Low (<1%).

## Executive Report
- Synthesizes insights from all analyses, with strong emphasis on velocity and momentum.
- Highlights the `Niche Velocity Score`, identifies "Breakout" topics from Keyword Analysis, and flags "Rising Star" competitors.
- Includes direct links to Rising Star channels for immediate competitive analysis.