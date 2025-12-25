# Cloud Servers

Cloud servers are virtual machine instances that run your applications and workloads in Aruba Cloud.

## Command Syntax

```bash
acloud compute cloudserver <command> [flags] [arguments]
```

## Available Commands

### `create`

Create a new cloud server instance.

**Syntax:**
```bash
acloud compute cloudserver create [flags]
```

**Required Flags:**
- `--name <string>` - Name for the cloud server
- `--region <string>` - Region code (e.g., `ITBG-Bergamo`)
- `--flavor <string>` - Flavor name (e.g., `small`, `medium`, `large`)
- `--image <string>` - Image ID or name

**Optional Flags:**
- `--project-id <string>` - Project ID (uses context if not specified)
- `--keypair <string>` - Key pair name for SSH access
- `--tags <stringSlice>` - Tags (comma-separated)

**Example:**
```bash
acloud compute cloudserver create \
  --name "web-server" \
  --region "ITBG-Bergamo" \
  --flavor "small" \
  --image "ubuntu-22.04" \
  --keypair "my-keypair" \
  --tags "production,web"
```

### `list`

List all cloud servers in the project.

**Syntax:**
```bash
acloud compute cloudserver list [flags]
```

**Optional Flags:**
- `--project-id <string>` - Project ID (uses context if not specified)

**Example:**
```bash
acloud compute cloudserver list
```

**Output:**
The command displays a table with the following columns:
- ID
- NAME
- LOCATION
- FLAVOR
- CPU
- RAM(GB)
- HD(GB)
- STATUS

### `get`

Get detailed information about a specific cloud server.

**Syntax:**
```bash
acloud compute cloudserver get <server-id> [flags]
```

**Arguments:**
- `<server-id>` - The ID of the cloud server

**Optional Flags:**
- `--project-id <string>` - Project ID (uses context if not specified)
- `--verbose` - Show detailed JSON output

**Example:**
```bash
acloud compute cloudserver get 69495ef64d0cdc87949b71ec
```

**Output:**
The command displays detailed information including:
- ID and URI
- Name and region
- Flavor details (CPU, RAM, HD)
- Image information
- Key pair (if configured)
- Status
- Creation date and creator
- Tags

### `update`

Update a cloud server's properties (name, tags).

**Syntax:**
```bash
acloud compute cloudserver update <server-id> [flags]
```

**Arguments:**
- `<server-id>` - The ID of the cloud server

**Optional Flags:**
- `--project-id <string>` - Project ID (uses context if not specified)
- `--name <string>` - New name for the cloud server
- `--tags <stringSlice>` - New tags (comma-separated)

**Example:**
```bash
acloud compute cloudserver update 69495ef64d0cdc87949b71ec \
  --name "web-server-updated" \
  --tags "production,web,updated"
```

**Note:** At least one of `--name` or `--tags` must be provided.

### `delete`

Delete a cloud server instance.

**Syntax:**
```bash
acloud compute cloudserver delete <server-id> [flags]
```

**Arguments:**
- `<server-id>` - The ID of the cloud server

**Optional Flags:**
- `--project-id <string>` - Project ID (uses context if not specified)
- `--yes, -y` - Skip confirmation prompt

**Example:**
```bash
acloud compute cloudserver delete 69495ef64d0cdc87949b71ec --yes
```

## Auto-completion

The CLI provides auto-completion for cloud server IDs:

```bash
acloud compute cloudserver get <TAB>
acloud compute cloudserver update <TAB>
acloud compute cloudserver delete <TAB>
```

## Common Workflows

### Launching a New Server

1. **Create a key pair** (if needed):
   ```bash
   acloud compute keypair create --name "my-keypair" --public-key "$(cat ~/.ssh/id_rsa.pub)"
   ```

2. **Create the cloud server**:
   ```bash
   acloud compute cloudserver create \
     --name "app-server" \
     --region "ITBG-Bergamo" \
     --flavor "medium" \
     --image "your-image-id" \
     --keypair "my-keypair"
   ```

3. **Wait for the server to be ready** and check status:
   ```bash
   acloud compute cloudserver get <server-id>
   ```

### Updating Server Metadata

```bash
# Update server name and tags
acloud compute cloudserver update <server-id> \
  --name "new-name" \
  --tags "production,updated"
```

### Listing and Filtering Servers

```bash
# List all servers
acloud compute cloudserver list

# Use grep to filter by name or tags
acloud compute cloudserver list | grep "production"
```

## Best Practices

- **Naming**: Use descriptive names that indicate the server's purpose (e.g., `web-server-prod`, `db-server-staging`)
- **Tags**: Use tags to organize servers by environment, project, or team
- **Flavors**: Choose appropriate flavors based on your workload requirements
- **Key Pairs**: Always use key pairs for SSH access instead of passwords
- **Monitoring**: Check server status before performing operations
- **Cleanup**: Delete unused servers to avoid unnecessary costs

## Related Resources

- [Key Pairs](./keypair.md) - Manage SSH key pairs for server access
- [Network Resources](../network.md) - Configure networking and security groups
- [Storage Resources](../storage.md) - Attach block storage volumes

