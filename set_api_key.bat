@echo off
echo YouTube API Key Configuration
echo =============================
echo.

set /p API_KEY="Enter your YouTube API key: "

if "%API_KEY%"=="" (
    echo Error: Key cannot be empty.
    pause
    exit /b 1
)

echo.
echo Configuring API key...
setx YOUTUBE_API_KEY "%API_KEY%"

if errorlevel 1 (
    echo Error configuring API key.
    pause
    exit /b 1
)

echo.
echo API key configured successfully!
echo.
echo IMPORTANT: Close and reopen terminal for variable to be recognized.
echo.
pause
