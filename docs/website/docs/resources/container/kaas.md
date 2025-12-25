# KaaS (Kubernetes as a Service)

KaaS provides managed Kubernetes clusters for running containerized applications in Aruba Cloud.

## Command Syntax

```bash
acloud container kaas <command> [flags] [arguments]
```

## Available Commands

### `create`

Create a new KaaS cluster.

**Syntax:**
```bash
acloud container kaas create [flags]
```

**Required Flags:**
- `--name <string>` - Name for the KaaS cluster
- `--region <string>` - Region code (e.g., `ITBG-Bergamo`)
- `--version <string>` - Kubernetes version (e.g., `1.28.0`)

**Optional Flags:**
- `--project-id <string>` - Project ID (uses context if not specified)
- `--tags <stringSlice>` - Tags (comma-separated)

**Example:**
```bash
acloud container kaas create \
  --name "production-cluster" \
  --region "ITBG-Bergamo" \
  --version "1.28.0" \
  --tags "production,kubernetes"
```

### `list`

List all KaaS clusters in the project.

**Syntax:**
```bash
acloud container kaas list [flags]
```

**Optional Flags:**
- `--project-id <string>` - Project ID (uses context if not specified)

**Example:**
```bash
acloud container kaas list
```

**Output:**
The command displays a table with the following columns:
- ID
- NAME
- VERSION
- REGION
- STATUS

### `get`

Get detailed information about a specific KaaS cluster.

**Syntax:**
```bash
acloud container kaas get <cluster-id> [flags]
```

**Arguments:**
- `<cluster-id>` - The ID of the KaaS cluster

**Optional Flags:**
- `--project-id <string>` - Project ID (uses context if not specified)
- `--verbose` - Show detailed JSON output

**Example:**
```bash
acloud container kaas get 69495ef64d0cdc87949b71ec
```

**Output:**
The command displays detailed information including:
- ID and URI
- Name and region
- Kubernetes version
- Status
- Creation date and creator
- Tags

### `update`

Update a KaaS cluster's properties (name, tags).

**Syntax:**
```bash
acloud container kaas update <cluster-id> [flags]
```

**Arguments:**
- `<cluster-id>` - The ID of the KaaS cluster

**Optional Flags:**
- `--project-id <string>` - Project ID (uses context if not specified)
- `--name <string>` - New name for the KaaS cluster
- `--tags <stringSlice>` - New tags (comma-separated)

**Example:**
```bash
acloud container kaas update 69495ef64d0cdc87949b71ec \
  --name "production-cluster-updated" \
  --tags "production,kubernetes,updated"
```

**Note:** At least one of `--name` or `--tags` must be provided.

### `delete`

Delete a KaaS cluster.

**Syntax:**
```bash
acloud container kaas delete <cluster-id> [flags]
```

**Arguments:**
- `<cluster-id>` - The ID of the KaaS cluster

**Optional Flags:**
- `--project-id <string>` - Project ID (uses context if not specified)
- `--yes, -y` - Skip confirmation prompt

**Example:**
```bash
acloud container kaas delete 69495ef64d0cdc87949b71ec --yes
```

**Warning:** Deleting a cluster will remove all workloads and data. Ensure you have backups if needed.

## Auto-completion

The CLI provides auto-completion for KaaS cluster IDs:

```bash
acloud container kaas get <TAB>
acloud container kaas update <TAB>
acloud container kaas delete <TAB>
```

## Common Workflows

### Creating a New Cluster

1. **Create the KaaS cluster**:
   ```bash
   acloud container kaas create \
     --name "my-k8s-cluster" \
     --region "ITBG-Bergamo" \
     --version "1.28.0"
   ```

2. **Wait for the cluster to be ready**:
   ```bash
   # Check cluster status
   acloud container kaas get <cluster-id>
   # Wait until status is "Active"
   ```

3. **Configure kubectl** (after cluster is ready):
   ```bash
   # Get kubeconfig from Aruba Cloud console or API
   # Export kubeconfig:
   export KUBECONFIG=/path/to/kubeconfig
   
   # Verify connection:
   kubectl get nodes
   ```

### Updating Cluster Metadata

```bash
# Update cluster name and tags
acloud container kaas update <cluster-id> \
  --name "production-cluster" \
  --tags "production,kubernetes,updated"
```

### Listing and Filtering Clusters

```bash
# List all clusters
acloud container kaas list

# Use grep to filter by name or tags
acloud container kaas list | grep "production"
```

## Best Practices

- **Naming**: Use descriptive names that indicate the cluster's purpose (e.g., `prod-k8s-cluster`, `staging-cluster`)
- **Tags**: Use tags to organize clusters by environment, project, or team
- **Versioning**: 
  - Keep Kubernetes versions up to date for security and features
  - Test version upgrades in staging before production
- **Resource Planning**: Plan cluster resources based on workload requirements
- **Monitoring**: Monitor cluster status, health, and resource usage
- **Security**: 
  - Follow Kubernetes security best practices
  - Use RBAC for access control
  - Enable network policies
  - Regularly update cluster components
- **Backup**: Implement backup strategies for persistent volumes and application data
- **Cleanup**: Delete unused clusters to avoid unnecessary costs

## Related Resources

- [Network Resources](../network.md) - Configure networking for container workloads
- [Storage Resources](../storage.md) - Persistent volumes for container applications
- [Security Resources](../security.md) - Security policies and key management
- [Compute Resources](../compute.md) - Cloud servers that can run container workloads

