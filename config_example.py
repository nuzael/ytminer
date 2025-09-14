"""
Example configuration file for YTMiner
Copy this file to config.py and configure your preferences
"""

# YouTube API settings
YOUTUBE_API_KEY = "your_api_key_here"

# Default settings
DEFAULT_MAX_RESULTS = 20
DEFAULT_ORDER = "relevance"
DEFAULT_REGION = "BR"  # Brazil

# Default filters
DEFAULT_MIN_VIEWS = 1000
DEFAULT_MIN_LIKES = 10

# Analysis settings
ENABLE_ANALYSIS = True
ANALYSIS_DAYS_BACK = 7

# Output settings
OUTPUT_FORMAT = "table"  # table, json, csv
SAVE_TO_FILE = False
OUTPUT_FILE = "ytminer_results.json"
