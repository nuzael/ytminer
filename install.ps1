# YTMiner Installation Script for Windows

Write-Host "Ì∫Ä YTMiner Installation Script" -ForegroundColor Green
Write-Host "===============================" -ForegroundColor Green

# Detect architecture
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
Write-Host "Ì≥ã Detected: Windows $arch" -ForegroundColor Blue

# Get latest release
try {
    $release = Invoke-RestMethod -Uri "https://api.github.com/repos/nuzael/ytminer/releases/latest"
    $version = $release.tag_name
    Write-Host "Ì≥¶ Downloading YTMiner $version..." -ForegroundColor Blue
} catch {
    Write-Host "‚ùå Failed to get latest release: $_" -ForegroundColor Red
    exit 1
}

# Set download URL
$binaryName = "ytminer-windows-$arch.exe"
$downloadUrl = "https://github.com/nuzael/ytminer/releases/download/$version/$binaryName"

# Create temp directory
$tempDir = [System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid().ToString()
New-Item -ItemType Directory -Path $tempDir | Out-Null

try {
    # Download binary
    $binaryPath = Join-Path $tempDir $binaryName
    Invoke-WebRequest -Uri $downloadUrl -OutFile $binaryPath
    
    # Install to current directory or Program Files
    $installPath = if ($args[0] -eq "--global") {
        "$env:ProgramFiles\YTMiner\ytminer.exe"
    } else {
        ".\ytminer.exe"
    }
    
    if ($args[0] -eq "--global") {
        New-Item -ItemType Directory -Path (Split-Path $installPath) -Force | Out-Null
    }
    
    Copy-Item $binaryPath $installPath -Force
    Write-Host "‚úÖ YTMiner installed successfully to: $installPath" -ForegroundColor Green
    Write-Host "ÌæØ Run 'ytminer --help' to get started" -ForegroundColor Yellow
    
} catch {
    Write-Host "‚ùå Installation failed: $_" -ForegroundColor Red
    exit 1
} finally {
    # Cleanup
    Remove-Item $tempDir -Recurse -Force -ErrorAction SilentlyContinue
}

Write-Host ""
Write-Host "Ì≥ù Usage examples:" -ForegroundColor Cyan
Write-Host "  .\ytminer.exe                    # Interactive mode" -ForegroundColor White
Write-Host "  .\ytminer.exe -k 'ai tools' -l quick  # CLI mode" -ForegroundColor White
Write-Host "  .\ytminer.exe --help            # Show help" -ForegroundColor White
