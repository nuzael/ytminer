# YTMiner Test Runner for PowerShell
# Uses gotestsum for beautiful test output with coverage

Write-Host "Running YTMiner Tests with gotestsum..." -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Green

# Check if gotestsum is available
if (!(Get-Command gotestsum -ErrorAction SilentlyContinue)) {
    Write-Host "gotestsum not found. Installing..." -ForegroundColor Yellow
    go install gotest.tools/gotestsum@latest
}

# Run tests with gotestsum and coverage (using testdox format)
gotestsum --format testdox -- -coverprofile="coverage.out" ./...

# Show coverage
Write-Host ""
Write-Host "Test Coverage:" -ForegroundColor Cyan
Write-Host "=================" -ForegroundColor Cyan
go tool cover -func="coverage.out" | Select-Object -Last 1

# Show coverage in browser (optional)
Write-Host ""
Write-Host "To see detailed coverage in browser:" -ForegroundColor Yellow
Write-Host "   go tool cover -html=coverage.out" -ForegroundColor Gray

# Clean up
Remove-Item -Path "coverage.out" -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "All tests completed!" -ForegroundColor Green
