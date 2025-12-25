# Container Resources

The `container` category provides commands for managing container resources in Aruba Cloud, including Kubernetes as a Service (KaaS) clusters.

## Available Resources

### [KaaS (Kubernetes as a Service)](container/kaas.md)

KaaS provides managed Kubernetes clusters for running containerized applications.

**Quick Commands:**
```bash
# List all KaaS clusters
acloud container kaas list

# Get KaaS cluster details
acloud container kaas get <cluster-id>

# Create a KaaS cluster
acloud container kaas create --name "my-cluster" --region "ITBG-Bergamo" --version "1.28.0"

# Update a KaaS cluster
acloud container kaas update <cluster-id> --name "new-name"

# Delete a KaaS cluster
acloud container kaas delete <cluster-id>
```

## Common Use Cases

### Creating a Kubernetes Cluster

1. **Create a KaaS cluster**:
   ```bash
   acloud container kaas create \
     --name "production-cluster" \
     --region "ITBG-Bergamo" \
     --version "1.28.0" \
     --tags "production,kubernetes"
   ```

2. **Wait for the cluster to be ready** and check status:
   ```bash
   acloud container kaas get <cluster-id>
   ```

3. **Configure kubectl** (after cluster is ready):
   ```bash
   # Get cluster kubeconfig from the Aruba Cloud console or API
   # Then configure kubectl:
   kubectl config use-context <cluster-context>
   ```

### Managing Cluster Metadata

```bash
# Update cluster name and tags
acloud container kaas update <cluster-id> \
  --name "production-cluster-updated" \
  --tags "production,kubernetes,updated"
```

## Best Practices

- **Naming**: Use descriptive names that indicate the cluster's purpose (e.g., `prod-k8s-cluster`, `staging-cluster`)
- **Tags**: Use tags to organize clusters by environment, project, or team
- **Versioning**: Keep Kubernetes versions up to date for security and features
- **Monitoring**: Monitor cluster status and health regularly
- **Resource Management**: Plan cluster resources based on workload requirements
- **Security**: Follow Kubernetes security best practices for your workloads

## Related Resources

- [Network Resources](./network.md) - Configure networking for container workloads
- [Storage Resources](./storage.md) - Persistent volumes for container applications
- [Security Resources](./security.md) - Security policies and key management

