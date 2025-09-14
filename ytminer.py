#!/usr/bin/env python3
"""
YTMiner - YouTube video analysis CLI tool
Search and analyze YouTube videos with advanced filters and growth pattern analysis
"""

import os
import sys
import json
import click
import pandas as pd
import re
from collections import Counter
from datetime import datetime, timedelta
from typing import List, Dict, Optional
from tabulate import tabulate
from colorama import init, Fore, Style
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

init()

try:
    from googleapiclient.discovery import build
    from googleapiclient.errors import HttpError
except ImportError:
    print(f"{Fore.RED}Error: Dependencies not installed. Run: pip install -r requirements.txt{Style.RESET_ALL}")
    sys.exit(1)


class YouTubeMiner:
    def __init__(self, api_key: str):
        self.api_key = api_key
        self.youtube = build('youtube', 'v3', developerKey=api_key)
        
    def search_videos(self, 
                     keyword: str,
                     max_results: int = 50,
                     order: str = 'relevance',
                     published_after: Optional[str] = None,
                     published_before: Optional[str] = None,
                     region_code: Optional[str] = None,
                     video_duration: Optional[str] = None) -> List[Dict]:
        try:
            search_params = {
                'part': 'snippet',
                'q': keyword,
                'type': 'video',
                'maxResults': min(max_results, 50),
                'order': order
            }
            
            if published_after:
                search_params['publishedAfter'] = published_after
            if published_before:
                search_params['publishedBefore'] = published_before
            if region_code:
                search_params['regionCode'] = region_code
            if video_duration:
                search_params['videoDuration'] = video_duration
                
            search_response = self.youtube.search().list(**search_params).execute()
            video_ids = [item['id']['videoId'] for item in search_response['items']]
            
            if not video_ids:
                return []
                
            videos_response = self.youtube.videos().list(
                part='snippet,statistics,contentDetails',
                id=','.join(video_ids)
            ).execute()
            
            videos = []
            for video in videos_response['items']:
                video_data = {
                    'id': video['id'],
                    'title': video['snippet']['title'],
                    'channel': video['snippet']['channelTitle'],
                    'published_at': video['snippet']['publishedAt'],
                    'description': video['snippet']['description'][:200] + '...' if len(video['snippet']['description']) > 200 else video['snippet']['description'],
                    'view_count': int(video['statistics'].get('viewCount', 0)),
                    'like_count': int(video['statistics'].get('likeCount', 0)),
                    'comment_count': int(video['statistics'].get('commentCount', 0)),
                    'duration': video['contentDetails']['duration'],
                    'url': f"https://www.youtube.com/watch?v={video['id']}"
                }
                videos.append(video_data)
                
            return videos
            
        except HttpError as e:
            print(f"{Fore.RED}YouTube API error: {e}{Style.RESET_ALL}")
            return []
        except Exception as e:
            print(f"{Fore.RED}Unexpected error: {e}{Style.RESET_ALL}")
            return []
    
    def analyze_growth_patterns(self, videos: List[Dict]) -> Dict:
        if not videos:
            return {}
            
        df = pd.DataFrame(videos)
        df['published_at'] = pd.to_datetime(df['published_at'])
        df['days_ago'] = (datetime.now() - df['published_at'].dt.tz_localize(None)).dt.days
        
        analysis = {
            'total_videos': len(videos),
            'avg_views': df['view_count'].mean(),
            'avg_likes': df['like_count'].mean(),
            'avg_engagement_rate': (df['like_count'] / df['view_count']).mean() * 100,
            'top_performers': df.nlargest(5, 'view_count')[['title', 'view_count', 'like_count', 'days_ago']].to_dict('records'),
            'recent_trends': df[df['days_ago'] <= 7].nlargest(3, 'view_count')[['title', 'view_count', 'days_ago']].to_dict('records'),
            'best_engagement': df.nlargest(3, 'like_count')[['title', 'like_count', 'view_count', 'days_ago']].to_dict('records')
        }
        
        return analysis

    def analyze_titles(self, videos: List[Dict]) -> Dict:
        """Analyze title patterns from successful videos"""
        if not videos:
            return {}
        
        # Filter top performers (top 20% by views)
        df = pd.DataFrame(videos)
        top_performers = df.nlargest(max(1, len(videos) // 5), 'view_count')
        
        titles = [video['title'] for video in top_performers.to_dict('records')]
        
        # Analyze title patterns
        analysis = {
            'total_titles_analyzed': len(titles),
            'avg_title_length': sum(len(title) for title in titles) / len(titles),
            'most_common_words': self._extract_common_words(titles),
            'most_common_phrases': self._extract_common_phrases(titles),
            'emoji_usage': self._analyze_emojis(titles),
            'title_patterns': self._analyze_title_patterns(titles),
            'successful_titles': titles[:5]  # Top 5 titles
        }
        
        return analysis
    
    def _extract_common_words(self, titles: List[str]) -> List[tuple]:
        """Extract most common words from titles"""
        # Clean and split titles
        all_words = []
        for title in titles:
            # Remove special characters and convert to lowercase
            clean_title = re.sub(r'[^\w\s]', ' ', title.lower())
            words = [word for word in clean_title.split() if len(word) > 2]
            all_words.extend(words)
        
        # Count and return top words
        word_count = Counter(all_words)
        return word_count.most_common(10)
    
    def _extract_common_phrases(self, titles: List[str]) -> List[tuple]:
        """Extract most common 2-word phrases from titles"""
        all_phrases = []
        for title in titles:
            clean_title = re.sub(r'[^\w\s]', ' ', title.lower())
            words = clean_title.split()
            # Create 2-word phrases
            for i in range(len(words) - 1):
                phrase = f"{words[i]} {words[i+1]}"
                if len(phrase) > 5:  # Filter short phrases
                    all_phrases.append(phrase)
        
        phrase_count = Counter(all_phrases)
        return phrase_count.most_common(8)
    
    def _analyze_emojis(self, titles: List[str]) -> Dict:
        """Analyze emoji usage in titles"""
        emoji_pattern = re.compile(r'[^\w\s]')
        emoji_count = 0
        emojis = []
        
        for title in titles:
            emojis_in_title = emoji_pattern.findall(title)
            emojis.extend(emojis_in_title)
            if emojis_in_title:
                emoji_count += 1
        
        return {
            'titles_with_emojis': emoji_count,
            'emoji_percentage': (emoji_count / len(titles)) * 100,
            'most_common_emojis': Counter(emojis).most_common(5)
        }
    
    def _analyze_title_patterns(self, titles: List[str]) -> Dict:
        """Analyze common title patterns"""
        patterns = {
            'tutorial_pattern': sum(1 for title in titles if any(word in title.lower() for word in ['tutorial', 'how to', 'como', 'aprenda', 'learn'])),
            'number_pattern': sum(1 for title in titles if re.search(r'\d+', title)),
            'question_pattern': sum(1 for title in titles if '?' in title),
            'exclamation_pattern': sum(1 for title in titles if '!' in title),
            'bracket_pattern': sum(1 for title in titles if '[' in title and ']' in title),
            'parentheses_pattern': sum(1 for title in titles if '(' in title and ')' in title)
        }
        
        # Convert to percentages
        total = len(titles)
        for pattern in patterns:
            patterns[pattern] = (patterns[pattern] / total) * 100 if total > 0 else 0
        
        return patterns

    def analyze_competitors(self, videos: List[Dict]) -> Dict:
        """Analyze competitor channels and market dominance"""
        if not videos:
            return {}
        
        df = pd.DataFrame(videos)
        
        # Channel analysis
        channel_stats = df.groupby('channel').agg({
            'view_count': ['count', 'sum', 'mean'],
            'like_count': ['sum', 'mean'],
            'title': 'count'
        }).round(2)
        
        # Flatten column names
        channel_stats.columns = ['video_count', 'total_views', 'avg_views', 'total_likes', 'avg_likes', 'title_count']
        channel_stats = channel_stats.reset_index()
        
        # Calculate engagement rate
        channel_stats['engagement_rate'] = (channel_stats['total_likes'] / channel_stats['total_views'] * 100).round(2)
        
        # Sort by total views
        channel_stats = channel_stats.sort_values('total_views', ascending=False)
        
        # Market share analysis
        total_views = channel_stats['total_views'].sum()
        channel_stats['market_share'] = (channel_stats['total_views'] / total_views * 100).round(2)
        
        # Identify top competitors
        top_competitors = channel_stats.head(10).to_dict('records')
        
        # Channel diversity analysis
        total_channels = len(channel_stats)
        top_5_share = channel_stats.head(5)['market_share'].sum()
        top_10_share = channel_stats.head(10)['market_share'].sum()
        
        # Content patterns by channel
        content_patterns = {}
        for channel in channel_stats.head(5)['channel']:
            channel_videos = df[df['channel'] == channel]
            if len(channel_videos) > 0:
                # Analyze title patterns for this channel
                titles = channel_videos['title'].tolist()
                patterns = {
                    'avg_title_length': sum(len(title) for title in titles) / len(titles),
                    'tutorial_ratio': sum(1 for title in titles if any(word in title.lower() for word in ['tutorial', 'how to', 'learn', 'course'])) / len(titles) * 100,
                    'number_ratio': sum(1 for title in titles if re.search(r'\d+', title)) / len(titles) * 100,
                    'emoji_ratio': sum(1 for title in titles if re.search(r'[^\w\s]', title)) / len(titles) * 100
                }
                content_patterns[channel] = patterns
        
        analysis = {
            'total_channels': total_channels,
            'market_concentration': {
                'top_5_share': top_5_share,
                'top_10_share': top_10_share,
                'concentration_level': 'High' if top_5_share > 60 else 'Medium' if top_5_share > 40 else 'Low'
            },
            'top_competitors': top_competitors,
            'content_patterns': content_patterns,
            'market_insights': self._generate_market_insights(channel_stats, total_views)
        }
        
        return analysis
    
    def _generate_market_insights(self, channel_stats, total_views):
        """Generate market insights and opportunities"""
        insights = []
        
        # Market concentration insights
        top_5_share = channel_stats.head(5)['market_share'].sum()
        if top_5_share > 70:
            insights.append("Market is highly concentrated - dominated by few channels")
        elif top_5_share > 50:
            insights.append("Market has moderate concentration - some dominant players")
        else:
            insights.append("Market is fragmented - many small players, opportunity for growth")
        
        # Engagement insights
        avg_engagement = channel_stats['engagement_rate'].mean()
        if avg_engagement > 5:
            insights.append(f"High engagement niche (avg {avg_engagement:.1f}%) - quality content matters")
        elif avg_engagement > 2:
            insights.append(f"Medium engagement niche (avg {avg_engagement:.1f}%) - room for improvement")
        else:
            insights.append(f"Low engagement niche (avg {avg_engagement:.1f}%) - focus on retention")
        
        # Content volume insights
        avg_videos_per_channel = channel_stats['video_count'].mean()
        if avg_videos_per_channel > 10:
            insights.append("High content volume - consistency is key")
        else:
            insights.append("Moderate content volume - opportunity for regular posting")
        
        # Market size insights
        if total_views > 100000000:  # 100M views
            insights.append("Large market - high competition but big rewards")
        elif total_views > 10000000:  # 10M views
            insights.append("Medium market - balanced competition and opportunity")
        else:
            insights.append("Smaller market - less competition, niche opportunity")
        
        return insights

    def analyze_temporal_patterns(self, videos: List[Dict]) -> Dict:
        """Analyze temporal patterns and optimal posting times"""
        if not videos:
            return {}
        
        df = pd.DataFrame(videos)
        df['published_at'] = pd.to_datetime(df['published_at'])
        
        # Extract time components
        df['hour'] = df['published_at'].dt.hour
        df['day_of_week'] = df['published_at'].dt.day_name()
        df['month'] = df['published_at'].dt.month
        df['year'] = df['published_at'].dt.year
        
        # Calculate days since publication
        df['days_ago'] = (datetime.now() - df['published_at'].dt.tz_localize(None)).dt.days
        
        # Calculate engagement rate first
        df['engagement_rate'] = (df['like_count'] / df['view_count'] * 100).round(2)
        
        # Hour analysis
        hour_performance = df.groupby('hour').agg({
            'view_count': ['count', 'mean', 'sum'],
            'like_count': ['mean', 'sum'],
            'engagement_rate': 'mean'
        }).round(2)
        
        # Flatten column names
        hour_performance.columns = ['video_count', 'avg_views', 'total_views', 'avg_likes', 'total_likes', 'engagement_rate']
        hour_performance = hour_performance.reset_index()
        
        # Day of week analysis
        day_performance = df.groupby('day_of_week').agg({
            'view_count': ['count', 'mean', 'sum'],
            'like_count': ['mean', 'sum']
        }).round(2)
        
        day_performance.columns = ['video_count', 'avg_views', 'total_views', 'avg_likes', 'total_likes']
        day_performance = day_performance.reset_index()
        
        # Calculate engagement rate for day analysis
        day_performance['engagement_rate'] = (day_performance['total_likes'] / day_performance['total_views'] * 100).round(2)
        
        # Sort days properly
        day_order = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday']
        day_performance['day_order'] = day_performance['day_of_week'].map({day: i for i, day in enumerate(day_order)})
        day_performance = day_performance.sort_values('day_order').drop('day_order', axis=1)
        
        # Recent vs older content analysis
        recent_videos = df[df['days_ago'] <= 30]  # Last 30 days
        older_videos = df[df['days_ago'] > 30]
        
        # Growth trends by time periods
        monthly_trends = df.groupby(['year', 'month']).agg({
            'view_count': ['count', 'mean', 'sum'],
            'like_count': ['mean', 'sum']
        }).round(2)
        
        monthly_trends.columns = ['video_count', 'avg_views', 'total_views', 'avg_likes', 'total_likes']
        monthly_trends = monthly_trends.reset_index()
        monthly_trends['engagement_rate'] = (monthly_trends['total_likes'] / monthly_trends['total_views'] * 100).round(2)
        
        # Find optimal posting times
        best_hours = hour_performance.nlargest(3, 'avg_views')[['hour', 'avg_views', 'engagement_rate']].to_dict('records')
        best_days = day_performance.nlargest(3, 'avg_views')[['day_of_week', 'avg_views', 'engagement_rate']].to_dict('records')
        
        # Generate temporal insights
        insights = self._generate_temporal_insights(hour_performance, day_performance, recent_videos, older_videos)
        
        analysis = {
            'hour_performance': hour_performance.to_dict('records'),
            'day_performance': day_performance.to_dict('records'),
            'monthly_trends': monthly_trends.to_dict('records'),
            'best_posting_times': {
                'best_hours': best_hours,
                'best_days': best_days
            },
            'recent_vs_older': {
                'recent_videos_count': len(recent_videos),
                'older_videos_count': len(older_videos),
                'recent_avg_views': recent_videos['view_count'].mean() if len(recent_videos) > 0 else 0,
                'older_avg_views': older_videos['view_count'].mean() if len(older_videos) > 0 else 0
            },
            'insights': insights
        }
        
        return analysis
    
    def _generate_temporal_insights(self, hour_performance, day_performance, recent_videos, older_videos):
        """Generate insights about temporal patterns"""
        insights = []
        
        # Best posting hour
        best_hour = hour_performance.loc[hour_performance['avg_views'].idxmax()]
        insights.append(f"Best posting hour: {int(best_hour['hour'])}:00 ({best_hour['avg_views']:,.0f} avg views)")
        
        # Best posting day
        best_day = day_performance.loc[day_performance['avg_views'].idxmax()]
        insights.append(f"Best posting day: {best_day['day_of_week']} ({best_day['avg_views']:,.0f} avg views)")
        
        # Peak hours analysis
        peak_hours = hour_performance[hour_performance['avg_views'] > hour_performance['avg_views'].mean() * 1.2]
        if len(peak_hours) > 0:
            peak_hour_list = [f"{int(h)}:00" for h in peak_hours['hour'].tolist()]
            insights.append(f"Peak performance hours: {', '.join(peak_hour_list)}")
        
        # Weekend vs weekday analysis
        weekday_avg = day_performance[day_performance['day_of_week'].isin(['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday'])]['avg_views'].mean()
        weekend_avg = day_performance[day_performance['day_of_week'].isin(['Saturday', 'Sunday'])]['avg_views'].mean()
        
        if weekend_avg > weekday_avg * 1.1:
            insights.append("Weekend content performs better - consider Saturday/Sunday posting")
        elif weekday_avg > weekend_avg * 1.1:
            insights.append("Weekday content performs better - focus on Monday-Friday")
        else:
            insights.append("No significant difference between weekday and weekend performance")
        
        # Recent vs older content
        if len(recent_videos) > 0 and len(older_videos) > 0:
            recent_avg = recent_videos['view_count'].mean()
            older_avg = older_videos['view_count'].mean()
            
            if recent_avg > older_avg * 1.2:
                insights.append("Recent content performs significantly better - focus on current trends")
            elif older_avg > recent_avg * 1.2:
                insights.append("Older content performs better - evergreen content works well")
            else:
                insights.append("Consistent performance across time periods")
        
        # Content volume insights
        total_videos = len(hour_performance)
        if total_videos < 10:
            insights.append("Limited data - analyze more videos for better insights")
        elif total_videos < 30:
            insights.append("Moderate data - patterns are emerging but more analysis needed")
        else:
            insights.append("Good data coverage - patterns are reliable")
        
        return insights

    def analyze_keywords(self, videos: List[Dict], search_keyword: str) -> Dict:
        """Analyze keywords and SEO opportunities"""
        if not videos:
            return {}
        
        df = pd.DataFrame(videos)
        
        # Extract all words from titles
        all_titles = ' '.join(df['title'].str.lower().tolist())
        
        # Clean and tokenize words
        import re
        words = re.findall(r'\b[a-zA-Z]{3,}\b', all_titles)
        
        # Count word frequency
        word_freq = Counter(words)
        
        # Filter out common stop words
        stop_words = {
            'the', 'and', 'for', 'are', 'but', 'not', 'you', 'all', 'can', 'had', 'her', 'was', 'one', 'our', 'out', 'day', 'get', 'has', 'him', 'his', 'how', 'its', 'may', 'new', 'now', 'old', 'see', 'two', 'who', 'boy', 'did', 'man', 'men', 'put', 'say', 'she', 'too', 'use', 'way', 'will', 'with', 'this', 'that', 'from', 'they', 'know', 'want', 'been', 'good', 'much', 'some', 'time', 'very', 'when', 'come', 'here', 'just', 'like', 'long', 'make', 'many', 'over', 'such', 'take', 'than', 'them', 'well', 'were', 'what', 'year', 'your', 'about', 'after', 'again', 'before', 'could', 'every', 'first', 'great', 'might', 'never', 'other', 'place', 'right', 'should', 'small', 'still', 'think', 'where', 'which', 'while', 'world', 'would', 'write', 'being', 'those', 'under', 'water', 'where', 'which', 'while', 'world', 'would', 'write'
        }
        
        # Filter meaningful words
        meaningful_words = {word: count for word, count in word_freq.items() 
                          if word not in stop_words and len(word) > 2}
        
        # Sort by frequency
        top_words = dict(sorted(meaningful_words.items(), key=lambda x: x[1], reverse=True)[:20])
        
        # Analyze keyword variations
        keyword_variations = self._find_keyword_variations(search_keyword, top_words)
        
        # Analyze long-tail keywords
        long_tail_keywords = self._extract_long_tail_keywords(df, search_keyword)
        
        # Analyze keyword performance
        keyword_performance = self._analyze_keyword_performance(df, top_words)
        
        # Generate SEO suggestions
        seo_suggestions = self._generate_seo_suggestions(search_keyword, top_words, keyword_variations)
        
        # Analyze trending keywords
        trending_keywords = self._analyze_trending_keywords(df)
        
        analysis = {
            'search_keyword': search_keyword,
            'top_keywords': top_words,
            'keyword_variations': keyword_variations,
            'long_tail_keywords': long_tail_keywords,
            'keyword_performance': keyword_performance,
            'seo_suggestions': seo_suggestions,
            'trending_keywords': trending_keywords,
            'total_words_analyzed': len(words),
            'unique_words': len(meaningful_words)
        }
        
        return analysis
    
    def _find_keyword_variations(self, search_keyword: str, top_words: Dict) -> List[Dict]:
        """Find keyword variations and related terms"""
        variations = []
        
        # Split search keyword
        keyword_parts = search_keyword.lower().split()
        
        # Find related words
        for word, count in top_words.items():
            if word not in keyword_parts:
                # Check if word is related to any part of the keyword
                for part in keyword_parts:
                    if len(part) > 3 and (word in part or part in word or 
                                         self._calculate_similarity(word, part) > 0.6):
                        variations.append({
                            'word': word,
                            'frequency': count,
                            'related_to': part,
                            'suggestion': f"{search_keyword} {word}"
                        })
                        break
        
        # Sort by frequency
        variations.sort(key=lambda x: x['frequency'], reverse=True)
        return variations[:10]
    
    def _extract_long_tail_keywords(self, df: pd.DataFrame, search_keyword: str) -> List[Dict]:
        """Extract long-tail keywords from titles"""
        long_tail = []
        
        for title in df['title'].tolist():
            title_lower = title.lower()
            if search_keyword.lower() in title_lower:
                # Extract phrases containing the keyword
                words = title_lower.split()
                keyword_index = -1
                for i, word in enumerate(words):
                    if search_keyword.lower() in word:
                        keyword_index = i
                        break
                
                if keyword_index >= 0:
                    # Extract 2-4 word phrases around the keyword
                    start = max(0, keyword_index - 1)
                    end = min(len(words), keyword_index + 3)
                    phrase = ' '.join(words[start:end])
                    
                    if len(phrase.split()) >= 2 and phrase not in [lt['phrase'] for lt in long_tail]:
                        long_tail.append({
                            'phrase': phrase,
                            'title': title[:50] + '...' if len(title) > 50 else title
                        })
        
        return long_tail[:15]
    
    def _analyze_keyword_performance(self, df: pd.DataFrame, top_words: Dict) -> Dict:
        """Analyze performance of different keywords"""
        performance = {}
        
        for word, count in list(top_words.items())[:10]:
            # Find videos with this word in title
            word_videos = df[df['title'].str.contains(word, case=False, na=False)]
            
            if len(word_videos) > 0:
                avg_views = word_videos['view_count'].mean()
                avg_likes = word_videos['like_count'].mean()
                engagement = (word_videos['like_count'] / word_videos['view_count'] * 100).mean()
                
                performance[word] = {
                    'frequency': count,
                    'avg_views': int(avg_views),
                    'avg_likes': int(avg_likes),
                    'engagement_rate': round(engagement, 2),
                    'video_count': len(word_videos)
                }
        
        return performance
    
    def _generate_seo_suggestions(self, search_keyword: str, top_words: Dict, variations: List[Dict]) -> List[str]:
        """Generate SEO suggestions based on analysis"""
        suggestions = []
        
        # High-frequency word suggestions
        high_freq_words = [word for word, count in list(top_words.items())[:5] 
                          if count > 2 and word not in search_keyword.lower()]
        
        if high_freq_words:
            suggestions.append(f"Consider using high-frequency words: {', '.join(high_freq_words[:3])}")
        
        # Keyword variation suggestions
        if variations:
            top_variations = [v['word'] for v in variations[:3]]
            suggestions.append(f"Try keyword variations: {', '.join(top_variations)}")
        
        # Long-tail keyword suggestions
        suggestions.append("Focus on long-tail keywords for less competition")
        suggestions.append("Include numbers in titles (tutorials, tips, etc.)")
        suggestions.append("Use action words: 'learn', 'master', 'build', 'create'")
        
        # Niche-specific suggestions
        if 'tutorial' in search_keyword.lower():
            suggestions.append("Add skill level indicators: 'beginner', 'advanced', 'complete'")
        if 'programming' in search_keyword.lower() or 'coding' in search_keyword.lower():
            suggestions.append("Include programming language names")
            suggestions.append("Use project-based keywords: 'build', 'create', 'project'")
        
        return suggestions
    
    def _analyze_trending_keywords(self, df: pd.DataFrame) -> List[Dict]:
        """Analyze trending keywords based on recent content"""
        # Sort by publication date (most recent first)
        recent_videos = df.sort_values('published_at', ascending=False).head(10)
        
        trending = []
        for title in recent_videos['title'].tolist():
            # Extract words that might be trending
            words = re.findall(r'\b[a-zA-Z]{4,}\b', title.lower())
            for word in words:
                if word not in ['tutorial', 'programming', 'learn', 'coding', 'python', 'javascript']:
                    trending.append({
                        'word': word,
                        'title': title[:40] + '...' if len(title) > 40 else title
                    })
        
        # Count frequency and return top trending
        trending_freq = Counter([item['word'] for item in trending])
        trending_words = [{'word': word, 'count': count} for word, count in trending_freq.most_common(5)]
        
        return trending_words
    
    def _calculate_similarity(self, word1: str, word2: str) -> float:
        """Calculate similarity between two words"""
        if word1 == word2:
            return 1.0
        
        # Simple character-based similarity
        common_chars = set(word1) & set(word2)
        total_chars = set(word1) | set(word2)
        
        if not total_chars:
            return 0.0
        
        return len(common_chars) / len(total_chars)

    def generate_executive_report(self, videos: List[Dict], search_keyword: str) -> Dict:
        """Generate executive summary and actionable insights"""
        if not videos:
            return {}
        
        # Run all analyses
        growth_data = self.analyze_growth_patterns(videos)
        title_data = self.analyze_titles(videos)
        competitor_data = self.analyze_competitors(videos)
        temporal_data = self.analyze_temporal_patterns(videos)
        keyword_data = self.analyze_keywords(videos, search_keyword)
        
        # Generate executive summary
        executive_summary = self._generate_executive_summary(videos, search_keyword, growth_data, competitor_data)
        
        # Generate actionable insights
        actionable_insights = self._generate_actionable_insights(
            growth_data, title_data, competitor_data, temporal_data, keyword_data
        )
        
        # Generate content strategy recommendations
        content_strategy = self._generate_content_strategy(
            title_data, keyword_data, competitor_data, temporal_data
        )
        
        # Generate competitive intelligence
        competitive_intel = self._generate_competitive_intelligence(competitor_data, keyword_data)
        
        # Generate performance benchmarks
        performance_benchmarks = self._generate_performance_benchmarks(videos, growth_data)
        
        # Generate next steps
        next_steps = self._generate_next_steps(actionable_insights, content_strategy)
        
        report = {
            'executive_summary': executive_summary,
            'actionable_insights': actionable_insights,
            'content_strategy': content_strategy,
            'competitive_intelligence': competitive_intel,
            'performance_benchmarks': performance_benchmarks,
            'next_steps': next_steps,
            'report_metadata': {
                'keyword': search_keyword,
                'videos_analyzed': len(videos),
                'generated_at': datetime.now().strftime('%Y-%m-%d %H:%M:%S'),
                'analysis_depth': 'comprehensive'
            }
        }
        
        return report
    
    def _generate_executive_summary(self, videos: List[Dict], search_keyword: str, growth_data: Dict, competitor_data: Dict) -> Dict:
        """Generate executive summary of findings"""
        df = pd.DataFrame(videos)
        total_views = df['view_count'].sum()
        avg_views = df['view_count'].mean()
        avg_engagement = (df['like_count'] / df['view_count'] * 100).mean()
        
        # Market size and opportunity
        market_size = "Large" if total_views > 100000000 else "Medium" if total_views > 10000000 else "Small"
        
        # Competition level
        concentration = competitor_data.get('market_concentration', {}).get('concentration_level', 'Unknown')
        competition_level = "High" if concentration == "High" else "Medium" if concentration == "Medium" else "Low"
        
        # Opportunity assessment
        opportunity_score = self._calculate_opportunity_score(competitor_data, avg_views, avg_engagement)
        
        summary = {
            'market_overview': {
                'keyword': search_keyword,
                'market_size': market_size,
                'total_views_analyzed': f"{total_views:,}",
                'average_views': f"{avg_views:,.0f}",
                'competition_level': competition_level,
                'opportunity_score': f"{opportunity_score}/10"
            },
            'key_findings': [
                f"Market has {len(competitor_data.get('top_competitors', []))} major players",
                f"Average engagement rate: {avg_engagement:.1f}%",
                f"Top performer has {df['view_count'].max():,} views",
                f"Content freshness: {len(df[df['published_at'] > (datetime.now() - pd.Timedelta(days=30)).strftime('%Y-%m-%d')])} recent videos"
            ],
            'strategic_recommendation': self._get_strategic_recommendation(opportunity_score, competition_level, market_size)
        }
        
        return summary
    
    def _generate_actionable_insights(self, growth_data: Dict, title_data: Dict, competitor_data: Dict, temporal_data: Dict, keyword_data: Dict) -> List[Dict]:
        """Generate actionable insights for content creators"""
        insights = []
        
        # Title optimization insights
        if title_data.get('common_words'):
            top_words = list(title_data['common_words'].items())[:3]
            insights.append({
                'category': 'Title Optimization',
                'priority': 'High',
                'insight': f"Use these high-performing words: {', '.join([word for word, _ in top_words])}",
                'action': 'Include these words in your next 5 video titles',
                'expected_impact': '15-25% increase in discoverability'
            })
        
        # Keyword insights
        if keyword_data.get('top_keywords'):
            top_keyword = list(keyword_data['top_keywords'].items())[0]
            insights.append({
                'category': 'SEO Strategy',
                'priority': 'High',
                'insight': f"'{top_keyword[0]}' appears {top_keyword[1]} times in successful videos",
                'action': f'Create content series around "{top_keyword[0]}"',
                'expected_impact': '20-30% increase in search visibility'
            })
        
        # Timing insights
        if temporal_data.get('best_posting_times'):
            best_time = temporal_data['best_posting_times']['best_hours'][0]['hour']
            best_day = temporal_data['best_posting_times']['best_days'][0]['day_of_week']
            insights.append({
                'category': 'Posting Strategy',
                'priority': 'Medium',
                'insight': f'Best performance: {best_day} at {int(best_time)}:00',
                'action': 'Schedule your next 3 videos for this time slot',
                'expected_impact': '10-20% increase in initial views'
            })
        
        # Competitor insights
        if competitor_data.get('top_competitors'):
            top_competitor = competitor_data['top_competitors'][0]
            insights.append({
                'category': 'Competitive Analysis',
                'priority': 'High',
                'insight': f"'{top_competitor['channel']}' dominates with {top_competitor['market_share']:.1f}% market share",
                'action': 'Study their content strategy and identify gaps',
                'expected_impact': 'Better positioning and unique value proposition'
            })
        
        # Content format insights
        if title_data.get('patterns'):
            patterns = title_data['patterns']
            if patterns.get('tutorial_pattern', 0) > 50:
                insights.append({
                    'category': 'Content Format',
                    'priority': 'Medium',
                    'insight': f"Tutorial content performs well ({patterns['tutorial_pattern']:.0f}% of top videos)",
                    'action': 'Create more step-by-step tutorial content',
                    'expected_impact': 'Higher engagement and retention'
                })
        
        return insights
    
    def _generate_content_strategy(self, title_data: Dict, keyword_data: Dict, competitor_data: Dict, temporal_data: Dict) -> Dict:
        """Generate content strategy recommendations"""
        strategy = {
            'content_themes': [],
            'posting_schedule': {},
            'title_formulas': [],
            'keyword_strategy': {},
            'differentiation_tactics': []
        }
        
        # Content themes based on top keywords
        if keyword_data.get('top_keywords'):
            top_keywords = list(keyword_data['top_keywords'].items())[:5]
            strategy['content_themes'] = [
                f"Create content around '{word}' (appears {count} times)" 
                for word, count in top_keywords
            ]
        
        # Posting schedule
        if temporal_data.get('best_posting_times'):
            best_day = temporal_data['best_posting_times']['best_days'][0]['day_of_week']
            best_hour = temporal_data['best_posting_times']['best_hours'][0]['hour']
            strategy['posting_schedule'] = {
                'optimal_day': best_day,
                'optimal_hour': f"{int(best_hour)}:00",
                'frequency': 'Weekly' if len(competitor_data.get('top_competitors', [])) > 5 else 'Bi-weekly'
            }
        
        # Title formulas
        if title_data.get('patterns'):
            patterns = title_data['patterns']
            formulas = []
            
            if patterns.get('tutorial_pattern', 0) > 40:
                formulas.append("'[Skill] Tutorial for [Audience]'")
            if patterns.get('number_pattern', 0) > 50:
                formulas.append("'[Number] [Topic] Tips/Tricks'")
            if patterns.get('question_pattern', 0) > 20:
                formulas.append("'How to [Action] [Topic]?'")
            
            strategy['title_formulas'] = formulas
        
        # Keyword strategy
        if keyword_data.get('long_tail_keywords'):
            strategy['keyword_strategy'] = {
                'primary_keywords': list(keyword_data['top_keywords'].keys())[:3],
                'long_tail_opportunities': [lt['phrase'] for lt in keyword_data['long_tail_keywords'][:5]],
                'trending_keywords': [t['word'] for t in keyword_data.get('trending_keywords', [])]
            }
        
        # Differentiation tactics
        if competitor_data.get('market_concentration'):
            concentration = competitor_data['market_concentration']['concentration_level']
            if concentration == 'High':
                strategy['differentiation_tactics'] = [
                    "Focus on niche sub-topics with less competition",
                    "Create unique content formats (series, challenges)",
                    "Target underserved audience segments",
                    "Develop signature teaching style"
                ]
            else:
                strategy['differentiation_tactics'] = [
                    "Build consistent content schedule",
                    "Focus on quality over quantity",
                    "Engage with community regularly",
                    "Create evergreen content"
                ]
        
        return strategy
    
    def _generate_competitive_intelligence(self, competitor_data: Dict, keyword_data: Dict) -> Dict:
        """Generate competitive intelligence report"""
        intel = {
            'market_leaders': [],
            'content_gaps': [],
            'opportunity_areas': [],
            'threat_assessment': {}
        }
        
        # Market leaders analysis
        if competitor_data.get('top_competitors'):
            for i, competitor in enumerate(competitor_data['top_competitors'][:3], 1):
                intel['market_leaders'].append({
                    'rank': i,
                    'channel': competitor['channel'],
                    'market_share': f"{competitor['market_share']:.1f}%",
                    'strength': f"{competitor['avg_views']:,.0f} avg views",
                    'strategy': 'High volume' if competitor['video_count'] > 5 else 'Quality focused'
                })
        
        # Content gaps
        if keyword_data.get('trending_keywords'):
            trending = [t['word'] for t in keyword_data['trending_keywords']]
            intel['content_gaps'] = [
                f"Underutilized trending keyword: '{word}'" for word in trending[:3]
            ]
        
        # Opportunity areas
        if competitor_data.get('market_concentration'):
            concentration = competitor_data['market_concentration']['concentration_level']
            if concentration == 'High':
                intel['opportunity_areas'] = [
                    "Niche sub-topics with less competition",
                    "Different content formats (shorts, live streams)",
                    "Target specific skill levels (absolute beginner, advanced)",
                    "Regional or language-specific content"
                ]
            else:
                intel['opportunity_areas'] = [
                    "Consistent content creation",
                    "Community building and engagement",
                    "Cross-platform content strategy",
                    "Collaboration opportunities"
                ]
        
        # Threat assessment
        intel['threat_assessment'] = {
            'competition_level': competitor_data.get('market_concentration', {}).get('concentration_level', 'Unknown'),
            'barrier_to_entry': 'High' if concentration == 'High' else 'Medium',
            'market_saturation': 'High' if len(competitor_data.get('top_competitors', [])) > 10 else 'Medium',
            'recommended_approach': 'Differentiation' if concentration == 'High' else 'Consistency'
        }
        
        return intel
    
    def _generate_performance_benchmarks(self, videos: List[Dict], growth_data: Dict) -> Dict:
        """Generate performance benchmarks for comparison"""
        df = pd.DataFrame(videos)
        
        benchmarks = {
            'view_benchmarks': {
                'excellent': int(df['view_count'].quantile(0.9)),
                'good': int(df['view_count'].quantile(0.7)),
                'average': int(df['view_count'].mean()),
                'poor': int(df['view_count'].quantile(0.3))
            },
            'engagement_benchmarks': {
                'excellent': float(df['like_count'].sum() / df['view_count'].sum() * 100 * 2),
                'good': float(df['like_count'].sum() / df['view_count'].sum() * 100 * 1.5),
                'average': float(df['like_count'].sum() / df['view_count'].sum() * 100),
                'poor': float(df['like_count'].sum() / df['view_count'].sum() * 100 * 0.5)
            },
            'content_frequency': {
                'top_performers': len(df[df['view_count'] > df['view_count'].quantile(0.8)]),
                'total_analyzed': len(df)
            }
        }
        
        return benchmarks
    
    def _generate_next_steps(self, actionable_insights: List[Dict], content_strategy: Dict) -> List[Dict]:
        """Generate prioritized next steps"""
        next_steps = []
        
        # High priority actions
        high_priority = [insight for insight in actionable_insights if insight.get('priority') == 'High']
        for i, insight in enumerate(high_priority[:3], 1):
            next_steps.append({
                'step': i,
                'priority': 'High',
                'action': insight['action'],
                'timeline': 'This week',
                'category': insight['category']
            })
        
        # Medium priority actions
        medium_priority = [insight for insight in actionable_insights if insight.get('priority') == 'Medium']
        for i, insight in enumerate(medium_priority[:2], len(next_steps) + 1):
            next_steps.append({
                'step': i,
                'priority': 'Medium',
                'action': insight['action'],
                'timeline': 'Next 2 weeks',
                'category': insight['category']
            })
        
        # Content strategy actions
        if content_strategy.get('content_themes'):
            next_steps.append({
                'step': len(next_steps) + 1,
                'priority': 'Medium',
                'action': f"Plan content around: {content_strategy['content_themes'][0]}",
                'timeline': 'Next month',
                'category': 'Content Planning'
            })
        
        return next_steps
    
    def _calculate_opportunity_score(self, competitor_data: Dict, avg_views: float, avg_engagement: float) -> int:
        """Calculate opportunity score from 1-10"""
        score = 5  # Base score
        
        # Market concentration factor
        concentration = competitor_data.get('market_concentration', {}).get('concentration_level', 'Medium')
        if concentration == 'Low':
            score += 2
        elif concentration == 'High':
            score -= 1
        
        # Engagement factor
        if avg_engagement > 5:
            score += 2
        elif avg_engagement > 2:
            score += 1
        else:
            score -= 1
        
        # Market size factor
        total_competitors = len(competitor_data.get('top_competitors', []))
        if total_competitors < 5:
            score += 1
        elif total_competitors > 15:
            score -= 1
        
        return max(1, min(10, score))
    
    def _get_strategic_recommendation(self, opportunity_score: int, competition_level: str, market_size: str) -> str:
        """Get strategic recommendation based on analysis"""
        if opportunity_score >= 8:
            return "Excellent opportunity - market is ripe for new entrants with high potential for growth"
        elif opportunity_score >= 6:
            return "Good opportunity - focus on differentiation and unique value proposition"
        elif opportunity_score >= 4:
            return "Moderate opportunity - requires careful strategy and consistent execution"
        else:
            return "Challenging market - consider niche sub-topics or different approach"


def interactive_mode():
    """Interactive mode that guides users through the analysis process"""
    print(f"\n{Fore.CYAN}🚀 Welcome to YTMiner Interactive Mode!{Style.RESET_ALL}")
    print(f"{Fore.YELLOW}I'll guide you through the analysis process step by step.{Style.RESET_ALL}\n")
    
    # Step 1: Get keyword
    keyword = click.prompt(f"{Fore.GREEN}1. What topic would you like to analyze?", type=str)
    
    # Step 2: Get number of results
    print(f"\n{Fore.GREEN}2. How many videos should I analyze?{Style.RESET_ALL}")
    print("   • 10-20: Quick analysis")
    print("   • 25-35: Detailed analysis") 
    print("   • 40-50: Comprehensive analysis")
    max_results = click.prompt("Number of videos", type=int, default=25)
    max_results = min(max_results, 50)  # Cap at 50
    
    # Step 3: Get user type
    print(f"\n{Fore.GREEN}3. What best describes you?{Style.RESET_ALL}")
    print("   1. Content Creator (YouTuber, streamer)")
    print("   2. Digital Marketer (Marketing professional)")
    print("   3. Researcher/Analyst (Academic or business research)")
    print("   4. Business Owner (Looking for market insights)")
    
    user_type = click.prompt("Choose your role", type=click.Choice(['1', '2', '3', '4']), default='1')
    
    # Step 4: Get analysis preferences based on user type
    analysis_options = get_analysis_preferences(user_type)
    
    # Step 5: Optional filters
    print(f"\n{Fore.GREEN}4. Any specific filters? (Press Enter to skip){Style.RESET_ALL}")
    
    # Time filter
    time_filter = click.prompt("   • Videos from last X days? (e.g., 7, 30, 90)", default="", type=str)
    days_back = int(time_filter) if time_filter.isdigit() else None
    
    # Region filter
    region = click.prompt("   • Specific country? (e.g., BR, US, UK)", default="", type=str)
    region = region.upper() if region else None
    
    # Duration filter
    print("   • Video duration?")
    print("     1. Short (under 4 minutes)")
    print("     2. Medium (4-20 minutes)")
    print("     3. Long (over 20 minutes)")
    print("     4. Any duration")
    duration_choice = click.prompt("Choose duration", type=click.Choice(['1', '2', '3', '4']), default='4')
    duration_map = {'1': 'short', '2': 'medium', '3': 'long', '4': 'any'}
    duration = duration_map[duration_choice]
    
    # Step 6: Run analysis
    print(f"\n{Fore.CYAN}🔍 Starting analysis...{Style.RESET_ALL}")
    print(f"   Keyword: {keyword}")
    print(f"   Videos: {max_results}")
    print(f"   Analysis: {', '.join(analysis_options)}")
    if days_back:
        print(f"   Time filter: Last {days_back} days")
    if region:
        print(f"   Region: {region}")
    print(f"   Duration: {duration}")
    
    # Execute the analysis
    run_analysis(keyword, max_results, days_back, region, duration, analysis_options)


def get_analysis_preferences(user_type):
    """Get analysis preferences based on user type"""
    preferences = {
        '1': {  # Content Creator
            'name': 'Content Creator',
            'recommended': ['title_analysis', 'temporal_analysis', 'competitor_analysis', 'keyword_analysis'],
            'description': 'Perfect for optimizing your content strategy!'
        },
        '2': {  # Digital Marketer
            'name': 'Digital Marketer',
            'recommended': ['competitor_analysis', 'keyword_analysis', 'executive_report'],
            'description': 'Great for market research and campaign planning!'
        },
        '3': {  # Researcher/Analyst
            'name': 'Researcher/Analyst',
            'recommended': ['analysis', 'temporal_analysis', 'executive_report'],
            'description': 'Ideal for comprehensive data analysis!'
        },
        '4': {  # Business Owner
            'name': 'Business Owner',
            'recommended': ['executive_report', 'competitor_analysis', 'keyword_analysis'],
            'description': 'Perfect for strategic business insights!'
        }
    }
    
    user_pref = preferences[user_type]
    print(f"\n{Fore.YELLOW}Great choice! {user_pref['description']}{Style.RESET_ALL}")
    
    # Show recommended analyses
    print(f"\n{Fore.GREEN}Recommended analyses for {user_pref['name']}:{Style.RESET_ALL}")
    analysis_descriptions = {
        'analysis': '📈 Growth Pattern Analysis - Performance trends and top performers',
        'title_analysis': '📝 Title Analysis - Winning title formulas and patterns',
        'competitor_analysis': '🏆 Competitor Analysis - Market leaders and positioning',
        'temporal_analysis': '⏰ Temporal Analysis - Best posting times and schedules',
        'keyword_analysis': '🔍 Keyword Analysis - SEO opportunities and trending terms',
        'executive_report': '📊 Executive Report - Comprehensive strategic insights'
    }
    
    for analysis in user_pref['recommended']:
        print(f"   ✓ {analysis_descriptions[analysis]}")
    
    # Ask if they want all recommended or custom selection
    print(f"\n{Fore.GREEN}Would you like to:{Style.RESET_ALL}")
    print("   1. Run all recommended analyses (recommended)")
    print("   2. Choose specific analyses")
    
    choice = click.prompt("Your choice", type=click.Choice(['1', '2']), default='1')
    
    if choice == '1':
        return user_pref['recommended']
    else:
        return custom_analysis_selection()


def custom_analysis_selection():
    """Allow users to select specific analyses"""
    print(f"\n{Fore.GREEN}Select the analyses you want:{Style.RESET_ALL}")
    print("   1. 📈 Growth Pattern Analysis")
    print("   2. 📝 Title Pattern Analysis") 
    print("   3. 🏆 Competitor Analysis")
    print("   4. ⏰ Temporal Analysis")
    print("   5. 🔍 Keyword Analysis")
    print("   6. 📊 Executive Report")
    
    selections = click.prompt("Enter numbers separated by commas (e.g., 1,3,5)", type=str)
    
    analysis_map = {
        '1': 'analysis',
        '2': 'title_analysis', 
        '3': 'competitor_analysis',
        '4': 'temporal_analysis',
        '5': 'keyword_analysis',
        '6': 'executive_report'
    }
    
    selected = []
    for num in selections.split(','):
        num = num.strip()
        if num in analysis_map:
            selected.append(analysis_map[num])
    
    return selected if selected else ['analysis']  # Default to basic analysis


def run_analysis(keyword, max_results, days_back, region, duration, analysis_options):
    """Run the analysis with the selected options"""
    api_key = os.getenv('YOUTUBE_API_KEY')
    miner = YouTubeMiner(api_key)
    
    # Set up search parameters
    published_after = None
    if days_back:
        published_after = (datetime.now() - timedelta(days=days_back)).isoformat() + 'Z'
    
    # Search for videos
    print(f"\n{Fore.CYAN}Searching for videos about '{keyword}'...{Style.RESET_ALL}")
    
    videos = miner.search_videos(
        keyword=keyword,
        max_results=max_results,
        order='relevance',
        published_after=published_after,
        region_code=region,
        video_duration=duration if duration != 'any' else None
    )
    
    if not videos:
        print(f"{Fore.RED}No videos found. Try adjusting your search criteria.{Style.RESET_ALL}")
        return
    
    print(f"{Fore.GREEN}Found {len(videos)} videos!{Style.RESET_ALL}")
    
    # Display basic results
    display_basic_results(videos)
    
    # Run selected analyses
    for analysis in analysis_options:
        if analysis == 'analysis':
            run_growth_analysis(miner, videos)
        elif analysis == 'title_analysis':
            run_title_analysis(miner, videos)
        elif analysis == 'competitor_analysis':
            run_competitor_analysis(miner, videos)
        elif analysis == 'temporal_analysis':
            run_temporal_analysis(miner, videos)
        elif analysis == 'keyword_analysis':
            run_keyword_analysis(miner, videos, keyword)
        elif analysis == 'executive_report':
            run_executive_report(miner, videos, keyword)
    
    print(f"\n{Fore.GREEN}✅ Analysis complete!{Style.RESET_ALL}")
    print(f"{Fore.YELLOW}💡 Tip: Use specific flags for faster analysis next time!{Style.RESET_ALL}")


def display_basic_results(videos):
    """Display basic video results"""
    print(f"\n{Fore.CYAN}📊 Top Videos:{Style.RESET_ALL}")
    
    for i, video in enumerate(videos[:5], 1):
        print(f"   {i}. {video['title'][:60]}{'...' if len(video['title']) > 60 else ''}")
        print(f"      👀 {video['view_count']:,} views | 👍 {video['like_count']:,} likes | 📅 {video['published_at'][:10]}")
        print(f"      🔗 {video['url']}")
        print()


def run_growth_analysis(miner, videos):
    """Run growth pattern analysis"""
    print(f"\n{Fore.CYAN}📈 Growth Pattern Analysis:{Style.RESET_ALL}")
    growth_data = miner.analyze_growth_patterns(videos)
    
    print(f"   • Total videos: {len(videos)}")
    print(f"   • Average views: {growth_data.get('avg_views', 0):,.0f}")
    print(f"   • Average engagement: {growth_data.get('avg_engagement', 0):.1f}%")
    
    if growth_data.get('top_performers'):
        print(f"\n   🏆 Top Performers:")
        for i, performer in enumerate(growth_data['top_performers'][:3], 1):
            print(f"      {i}. {performer['title'][:50]}...")
            print(f"         Views: {performer['view_count']:,} | Likes: {performer['like_count']:,}")


def run_title_analysis(miner, videos):
    """Run title pattern analysis"""
    print(f"\n{Fore.CYAN}📝 Title Pattern Analysis:{Style.RESET_ALL}")
    title_data = miner.analyze_titles(videos)
    
    if title_data.get('common_words'):
        print(f"   🔤 Most common words:")
        for word, count in list(title_data['common_words'].items())[:5]:
            print(f"      • {word} ({count} times)")
    
    if title_data.get('patterns'):
        patterns = title_data['patterns']
        print(f"\n   📊 Title Patterns:")
        print(f"      • Tutorial content: {patterns.get('tutorial_pattern', 0):.0f}%")
        print(f"      • Numbers in titles: {patterns.get('number_pattern', 0):.0f}%")
        print(f"      • Questions: {patterns.get('question_pattern', 0):.0f}%")


def run_competitor_analysis(miner, videos):
    """Run competitor analysis"""
    print(f"\n{Fore.CYAN}🏆 Competitor Analysis:{Style.RESET_ALL}")
    competitor_data = miner.analyze_competitors(videos)
    
    print(f"   • Total channels: {competitor_data.get('total_channels', 0)}")
    print(f"   • Market concentration: {competitor_data.get('market_concentration', {}).get('concentration_level', 'Unknown')}")
    
    if competitor_data.get('top_competitors'):
        print(f"\n   🏅 Top Competitors:")
        for i, competitor in enumerate(competitor_data['top_competitors'][:3], 1):
            print(f"      {i}. {competitor['channel']}")
            print(f"         Market share: {competitor['market_share']:.1f}% | Videos: {competitor['video_count']}")


def run_temporal_analysis(miner, videos):
    """Run temporal analysis"""
    print(f"\n{Fore.CYAN}⏰ Temporal Analysis:{Style.RESET_ALL}")
    temporal_data = miner.analyze_temporal_patterns(videos)
    
    if temporal_data.get('best_posting_times'):
        best_time = temporal_data['best_posting_times']['best_hours'][0]['hour']
        best_day = temporal_data['best_posting_times']['best_days'][0]['day_of_week']
        print(f"   🎯 Best posting time: {best_day} at {int(best_time)}:00")
    
    print(f"   📊 Performance by day:")
    for day_data in temporal_data.get('day_performance', [])[:3]:
        print(f"      • {day_data['day_of_week']}: {day_data['avg_views']:,.0f} avg views")


def run_keyword_analysis(miner, videos, keyword):
    """Run keyword analysis"""
    print(f"\n{Fore.CYAN}🔍 Keyword Analysis:{Style.RESET_ALL}")
    keyword_data = miner.analyze_keywords(videos, keyword)
    
    if keyword_data.get('top_keywords'):
        print(f"   🔤 Top keywords:")
        for word, count in list(keyword_data['top_keywords'].items())[:5]:
            print(f"      • {word} ({count} times)")
    
    if keyword_data.get('seo_suggestions'):
        print(f"\n   💡 SEO Suggestions:")
        for suggestion in keyword_data['seo_suggestions'][:3]:
            print(f"      • {suggestion}")


def run_executive_report(miner, videos, keyword):
    """Run executive report"""
    print(f"\n{Fore.CYAN}📊 Executive Report:{Style.RESET_ALL}")
    report = miner.generate_executive_report(videos, keyword)
    
    summary = report['executive_summary']
    print(f"   🎯 Opportunity Score: {summary['market_overview']['opportunity_score']}")
    print(f"   📈 Market Size: {summary['market_overview']['market_size']}")
    print(f"   🏆 Competition: {summary['market_overview']['competition_level']}")
    
    if report.get('actionable_insights'):
        print(f"\n   🚀 Top Insights:")
        for insight in report['actionable_insights'][:2]:
            print(f"      • {insight['insight']}")
            print(f"        Action: {insight['action']}")


@click.command()
@click.option('--keyword', '-k', help='Search keyword')
@click.option('--max-results', '-n', default=20, help='Maximum number of results (default: 20)')
@click.option('--order', '-o', 
              type=click.Choice(['relevance', 'date', 'rating', 'viewCount', 'title']),
              default='relevance', help='Result ordering')
@click.option('--days-back', '-d', type=int, help='Number of days back (e.g., 7 for last week)')
@click.option('--region', '-r', help='Country code (e.g., BR, US)')
@click.option('--duration', 
              type=click.Choice(['short', 'medium', 'long', 'any']),
              help='Video duration')
@click.option('--min-views', type=int, help='Minimum view count')
@click.option('--min-likes', type=int, help='Minimum like count')
@click.option('--output', '-f', help='Output file for results (JSON)')
@click.option('--analysis', '-a', is_flag=True, help='Show growth pattern analysis')
@click.option('--title-analysis', '-t', is_flag=True, help='Show title pattern analysis for creators')
@click.option('--competitor-analysis', '-c', is_flag=True, help='Show competitor and market analysis')
@click.option('--temporal-analysis', '-time', is_flag=True, help='Show temporal patterns and optimal posting times')
@click.option('--keyword-analysis', '-kwd', is_flag=True, help='Show keyword analysis and SEO opportunities')
@click.option('--executive-report', '-rpt', is_flag=True, help='Generate comprehensive executive report with actionable insights')
@click.option('--interactive', '-i', is_flag=True, help='Run in interactive mode (guided experience)')
def main(keyword, max_results, order, days_back, region, duration, min_views, min_likes, output, analysis, title_analysis, competitor_analysis, temporal_analysis, keyword_analysis, executive_report, interactive):
    """YTMiner - YouTube video analysis CLI tool"""
    
    api_key = os.getenv('YOUTUBE_API_KEY')
    if not api_key:
        print(f"{Fore.RED}Error: YouTube API key not found.{Style.RESET_ALL}")
        print(f"{Fore.YELLOW}Set the YOUTUBE_API_KEY environment variable{Style.RESET_ALL}")
        print(f"{Fore.CYAN}Example: set YOUTUBE_API_KEY=your_key_here{Style.RESET_ALL}")
        return
    
    # Interactive mode
    if interactive:
        interactive_mode()
        return
    
    # Validate required parameters for non-interactive mode
    if not keyword:
        print(f"{Fore.RED}Error: Keyword is required for non-interactive mode.{Style.RESET_ALL}")
        print("Use --interactive (-i) for guided experience or provide --keyword (-k)")
        return
    
    published_after = None
    if days_back:
        published_after = (datetime.now() - timedelta(days=days_back)).isoformat() + 'Z'
    
    miner = YouTubeMiner(api_key)
    
    print(f"{Fore.CYAN}Searching for videos about '{keyword}'...{Style.RESET_ALL}")
    
    videos = miner.search_videos(
        keyword=keyword,
        max_results=max_results,
        order=order,
        published_after=published_after,
        region_code=region,
        video_duration=duration
    )
    
    if not videos:
        print(f"{Fore.RED}No videos found with the specified criteria.{Style.RESET_ALL}")
        return
    
    if min_views:
        videos = [v for v in videos if v['view_count'] >= min_views]
    if min_likes:
        videos = [v for v in videos if v['like_count'] >= min_likes]
    
    if not videos:
        print(f"{Fore.RED}No videos meet the view/like filters.{Style.RESET_ALL}")
        return
    
    videos.sort(key=lambda x: x['view_count'], reverse=True)
    
    print(f"\n{Fore.GREEN}Found {len(videos)} videos:{Style.RESET_ALL}\n")
    
    table_data = []
    for i, video in enumerate(videos, 1):
        table_data.append([
            i,
            video['title'][:40] + '...' if len(video['title']) > 40 else video['title'],
            video['channel'][:20] + '...' if len(video['channel']) > 20 else video['channel'],
            f"{video['view_count']:,}",
            f"{video['like_count']:,}",
            video['published_at'][:10],
            video['url']
        ])
    
    headers = ['#', 'Title', 'Channel', 'Views', 'Likes', 'Date', 'URL']
    print(tabulate(table_data, headers=headers, tablefmt='grid'))
    
    print(f"\n{Fore.YELLOW}Top 5 videos:{Style.RESET_ALL}")
    for i, video in enumerate(videos[:5], 1):
        print(f"{i}. {video['title']}")
        print(f"   {video['url']}")
        print(f"   Views: {video['view_count']:,} | Likes: {video['like_count']:,}")
        print()
    
    if analysis:
        print(f"{Fore.CYAN}Growth Pattern Analysis:{Style.RESET_ALL}\n")
        
        analysis_data = miner.analyze_growth_patterns(videos)
        
        print(f"General Statistics:")
        print(f"   • Total videos: {analysis_data['total_videos']}")
        print(f"   • Average views: {analysis_data['avg_views']:,.0f}")
        print(f"   • Average likes: {analysis_data['avg_likes']:,.0f}")
        print(f"   • Engagement rate: {analysis_data['avg_engagement_rate']:.2f}%")
        
        print(f"\nTop Performers (most views):")
        for i, video in enumerate(analysis_data['top_performers'], 1):
            print(f"   {i}. {video['title'][:60]}...")
            print(f"      Views: {video['view_count']:,} | Likes: {video['like_count']:,} | {video['days_ago']} days ago")
        
        if analysis_data['recent_trends']:
            print(f"\nRecent Trends (last 7 days):")
            for i, video in enumerate(analysis_data['recent_trends'], 1):
                print(f"   {i}. {video['title'][:60]}...")
                print(f"      Views: {video['view_count']:,} | {video['days_ago']} days ago")
        
        print(f"\nBest Engagement (most likes):")
        for i, video in enumerate(analysis_data['best_engagement'], 1):
            print(f"   {i}. {video['title'][:60]}...")
            print(f"      Likes: {video['like_count']:,} | Views: {video['view_count']:,} | {video['days_ago']} days ago")
    
    if title_analysis:
        print(f"\n{Fore.CYAN}Title Pattern Analysis (Creator Insights):{Style.RESET_ALL}\n")
        
        title_data = miner.analyze_titles(videos)
        
        print(f"📊 Title Statistics:")
        print(f"   • Titles analyzed: {title_data['total_titles_analyzed']}")
        print(f"   • Average title length: {title_data['avg_title_length']:.0f} characters")
        
        print(f"\n🔤 Most Common Words:")
        for word, count in title_data['most_common_words'][:5]:
            print(f"   • '{word}' ({count} times)")
        
        print(f"\n📝 Most Common Phrases:")
        for phrase, count in title_data['most_common_phrases'][:5]:
            print(f"   • '{phrase}' ({count} times)")
        
        print(f"\n😀 Emoji Usage:")
        emoji_data = title_data['emoji_usage']
        print(f"   • Titles with emojis: {emoji_data['titles_with_emojis']} ({emoji_data['emoji_percentage']:.1f}%)")
        if emoji_data['most_common_emojis']:
            print(f"   • Most common emojis: {', '.join([emoji for emoji, _ in emoji_data['most_common_emojis'][:3]])}")
        
        print(f"\n🎯 Title Patterns:")
        patterns = title_data['title_patterns']
        print(f"   • Tutorial/How-to: {patterns['tutorial_pattern']:.1f}%")
        print(f"   • Numbers in title: {patterns['number_pattern']:.1f}%")
        print(f"   • Questions (?): {patterns['question_pattern']:.1f}%")
        print(f"   • Exclamations (!): {patterns['exclamation_pattern']:.1f}%")
        print(f"   • Brackets [...]: {patterns['bracket_pattern']:.1f}%")
        
        print(f"\n🏆 Successful Title Examples:")
        for i, title in enumerate(title_data['successful_titles'], 1):
            print(f"   {i}. {title}")
        
        print(f"\n💡 Creator Tips:")
        print(f"   • Use {title_data['avg_title_length']:.0f} characters on average")
        if emoji_data['emoji_percentage'] > 30:
            print(f"   • Emojis are popular in this niche ({emoji_data['emoji_percentage']:.1f}% usage)")
        if patterns['tutorial_pattern'] > 40:
            print(f"   • Tutorial content performs well ({patterns['tutorial_pattern']:.1f}% of top videos)")
        if patterns['number_pattern'] > 50:
            print(f"   • Numbers in titles are effective ({patterns['number_pattern']:.1f}% usage)")
    
    if competitor_analysis:
        print(f"\n{Fore.CYAN}Competitor & Market Analysis:{Style.RESET_ALL}\n")
        
        competitor_data = miner.analyze_competitors(videos)
        
        print(f"🏢 Market Overview:")
        print(f"   • Total channels: {competitor_data['total_channels']}")
        print(f"   • Market concentration: {competitor_data['market_concentration']['concentration_level']}")
        print(f"   • Top 5 channels share: {competitor_data['market_concentration']['top_5_share']:.1f}%")
        print(f"   • Top 10 channels share: {competitor_data['market_concentration']['top_10_share']:.1f}%")
        
        print(f"\n🏆 Top Competitors:")
        for i, competitor in enumerate(competitor_data['top_competitors'][:5], 1):
            print(f"   {i}. {competitor['channel']}")
            print(f"      Videos: {competitor['video_count']} | Views: {competitor['total_views']:,} | Market Share: {competitor['market_share']:.1f}%")
            print(f"      Avg Views: {competitor['avg_views']:,} | Engagement: {competitor['engagement_rate']:.1f}%")
            print()
        
        print(f"📊 Content Patterns by Top Channels:")
        for channel, patterns in list(competitor_data['content_patterns'].items())[:3]:
            print(f"   • {channel[:30]}...")
            print(f"     Title length: {patterns['avg_title_length']:.0f} chars | Tutorial: {patterns['tutorial_ratio']:.0f}% | Numbers: {patterns['number_ratio']:.0f}% | Emojis: {patterns['emoji_ratio']:.0f}%")
        
        print(f"\n💡 Market Insights:")
        for insight in competitor_data['market_insights']:
            print(f"   • {insight}")
        
        print(f"\n🎯 Strategic Recommendations:")
        concentration = competitor_data['market_concentration']['concentration_level']
        if concentration == 'High':
            print(f"   • High competition - focus on unique value proposition")
            print(f"   • Study top 3 competitors closely for content gaps")
            print(f"   • Consider niche sub-topics with less competition")
        elif concentration == 'Medium':
            print(f"   • Moderate competition - good opportunity for growth")
            print(f"   • Identify underserved audience segments")
            print(f"   • Build consistent content schedule")
        else:
            print(f"   • Low competition - great opportunity to establish dominance")
            print(f"   • Focus on content quality and consistency")
            print(f"   • Build strong community engagement")
    
    if temporal_analysis:
        print(f"\n{Fore.CYAN}Temporal Analysis (Optimal Posting Times):{Style.RESET_ALL}\n")
        
        temporal_data = miner.analyze_temporal_patterns(videos)
        
        print(f"⏰ Best Posting Times:")
        best_times = temporal_data['best_posting_times']
        hours_str = ', '.join([f"{int(h['hour'])}:00" for h in best_times['best_hours']])
        days_str = ', '.join([h['day_of_week'] for h in best_times['best_days']])
        print(f"   • Top 3 Hours: {hours_str}")
        print(f"   • Top 3 Days: {days_str}")
        
        print(f"\n📊 Performance by Hour:")
        for hour_data in temporal_data['hour_performance'][:5]:
            print(f"   • {int(hour_data['hour']):02d}:00 - {hour_data['video_count']} videos, {hour_data['avg_views']:,.0f} avg views, {hour_data['engagement_rate']:.1f}% engagement")
        
        print(f"\n📅 Performance by Day:")
        for day_data in temporal_data['day_performance']:
            print(f"   • {day_data['day_of_week']} - {day_data['video_count']} videos, {day_data['avg_views']:,.0f} avg views, {day_data['engagement_rate']:.1f}% engagement")
        
        print(f"\n📈 Recent vs Older Content:")
        recent_older = temporal_data['recent_vs_older']
        print(f"   • Recent videos (last 30 days): {recent_older['recent_videos_count']} videos, {recent_older['recent_avg_views']:,.0f} avg views")
        print(f"   • Older videos: {recent_older['older_videos_count']} videos, {recent_older['older_avg_views']:,.0f} avg views")
        
        print(f"\n💡 Temporal Insights:")
        for insight in temporal_data['insights']:
            print(f"   • {insight}")
        
        print(f"\n🎯 Posting Strategy Recommendations:")
        best_hour = best_times['best_hours'][0]['hour']
        best_day = best_times['best_days'][0]['day_of_week']
        print(f"   • Optimal posting: {best_day} at {int(best_hour)}:00")
        
        # Additional recommendations based on data
        if len(temporal_data['hour_performance']) > 5:
            print(f"   • Consider posting during peak hours for maximum reach")
        if recent_older['recent_avg_views'] > recent_older['older_avg_views'] * 1.2:
            print(f"   • Focus on current trends - recent content performs better")
        elif recent_older['older_avg_views'] > recent_older['recent_avg_views'] * 1.2:
            print(f"   • Evergreen content works well - focus on timeless topics")
    
    if keyword_analysis:
        print(f"\n{Fore.CYAN}Keyword Analysis & SEO Opportunities:{Style.RESET_ALL}\n")
        
        keyword_data = miner.analyze_keywords(videos, keyword)
        
        print(f"🔍 Keyword Analysis for: '{keyword}'")
        print(f"   • Total words analyzed: {keyword_data['total_words_analyzed']}")
        print(f"   • Unique meaningful words: {keyword_data['unique_words']}")
        
        print(f"\n📊 Top Keywords by Frequency:")
        for i, (word, count) in enumerate(list(keyword_data['top_keywords'].items())[:10], 1):
            print(f"   {i:2d}. {word:<15} ({count} times)")
        
        print(f"\n🔄 Keyword Variations & Suggestions:")
        for variation in keyword_data['keyword_variations'][:5]:
            print(f"   • {variation['word']} (freq: {variation['frequency']}) → '{variation['suggestion']}'")
        
        print(f"\n🎯 Long-tail Keywords:")
        for lt in keyword_data['long_tail_keywords'][:8]:
            print(f"   • '{lt['phrase']}'")
            print(f"     Example: {lt['title']}")
        
        print(f"\n📈 Keyword Performance Analysis:")
        for word, perf in list(keyword_data['keyword_performance'].items())[:5]:
            print(f"   • {word}: {perf['avg_views']:,} avg views, {perf['engagement_rate']:.1f}% engagement ({perf['video_count']} videos)")
        
        print(f"\n🔥 Trending Keywords:")
        for trend in keyword_data['trending_keywords']:
            print(f"   • {trend['word']} (appears {trend['count']} times in recent videos)")
        
        print(f"\n💡 SEO Suggestions:")
        for suggestion in keyword_data['seo_suggestions']:
            print(f"   • {suggestion}")
        
        print(f"\n🚀 Action Items for Better SEO:")
        print(f"   • Use high-frequency keywords in your titles")
        print(f"   • Create content around long-tail keyword phrases")
        print(f"   • Monitor trending keywords for timely content")
        print(f"   • A/B test different keyword combinations")
        print(f"   • Focus on keywords with high engagement rates")
    
    if executive_report:
        print(f"\n{Fore.CYAN}📊 EXECUTIVE REPORT - Comprehensive Analysis:{Style.RESET_ALL}\n")
        
        report = miner.generate_executive_report(videos, keyword)
        
        # Executive Summary
        summary = report['executive_summary']
        print(f"🎯 EXECUTIVE SUMMARY")
        print(f"   Keyword: {summary['market_overview']['keyword']}")
        print(f"   Market Size: {summary['market_overview']['market_size']}")
        print(f"   Competition Level: {summary['market_overview']['competition_level']}")
        print(f"   Opportunity Score: {summary['market_overview']['opportunity_score']}")
        print(f"   Total Views Analyzed: {summary['market_overview']['total_views_analyzed']}")
        print(f"   Average Views: {summary['market_overview']['average_views']}")
        
        print(f"\n📋 KEY FINDINGS:")
        for finding in summary['key_findings']:
            print(f"   • {finding}")
        
        print(f"\n💡 STRATEGIC RECOMMENDATION:")
        print(f"   {summary['strategic_recommendation']}")
        
        # Actionable Insights
        print(f"\n🚀 ACTIONABLE INSIGHTS:")
        for i, insight in enumerate(report['actionable_insights'], 1):
            priority_icon = "🔴" if insight['priority'] == 'High' else "🟡" if insight['priority'] == 'Medium' else "🟢"
            print(f"   {i}. {priority_icon} {insight['category']}: {insight['insight']}")
            print(f"      Action: {insight['action']}")
            print(f"      Expected Impact: {insight['expected_impact']}")
            print()
        
        # Content Strategy
        strategy = report['content_strategy']
        print(f"📝 CONTENT STRATEGY:")
        
        if strategy.get('content_themes'):
            print(f"   🎬 Content Themes:")
            for theme in strategy['content_themes'][:3]:
                print(f"      • {theme}")
        
        if strategy.get('posting_schedule'):
            schedule = strategy['posting_schedule']
            print(f"   ⏰ Posting Schedule:")
            print(f"      • Optimal Day: {schedule['optimal_day']}")
            print(f"      • Optimal Time: {schedule['optimal_hour']}")
            print(f"      • Frequency: {schedule['frequency']}")
        
        if strategy.get('title_formulas'):
            print(f"   📝 Title Formulas:")
            for formula in strategy['title_formulas']:
                print(f"      • {formula}")
        
        # Competitive Intelligence
        intel = report['competitive_intelligence']
        print(f"\n🏆 COMPETITIVE INTELLIGENCE:")
        
        if intel.get('market_leaders'):
            print(f"   Market Leaders:")
            for leader in intel['market_leaders']:
                print(f"      {leader['rank']}. {leader['channel']} ({leader['market_share']}) - {leader['strategy']}")
        
        if intel.get('opportunity_areas'):
            print(f"   Opportunity Areas:")
            for area in intel['opportunity_areas'][:3]:
                print(f"      • {area}")
        
        # Performance Benchmarks
        benchmarks = report['performance_benchmarks']
        print(f"\n📊 PERFORMANCE BENCHMARKS:")
        print(f"   View Benchmarks:")
        print(f"      • Excellent: {benchmarks['view_benchmarks']['excellent']:,} views")
        print(f"      • Good: {benchmarks['view_benchmarks']['good']:,} views")
        print(f"      • Average: {benchmarks['view_benchmarks']['average']:,} views")
        print(f"   Engagement Benchmarks:")
        print(f"      • Excellent: {benchmarks['engagement_benchmarks']['excellent']:.1f}%")
        print(f"      • Good: {benchmarks['engagement_benchmarks']['good']:.1f}%")
        print(f"      • Average: {benchmarks['engagement_benchmarks']['average']:.1f}%")
        
        # Next Steps
        print(f"\n✅ NEXT STEPS (Prioritized):")
        for step in report['next_steps']:
            priority_icon = "🔴" if step['priority'] == 'High' else "🟡" if step['priority'] == 'Medium' else "🟢"
            print(f"   {step['step']}. {priority_icon} {step['action']}")
            print(f"      Timeline: {step['timeline']} | Category: {step['category']}")
            print()
        
        # Report Metadata
        metadata = report['report_metadata']
        print(f"📄 Report Generated: {metadata['generated_at']}")
        print(f"   Videos Analyzed: {metadata['videos_analyzed']}")
        print(f"   Analysis Depth: {metadata['analysis_depth']}")
    
    if output:
        with open(output, 'w', encoding='utf-8') as f:
            json.dump(videos, f, ensure_ascii=False, indent=2, default=str)
        print(f"\n{Fore.GREEN}Data saved to: {output}{Style.RESET_ALL}")


if __name__ == '__main__':
    main()
