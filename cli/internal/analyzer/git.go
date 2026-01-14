package analyzer

import (
	"path/filepath"
	"sort"

	"github.com/go-git/go-git/v5"
	"github.com/stripsior/loc.cli/internal/models"
)

// AnalyzeAuthors performs git blame analysis to calculate author contributions
func AnalyzeAuthors(rootPath string, files []models.FileStats) ([]models.AuthorStats, error) {
	// Open the git repository
	repo, err := git.PlainOpen(rootPath)
	if err != nil {
		// Not a git repository or error opening
		return nil, err
	}

	// Get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	// Get the commit object
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	authorMap := make(map[string]*models.AuthorStats)
	filesProcessed := make(map[string]bool)

	// Iterate through files and run blame
	for _, fileStats := range files {
		if filesProcessed[fileStats.Path] {
			continue
		}
		filesProcessed[fileStats.Path] = true

		// Get relative path
		relPath, err := filepath.Rel(rootPath, fileStats.Path)
		if err != nil {
			continue
		}

		// Get blame for the file
		blame, err := git.Blame(commit, relPath)
		if err != nil {
			// File might not be in git, skip
			continue
		}

		// Track which authors contributed to this file
		fileAuthors := make(map[string]bool)

		// Process blame lines
		for _, line := range blame.Lines {
			authorName := line.Author
			if authorName == "" {
				authorName = "Unknown"
			}

			fileAuthors[authorName] = true

			if author, ok := authorMap[authorName]; ok {
				author.Lines++
				// Classify line based on original file stats
				// This is a simplified approach
				if fileStats.TotalLines > 0 {
					codeRatio := float64(fileStats.CodeLines) / float64(fileStats.TotalLines)
					commentRatio := float64(fileStats.CommentLines) / float64(fileStats.TotalLines)

					author.CodeLines += int(codeRatio)
					author.CommentLines += int(commentRatio)
					author.BlankLines += int(1 - codeRatio - commentRatio)
				}
			} else {
				authorMap[authorName] = &models.AuthorStats{
					Name:         authorName,
					Lines:        1,
					CodeLines:    0,
					CommentLines: 0,
					BlankLines:   0,
					FileCount:    0,
				}
			}
		}

		// Update file counts for authors
		for authorName := range fileAuthors {
			if author, ok := authorMap[authorName]; ok {
				author.FileCount++
			}
		}
	}

	// Convert map to slice
	var authors []models.AuthorStats
	totalLines := 0
	for _, author := range authorMap {
		totalLines += author.Lines
		authors = append(authors, *author)
	}

	// Calculate percentages
	for i := range authors {
		if totalLines > 0 {
			authors[i].Percentage = float64(authors[i].Lines) / float64(totalLines) * 100
		}
	}

	// Sort by lines (descending)
	sort.Slice(authors, func(i, j int) bool {
		return authors[i].Lines > authors[j].Lines
	})

	return authors, nil
}
