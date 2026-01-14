# CodeLoc CLI

Go implementation of the CodeLoc lines of code counter.

## Development

### Prerequisites

- Go 1.20 or higher
- Make (optional, but recommended)

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Install locally
make install
```

### Running

```bash
# Run directly
go run ./cmd/codeloc .

# Or use make
make run
```

### Testing

```bash
make test
```

### Dependencies

- **go-enry/go-enry**: Language detection (GitHub's Linguist)
- **go-git/go-git**: Git operations
- **charmbracelet/bubbletea**: TUI framework
- **charmbracelet/lipgloss**: Styling
- **charmbracelet/bubbles**: TUI components
- **spf13/cobra**: CLI framework

## Architecture

```
cmd/codeloc/main.go         # Entry point and CLI setup
internal/
  ├── analyzer/
  │   ├── analyzer.go       # Core analysis engine
  │   ├── counter.go        # Line counting logic
  │   ├── language.go       # Language detection
  │   └── git.go            # Git blame/author tracking
  ├── repository/
  │   ├── local.go          # Local directory handling
  │   └── remote.go         # Remote repo cloning
  ├── models/
  │   └── models.go         # Data structures
  └── ui/
      ├── interactive.go    # Bubbletea TUI
      └── static.go         # Static output
```

## Release Process

1. Update version in `cmd/codeloc/main.go`
2. Build all platform binaries: `make build-all`
3. Create GitHub release and upload binaries
4. Update npm package version
5. Publish npm package
