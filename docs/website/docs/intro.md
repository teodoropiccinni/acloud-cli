# Welcome to Aruba Cloud CLI Documentation

Aruba Cloud CLI (`acloud`) is a command-line interface for managing your Aruba Cloud resources.

## Quick Start

```bash
# Configure credentials
acloud config set

# Set up a context (optional but recommended)
acloud context set my-prod --project-id "your-project-id"
acloud context use my-prod

# List projects
acloud management project list

# List storage resources
acloud storage blockstorage list
acloud storage snapshot list

# List network resources
acloud network vpc list
acloud network elasticip list
```

## Features

- **Full CRUD Operations**: Create, read, update, and delete operations for all supported resources
- **Context Management**: Switch between different projects and configurations easily
- **Auto-completion**: Shell completion for resource IDs and commands
- **Debug Mode**: Detailed logging for troubleshooting

## Installation

See the [Getting Started](getting-started.md) guide for installation instructions for your platform.

## Resources

- [Management Resources](resources/management.md) - Projects
- [Storage Resources](resources/storage.md) - Block Storage, Snapshots, Backups, Restores
- [Network Resources](resources/network.md) - VPCs, Subnets, Security Groups, Elastic IPs, and more

