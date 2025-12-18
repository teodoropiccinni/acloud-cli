# Aruba Cloud CLI Documentation

Welcome to the Aruba Cloud CLI (`acloud`) documentation. This CLI provides a powerful command-line interface for managing your Aruba Cloud resources.

## Table of Contents

- [Getting Started](getting-started.md)
  - [Installation](getting-started.md#installation)
  - [Authentication](getting-started.md#authentication)
  - [Auto-completion](getting-started.md#auto-completion)
- [Resources](resources/)
  - [Management](resources/management.md)
    - [Projects](resources/management/projects.md)

## Quick Start

```bash
# Configure credentials
acloud config set

# List projects
acloud management project list

# Get project details
acloud management project get <project-id>
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
