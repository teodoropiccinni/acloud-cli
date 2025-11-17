# acloud-cli

ArubaCloud Command Line Interface - A CLI tool for interacting with Aruba Cloud APIs.

## Installation

### Windows

```powershell
# Download the latest release
# Extract and add to your PATH
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
