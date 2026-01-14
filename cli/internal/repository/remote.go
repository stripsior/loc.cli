package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// CloneRepository clones a remote repository to a temporary directory
func CloneRepository(url string, branch string) (string, func(), error) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "codeloc-*")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	// Clone options
	cloneOpts := &git.CloneOptions{
		URL:      url,
		Progress: nil, // Silent clone
		Depth:    1,   // Shallow clone for speed
	}

	// If branch is specified, use it
	if branch != "" {
		cloneOpts.ReferenceName = plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))
		cloneOpts.SingleBranch = true
	}

	// Clone the repository
	_, err = git.PlainClone(tempDir, false, cloneOpts)
	if err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to clone repository: %w", err)
	}

	return tempDir, cleanup, nil
}

// IsRemoteURL checks if a string looks like a remote repository URL
func IsRemoteURL(path string) bool {
	// Check for common URL patterns
	return strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "git@") ||
		strings.HasPrefix(path, "ssh://")
}

// NormalizeRepoURL normalizes GitHub-style "owner/repo" to full URL
func NormalizeRepoURL(input string) string {
	// If already a URL, return as-is
	if IsRemoteURL(input) {
		return input
	}

	// Check if it looks like "owner/repo" format
	parts := strings.Split(input, "/")
	if len(parts) == 2 && !strings.Contains(input, string(filepath.Separator)) {
		return fmt.Sprintf("https://github.com/%s", input)
	}

	return input
}
