#!/bin/bash

# YTMiner Test Runner
# Uses gotestsum for beautiful test output with coverage

echo "� Running YTMiner Tests with gotestsum..."
echo "=========================================="

# Check if gotestsum is available
if ! command -v gotestsum &> /dev/null; then
    echo "❌ gotestsum not found. Installing..."
    go install gotest.tools/gotestsum@latest
fi

# Run tests with gotestsum and coverage (using testdox format)
gotestsum --format testdox -- -coverprofile=coverage.out ./...

# Show coverage
echo ""
echo "� Test Coverage:"
echo "================="
go tool cover -func=coverage.out | tail -1

# Show coverage in browser (optional)
echo ""
echo "� To see detailed coverage in browser:"
echo "   go tool cover -html=coverage.out"

# Clean up
rm -f coverage.out

echo ""
echo "✅ All tests completed!"
