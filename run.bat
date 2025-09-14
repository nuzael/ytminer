@echo off
echo YTMiner - YouTube Analysis CLI Tool
echo ====================================
echo.

python --version >nul 2>&1
if errorlevel 1 (
    echo Error: Python not found. Install Python 3.7+ first.
    pause
    exit /b 1
)

python -c "import googleapiclient" >nul 2>&1
if errorlevel 1 (
    echo Installing dependencies...
    pip install -r requirements.txt
    if errorlevel 1 (
        echo Error installing dependencies.
        pause
        exit /b 1
    )
)

if "%YOUTUBE_API_KEY%"=="" (
    echo.
    echo WARNING: YouTube API key not configured!
    echo.
    echo To configure:
    echo 1. Get a key at: https://console.cloud.google.com/
    echo 2. Set environment variable:
    echo    set YOUTUBE_API_KEY=your_key_here
    echo.
    echo Or run: set_api_key.bat
    echo.
    pause
    exit /b 1
)

python ytminer.py %*

pause
