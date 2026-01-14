package analyzer

import (
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/stripsior/loc.cli/internal/models"
)

const numWorkers = 10

// Analyzer handles the analysis of a directory
type Analyzer struct {
	rootPath string
	files    chan string
	results  chan *models.FileStats
	wg       sync.WaitGroup
}

// NewAnalyzer creates a new analyzer
func NewAnalyzer(rootPath string) *Analyzer {
	return &Analyzer{
		rootPath: rootPath,
		files:    make(chan string, 100),
		results:  make(chan *models.FileStats, 100),
	}
}

// Analyze performs the analysis on the directory
func (a *Analyzer) Analyze() (*models.AnalysisResult, error) {
	startTime := time.Now()

	// Start worker pool
	for i := 0; i < numWorkers; i++ {
		a.wg.Add(1)
		go a.worker()
	}

	// Collect results in a separate goroutine
	var allFiles []*models.FileStats
	var collectorWg sync.WaitGroup
	collectorWg.Add(1)
	go func() {
		defer collectorWg.Done()
		for stats := range a.results {
			if stats != nil {
				allFiles = append(allFiles, stats)
			}
		}
	}()

	// Walk the directory tree
	err := filepath.Walk(a.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files with errors
		}

		if info.IsDir() {
			// Check if directory should be ignored
			if ShouldIgnoreFile(path) && path != a.rootPath {
				return filepath.SkipDir
			}
			return nil
		}

		// Send file for processing
		a.files <- path
		return nil
	})

	// Close files channel and wait for workers
	close(a.files)
	a.wg.Wait()

	// Close results channel and wait for collector
	close(a.results)
	collectorWg.Wait()

	if err != nil {
		return nil, err
	}

	// Aggregate results
	result := a.aggregateResults(allFiles)
	result.Path = a.rootPath
	result.ProcessingTime = time.Since(startTime)

	return result, nil
}

// worker processes files from the channel
func (a *Analyzer) worker() {
	defer a.wg.Done()

	for filePath := range a.files {
		stats, err := CountLines(filePath)
		if err != nil {
			// Skip files with errors
			continue
		}
		a.results <- stats
	}
}

// aggregateResults aggregates file statistics into the final result
func (a *Analyzer) aggregateResults(files []*models.FileStats) *models.AnalysisResult {
	// Convert to non-pointer slice
	fileSlice := make([]models.FileStats, len(files))
	for i, f := range files {
		if f != nil {
			fileSlice[i] = *f
		}
	}

	result := &models.AnalysisResult{
		Files: fileSlice,
	}

	languageMap := make(map[string]*models.LanguageStats)

	for _, file := range fileSlice {
		result.TotalFiles++
		result.TotalLines += file.TotalLines
		result.CodeLines += file.CodeLines
		result.CommentLines += file.CommentLines
		result.BlankLines += file.BlankLines

		// Aggregate by language
		if langStats, ok := languageMap[file.Language]; ok {
			langStats.FileCount++
			langStats.TotalLines += file.TotalLines
			langStats.CodeLines += file.CodeLines
			langStats.CommentLines += file.CommentLines
			langStats.BlankLines += file.BlankLines
			langStats.Bytes += file.Size
		} else {
			languageMap[file.Language] = &models.LanguageStats{
				Language:     file.Language,
				FileCount:    1,
				TotalLines:   file.TotalLines,
				CodeLines:    file.CodeLines,
				CommentLines: file.CommentLines,
				BlankLines:   file.BlankLines,
				Bytes:        file.Size,
			}
		}
	}

	// Convert map to slice and calculate percentages
	totalBytes := int64(0)
	for _, langStats := range languageMap {
		totalBytes += langStats.Bytes
		result.Languages = append(result.Languages, *langStats)
	}

	// Calculate percentages
	for i := range result.Languages {
		if totalBytes > 0 {
			result.Languages[i].Percentage = float64(result.Languages[i].Bytes) / float64(totalBytes) * 100
		}
	}

	// Sort languages by bytes (descending)
	sort.Slice(result.Languages, func(i, j int) bool {
		return result.Languages[i].Bytes > result.Languages[j].Bytes
	})

	// Sort files by total lines (descending)
	sort.Slice(result.Files, func(i, j int) bool {
		return result.Files[i].TotalLines > result.Files[j].TotalLines
	})

	return result
}
