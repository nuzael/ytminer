package tui

import (
"fmt"
"strings"

"ytminer/analysis"

"github.com/charmbracelet/bubbles/table"
tea "github.com/charmbracelet/bubbletea"
"github.com/charmbracelet/lipgloss"
)

// ShowViewer opens a responsive TUI to navigate analysis results via tabs
func ShowViewer(g analysis.GrowthPattern, t analysis.TitleAnalysis, c analysis.CompetitorAnalysis, temp analysis.TemporalAnalysis, k analysis.KeywordAnalysis, o []analysis.OpportunityItem, e analysis.ExecutiveReport) error {
m := newModel(g, t, c, temp, k, o, e)
p := tea.NewProgram(m, tea.WithAltScreen())
_, err := p.Run()
return err
}

type tab int

const (
tabGrowth tab = iota
tabTitles
tabCompetitors
tabTemporal
tabKeywords
tabOpportunity
tabExecutive
)

type model struct {
active tab
width  int
height int

// Tables for each tab
growthMain      table.Model
growthInsights  table.Model
titlesWords     table.Model
titlesPhrases   table.Model
competitorsTop  table.Model
competitorsInfo table.Model
temporalHours   table.Model
temporalDays    table.Model
keywordsTrending table.Model
keywordsCore    table.Model
opportunityTop  table.Model
executiveSummary table.Model

// Focus management
subFocus int // 0 = main table, 1 = insights table
}

var (
// Styles
headerStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
tabStyle       = lipgloss.NewStyle().Padding(0, 1)
activeTabStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Padding(0, 1)
titleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
subtitleStyle  = lipgloss.NewStyle().Faint(true)
footerStyle    = lipgloss.NewStyle().Faint(true)
textStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))

// Table border style like in documentation
baseStyle = lipgloss.NewStyle().
BorderStyle(lipgloss.NormalBorder()).
BorderForeground(lipgloss.Color("240"))
)

// ============================================================================
// GROWTH PATTERN ANALYSIS - ✅ IMPLEMENTED
// ============================================================================

func createGrowthMainTable(g analysis.GrowthPattern) table.Model {
columns := []table.Column{
{Title: "Title", Width: 40},
{Title: "Channel", Width: 18},
{Title: "Views", Width: 12},
{Title: "VPD", Width: 8},
{Title: "Eng%", Width: 6},
}

rows := make([]table.Row, len(g.TopPerformers))
for i, v := range g.TopPerformers {
rows[i] = table.Row{
v.Title,
v.Channel,
formatInt(v.Views),
formatFloat(v.VPD),
formatPercent(v.Engagement),
}
}

t := table.New(
table.WithColumns(columns),
table.WithRows(rows),
table.WithFocused(true),
table.WithHeight(5),
)
t.SetStyles(createTableStyles())
return t
}

func createGrowthInsightsTable(g analysis.GrowthPattern) table.Model {
columns := []table.Column{
{Title: "Insights", Width: 80},
}

rows := make([]table.Row, len(g.Insights))
for i, insight := range g.Insights {
rows[i] = table.Row{insight}
}

t := table.New(
table.WithColumns(columns),
table.WithRows(rows),
table.WithFocused(false),
table.WithHeight(5),
)
t.SetStyles(createTableStyles())
return t
}

func (m model) renderGrowthContent() string {
var content strings.Builder

// Summary stats
content.WriteString(textStyle.Render(fmt.Sprintf("Total Videos: %d | Avg Views: %.0f | Niche Velocity: %.2f", 
len(m.growthMain.Rows()), 0.0, 0.0)))
content.WriteString("\n\n")

// Main table
content.WriteString(subtitleStyle.Render("Top Performers"))
content.WriteString("\n")
content.WriteString(baseStyle.Render(m.growthMain.View()))
content.WriteString("\n\n")

// Insights table
content.WriteString(subtitleStyle.Render("Insights"))
content.WriteString("\n")
content.WriteString(baseStyle.Render(m.growthInsights.View()))

return content.String()
}

// ============================================================================
// TITLE PATTERN ANALYSIS - ✅ IMPLEMENTED
// ============================================================================

func createTitlesWordsTable(t analysis.TitleAnalysis) table.Model {
columns := []table.Column{
{Title: "Word", Width: 20},
{Title: "Count", Width: 8},
}

rows := make([]table.Row, len(t.CommonWords))
for i, wc := range t.CommonWords {
rows[i] = table.Row{wc.Word, fmt.Sprintf("%d", wc.Count)}
}

table := table.New(
table.WithColumns(columns),
table.WithRows(rows),
table.WithFocused(false),
table.WithHeight(5),
)
table.SetStyles(createTableStyles())
return table
}

func createTitlesPhrasesTable(t analysis.TitleAnalysis) table.Model {
columns := []table.Column{
{Title: "Phrase", Width: 36},
{Title: "Count", Width: 8},
}

rows := make([]table.Row, len(t.CommonPhrases))
for i, pc := range t.CommonPhrases {
rows[i] = table.Row{pc.Phrase, fmt.Sprintf("%d", pc.Count)}
}

table := table.New(
table.WithColumns(columns),
table.WithRows(rows),
table.WithFocused(false),
table.WithHeight(5),
)
table.SetStyles(createTableStyles())
return table
}

func (m model) renderTitlesContent() string {
var content strings.Builder

// Words table
content.WriteString(subtitleStyle.Render("Common Words"))
content.WriteString("\n")
content.WriteString(baseStyle.Render(m.titlesWords.View()))
content.WriteString("\n\n")

// Phrases table
content.WriteString(subtitleStyle.Render("Common Phrases"))
content.WriteString("\n")
content.WriteString(baseStyle.Render(m.titlesPhrases.View()))

return content.String()
}

// ============================================================================
// COMPETITOR ANALYSIS - ✅ IMPLEMENTED
// ============================================================================

func createCompetitorsTopTable(c analysis.CompetitorAnalysis) table.Model {
columns := []table.Column{
{Title: "Channel", Width: 24},
{Title: "AvgVPD", Width: 10},
{Title: "AvgViews", Width: 12},
{Title: "Eng%", Width: 8},
{Title: "Rising", Width: 8},
}

rows := make([]table.Row, len(c.TopChannels))
for i, ch := range c.TopChannels {
rows[i] = table.Row{
ch.Channel,
formatFloat(ch.AvgVPD),
formatFloat(ch.AvgViews),
formatPercent(ch.Engagement),
boolMark(ch.IsRisingStar),
}
}

table := table.New(
table.WithColumns(columns),
table.WithRows(rows),
table.WithFocused(false),
table.WithHeight(5),
)
table.SetStyles(createTableStyles())
return table
}

func createCompetitorsInfoTable(c analysis.CompetitorAnalysis) table.Model {
columns := []table.Column{
{Title: "Insights", Width: 80},
}

rows := make([]table.Row, len(c.Insights))
for i, insight := range c.Insights {
rows[i] = table.Row{insight}
}

table := table.New(
table.WithColumns(columns),
table.WithRows(rows),
table.WithFocused(false),
table.WithHeight(5),
)
table.SetStyles(createTableStyles())
return table
}

func (m model) renderCompetitorsContent() string {
var content strings.Builder

// Top channels table
content.WriteString(subtitleStyle.Render("Top Channels"))
content.WriteString("\n")
content.WriteString(baseStyle.Render(m.competitorsTop.View()))
content.WriteString("\n\n")

// Insights table
content.WriteString(subtitleStyle.Render("Insights"))
content.WriteString("\n")
content.WriteString(baseStyle.Render(m.competitorsInfo.View()))

return content.String()
}

// ============================================================================
// TEMPORAL ANALYSIS - ✅ IMPLEMENTED
// ============================================================================

func createTemporalHoursTable(t analysis.TemporalAnalysis) table.Model {
columns := []table.Column{
{Title: "Hour", Width: 8},
{Title: "Avg Views", Width: 12},
{Title: "Avg Likes", Width: 12},
{Title: "Engagement", Width: 12},
}

rows := make([]table.Row, len(t.BestHours))
for i, h := range t.BestHours {
rows[i] = table.Row{
fmt.Sprintf("%d:00", h.Hour),
formatFloat(h.AvgViews),
formatFloat(h.AvgLikes),
formatPercent(h.Engagement),
}
}

table := table.New(
table.WithColumns(columns),
table.WithRows(rows),
table.WithFocused(false),
table.WithHeight(5),
)
table.SetStyles(createTableStyles())
return table
}

func createTemporalDaysTable(t analysis.TemporalAnalysis) table.Model {
columns := []table.Column{
{Title: "Day", Width: 12},
{Title: "Avg Views", Width: 12},
{Title: "Avg Likes", Width: 12},
{Title: "Engagement", Width: 12},
}

rows := make([]table.Row, len(t.BestDays))
for i, d := range t.BestDays {
rows[i] = table.Row{
d.Day,
formatFloat(d.AvgViews),
formatFloat(d.AvgLikes),
formatPercent(d.Engagement),
}
}

table := table.New(
table.WithColumns(columns),
table.WithRows(rows),
table.WithFocused(false),
table.WithHeight(5),
)
table.SetStyles(createTableStyles())
return table
}

func (m model) renderTemporalContent() string {
var content strings.Builder

// Hours table
content.WriteString(subtitleStyle.Render("Best Posting Hours"))
content.WriteString("\n")
content.WriteString(baseStyle.Render(m.temporalHours.View()))
content.WriteString("\n\n")

// Days table
content.WriteString(subtitleStyle.Render("Best Posting Days"))
content.WriteString("\n")
content.WriteString(baseStyle.Render(m.temporalDays.View()))

return content.String()
}

// ============================================================================
// KEYWORD ANALYSIS - ✅ IMPLEMENTED
// ============================================================================

func createKeywordsTrendingTable(k analysis.KeywordAnalysis) table.Model {
columns := []table.Column{
{Title: "Keyword", Width: 20},
{Title: "Freq", Width: 6},
{Title: "Avg Views", Width: 12},
{Title: "Avg VPD", Width: 10},
{Title: "Eng%", Width: 8},
}

rows := make([]table.Row, len(k.TrendingKeywords))
for i, kw := range k.TrendingKeywords {
rows[i] = table.Row{
kw.Keyword,
fmt.Sprintf("%d", kw.Frequency),
formatFloat(kw.AvgViews),
formatFloat(kw.AvgVPD),
formatPercent(kw.Engagement),
}
}

table := table.New(
table.WithColumns(columns),
table.WithRows(rows),
table.WithFocused(false),
table.WithHeight(5),
)
table.SetStyles(createTableStyles())
return table
}

func createKeywordsCoreTable(k analysis.KeywordAnalysis) table.Model {
columns := []table.Column{
{Title: "Keyword", Width: 20},
{Title: "Freq", Width: 6},
{Title: "Avg Views", Width: 12},
{Title: "Avg VPD", Width: 10},
{Title: "Eng%", Width: 8},
}

rows := make([]table.Row, len(k.CoreKeywords))
for i, kw := range k.CoreKeywords {
rows[i] = table.Row{
kw.Keyword,
fmt.Sprintf("%d", kw.Frequency),
formatFloat(kw.AvgViews),
formatFloat(kw.AvgVPD),
formatPercent(kw.Engagement),
}
}

table := table.New(
table.WithColumns(columns),
table.WithRows(rows),
table.WithFocused(false),
table.WithHeight(5),
)
table.SetStyles(createTableStyles())
return table
}

func (m model) renderKeywordsContent() string {
var content strings.Builder

// Trending keywords table
content.WriteString(subtitleStyle.Render("Trending Keywords (by VPD)"))
content.WriteString("\n")
content.WriteString(baseStyle.Render(m.keywordsTrending.View()))
content.WriteString("\n\n")

// Core keywords table
content.WriteString(subtitleStyle.Render("Core Keywords (by Frequency)"))
content.WriteString("\n")
content.WriteString(baseStyle.Render(m.keywordsCore.View()))

return content.String()
}

// ============================================================================
// OPPORTUNITY SCORE ANALYSIS - ✅ IMPLEMENTED
// ============================================================================

func createOpportunityTopTable(o []analysis.OpportunityItem) table.Model {
columns := []table.Column{
{Title: "Title", Width: 35},
{Title: "Channel", Width: 15},
{Title: "Score", Width: 8},
{Title: "VPD", Width: 8},
{Title: "Like%", Width: 8},
{Title: "Age", Width: 6},
}

rows := make([]table.Row, len(o))
for i, item := range o {
rows[i] = table.Row{
item.Title,
item.Channel,
formatFloat(item.Score),
formatFloat(item.VPD),
formatPercent(item.LikeRate),
fmt.Sprintf("%dd", item.AgeDays),
}
}

table := table.New(
table.WithColumns(columns),
table.WithRows(rows),
table.WithFocused(false),
table.WithHeight(5),
)
table.SetStyles(createTableStyles())
return table
}

func (m model) renderOpportunityContent() string {
var content strings.Builder

// Opportunity scores table
content.WriteString(subtitleStyle.Render("Top Opportunity Scores"))
content.WriteString("\n")
content.WriteString(baseStyle.Render(m.opportunityTop.View()))

return content.String()
}

// ============================================================================
// EXECUTIVE REPORT ANALYSIS - ✅ IMPLEMENTED
// ============================================================================

func createExecutiveSummaryTable(e analysis.ExecutiveReport) table.Model {
columns := []table.Column{
{Title: "Section", Width: 20},
{Title: "Content", Width: 60},
}

var rows []table.Row

// Summary
rows = append(rows, table.Row{"Summary", e.Summary})

// Key Insights
for _, insight := range e.KeyInsights {
rows = append(rows, table.Row{"Key Insight", insight})
}

// Recommendations
for _, rec := range e.Recommendations {
rows = append(rows, table.Row{"Recommendation", rec})
}

table := table.New(
table.WithColumns(columns),
table.WithRows(rows),
table.WithFocused(false),
table.WithHeight(5),
)
table.SetStyles(createTableStyles())
return table
}

func (m model) renderExecutiveContent() string {
var content strings.Builder

// Executive summary table
content.WriteString(subtitleStyle.Render("Executive Summary"))
content.WriteString("\n")
content.WriteString(baseStyle.Render(m.executiveSummary.View()))

return content.String()
}

// ============================================================================
// SHARED UTILITIES
// ============================================================================

func createTableStyles() table.Styles {
s := table.DefaultStyles()
s.Header = s.Header.
BorderStyle(lipgloss.NormalBorder()).
BorderForeground(lipgloss.Color("240")).
BorderBottom(true).
Bold(false)
s.Selected = s.Selected.
Foreground(lipgloss.Color("229")).
Background(lipgloss.Color("57")).
Bold(false)
return s
}

func newModel(g analysis.GrowthPattern, t analysis.TitleAnalysis, c analysis.CompetitorAnalysis, temp analysis.TemporalAnalysis, k analysis.KeywordAnalysis, o []analysis.OpportunityItem, e analysis.ExecutiveReport) model {
// Growth tables
growthMain := createGrowthMainTable(g)
growthInsights := createGrowthInsightsTable(g)

// Titles tables
titlesWords := createTitlesWordsTable(t)
titlesPhrases := createTitlesPhrasesTable(t)

// Competitors tables
competitorsTop := createCompetitorsTopTable(c)
competitorsInfo := createCompetitorsInfoTable(c)

// Temporal tables
temporalHours := createTemporalHoursTable(temp)
temporalDays := createTemporalDaysTable(temp)

// Keywords tables
keywordsTrending := createKeywordsTrendingTable(k)
keywordsCore := createKeywordsCoreTable(k)

// Opportunity table
opportunityTop := createOpportunityTopTable(o)

// Executive table
executiveSummary := createExecutiveSummaryTable(e)

return model{
active:           tabGrowth,
growthMain:       growthMain,
growthInsights:   growthInsights,
titlesWords:      titlesWords,
titlesPhrases:    titlesPhrases,
competitorsTop:   competitorsTop,
competitorsInfo:  competitorsInfo,
temporalHours:    temporalHours,
temporalDays:     temporalDays,
keywordsTrending: keywordsTrending,
keywordsCore:     keywordsCore,
opportunityTop:   opportunityTop,
executiveSummary: executiveSummary,
subFocus:         0,
}
}

// ============================================================================
// TUI CORE LOGIC
// ============================================================================

func (m model) Init() tea.Cmd {
return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
var cmd tea.Cmd
var cmds []tea.Cmd

switch msg := msg.(type) {
case tea.WindowSizeMsg:
m.width = msg.Width
m.height = msg.Height
return m, nil

case tea.KeyMsg:
switch msg.String() {
case "ctrl+c", "q":
return m, tea.Quit
case "1":
m.active = tabGrowth
m.subFocus = 0
m.setFocus()
case "2":
m.active = tabTitles
m.subFocus = 0
m.setFocus()
case "3":
m.active = tabCompetitors
m.subFocus = 0
m.setFocus()
case "4":
m.active = tabTemporal
m.subFocus = 0
m.setFocus()
case "5":
m.active = tabKeywords
m.subFocus = 0
m.setFocus()
case "6":
m.active = tabOpportunity
m.subFocus = 0
m.setFocus()
case "7":
m.active = tabExecutive
m.subFocus = 0
m.setFocus()
case "tab":
m.toggleSubFocus()
default:
// Delegate to active table
switch m.active {
case tabGrowth:
if m.subFocus == 0 {
m.growthMain, cmd = m.growthMain.Update(msg)
} else {
m.growthInsights, cmd = m.growthInsights.Update(msg)
}
case tabTitles:
if m.subFocus == 0 {
m.titlesWords, cmd = m.titlesWords.Update(msg)
} else {
m.titlesPhrases, cmd = m.titlesPhrases.Update(msg)
}
case tabCompetitors:
if m.subFocus == 0 {
m.competitorsTop, cmd = m.competitorsTop.Update(msg)
} else {
m.competitorsInfo, cmd = m.competitorsInfo.Update(msg)
}
case tabTemporal:
if m.subFocus == 0 {
m.temporalHours, cmd = m.temporalHours.Update(msg)
} else {
m.temporalDays, cmd = m.temporalDays.Update(msg)
}
case tabKeywords:
if m.subFocus == 0 {
m.keywordsTrending, cmd = m.keywordsTrending.Update(msg)
} else {
m.keywordsCore, cmd = m.keywordsCore.Update(msg)
}
case tabOpportunity:
m.opportunityTop, cmd = m.opportunityTop.Update(msg)
case tabExecutive:
m.executiveSummary, cmd = m.executiveSummary.Update(msg)
}
cmds = append(cmds, cmd)
}
}

return m, tea.Batch(cmds...)
}

func (m *model) toggleSubFocus() {
m.subFocus = 1 - m.subFocus
m.setFocus()
}

func (m *model) setFocus() {
// Blur all tables
m.growthMain.Blur()
m.growthInsights.Blur()
m.titlesWords.Blur()
m.titlesPhrases.Blur()
m.competitorsTop.Blur()
m.competitorsInfo.Blur()
m.temporalHours.Blur()
m.temporalDays.Blur()
m.keywordsTrending.Blur()
m.keywordsCore.Blur()
m.opportunityTop.Blur()
m.executiveSummary.Blur()

// Focus active table
switch m.active {
case tabGrowth:
if m.subFocus == 0 {
m.growthMain.Focus()
} else {
m.growthInsights.Focus()
}
case tabTitles:
if m.subFocus == 0 {
m.titlesWords.Focus()
} else {
m.titlesPhrases.Focus()
}
case tabCompetitors:
if m.subFocus == 0 {
m.competitorsTop.Focus()
} else {
m.competitorsInfo.Focus()
}
case tabTemporal:
if m.subFocus == 0 {
m.temporalHours.Focus()
} else {
m.temporalDays.Focus()
}
case tabKeywords:
if m.subFocus == 0 {
m.keywordsTrending.Focus()
} else {
m.keywordsCore.Focus()
}
case tabOpportunity:
m.opportunityTop.Focus()
case tabExecutive:
m.executiveSummary.Focus()
}
}

func (m model) View() string {
if m.width <= 0 || m.height <= 0 {
return "Loading..."
}

var content strings.Builder

// Header with tabs
content.WriteString(m.renderHeader())
content.WriteString("\n\n")

// Content based on active tab
content.WriteString(m.renderContent())

// Footer
content.WriteString("\n")
content.WriteString(m.renderFooter())

return content.String()
}

func (m model) renderHeader() string {
var header strings.Builder

// Tabs
tabs := []string{
m.renderTab("1 Growth", m.active == tabGrowth),
m.renderTab("2 Titles", m.active == tabTitles),
m.renderTab("3 Competitors", m.active == tabCompetitors),
m.renderTab("4 Temporal", m.active == tabTemporal),
m.renderTab("5 Keywords", m.active == tabKeywords),
m.renderTab("6 Opportunity", m.active == tabOpportunity),
m.renderTab("7 Executive", m.active == tabExecutive),
}
header.WriteString(strings.Join(tabs, "  "))

// Title
title := m.getTabTitle()
if title != "" {
header.WriteString("  ")
header.WriteString(titleStyle.Render(title))
}

return header.String()
}

func (m model) renderTab(label string, active bool) string {
if active {
return activeTabStyle.Render("[" + label + "]")
}
return tabStyle.Render(label)
}

func (m model) getTabTitle() string {
switch m.active {
case tabGrowth:
return "Growth Pattern Analysis"
case tabTitles:
return "Title Pattern Analysis"
case tabCompetitors:
return "Competitor Analysis"
case tabTemporal:
return "Temporal Analysis"
case tabKeywords:
return "Keyword Analysis"
case tabOpportunity:
return "Opportunity Score Analysis"
case tabExecutive:
return "Executive Report Analysis"
default:
return ""
}
}

func (m model) renderContent() string {
switch m.active {
case tabGrowth:
return m.renderGrowthContent()
case tabTitles:
return m.renderTitlesContent()
case tabCompetitors:
return m.renderCompetitorsContent()
case tabTemporal:
return m.renderTemporalContent()
case tabKeywords:
return m.renderKeywordsContent()
case tabOpportunity:
return m.renderOpportunityContent()
case tabExecutive:
return m.renderExecutiveContent()
default:
return ""
}
}

func (m model) renderFooter() string {
return footerStyle.Render("TAB: switch focus • ↑/↓: navigate • 1-7: tabs • q: quit")
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func formatInt(v int64) string       { return fmt.Sprintf("%d", v) }
func formatFloat(v float64) string   { return fmt.Sprintf("%.2f", v) }
func formatPercent(v float64) string { return fmt.Sprintf("%.1f%%", v) }
func boolMark(b bool) string {
if b {
return "✔"
}
return ""
}
