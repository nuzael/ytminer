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


@click.command()
@click.option('--keyword', '-k', required=True, help='Search keyword')
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
def main(keyword, max_results, order, days_back, region, duration, min_views, min_likes, output, analysis):
    """YTMiner - YouTube video analysis CLI tool"""
    
    api_key = os.getenv('YOUTUBE_API_KEY')
    if not api_key:
        print(f"{Fore.RED}Error: YouTube API key not found.{Style.RESET_ALL}")
        print(f"{Fore.YELLOW}Set the YOUTUBE_API_KEY environment variable{Style.RESET_ALL}")
        print(f"{Fore.CYAN}Example: set YOUTUBE_API_KEY=your_key_here{Style.RESET_ALL}")
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
            video['title'][:50] + '...' if len(video['title']) > 50 else video['title'],
            video['channel'],
            f"{video['view_count']:,}",
            f"{video['like_count']:,}",
            video['published_at'][:10]
        ])
    
    headers = ['#', 'Title', 'Channel', 'Views', 'Likes', 'Date']
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
    
    if output:
        with open(output, 'w', encoding='utf-8') as f:
            json.dump(videos, f, ensure_ascii=False, indent=2, default=str)
        print(f"\n{Fore.GREEN}Data saved to: {output}{Style.RESET_ALL}")


if __name__ == '__main__':
    main()
