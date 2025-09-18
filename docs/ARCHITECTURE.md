# YTMiner Architecture

Technical overview of YTMiner's architecture and design principles.

## Architecture Overview

YTMiner follows **Clean Architecture** principles with clear separation of concerns:

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
└── utils/            # Utilities
```

## Design Principles

### 1. Functional Core, Imperative Shell
- **Core**: Pure functions with no side effects
- **Shell**: I/O operations and external dependencies
- **Benefits**: Testable, predictable, maintainable

### 2. Ports & Adapters (Hexagonal Architecture)
- **Ports**: Interfaces defining contracts
- **Adapters**: Implementations of external services
- **Benefits**: Testable, swappable, flexible

### 3. Domain-Driven Design
- **Domain Layer**: Business logic and rules
- **Platform Layer**: External integrations
- **Application Layer**: Use cases and orchestration
- **Benefits**: Clear boundaries, focused responsibilities

## Layer Details

### Domain Layer (`domain/`)

#### Metrics (`domain/metrics/`)
Pure functions for core business metrics:

```go
// VPD - Views Per Day
func VPD(views int64, publishedAt time.Time, now time.Time) float64

// VPDWindow - Windowed velocity
func VPDWindow(views int64, publishedAt time.Time, now time.Time, windowDays int) float64

// SlopeVPD - Velocity acceleration
func SlopeVPD(v7 float64, v30 float64) float64

// LikeRatePerThousand - Engagement rate
func LikeRatePerThousand(views int64, likes int64) float64

// FreshnessFromAges - Content freshness
func FreshnessFromAges(allAges []float64, ageDays int) float64

// NormalizeSaturation - Market saturation
func NormalizeSaturation(freq int, sampleSize int) float64
```

**Characteristics**:
- Pure functions (no side effects)
- Deterministic (same input = same output)
- Testable (easy to unit test)
- Reusable (used across different contexts)

#### Score (`domain/score/`)
Opportunity Score computation:

```go
// Compute - Main scoring function
func Compute(videos []youtube.Video, w Weights, now time.Time) []Item

// Weights - Configuration for scoring
type Weights struct {
    VPD   float64
    Like  float64
    Fresh float64
    Sat   float64
    Slope float64
}

// Item - Scored video result
type Item struct {
    Title      string
    Channel    string
    URL        string
    Score      float64
    VPD        float64
    VPD7       float64
    VPD30      float64
    Slope      float64
    LikeRate   float64
    AgeDays    int
    Saturation float64
    Reasons    []string
}
```

**Characteristics**:
- Pure computation (no I/O)
- Configurable weights
- Detailed reasoning
- Z-score normalization

### Platform Layer (`platform/`)

#### YouTube API (`platform/ytapi/`)
Interface for YouTube data access:

```go
// Client - Port interface
type Client interface {
    SearchVideos(opts youtube.SearchOptions) ([]youtube.Video, error)
    GetTranscript(ctx context.Context, videoID string) (*transcripts.Transcript, error)
}

// adapter - Adapter implementation
type adapter struct {
    c        *youtube.Client
    fetcher  transcripts.Fetcher
}
```

**Characteristics**:
- Interface-based design
- Swappable implementations
- Error handling
- Context support

#### Transcripts (`platform/transcripts/`)
Transcript fetching and caching:

```go
// Fetcher - Port interface
type Fetcher interface {
    Get(ctx context.Context, videoID string) (*Transcript, error)
}

// DefaultFetcher - Adapter implementation
type DefaultFetcher struct{}

// OAuth2Fetcher - Future OAuth 2.0 implementation
type OAuth2Fetcher struct{}
```

**Characteristics**:
- Multiple implementations
- Caching support
- Language fallback
- Error handling

### Application Layer (`analysis/`)

#### Analyzer (`analysis/analyzer.go`)
Orchestrates different analysis types:

```go
type Analyzer struct {
    videos []youtube.Video
    cfg    config.AppConfig
}

// Analysis methods
func (a *Analyzer) AnalyzeGrowthPatterns() GrowthResult
func (a *Analyzer) AnalyzeTitles() TitleResult
func (a *Analyzer) AnalyzeCompetitors() CompetitorResult
func (a *Analyzer) AnalyzeTemporal() TemporalResult
func (a *Analyzer) AnalyzeKeywords() KeywordResult
func (a *Analyzer) AnalyzeOpportunityScore() []OpportunityItem
func (a *Analyzer) GenerateExecutiveReport() ExecutiveResult
```

**Characteristics**:
- Use case orchestration
- Configuration-driven
- Result aggregation
- Error handling

### Configuration Layer (`config/`)

#### AppConfig (`config/config.go`)
Application configuration management:

```go
type AppConfig struct {
    // Search parameters
    DefaultRegion    string
    DefaultDuration  string
    DefaultTimeRange string
    DefaultOrder     string
    
    // Transcript settings
    WithTranscripts  bool
    TranscriptLangs  string
    CacheDir         string
    
    // Opportunity Score weights
    OppWeightVPD     float64
    OppWeightLike    float64
    OppWeightFresh   float64
    OppWeightSatPen  float64
    OppWeightSlope   float64
}
```

**Characteristics**:
- Environment variable support
- Default values
- Validation
- Profile support

#### Profiles (`config/profiles.go`)
Predefined weight profiles:

```go
type WeightProfile struct {
    Name        string
    Description string
    Weights     ProfileWeights
}

var WeightProfiles = map[string]WeightProfile{
    "exploration": {...},
    "evergreen":   {...},
    "trending":    {...},
    "balanced":    {...},
}
```

**Characteristics**:
- Predefined strategies
- Easy switching
- Documentation
- Extensible

### UI Layer (`ui/`)

#### Display (`ui/display.go`)
User interface components:

```go
// Display functions
func DisplayInfo(message string)
func DisplayWarning(message string)
func DisplayError(message string)
func DisplayLoading(message string) func()
func DisplayOpportunityScore(items []analysis.OpportunityItem)
func DisplayGrowthPatterns(result analysis.GrowthResult)
// ... other display functions
```

**Characteristics**:
- Consistent formatting
- Color coding
- Progress indicators
- Error handling

## Data Flow

### 1. Input Processing
```
User Input → CLI Parser → Configuration → Search Options
```

### 2. Data Fetching
```
Search Options → YouTube API → Video Data → Transcript Fetcher → Enriched Data
```

### 3. Analysis Processing
```
Enriched Data → Domain Functions → Analysis Results → UI Display
```

### 4. Caching
```
Transcript Data → Cache Storage → Future Requests
```

## Error Handling

### Strategy
1. **Fail fast**: Validate inputs early
2. **Graceful degradation**: Continue with available data
3. **User feedback**: Clear error messages
4. **Logging**: Detailed error information

### Error Types
- **Configuration errors**: Invalid settings
- **API errors**: YouTube API failures
- **Network errors**: Connection issues
- **Data errors**: Invalid video data

## Testing Strategy

### Unit Tests
- **Domain functions**: Pure function testing
- **Configuration**: Settings validation
- **Utilities**: Helper function testing

### Integration Tests
- **Platform adapters**: External service integration
- **Analysis orchestration**: End-to-end workflows
- **UI components**: Display functionality

### Golden Tests
- **Opportunity Score**: Deterministic output validation
- **Analysis results**: Consistent formatting
- **CLI output**: User experience validation

## Performance Considerations

### Memory Usage
- **Video data**: ~1KB per video
- **Transcripts**: ~100KB per video
- **Analysis results**: ~10KB per result

### API Quotas
- **Search requests**: 100 units each
- **Video details**: 1 unit each
- **Transcripts**: No quota (rate limited)

### Caching
- **Transcript cache**: Local file storage
- **Analysis cache**: In-memory during session
- **Configuration cache**: Loaded once per session

## Extensibility

### Adding New Metrics
1. Add function to `domain/metrics/`
2. Add to Opportunity Score calculation
3. Add to analysis results
4. Add to UI display

### Adding New Analysis Types
1. Add method to `analysis/analyzer.go`
2. Add result type
3. Add UI display function
4. Add CLI option

### Adding New Data Sources
1. Create new platform adapter
2. Implement port interface
3. Add to application layer
4. Update configuration

## Security Considerations

### API Key Management
- **Environment variables**: Secure storage
- **No hardcoding**: Keys not in source code
- **Rotation**: Regular key updates

### Data Privacy
- **Local processing**: No data sent to external services
- **Cache security**: Local file permissions
- **Logging**: No sensitive data in logs

### Input Validation
- **Sanitization**: Clean user inputs
- **Bounds checking**: Validate numeric inputs
- **Type safety**: Strong typing throughout

## Future Considerations

### OAuth 2.0 Integration
- **YouTube Data API v3**: Official transcript access
- **Rate limits**: Higher quotas
- **Authentication**: User consent flow

### Scalability
- **Concurrent processing**: Parallel analysis
- **Database storage**: Persistent data
- **API optimization**: Batch requests

### Advanced Features
- **Machine learning**: Predictive analytics
- **Real-time updates**: Live data feeds
- **Collaboration**: Multi-user support

## Development Guidelines

### Code Organization
- **Single responsibility**: One purpose per function
- **Clear interfaces**: Well-defined contracts
- **Error handling**: Comprehensive error management
- **Documentation**: Clear comments and docs

### Testing Requirements
- **Unit tests**: All domain functions
- **Integration tests**: Critical workflows
- **Golden tests**: Stable outputs
- **Coverage**: 80%+ for critical paths

### Performance Standards
- **Response time**: <3 minutes for deep analysis
- **Memory usage**: <100MB for typical analysis
- **API efficiency**: Minimize quota usage
- **Cache hit rate**: >80% for repeated analysis
