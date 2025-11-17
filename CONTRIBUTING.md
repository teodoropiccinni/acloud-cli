# Contributing to acloud-cli

Thank you for your interest in contributing to acloud-cli!

## Development Setup

### Prerequisites

- Go 1.24.2 or higher
- Git
- Access to Aruba Cloud API credentials for testing

### Getting Started

1. **Clone the repository**

```bash
git clone https://github.com/Arubacloud/acloud-cli.git
cd acloud-cli
```

2. **Install dependencies**

```bash
go mod download
```

3. **Build the project**

```bash
go build -o acloud
```

4. **Run locally**

```bash
./acloud --help
```

## Building

### Build for current platform

```bash
go build -o acloud
```

### Build for all platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o acloud-linux-amd64

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o acloud-darwin-amd64

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o acloud-darwin-arm64

# Windows
GOOS=windows GOARCH=amd64 go build -o acloud-windows-amd64.exe
```

## Testing

### Run tests

```bash
go test ./...
```

### Run tests with coverage

```bash
go test -cover ./...
```

### Run tests with verbose output

```bash
go test -v ./...
```

## Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Use meaningful variable and function names
- Add comments for exported functions and types

### Format code

```bash
go fmt ./...
```

### Lint code

```bash
go vet ./...
```

## Project Structure

```
acloud-cli/
├── cmd/           # Command implementations
│   ├── root.go    # Root command
│   └── config.go  # Config command
├── main.go        # Entry point
├── go.mod         # Go module definition
├── go.sum         # Go dependencies checksums
├── README.md      # User documentation
└── CONTRIBUTING.md # This file
```

## Adding New Commands

1. Create a new file in the `cmd/` directory (e.g., `cmd/newcommand.go`)
2. Define your command using Cobra:

```go
package cmd

import (
    "github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
    Use:   "new",
    Short: "Short description",
    Long:  `Long description`,
    Run: func(cmd *cobra.Command, args []string) {
        // Implementation
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
}
```

3. Build and test your changes
4. Update documentation as needed

## Release Process

Releases are automated via GitHub Actions pipeline:

1. Push changes to `main` branch
2. Create and push a new tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
3. GitHub Actions will automatically:
   - Build binaries for Linux, macOS, and Windows
   - Create a GitHub release
   - Upload binaries as release assets

## Getting Help

- Open an issue for bug reports or feature requests
- Check existing issues before creating a new one
- Provide detailed information about your environment and the problem

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow

Thank you for contributing!
