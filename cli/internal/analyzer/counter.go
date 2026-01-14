package analyzer

import (
	"bufio"
	"os"
	"strings"

	"github.com/stripsior/loc.cli/internal/models"
)

// CountLines analyzes a file and returns statistics
func CountLines(filePath string) (*models.FileStats, error) {
	// Check if file should be ignored
	if ShouldIgnoreFile(filePath) {
		return nil, nil
	}

	// Check if binary
	if IsBinaryFile(filePath) {
		return nil, nil
	}

	// Get file info for size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Detect language
	language := DetectLanguage(filePath, content)
	if language == "" {
		// Unknown language, skip or mark as unknown
		return nil, nil
	}

	// Get comment pattern for the language
	commentPattern, hasPattern := GetCommentPattern(language)

	stats := &models.FileStats{
		Path:     filePath,
		Language: language,
		Size:     fileInfo.Size(),
	}

	// If no comment pattern, treat all non-blank lines as code
	if !hasPattern {
		stats.TotalLines, stats.CodeLines, stats.BlankLines = countBasicLines(content)
		return stats, nil
	}

	// Count lines with comment detection
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	inBlockComment := false
	blockCommentStart := commentPattern.BlockCommentStart
	blockCommentEnd := commentPattern.BlockCommentEnd

	for scanner.Scan() {
		line := scanner.Text()
		stats.TotalLines++

		trimmedLine := strings.TrimSpace(line)

		// Check for blank lines
		if trimmedLine == "" {
			stats.BlankLines++
			continue
		}

		isComment := false

		// Handle block comments
		if blockCommentStart != "" && blockCommentEnd != "" {
			// Check if we're entering a block comment
			if !inBlockComment && strings.Contains(trimmedLine, blockCommentStart) {
				inBlockComment = true
				isComment = true

				// Check if block comment ends on the same line
				startIdx := strings.Index(trimmedLine, blockCommentStart)
				endIdx := strings.Index(trimmedLine[startIdx+len(blockCommentStart):], blockCommentEnd)
				if endIdx != -1 {
					inBlockComment = false
					// Check if there's code after the block comment
					remaining := strings.TrimSpace(trimmedLine[startIdx+len(blockCommentStart)+endIdx+len(blockCommentEnd):])
					if remaining != "" {
						isComment = false // Mixed line
					}
				}
			} else if inBlockComment {
				isComment = true
				// Check if block comment ends on this line
				if strings.Contains(trimmedLine, blockCommentEnd) {
					inBlockComment = false
				}
			}
		}

		// Check for line comments (if not already in block comment)
		if !inBlockComment && !isComment {
			if commentPattern.LineComment != "" && strings.HasPrefix(trimmedLine, commentPattern.LineComment) {
				isComment = true
			} else {
				// Check alternative line comments
				for _, altComment := range commentPattern.AltLineComments {
					if strings.HasPrefix(trimmedLine, altComment) {
						isComment = true
						break
					}
				}
			}
		}

		// Classify the line
		if isComment {
			stats.CommentLines++
		} else {
			stats.CodeLines++
		}
	}

	return stats, nil
}

// countBasicLines counts lines without comment detection
func countBasicLines(content []byte) (total, code, blank int) {
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		total++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			blank++
		} else {
			code++
		}
	}
	return
}
