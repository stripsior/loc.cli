package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stripsior/loc.cli/internal/models"
)

type state int

const (
	loadingState state = iota
	resultState
)

type view int

const (
	summaryView view = iota
	languagesView
	filesView
	authorsView
)

type model struct {
	state          state
	result         *models.AnalysisResult
	currentView    view
	table          table.Model
	spinner        spinner.Model
	progress       progress.Model
	progressValue  float64
	width          int
	height         int
	err            error
	analysisFunc   func() (*models.AnalysisResult, error)
	showTransition bool
	transitionStep int
}

type analysisCompleteMsg struct {
	result *models.AnalysisResult
	err    error
}

type progressMsg float64

type tickMsg time.Time

var (
	// Minimalistic color scheme
	primaryColor = lipgloss.Color("#e95268")
	grayDark     = lipgloss.Color("#4B5563")
	grayMedium   = lipgloss.Color("#6B7280")
	grayLight    = lipgloss.Color("#9CA3AF")
	grayLighter  = lipgloss.Color("#D1D5DB")
	white        = lipgloss.Color("#FFFFFF")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(white).
			Bold(true).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(grayLight).
			MarginBottom(2)

	tabStyle = lipgloss.NewStyle().
			Foreground(grayMedium).
			Padding(0, 2, 1, 2)

	activeTabStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 2).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(primaryColor)

	tabSeparator = lipgloss.NewStyle().
			Foreground(grayDark).
			SetString("│")

	dividerStyle = lipgloss.NewStyle().
			Foreground(grayDark).
			MarginTop(1).
			MarginBottom(1)

	labelStyle = lipgloss.NewStyle().
			Foreground(grayLight).
			Width(20)

	valueStyle = lipgloss.NewStyle().
			Foreground(white).
			Bold(true)

	dimValueStyle = lipgloss.NewStyle().
			Foreground(grayLighter)

	helpStyle = lipgloss.NewStyle().
			Foreground(grayMedium).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)
)

func InitialModel(analysisFunc func() (*models.AnalysisResult, error)) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primaryColor)

	p := progress.New(
		progress.WithSolidFill(string(primaryColor)),
		progress.WithoutPercentage(),
	)

	return model{
		state:        loadingState,
		spinner:      s,
		progress:     p,
		analysisFunc: analysisFunc,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.doAnalysis(),
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) doAnalysis() tea.Cmd {
	return func() tea.Msg {
		result, err := m.analysisFunc()
		return analysisCompleteMsg{result: result, err: err}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = min(50, m.width-4)
		if m.state == resultState {
			m.updateTable()
		}
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}

		if m.state == resultState {
			switch msg.String() {
			case "tab", "right", "l":
				m.currentView = (m.currentView + 1) % view(4)
				if m.currentView == authorsView && len(m.result.Authors) == 0 {
					m.currentView = summaryView
				}
				m.updateTable()

			case "shift+tab", "left", "h":
				m.currentView = (m.currentView - 1 + 4) % view(4)
				if m.currentView == authorsView && len(m.result.Authors) == 0 {
					m.currentView = filesView
				}
				m.updateTable()

			case "1":
				m.currentView = summaryView
			case "2":
				m.currentView = languagesView
				m.updateTable()
			case "3":
				m.currentView = filesView
				m.updateTable()
			case "4":
				if len(m.result.Authors) > 0 {
					m.currentView = authorsView
					m.updateTable()
				}
			}
		}

	case analysisCompleteMsg:
		m.state = resultState
		m.result = msg.result
		m.err = msg.err
		if m.result != nil {
			m.updateTable()
		}
		return m, nil

	case tickMsg:
		if m.state == loadingState {
			m.progressValue += 0.02
			if m.progressValue > 1.0 {
				m.progressValue = 0.1
			}
			return m, tickCmd()
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	if m.state == resultState && m.table.Focused() {
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return ""
	}

	var content string

	if m.state == loadingState {
		content = m.renderLoading()
	} else if m.err != nil {
		content = m.renderError()
	} else {
		content = m.renderResults()
	}

	return m.centerContent(content)
}

func (m model) renderLoading() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("CodeLoc") + "\n")
	s.WriteString(subtitleStyle.Render("Analyzing codebase") + "\n")
	s.WriteString(fmt.Sprintf("%s  Scanning files and counting lines\n\n", m.spinner.View()))
	s.WriteString(m.progress.ViewAs(m.progressValue))

	return s.String()
}

func (m model) renderError() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("CodeLoc") + "\n")
	s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))

	return s.String()
}

func (m model) renderResults() string {
	var s strings.Builder

	// Title
	s.WriteString(titleStyle.Render("CodeLoc") + "\n")

	// Tabs
	tabs := m.renderTabs()
	s.WriteString(tabs + "\n")
	s.WriteString(dividerStyle.Render(strings.Repeat("─", min(100, m.width-4))) + "\n\n")

	// Content
	var contentView string
	switch m.currentView {
	case summaryView:
		contentView = m.renderSummary()
	case languagesView, filesView, authorsView:
		contentView = m.renderTable()
	}

	s.WriteString(contentView)

	// Help
	s.WriteString("\n" + m.renderHelp())

	return s.String()
}

func (m model) renderTabs() string {
	tabs := []string{"Summary", "Languages", "Files"}
	if len(m.result.Authors) > 0 {
		tabs = append(tabs, "Authors")
	}

	var renderedTabs []string
	for i, t := range tabs {
		style := tabStyle
		if view(i) == m.currentView {
			style = activeTabStyle
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	return lipgloss.JoinHorizontal(lipgloss.Bottom, renderedTabs...)
}

func (m model) renderSummary() string {
	var s strings.Builder

	// Statistics in clean layout
	stats := []struct {
		label string
		value string
	}{
		{"Path", m.result.Path},
		{"Total Files", fmt.Sprintf("%d", m.result.TotalFiles)},
		{"Total Lines", fmt.Sprintf("%d", m.result.TotalLines)},
		{"Code Lines", fmt.Sprintf("%d", m.result.CodeLines)},
		{"Comment Lines", fmt.Sprintf("%d", m.result.CommentLines)},
		{"Blank Lines", fmt.Sprintf("%d", m.result.BlankLines)},
		{"Languages", fmt.Sprintf("%d", len(m.result.Languages))},
	}

	for _, stat := range stats {
		line := fmt.Sprintf("%s %s\n",
			labelStyle.Render(stat.label),
			valueStyle.Render(stat.value))
		s.WriteString(line)
	}

	// Resource usage section
	s.WriteString("\n")
	s.WriteString(dividerStyle.Render(strings.Repeat("─", 60)) + "\n")
	s.WriteString(labelStyle.Render("Performance") + "\n\n")

	// Calculate metrics
	seconds := m.result.ProcessingTime.Seconds()
	if seconds == 0 {
		seconds = 0.001 // Avoid division by zero
	}
	filesPerSec := float64(m.result.TotalFiles) / seconds
	linesPerSec := float64(m.result.TotalLines) / seconds

	perfStats := []struct {
		label string
		value string
	}{
		{"Processing Time", m.result.ProcessingTime.String()},
		{"Files/sec", fmt.Sprintf("%.0f", filesPerSec)},
		{"Lines/sec", fmt.Sprintf("%.0f", linesPerSec)},
	}

	for _, stat := range perfStats {
		line := fmt.Sprintf("  %s %s\n",
			labelStyle.Render(stat.label),
			dimValueStyle.Render(stat.value))
		s.WriteString(line)
	}

	// Language breakdown
	if len(m.result.Languages) > 0 {
		s.WriteString("\n")
		s.WriteString(dividerStyle.Render(strings.Repeat("─", 60)) + "\n")
		s.WriteString(labelStyle.Render("Language Distribution") + "\n\n")

		for i, lang := range m.result.Languages {
			if i >= 5 {
				break
			}

			// Calculate bar width (max 40 chars)
			barWidth := int(lang.Percentage * 0.4)
			if barWidth < 1 {
				barWidth = 1
			}

			// Build bar with color based on index
			var bar string
			if i == 0 {
				bar = lipgloss.NewStyle().Foreground(primaryColor).Render(strings.Repeat("█", barWidth))
			} else {
				grayShade := lipgloss.Color(fmt.Sprintf("#%x", 0x9CA3AF-(i*0x0A0A0A)))
				bar = lipgloss.NewStyle().Foreground(grayShade).Render(strings.Repeat("█", barWidth))
			}

			// Format line
			langName := fmt.Sprintf("%-15s", lang.Language)
			percentage := dimValueStyle.Render(fmt.Sprintf("%5.1f%%", lang.Percentage))
			files := dimValueStyle.Render(fmt.Sprintf("%d files", lang.FileCount))

			s.WriteString(fmt.Sprintf("  %s %s %s %s\n", langName, bar, percentage, files))
		}
	}

	return s.String()
}

func (m model) renderTable() string {
	return m.table.View()
}

func (m model) renderHelp() string {
	// Dynamic help based on available views
	maxView := "3"
	if len(m.result.Authors) > 0 {
		maxView = "4"
	}

	helps := []string{
		"tab/shift+tab: switch",
		"↑/↓: scroll",
		fmt.Sprintf("1-%s: jump", maxView),
		"q: quit",
	}
	return helpStyle.Render(strings.Join(helps, "  •  "))
}

func (m *model) updateTable() {
	columns := []table.Column{}
	rows := []table.Row{}

	switch m.currentView {
	case languagesView:
		columns = []table.Column{
			{Title: "Language", Width: 20},
			{Title: "Files", Width: 10},
			{Title: "Lines", Width: 12},
			{Title: "Code", Width: 12},
			{Title: "Comments", Width: 12},
			{Title: "Blanks", Width: 12},
			{Title: "Percent", Width: 10},
		}

		for _, lang := range m.result.Languages {
			rows = append(rows, table.Row{
				lang.Language,
				fmt.Sprintf("%d", lang.FileCount),
				fmt.Sprintf("%d", lang.TotalLines),
				fmt.Sprintf("%d", lang.CodeLines),
				fmt.Sprintf("%d", lang.CommentLines),
				fmt.Sprintf("%d", lang.BlankLines),
				fmt.Sprintf("%.1f%%", lang.Percentage),
			})
		}

	case filesView:
		columns = []table.Column{
			{Title: "File", Width: 55},
			{Title: "Language", Width: 15},
			{Title: "Lines", Width: 10},
			{Title: "Code", Width: 10},
			{Title: "Comments", Width: 10},
		}

		for _, file := range m.result.Files {
			path := file.Path
			if len(path) > 53 {
				path = "..." + path[len(path)-50:]
			}

			rows = append(rows, table.Row{
				path,
				file.Language,
				fmt.Sprintf("%d", file.TotalLines),
				fmt.Sprintf("%d", file.CodeLines),
				fmt.Sprintf("%d", file.CommentLines),
			})
		}

	case authorsView:
		columns = []table.Column{
			{Title: "Author", Width: 35},
			{Title: "Lines", Width: 10},
			{Title: "Files", Width: 10},
			{Title: "Code", Width: 10},
			{Title: "Comments", Width: 10},
			{Title: "Percent", Width: 10},
		}

		for _, author := range m.result.Authors {
			rows = append(rows, table.Row{
				author.Name,
				fmt.Sprintf("%d", author.Lines),
				fmt.Sprintf("%d", author.FileCount),
				fmt.Sprintf("%d", author.CodeLines),
				fmt.Sprintf("%d", author.CommentLines),
				fmt.Sprintf("%.1f%%", author.Percentage),
			})
		}
	}

	tableHeight := min(20, m.height-15)
	if tableHeight < 5 {
		tableHeight = 5
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(grayDark).
		BorderBottom(true).
		Bold(false).
		Foreground(grayLight)
	s.Selected = s.Selected.
		Foreground(white).
		Background(primaryColor).
		Bold(false)
	s.Cell = s.Cell.
		Foreground(grayLighter)

	t.SetStyles(s)
	m.table = t
}

func (m model) centerContent(content string) string {
	lines := strings.Split(content, "\n")
	var maxWidth int
	for _, line := range lines {
		width := lipgloss.Width(line)
		if width > maxWidth {
			maxWidth = width
		}
	}

	verticalPadding := (m.height - len(lines)) / 2
	if verticalPadding < 0 {
		verticalPadding = 0
	}

	horizontalPadding := (m.width - maxWidth) / 2
	if horizontalPadding < 0 {
		horizontalPadding = 0
	}

	style := lipgloss.NewStyle().
		Padding(verticalPadding, 0, 0, horizontalPadding)

	return style.Render(content)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RenderInteractive starts the interactive TUI with analysis
func RenderInteractive(analysisFunc func() (*models.AnalysisResult, error)) error {
	p := tea.NewProgram(
		InitialModel(analysisFunc),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	_, err := p.Run()
	return err
}
