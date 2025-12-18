# Aruba Cloud CLI Documentation

Welcome to the Aruba Cloud CLI (`acloud`) documentation. This CLI provides a powerful command-line interface for managing your Aruba Cloud resources.

## Table of Contents

- [Getting Started](getting-started.md)
  - [Installation](getting-started.md#installation)
  - [Authentication](getting-started.md#authentication)
  - [Context Management](getting-started.md#context-management)
  - [Auto-completion](getting-started.md#auto-completion)
- [Resources](resources/)
  - [Management](resources/management.md)
    - [Projects](resources/management/projects.md)
  - [Storage](resources/storage.md)
    - Block Storage
    - Snapshots

## Quick Start

```bash
# Configure credentials
acloud config set

# Set up a context (optional but recommended)
acloud context set my-prod --project-id "your-project-id"
acloud context use my-prod

# List projects
acloud management project list

# List storage resources (uses context)
acloud storage blockstorage list
acloud storage snapshot list
```

## Getting Help

For help with any command, use the `--help` flag:

```bash
acloud --help
acloud management --help
acloud management project --help
acloud management project create --help
```

## Additional Resources

- [GitHub Repository](https://github.com/Arubacloud/acloud-cli)
- [Aruba Cloud Documentation](https://kb.arubacloud.com/en/home.aspx)
- [SDK Documentation](https://github.com/Arubacloud/sdk-go)
