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

### Debug Mode

Enable verbose logging to troubleshoot issues:

```bash
# Enable debug logging (shows HTTP requests/responses)
acloud --debug network securityrule update <vpc-id> <securitygroup-id> <securityrule-id> --tags test

# Short form
acloud -d network vpc list
```

The `--debug` flag enables:
- HTTP request/response logging from the SDK
- Detailed request payloads (JSON formatted)
- Full error response details

Debug output is sent to `stderr`, so it won't interfere with normal command output.

## Documentation

For comprehensive documentation, including detailed guides, resource references, and examples, see the [Documentation](docs/) folder.

- **[Getting Started Guide](docs/getting-started.md)** - Installation, authentication, and basic usage
- **[Resource Documentation](docs/resources/)** - Detailed guides for all resource types:
  - [Management Resources](docs/resources/management.md) - Projects
  - [Storage Resources](docs/resources/storage.md) - Block Storage, Snapshots, Backups, Restores
  - [Network Resources](docs/resources/network.md) - VPCs, Subnets, Security Groups, Elastic IPs, Load Balancers, VPN Tunnels, and more

## Testing

End-to-end (E2E) tests are available to validate CRUD operations across all resource categories:

- **[E2E Tests Documentation](e2e/README.md)** - Comprehensive guide to running E2E tests
- **[Management Tests](e2e/management/test.sh)** - Test projects and organization resources
- **[Storage Tests](e2e/storage/test.sh)** - Test block storage, snapshots, backups, and restores
- **[Network Tests](e2e/network/test.sh)** - Test VPCs, subnets, security groups, VPN tunnels, and more

To run E2E tests:

```bash
# Set required environment variables
export ACLOUD_PROJECT_ID="your-project-id"
export ACLOUD_REGION="ITBG-Bergamo"

# Run tests for a specific category
./e2e/management/test.sh
./e2e/storage/test.sh
./e2e/network/test.sh
```

See the [E2E Tests README](e2e/README.md) for detailed instructions and prerequisites.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## License

See [LICENSE](LICENSE) file for details.
