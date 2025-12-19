# acloud-cli

[![GitHub release](https://img.shields.io/github/tag/Arubacloud/acloud-cli.svg?label=release)](https://github.com/Arubacloud/acloud-cli/releases/latest)

ArubaCloud Command Line Interface - A CLI tool for interacting with Aruba Cloud APIs.

> **⚠️ Development Status**: This CLI is currently under active development and is **not production-ready yet**. 

## Installation

### Windows

```powershell
 # Download the latest release binary
 Invoke-WebRequest -Uri "https://github.com/Arubacloud/acloud-cli/releases/latest/download/acloud-windows-amd64.exe" -OutFile "acloud.exe"

 # (Optional) Move to a folder in your PATH, e.g. C:\acloud
 # Move-Item -Path .\acloud.exe -Destination C:\acloud\acloud.exe
 # Optionally, add C:\acloud to your PATH environment variable

 # Run from the command prompt:
 acloud.exe --help
```

### Linux

```bash
# Download the latest release
curl -LO https://github.com/Arubacloud/acloud-cli/releases/latest/download/acloud-linux-amd64
chmod +x acloud-linux-amd64
sudo mv acloud-linux-amd64 /usr/local/bin/acloud
```

### macOS

```bash
# Download the latest release
curl -LO https://github.com/Arubacloud/acloud-cli/releases/latest/download/acloud-darwin-amd64
chmod +x acloud-darwin-amd64
sudo mv acloud-darwin-amd64 /usr/local/bin/acloud
```

## Configuration

Before using acloud, you need to configure your Aruba Cloud API credentials.

### Set Credentials

```bash
# Set both client ID and client secret
acloud config set --client-id YOUR_CLIENT_ID --client-secret YOUR_CLIENT_SECRET

# Set individual values
acloud config set --client-id YOUR_CLIENT_ID
acloud config set --client-secret YOUR_CLIENT_SECRET
```

### View Configuration

```bash
acloud config show
```

Configuration is stored in `~/.acloud.yaml` with secure file permissions.

## Quick Start

### 1. Configure Credentials

```bash
acloud config set --client-id YOUR_CLIENT_ID --client-secret YOUR_CLIENT_SECRET
```

### 2. Set up a Context (Optional but Recommended)

Contexts allow you to work with a specific project without repeatedly passing `--project-id`:

```bash
# Create a context with your project ID
acloud context set my-prod --project-id "66a10244f62b99c686572a9f"

# Switch to that context
acloud context use my-prod

# Now commands use the context project ID automatically
acloud storage blockstorage list
```

### 3. Explore Resources

```bash
# List projects
acloud management project list

# List block storage volumes
acloud storage blockstorage list

# List snapshots
acloud storage snapshot list
```

## Context Management

Manage multiple project contexts to simplify your workflow:

```bash
# Set contexts for different environments
acloud context set prod --project-id "prod-project-id"
acloud context set dev --project-id "dev-project-id"
acloud context set staging --project-id "staging-project-id"

# Switch between contexts
acloud context use prod
acloud context use dev

# View current context
acloud context current

# List all contexts
acloud context list

# Delete a context
acloud context delete staging
```

## Usage

```bash
# View all available commands
acloud --help

# View config command options
acloud config --help
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## License

See [LICENSE](LICENSE) file for details.
