# Contributing to YTMiner

Thank you for your interest in contributing to YTMiner! This guide will help you get started.

## Development Setup

### Prerequisites
- Go 1.24+ installed
- Git installed
- YouTube Data API v3 key (for testing)

### Getting Started
1. **Fork the repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/yourusername/ytminer.git
   cd ytminer
   ```

2. **Set up environment**
   ```bash
   # Copy environment template
   cp env.example .env
   
   # Add your YouTube API key
   echo "YOUTUBE_API_KEY=your_api_key_here" >> .env
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Build the project**
   ```bash
   go build -o ytminer.exe .
   ```

5. **Run tests**
   ```bash
   go test ./...
   ```

## Project Structure

```
ytminer/
├── domain/           # Business logic (pure functions)
│   ├── metrics/      # Core metrics (VPD, Slope, etc.)
│   └── score/        # Opportunity Score computation
├── platform/         # External integrations
│   ├── ytapi/        # YouTube API adapter
│   └── transcripts/  # Transcript fetching & caching
├── analysis/         # Analysis orchestration
├── config/           # Configuration management
├── ui/               # User interface
├── utils/            # Utilities
├── docs/             # Documentation
└── e2e/              # End-to-end tests
```

## Architecture Principles

### 1. Clean Architecture
- **Domain Layer**: Pure business logic, no dependencies
- **Platform Layer**: External integrations (APIs, databases)
- **Application Layer**: Use cases and orchestration
- **UI Layer**: User interface and presentation

### 2. Functional Core, Imperative Shell
- **Core**: Pure functions with no side effects
- **Shell**: I/O operations and external dependencies
- **Benefits**: Testable, predictable, maintainable

### 3. Ports & Adapters
- **Ports**: Interfaces defining contracts
- **Adapters**: Implementations of external services
- **Benefits**: Testable, swappable, flexible

## Development Workflow

### 1. Create a Feature Branch
```bash
git checkout -b feature/your-feature-name
```

### 2. Make Changes
- Follow the architecture principles
- Write tests for new functionality
- Update documentation as needed

### 3. Test Your Changes
```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./domain/metrics

# Run with verbose output
go test ./... -v

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### 4. Format Code
```bash
# Format all Go files
gofmt -s -w .

# Check formatting
gofmt -s -l .
```

### 5. Lint Code
```bash
# Run go vet
go vet ./...

# Run go mod tidy
go mod tidy
```

### 6. Commit Changes
```bash
# Add changes
git add .

# Commit with descriptive message
git commit -m "feat: add new metric calculation"

# Push to your fork
git push origin feature/your-feature-name
```

### 7. Create Pull Request
- Go to your fork on GitHub
- Click "New Pull Request"
- Fill out the PR template
- Request review from maintainers

## Coding Standards

### Go Style
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Write clear, concise comments
- Keep functions small and focused

### Error Handling
```go
// Good: Handle errors explicitly
result, err := someFunction()
if err != nil {
    return fmt.Errorf("failed to process: %w", err)
}

// Bad: Ignore errors
result, _ := someFunction()
```

### Testing
```go
// Good: Test with clear assertions
func TestVPD(t *testing.T) {
    result := metrics.VPD(1000, time.Now().AddDate(0, 0, -10), time.Now())
    expected := 100.0
    if result != expected {
        t.Fatalf("expected %v, got %v", expected, result)
    }
}

// Bad: Test without clear assertions
func TestVPD(t *testing.T) {
    result := metrics.VPD(1000, time.Now().AddDate(0, 0, -10), time.Now())
    // No assertion
}
```

### Documentation
```go
// Good: Clear function documentation
// VPD computes Views Per Day given total views and publication time.
// It normalizes total views by the video's age, allowing comparison
// of new and older videos on a level playing field.
func VPD(views int64, publishedAt time.Time, now time.Time) float64 {
    // Implementation
}

// Bad: No documentation
func VPD(views int64, publishedAt time.Time, now time.Time) float64 {
    // Implementation
}
```

## Testing Guidelines

### Unit Tests
- Test all domain functions
- Use table-driven tests for multiple cases
- Test edge cases and error conditions
- Aim for good test coverage on critical paths

### Integration Tests
- Test platform adapters
- Test analysis orchestration
- Use real API keys when possible
- Mock external services when needed

### Golden Tests
- Use for stable outputs
- Test CLI output formatting
- Test analysis result consistency
- Update when expected output changes

### Example Test Structure
```go
func TestVPD_Basic(t *testing.T) {
    tests := []struct {
        name     string
        views    int64
        published time.Time
        now      time.Time
        expected float64
    }{
        {
            name:     "basic calculation",
            views:    1000,
            published: time.Now().AddDate(0, 0, -10),
            now:      time.Now(),
            expected: 100.0,
        },
        {
            name:     "zero views",
            views:    0,
            published: time.Now().AddDate(0, 0, -10),
            now:      time.Now(),
            expected: 0.0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := metrics.VPD(tt.views, tt.published, tt.now)
            if result != tt.expected {
                t.Fatalf("expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

## Pull Request Guidelines

### PR Title
Use conventional commit format:
- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `test:` - Test additions/changes
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `ci:` - CI/CD changes

### PR Description
Include:
- **What**: What changes were made
- **Why**: Why the changes were needed
- **How**: How the changes work
- **Testing**: How the changes were tested
- **Breaking changes**: Any breaking changes

### Example PR Description
```markdown
## What
Add VPD7 and VPD30 metrics for better velocity analysis.

## Why
Current VPD calculation doesn't distinguish between recent and overall performance.
VPD7/VPD30 provide better insights into momentum and acceleration.

## How
- Added VPDWindow function to calculate windowed velocity
- Added SlopeVPD function to calculate acceleration
- Integrated new metrics into Opportunity Score
- Added comprehensive tests

## Testing
- Unit tests for all new functions
- Integration tests with Opportunity Score
- Golden tests for consistent output
- Manual testing with real data

## Breaking Changes
None - all changes are additive.
```

## Code Review Process

### For Contributors
1. **Self-review**: Review your own PR before submitting
2. **Address feedback**: Respond to review comments
3. **Update tests**: Add tests for new functionality
4. **Update docs**: Update documentation as needed

### For Reviewers
1. **Check functionality**: Ensure code works as intended
2. **Check tests**: Verify adequate test coverage
3. **Check style**: Ensure code follows standards
4. **Check docs**: Verify documentation is updated

## Common Issues

### Import Errors
```bash
# Clean module cache
go clean -modcache

# Download dependencies
go mod download

# Tidy modules
go mod tidy
```

### Test Failures
```bash
# Run tests with verbose output
go test ./... -v

# Run specific test
go test -run TestVPD ./domain/metrics

# Check test coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Build Errors
```bash
# Check Go version
go version

# Clean build cache
go clean -cache

# Rebuild
go build -o ytminer.exe .
```

## Getting Help

### Documentation
- [Usage Guide](USAGE.md) - How to use YTMiner
- [Analysis Types](ANALYSIS.md) - Detailed analysis explanations
- [Configuration](CONFIGURATION.md) - Configuration options
- [Architecture](ARCHITECTURE.md) - Technical architecture

### Community
- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: Questions and general discussion
- **Pull Requests**: Code contributions

### Contact
- **Maintainer**: [@nuzael](https://github.com/nuzael)
- **Repository**: [ytminer](https://github.com/nuzael/ytminer)

## License

By contributing to YTMiner, you agree that your contributions will be licensed under the MIT License.
