# End-to-End (E2E) Testing

This directory contains end-to-end test scripts for validating CRUD operations across all resource categories in the Aruba Cloud CLI.

## Overview

The E2E tests are organized by resource category:
- **[Management](management/)** - Organization-level resources (Projects)
- **[Storage](storage/)** - Storage resources (Block Storage, Snapshots, Backups, Restores)
- **[Network](network/)** - Network resources (VPCs, Subnets, Security Groups, Elastic IPs, VPN Tunnels, etc.)
- **[Database](database/)** - Database resources (DBaaS, Databases, Users, Backups)
- **[Schedule](schedule/)** - Scheduled jobs (OneShot and Recurring jobs)
- **[Security](security/)** - Security resources (KMS Keys)
- **[Compute](compute/)** - Compute resources (Cloud Servers, Key Pairs)
- **[Container](container/)** - Container resources (KaaS clusters)

## Prerequisites

Before running the tests, ensure you have:

1. **Configured CLI credentials:**
   ```bash
   acloud config set --client-id YOUR_CLIENT_ID --client-secret YOUR_CLIENT_SECRET
   ```

2. **Set required environment variables** (see each test script for specific requirements):
   ```bash
   export ACLOUD_PROJECT_ID="your-project-id"
   export ACLOUD_REGION="ITBG-Bergamo"
   # ... other variables as needed
   ```

3. **Built the CLI:**
   ```bash
   go build -o acloud .
   ```

## Running Tests

### Run All Tests

To run all E2E tests:

```bash
# From the project root
./e2e/management/test.sh
./e2e/storage/test.sh
./e2e/network/test.sh
./e2e/database/test.sh
./e2e/schedule/test.sh
./e2e/security/test.sh
./e2e/compute/test.sh
./e2e/container/test.sh
```

### Run Individual Category Tests

Each category has its own test script:

```bash
# Management resources
./e2e/management/test.sh

# Storage resources
./e2e/storage/test.sh

# Network resources
./e2e/network/test.sh

# Database resources
./e2e/database/test.sh

# Schedule resources
./e2e/schedule/test.sh

# Security resources
./e2e/security/test.sh

# Compute resources
./e2e/compute/test.sh

# Container resources
./e2e/container/test.sh
```

## Test Structure

Each test script follows a consistent structure:

1. **Configuration** - Sets up environment variables and validates prerequisites
2. **CREATE** - Creates test resources
3. **LIST** - Lists created resources
4. **GET** - Retrieves detailed information about resources
5. **UPDATE** - Updates resource properties (name, tags, etc.)
6. **DELETE** - Cleans up test resources

## Test Scripts

### Management Tests

Tests organization-level resources:
- **Projects** - Create, list, get, update, delete projects

See [management/test.sh](management/test.sh) for details.

### Storage Tests

Tests storage resources:
- **Block Storage** - Volume creation, updates, deletion
- **Snapshots** - Snapshot creation from volumes
- **Backups** - Backup operations
- **Restores** - Restore operations from backups

See [storage/test.sh](storage/test.sh) for details.

### Network Tests

Tests network resources:
- **VPC** - Virtual Private Cloud management
- **Subnet** - Subnet creation and management
- **Security Group** - Security group operations
- **Security Rule** - Firewall rule management
- **Elastic IP** - Public IP address management
- **VPC Peering** - VPC peering connections
- **VPC Peering Route** - Peering route management
- **VPN Tunnel** - VPN tunnel creation and management
- **VPN Route** - VPN route configuration

See [network/test.sh](network/test.sh) for details.

### Database Tests

Tests database resources:
- **DBaaS** - Database as a Service instance management
- **DBaaS Databases** - Database creation and management within DBaaS
- **DBaaS Users** - User management for DBaaS instances
- **Database Backups** - Backup operations for databases

See [database/test.sh](database/test.sh) for details.

### Schedule Tests

Tests scheduled job resources:
- **OneShot Jobs** - One-time scheduled jobs
- **Recurring Jobs** - Recurring scheduled jobs with CRON expressions

See [schedule/test.sh](schedule/test.sh) for details.

### Security Tests

Tests security resources:
- **KMS Keys** - Key Management System key operations

See [security/test.sh](security/test.sh) for details.

### Compute Tests

Tests compute resources:
- **Cloud Servers** - Virtual machine instances
- **Key Pairs** - SSH key pair management

See [compute/test.sh](compute/test.sh) for details.

### Container Tests

Tests container resources:
- **KaaS Clusters** - Kubernetes as a Service cluster management
  - Create, list, get, update, delete KaaS clusters
  - Connect to clusters and configure kubectl (requires kubectl installed)

**Required Environment Variables:**
- `ACLOUD_VPC_URI` - VPC URI for the cluster (e.g., `/projects/{project-id}/providers/Aruba.Network/vpcs/{vpc-id}`)
- `ACLOUD_SUBNET_URI` - Subnet URI for the cluster (e.g., `/projects/{project-id}/providers/Aruba.Network/subnets/{subnet-id}`)
- `ACLOUD_NODE_POOL_INSTANCE` - Instance configuration name for nodes
- `ACLOUD_NODE_POOL_ZONE` - Datacenter/zone code for nodes

**Optional Environment Variables:**
- `ACLOUD_NODE_CIDR` - Node CIDR address (default: `10.0.0.0/16`)
- `ACLOUD_NODE_CIDR_NAME` - Node CIDR name (default: `node-cidr`)
- `ACLOUD_SECURITY_GROUP_NAME` - Security group name (default: `kaas-sg`)
- `ACLOUD_NODE_POOL_NAME` - Node pool name (default: `default-pool`)
- `ACLOUD_NODE_POOL_NODES` - Number of nodes (default: `1`)
- `ACLOUD_K8S_VERSION` - Kubernetes version (default: `1.28.0`)

**Note:** The connect test requires `kubectl` to be installed and available in PATH. If kubectl is not found, the connect test will be skipped.

See [container/test.sh](container/test.sh) for details.

## Environment Variables

Common environment variables used across tests:

| Variable | Description | Example |
|----------|-------------|---------|
| `ACLOUD_PROJECT_ID` | Project ID for resources | `66a10244f62b99c686572a9f` |
| `ACLOUD_REGION` | Region code | `ITBG-Bergamo` |
| `ACLOUD_VPC_ID` | VPC ID for network resources | `69495ef64d0cdc87949b71ec` |
| `ACLOUD_PEER_VPC_ID` | Peer VPC ID for peering | `69485a584d0cdc87949b6ff8` |
| `ACLOUD_VPC_URI` | VPC URI for KaaS clusters | `/projects/{id}/providers/Aruba.Network/vpcs/{vpc-id}` |
| `ACLOUD_SUBNET_URI` | Subnet URI for KaaS clusters | `/projects/{id}/providers/Aruba.Network/subnets/{subnet-id}` |
| `ACLOUD_NODE_POOL_INSTANCE` | Instance type for KaaS node pool | `small` |
| `ACLOUD_NODE_POOL_ZONE` | Zone for KaaS node pool | `ITBG-Bergamo-A` |

See individual test scripts for category-specific variables.

## Debug Mode

To see detailed HTTP requests/responses during tests, use the `--debug` flag:

```bash
# The test scripts can be modified to add --debug to commands
# Or run individual commands with debug:
acloud --debug network vpc create --name test-vpc --region ITBG-Bergamo
```

## Test Output

Each test script provides:
- ✅ **Success indicators** - Green checkmarks for passed tests
- ❌ **Error messages** - Red text for failed operations
- 📊 **Resource information** - Tables showing created/updated resources

## Cleanup

Test scripts attempt to clean up created resources, but if a test fails:
1. Manually delete any remaining test resources
2. Check resource names/IDs in the test output
3. Use the CLI to list and delete orphaned resources

## Troubleshooting

### "Error: project ID not specified"
- Set `ACLOUD_PROJECT_ID` environment variable
- Or use `--project-id` flag in commands
- Or set up a context: `acloud context set my-prod --project-id <id>`

### "Error: Unable to determine region value"
- Ensure `ACLOUD_REGION` is set correctly
- Use the correct region format (e.g., `ITBG-Bergamo`)

### "Failed to create ... - Status: 400"
- Verify all required parameters are provided
- Check that dependent resources exist (e.g., VPC for subnet)
- Review error details in the output

### "Cannot update ... while it is in 'InCreation' state"
- Wait for resources to reach `Active` state before updating
- Some resources take time to provision

## Contributing

When adding new resources or operations:
1. Add test cases to the appropriate category script
2. Follow the existing test structure (CREATE → LIST → GET → UPDATE → DELETE)
3. Include proper cleanup in case of failures
4. Update this README with new test coverage

## See Also

- [Main Documentation](../docs/)
- [Getting Started Guide](../docs/getting-started.md)
- [Resource Documentation](../docs/resources/)

