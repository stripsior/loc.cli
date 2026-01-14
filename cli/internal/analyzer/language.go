package analyzer

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-enry/go-enry/v2"
)

// CommentPattern defines comment syntax for a programming language
type CommentPattern struct {
	LineComment       string   // Single-line comment prefix (e.g., "//" for Go)
	BlockCommentStart string   // Block comment start (e.g., "/*" for Go)
	BlockCommentEnd   string   // Block comment end (e.g., "*/" for Go)
	AltLineComments   []string // Alternative line comment prefixes
}

// languageCommentPatterns maps programming languages to their comment patterns
var languageCommentPatterns = map[string]CommentPattern{
	"Go": {
		LineComment:       "//",
		BlockCommentStart: "/*",
		BlockCommentEnd:   "*/",
	},
	"JavaScript": {
		LineComment:       "//",
		BlockCommentStart: "/*",
		BlockCommentEnd:   "*/",
	},
	"TypeScript": {
		LineComment:       "//",
		BlockCommentStart: "/*",
		BlockCommentEnd:   "*/",
	},
	"Java": {
		LineComment:       "//",
		BlockCommentStart: "/*",
		BlockCommentEnd:   "*/",
	},
	"C": {
		LineComment:       "//",
		BlockCommentStart: "/*",
		BlockCommentEnd:   "*/",
	},
	"C++": {
		LineComment:       "//",
		BlockCommentStart: "/*",
		BlockCommentEnd:   "*/",
	},
	"C#": {
		LineComment:       "//",
		BlockCommentStart: "/*",
		BlockCommentEnd:   "*/",
	},
	"Rust": {
		LineComment:       "//",
		BlockCommentStart: "/*",
		BlockCommentEnd:   "*/",
	},
	"Swift": {
		LineComment:       "//",
		BlockCommentStart: "/*",
		BlockCommentEnd:   "*/",
	},
	"Kotlin": {
		LineComment:       "//",
		BlockCommentStart: "/*",
		BlockCommentEnd:   "*/",
	},
	"Python": {
		LineComment:       "#",
		BlockCommentStart: "\"\"\"",
		BlockCommentEnd:   "\"\"\"",
		AltLineComments:   []string{"'''"},
	},
	"Ruby": {
		LineComment:       "#",
		BlockCommentStart: "=begin",
		BlockCommentEnd:   "=end",
	},
	"PHP": {
		LineComment:       "//",
		BlockCommentStart: "/*",
		BlockCommentEnd:   "*/",
		AltLineComments:   []string{"#"},
	},
	"Shell": {
		LineComment: "#",
	},
	"Bash": {
		LineComment: "#",
	},
	"HTML": {
		BlockCommentStart: "<!--",
		BlockCommentEnd:   "-->",
	},
	"CSS": {
		BlockCommentStart: "/*",
		BlockCommentEnd:   "*/",
	},
	"SQL": {
		LineComment:       "--",
		BlockCommentStart: "/*",
		BlockCommentEnd:   "*/",
	},
	"Lua": {
		LineComment:       "--",
		BlockCommentStart: "--[[",
		BlockCommentEnd:   "]]",
	},
	"R": {
		LineComment: "#",
	},
	"YAML": {
		LineComment: "#",
	},
	"TOML": {
		LineComment: "#",
	},
	"Dockerfile": {
		LineComment: "#",
	},
	"Makefile": {
		LineComment: "#",
	},
	"Perl": {
		LineComment:       "#",
		BlockCommentStart: "=pod",
		BlockCommentEnd:   "=cut",
	},
	"Haskell": {
		LineComment:       "--",
		BlockCommentStart: "{-",
		BlockCommentEnd:   "-}",
	},
	"Scala": {
		LineComment:       "//",
		BlockCommentStart: "/*",
		BlockCommentEnd:   "*/",
	},
	"Dart": {
		LineComment:       "//",
		BlockCommentStart: "/*",
		BlockCommentEnd:   "*/",
	},
}

// DetectLanguage detects the programming language of a file using go-enry
func DetectLanguage(filePath string, content []byte) string {
	// First try detection by filename
	language := enry.GetLanguage(filepath.Base(filePath), content)

	if language == "" {
		// Fallback: try detection by content only
		language, _ = enry.GetLanguageByContent(filePath, content)
	}

	return language
}

// GetCommentPattern returns the comment pattern for a given language
func GetCommentPattern(language string) (CommentPattern, bool) {
	pattern, ok := languageCommentPatterns[language]
	return pattern, ok
}

// IsBinaryFile checks if a file is binary (should skip analysis)
func IsBinaryFile(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Read first 512 bytes to check for binary content
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		return false
	}

	return enry.IsBinary(buf[:n])
}

// ShouldIgnoreFile checks if a file should be ignored
func ShouldIgnoreFile(filePath string) bool {
	name := filepath.Base(filePath)

	// Ignore common binary and generated files
	ignoredExtensions := []string{
		".exe", ".dll", ".so", ".dylib", ".a", ".o", ".obj",
		".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg", ".webp",
		".pdf", ".zip", ".tar", ".gz", ".7z", ".rar",
		".mp3", ".mp4", ".avi", ".mov", ".wav",
		".ttf", ".woff", ".woff2", ".eot",
		".lock", ".sum",
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	for _, ignoredExt := range ignoredExtensions {
		if ext == ignoredExt {
			return true
		}
	}

	// Ignore common directories
	ignoredDirs := []string{
		"node_modules", "vendor", ".git", ".svn", ".hg",
		"bin", "obj", "dist", "build", "target",
		".idea", ".vscode", "__pycache__", ".cache",
	}

	parts := strings.Split(filePath, string(filepath.Separator))
	for _, part := range parts {
		for _, ignored := range ignoredDirs {
			if part == ignored {
				return true
			}
		}
	}

	// Ignore hidden files (starting with .)
	if strings.HasPrefix(name, ".") && name != ".." && name != "." {
		return true
	}

	return false
}
