# CodeLoc

Lines of code counter with beautiful TUI - analyze code repositories with detailed statistics.

## Installation

```bash
# Run directly with npx (no installation needed)
npx codeloc

# Or install globally
npm install -g codeloc
```

## Usage

```bash
# Analyze current directory
codeloc

# Analyze a specific path
codeloc ./src

# Analyze a remote repository
codeloc https://github.com/user/repo
codeloc user/repo  # Short format for GitHub

# Enable author tracking (git blame)
codeloc --authors

# Use static output mode
codeloc --static

# JSON output for scripting
codeloc --static --output json
```

## Features

- 🚀 Fast concurrent analysis
- 🎨 Beautiful interactive TUI
- 📊 Detailed per-file, per-language statistics
- 👥 Git blame author contributions
- 🌐 Analyze remote repositories
- 📦 Zero native dependencies (pure Go binary)

## Options

```
--authors           Enable git blame author tracking
--static            Use static output instead of interactive TUI
--branch <name>     Specify branch for remote repos
--output <format>   Output format: table, json (requires --static)
```

## Interactive Controls

- **Tab**: Next view
- **Shift+Tab**: Previous view
- **↑/↓**: Scroll
- **Q**: Quit

## Platform Support

- Linux (x64, arm64)
- macOS (x64, arm64)
- Windows (x64)

## How It Works

The npm package automatically downloads the appropriate pre-built Go binary for your platform during installation. The binary is cached locally and executed when you run `codeloc`.

## License

MIT

## Repository

https://github.com/stripsior/loc.cli
