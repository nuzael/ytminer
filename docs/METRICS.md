# YTMiner Metrics and Methodology

This document explains how YTMiner computes each analysis. Our methodology is centered around **View Velocity** to provide insights into current momentum, not just historical popularity.

## Global Inputs and Filters
- Keyword, Region, Duration, Time Range, Order, Analysis Level.
- All analyses operate on the same filtered set of videos. The report header shows the parameters and the total sample size (N videos).

## Data Cleaning and Normalization
- Unicode normalization, tokenization, stopword removal, and emoji extraction are performed on titles.

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
- **Rising Stars**: New category to identify emerging competitors. A channel is flagged as "Rising Star" if its `Channel Average VPD` is significantly higher than the overall `Niche Average VPD` (e.g., > 1.5x).
- Ranking: The main competitor list is ranked by `TotalViews` to show established leaders, but the report highlights channels with highest `Channel Average VPD` and identified "Rising Stars".

## Temporal Analysis
- Methodology remains the same: aggregates views and engagement by hour and day of the week. This analysis is independent of velocity.

## Keyword Analysis
- **Trending Keywords (redefined)**: This analysis now identifies keywords driving the most momentum.
  - For each keyword (token), we calculate the **Average VPD** of all videos in the sample containing that keyword.
  - Keywords are ranked by this `Average VPD` to reveal which specific topics are currently "hot" and generating high-velocity views.
- **Core Keywords**: The previous frequency-based ranking is available as "Core Keywords", identifying the most common terms in the niche.
- **Long-tail Candidates**: The definition remains: low frequency but high `Engagement%`.

## Engagement Analysis
- **Engagement Rate**: `(Likes + Comments) / max(1, Views) * 100`
- Channel engagement is calculated as the average engagement across all videos from that channel in the sample.
- Engagement formatting adjusts precision based on value range for better readability.
- Engagement categories: Excellent (>10%), Very Good (5-10%), Good (2-5%), Average (1-2%), Low (<1%).

## Executive Report
- Synthesizes insights from all analyses, with strong emphasis on velocity and momentum.
- Highlights the `Niche Velocity Score`, identifies "Breakout" topics from Keyword Analysis, and flags "Rising Star" competitors.
- Includes direct links to Rising Star channels for immediate competitive analysis.