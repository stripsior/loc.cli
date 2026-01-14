# CodeLoc

A blazingly fast lines of code counter with a beautiful TUI, written in Go and distributed via npm.

## Features

- 🚀 **Fast Analysis**: Concurrent processing with worker pools
- 🎨 **Beautiful TUI**: Interactive terminal UI built with Charm's Bubbletea
- 📊 **Detailed Statistics**: Per-file, per-language, and per-author breakdowns
- 🌐 **Remote Repositories**: Analyze GitHub repos directly
- 🔍 **20+ Languages**: Accurate language detection and comment parsing
- 📦 **Easy Installation**: Install via npm, run via `npx`
- 🎯 **Multiple Output Modes**: Interactive TUI or static table/JSON output

## Quick Start

```bash
# Run directly with npx
npx codeloc

# Or install globally
npm install -g codeloc

# Analyze current directory
codeloc

# Analyze a specific path
codeloc ./src

# Analyze a remote repository
codeloc https://github.com/stripsior/loc.api
codeloc stripsior/loc.api  # Short format

# Enable author tracking (git blame)
codeloc --authors

# Use static output mode
codeloc --static

# JSON output for scripting
codeloc --static --output json
```

## Options

```
Flags:
  --authors           Enable git blame author tracking
  --static            Use static output instead of interactive TUI
  --branch <name>     Specify branch for remote repos (default: main/master)
  --output <format>   Output format: table, json (requires --static)
  -h, --help          Help for codeloc
  -v, --version       Version information
```

## Interactive TUI Controls

- **Tab** / **→** / **L**: Next view
- **Shift+Tab** / **←** / **H**: Previous view
- **↑ / ↓**: Scroll through tables
- **1-4**: Jump to specific view
- **Q** / **Ctrl+C**: Quit

## Views

1. **Summary**: Overview of total statistics
2. **Languages**: Breakdown by programming language
3. **Files**: Per-file statistics sorted by size
4. **Authors**: Contributor statistics (with `--authors` flag)

## Supported Languages

Go, JavaScript, TypeScript, Python, Java, C, C++, C#, Rust, Swift, Kotlin, Ruby, PHP, Shell, Bash, HTML, CSS, SQL, Lua, R, YAML, TOML, Dockerfile, and more.

## Development

This is a monorepo containing:

- **cli/**: Go CLI application
- **npm/**: NPM package wrapper

### Building from Source

```bash
# Build the Go CLI
cd cli
make build

# Run locally
make run

# Build for all platforms
make build-all
```

### Project Structure

```
loc.cli/
├── cli/                    # Go application
│   ├── cmd/codeloc/        # Main entry point
│   ├── internal/
│   │   ├── analyzer/       # Core analysis logic
│   │   ├── repository/     # Local/remote repo handling
│   │   ├── models/         # Data structures
│   │   └── ui/             # TUI and static output
│   └── Makefile
└── npm/                    # NPM wrapper
    ├── bin/codeloc.js      # Wrapper script
    └── scripts/            # Postinstall binary downloader
```

## Publishing (Maintainers)

This project uses GitHub Actions for automated publishing. The workflow handles version bumping, building, and releasing.

### Prerequisites

1. **NPM Token**: Create an automation token at [npmjs.com](https://www.npmjs.com/settings)
2. **GitHub Secret**: Add `NPM_TOKEN` to repository secrets
3. **Permissions**: Enable "Read and write permissions" in Actions settings

### Publish New Version

1. Go to **Actions** tab → **Publish CodeLoc**
2. Click **Run workflow**
3. Select bump type:
   - `patch` - Bug fixes (1.0.0 → 1.0.1)
   - `minor` - New features (1.0.0 → 1.1.0)
   - `major` - Breaking changes (1.0.0 → 2.0.0)
4. Click **Run workflow**

The workflow will:
- ✅ Bump version in `package.json` and `main.go`
- ✅ Commit and tag the release
- ✅ Build binaries for all platforms
- ✅ Create GitHub release
- ✅ Publish to NPM

### Verify Release

```bash
# Check NPM
npm view codeloc version

# Test installation
npx codeloc@latest .
```

## License

MIT

## Author

stripsior

## Repository

https://github.com/stripsior/loc.cli
