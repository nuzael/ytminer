# YTMiner Metrics and Methodology

This document explains how YTMiner computes each analysis and what parameters influence the results. Our goal is to provide transparency so users understand where insights come from and what sample sizes were used.

## Global Inputs and Filters
- Keyword: the search topic used to retrieve videos.
- Region: ISO2 country code or `any`.
- Duration: `short` (<4m), `medium` (4–20m), `long` (>20m), or `any`.
- Time Range: `any`, `7d`, `30d`, `90d`, `1y`.
- Order: `relevance` (default), `date`, `viewCount`, `rating`, `title`.
- Analysis Level: controls breadth/depth (number of searches/videos sampled).

All analyses operate on the same filtered set of videos retrieved from the YouTube Data API v3 using the parameters above. Each report header shows the parameters and the total sample size (N videos).

Sampling note: for Balanced/Deep Dive levels, the first search honors the user-selected `order/region/duration`, then additional searches introduce diversity (e.g., different orders/regions/durations) while deduplicating videos.

## Data Cleaning and Normalization
- Unicode normalization (NFKC-like) with accent removal before tokenization.
- Tokenization removes punctuation/symbols and collapses whitespace.
- Stopwords: small bilingual list (English and Portuguese) to remove common function words.
- Emoji extraction: Unicode grapheme cluster aware (ZWJ sequences, flags), not just non-ASCII characters.
- Outlier protection: engagement calculations guard against division-by-zero.

## Growth Pattern Analysis
- Metrics per video:
  - Engagement% = (likes + comments) / max(views, 1) × 100.
- Aggregation across videos (N):
  - Avg Views = sum(views)/N
  - Avg Likes = sum(likes)/N
  - Avg Comments = sum(comments)/N
- Top Performers: ranked by per-video Engagement%, top 5.
- Growth Trend (rate): simple linear regression slope of views vs. publish order (time proxy), expressed as % of the first video’s views: slope / firstViews × 100.
  - Interpretation: positive = upward trend; negative = downward.

## Title Pattern Analysis
- Words: tokenized from normalized titles, stopwords removed; ranked by frequency (top 10).
- Phrases: bigrams from tokenized words; ranked by frequency (top 5).
- Emojis: extracted via grapheme clusters; ranked by frequency (top 5).
- Patterns: simple rule matches (e.g., contains "tutorial", "how to", year markers like 2023/2024).

## Competitor Analysis
- For each channel within the sample:
  - VideoCount, TotalViews, AvgViews, Engagement% (avg of per-video engagement).
- Ranking by TotalViews; top channels displayed.
- Market Share: channel.TotalViews / sum(all channels’ views) × 100 within the sampled set.
- Note: The same global filters apply (time range, region, duration, order). Interpret as sample-level market share, not global.

## Temporal Analysis
- For each hour (0–23) and weekday:
  - Sum views, likes, and engagement across matching videos; track counts.
  - Compute averages per bucket: AvgViews, AvgLikes, Engagement%.
- Minimum sample size per bucket: N ≥ 5 (buckets below threshold are hidden).
- Rankings: by Engagement% descending.

## Keyword Analysis
- Tokenized keywords from titles (normalized, stopwords removed).
- For each keyword:
  - Frequency = count of videos where keyword appears.
  - AvgViews = mean of views across those videos.
  - Engagement% = mean of per-video engagement across those videos.
- Trending Keywords: top by Frequency (top 10).
- Long-tail Candidates: Frequency < 3 and Engagement% > 5%; sorted by Engagement%.

## Executive Report
- Synthesizes insights from all analyses, summarizing:
  - Key Insights, Recommendations, Content Strategy, Competitive Intel, Performance Benchmarks, Next Steps.
- Each section is derived from the metrics described above and the current sample.

## Sample Size Reporting
- Each analysis prints N = number of videos considered.
- Temporal buckets also imply per-bucket N with a minimum threshold; consider increasing the Analysis Level and time window to improve per-bucket coverage.

## Limitations
- API quota and sampling strategy may bias results; consider higher "Analysis Level" for broader coverage.
- Growth slope uses publish order as a time proxy; for precise time series, a full time-bucketed trend should be computed.
- Keyword stopword lists are minimal; specialized domains may require custom lists.

## Reproducibility
- The CLI and interactive flows document the filters used (including `order`) at the top of each report or preview.
- Results depend on current YouTube search results and may vary over time. 