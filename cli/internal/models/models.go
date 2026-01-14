package models

import "time"

// FileStats represents statistics for a single file
type FileStats struct {
	Path         string
	Language     string
	TotalLines   int
	CodeLines    int
	CommentLines int
	BlankLines   int
	Size         int64
}

// LanguageStats represents aggregated statistics for a programming language
type LanguageStats struct {
	Language     string
	FileCount    int
	TotalLines   int
	CodeLines    int
	CommentLines int
	BlankLines   int
	Bytes        int64
	Percentage   float64
}

// AuthorStats represents statistics for an author (git blame)
type AuthorStats struct {
	Name         string
	Lines        int
	CodeLines    int
	CommentLines int
	BlankLines   int
	FileCount    int
	Percentage   float64
}

// AnalysisResult represents the complete analysis result
type AnalysisResult struct {
	Path            string
	TotalFiles      int
	TotalLines      int
	CodeLines       int
	CommentLines    int
	BlankLines      int
	Languages       []LanguageStats
	Files           []FileStats
	Authors         []AuthorStats // Only populated with --authors flag
	ProcessingTime  time.Duration
}
