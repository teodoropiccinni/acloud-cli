# Getting Started with Aruba Cloud CLI

This guide will help you install, configure, and start using the Aruba Cloud CLI.

## Installation

### Download Pre-built Binary

Download the latest release for your platform from the [releases page](https://github.com/Arubacloud/acloud-cli/releases).

#### Linux
```bash
# Download and extract
wget https://github.com/Arubacloud/acloud-cli/releases/latest/download/acloud-linux-amd64.tar.gz
tar -xzf acloud-linux-amd64.tar.gz

# Move to PATH
sudo mv acloud /usr/local/bin/
sudo chmod +x /usr/local/bin/acloud
```

#### macOS
```bash
# Download and extract
wget https://github.com/Arubacloud/acloud-cli/releases/latest/download/acloud-darwin-amd64.tar.gz
tar -xzf acloud-darwin-amd64.tar.gz

# Move to PATH
sudo mv acloud /usr/local/bin/
sudo chmod +x /usr/local/bin/acloud
```

#### Windows
Download `acloud-windows-amd64.zip`, extract it, and add the location to your PATH.

### Build from Source

Requirements:
- Go 1.24 or later

```bash
git clone https://github.com/Arubacloud/acloud-cli.git
cd acloud-cli
go build -o acloud
```

## Authentication

The Aruba Cloud CLI requires API credentials to authenticate with Aruba Cloud services.

### Setting up Credentials

1. **Obtain API Credentials**: Get your Client ID and Client Secret from the Aruba Cloud console.

2. **Configure the CLI**:
   ```bash
   acloud config set
   ```

3. **Enter your credentials** when prompted:
   - Client ID
   - Client Secret

4. **Verify configuration**:
   ```bash
   acloud config show
   ```

### Configuration File

Credentials are stored in `~/.acloud.yaml`:

```yaml
clientId: your-client-id
clientSecret: your-client-secret
```

**Security Note**: Keep your credentials secure. The configuration file contains sensitive information.

### Environment Variables

You can also set credentials via environment variables:

```bash
export ACLOUD_CLIENT_ID="your-client-id"
export ACLOUD_CLIENT_SECRET="your-client-secret"
```

## Auto-completion

The CLI supports shell auto-completion for commands, flags, and resource IDs.

### Bash

#### Current Session
```bash
source <(acloud completion bash)
```

#### Permanent Installation

**Linux:**
```bash
acloud completion bash | sudo tee /etc/bash_completion.d/acloud
```

**macOS:**
```bash
acloud completion bash > $(brew --prefix)/etc/bash_completion.d/acloud
```

After installation, restart your shell or run:
```bash
source ~/.bashrc  # or ~/.bash_profile on macOS
```

### Zsh

Add to `~/.zshrc`:

```bash
# Enable completion
autoload -Uz compinit
compinit

# Load acloud completion
source <(acloud completion zsh)
```

Or for permanent installation:
```bash
acloud completion zsh > "${fpath[1]}/_acloud"
```

### Fish

```bash
acloud completion fish | source
```

Or for permanent installation:
```bash
acloud completion fish > ~/.config/fish/completions/acloud.fish
```

### PowerShell

Add to your PowerShell profile:

```powershell
acloud completion powershell | Out-String | Invoke-Expression
```

## Features of Auto-completion

The auto-completion system provides:

1. **Command completion**: Tab-complete commands and subcommands
   ```bash
   acloud man<TAB>  # completes to "management"
   ```

2. **Flag completion**: Tab-complete available flags
   ```bash
   acloud config set --<TAB>  # shows available flags
   ```

3. **Resource ID completion**: Tab-complete resource IDs with descriptions
   ```bash
   acloud management project get <TAB>
   # Shows:
   # 655b2822af30f667f826994e    defaultproject
   # 66a10244f62b99c686572a9f    develop
   # ...
   ```

## Verifying Installation

Test your installation:

```bash
# Check version
acloud --version

# View available commands
acloud --help

# Test API connectivity
acloud management project list
```

## Next Steps

- Learn about [Project Management](resources/management/projects.md)
- Explore other [Resources](resources/management.md)
- Read the [Command Reference](command-reference.md)

## Troubleshooting

### "Error initializing client"

This usually means credentials are not configured. Run:
```bash
acloud config set
```

### "No projects found"

Ensure your credentials have the correct permissions and you have projects in your account.

### Auto-completion not working

1. Ensure bash-completion is installed:
   ```bash
   # Ubuntu/Debian
   sudo apt-get install bash-completion
   
   # macOS
   brew install bash-completion
   ```

2. Reload your shell configuration:
   ```bash
   source ~/.bashrc
   ```
