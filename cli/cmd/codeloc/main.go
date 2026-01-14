package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/stripsior/loc.cli/internal/analyzer"
	"github.com/stripsior/loc.cli/internal/models"
	"github.com/stripsior/loc.cli/internal/repository"
	"github.com/stripsior/loc.cli/internal/ui"
)

var (
	enableAuthors bool
	staticMode    bool
	branch        string
	outputFormat  string
	version       = "1.0.1"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:     "codeloc [path/url]",
		Short:   "Count lines of code with beautiful visualizations",
		Long:    "CodeLoc analyzes code repositories and provides detailed statistics about lines of code, languages, and contributors.",
		Version: version,
		Args:    cobra.MaximumNArgs(1),
		RunE:    run,
	}

	rootCmd.Flags().BoolVar(&enableAuthors, "authors", false, "Enable git blame author tracking")
	rootCmd.Flags().BoolVar(&staticMode, "static", false, "Use static output instead of interactive TUI")
	rootCmd.Flags().StringVar(&branch, "branch", "", "Specify branch for remote repos (default: main/master)")
	rootCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Output format: table, json (requires --static)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Determine the target path
	targetPath := "."
	if len(args) > 0 {
		targetPath = args[0]
	}

	// Normalize URL if it's in "owner/repo" format
	targetPath = repository.NormalizeRepoURL(targetPath)

	var analysisPath string
	var cleanup func()

	// Check if it's a remote repository
	if repository.IsRemoteURL(targetPath) {
		if !staticMode {
			fmt.Printf("Cloning repository: %s\n", targetPath)
		}
		clonedPath, cleanupFunc, err := repository.CloneRepository(targetPath, branch)
		if err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
		analysisPath = clonedPath
		cleanup = cleanupFunc
		defer cleanup()
		if !staticMode {
			fmt.Println("Clone complete\n")
		}
	} else {
		// Local path
		if err := repository.ValidateLocalPath(targetPath); err != nil {
			return err
		}
		analysisPath = targetPath
	}

	// Create analysis function
	analysisFunc := func() (*models.AnalysisResult, error) {
		// Create analyzer and run analysis
		codeAnalyzer := analyzer.NewAnalyzer(analysisPath)

		result, err := codeAnalyzer.Analyze()
		if err != nil {
			return nil, fmt.Errorf("analysis failed: %w", err)
		}

		// Add author analysis if requested
		if enableAuthors {
			authors, err := analyzer.AnalyzeAuthors(analysisPath, result.Files)
			if err != nil {
				// Just log warning, don't fail
				fmt.Printf("Warning: Could not analyze authors: %v\n", err)
			} else {
				result.Authors = authors
			}
		}

		return result, nil
	}

	// Display results
	if staticMode {
		// For static mode, run analysis immediately
		fmt.Println("Analyzing repository...")
		result, err := analysisFunc()
		if err != nil {
			return err
		}
		return ui.RenderStatic(result, outputFormat)
	}

	// For interactive mode, pass the function to TUI
	return ui.RenderInteractive(analysisFunc)
}
