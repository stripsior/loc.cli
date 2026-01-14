package ui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/stripsior/loc.cli/internal/models"
)

var (
	staticTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12")).
			Padding(0, 1)

	staticHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("11"))

	staticStatStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	staticTableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("13"))
)

// RenderStatic renders the analysis results in static format
func RenderStatic(result *models.AnalysisResult, format string) error {
	if format == "json" {
		return renderJSON(result)
	}

	return renderTable(result)
}

func renderJSON(result *models.AnalysisResult) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func renderTable(result *models.AnalysisResult) error {
	// Print title
	fmt.Println(staticTitleStyle.Render("CodeLoc Analysis Results"))
	fmt.Println()

	// Print summary
	fmt.Println(staticHeaderStyle.Render("Summary"))
	fmt.Printf("  Path:           %s\n", result.Path)
	fmt.Printf("  Total Files:    %s\n", staticStatStyle.Render(fmt.Sprintf("%d", result.TotalFiles)))
	fmt.Printf("  Total Lines:    %s\n", staticStatStyle.Render(fmt.Sprintf("%d", result.TotalLines)))
	fmt.Printf("  Code Lines:     %s\n", staticStatStyle.Render(fmt.Sprintf("%d", result.CodeLines)))
	fmt.Printf("  Comment Lines:  %s\n", staticStatStyle.Render(fmt.Sprintf("%d", result.CommentLines)))
	fmt.Printf("  Blank Lines:    %s\n", staticStatStyle.Render(fmt.Sprintf("%d", result.BlankLines)))
	fmt.Printf("  Processing Time: %s\n", result.ProcessingTime)
	fmt.Println()

	// Print language breakdown
	if len(result.Languages) > 0 {
		fmt.Println(staticHeaderStyle.Render("Languages"))
		fmt.Println()

		// Table header
		fmt.Printf("  %-20s %10s %10s %10s %10s %10s %10s\n",
			staticTableHeaderStyle.Render("Language"),
			staticTableHeaderStyle.Render("Files"),
			staticTableHeaderStyle.Render("Lines"),
			staticTableHeaderStyle.Render("Code"),
			staticTableHeaderStyle.Render("Comments"),
			staticTableHeaderStyle.Render("Blanks"),
			staticTableHeaderStyle.Render("Percent"),
		)
		fmt.Println("  " + strings.Repeat("-", 90))

		// Table rows
		for _, lang := range result.Languages {
			fmt.Printf("  %-20s %10d %10d %10d %10d %10d %9.1f%%\n",
				lang.Language,
				lang.FileCount,
				lang.TotalLines,
				lang.CodeLines,
				lang.CommentLines,
				lang.BlankLines,
				lang.Percentage,
			)
		}
		fmt.Println()
	}

	// Print top files
	if len(result.Files) > 0 {
		fmt.Println(staticHeaderStyle.Render("Top Files (by lines)"))
		fmt.Println()

		// Show top 20 files
		limit := 20
		if len(result.Files) < limit {
			limit = len(result.Files)
		}

		// Table header
		fmt.Printf("  %-50s %12s %10s %10s %10s %10s\n",
			staticTableHeaderStyle.Render("File"),
			staticTableHeaderStyle.Render("Language"),
			staticTableHeaderStyle.Render("Lines"),
			staticTableHeaderStyle.Render("Code"),
			staticTableHeaderStyle.Render("Comments"),
			staticTableHeaderStyle.Render("Blanks"),
		)
		fmt.Println("  " + strings.Repeat("-", 110))

		// Table rows
		for i := 0; i < limit; i++ {
			file := result.Files[i]
			// Truncate path if too long
			path := file.Path
			if len(path) > 48 {
				path = "..." + path[len(path)-45:]
			}

			fmt.Printf("  %-50s %12s %10d %10d %10d %10d\n",
				path,
				file.Language,
				file.TotalLines,
				file.CodeLines,
				file.CommentLines,
				file.BlankLines,
			)
		}
		fmt.Println()
	}

	// Print author statistics if available
	if len(result.Authors) > 0 {
		fmt.Println(staticHeaderStyle.Render("Author Contributions"))
		fmt.Println()

		// Table header
		fmt.Printf("  %-30s %10s %10s %10s %10s %10s\n",
			staticTableHeaderStyle.Render("Author"),
			staticTableHeaderStyle.Render("Lines"),
			staticTableHeaderStyle.Render("Code"),
			staticTableHeaderStyle.Render("Comments"),
			staticTableHeaderStyle.Render("Blanks"),
			staticTableHeaderStyle.Render("Percent"),
		)
		fmt.Println("  " + strings.Repeat("-", 80))

		// Table rows
		for _, author := range result.Authors {
			fmt.Printf("  %-30s %10d %10d %10d %10d %9.1f%%\n",
				author.Name,
				author.Lines,
				author.CodeLines,
				author.CommentLines,
				author.BlankLines,
				author.Percentage,
			)
		}
		fmt.Println()
	}

	return nil
}
