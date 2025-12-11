# Aruba Cloud SDK Integration

This CLI is now integrated with the official Aruba Cloud Go SDK (`github.com/Arubacloud/sdk-go`).

## Configuration

Before using the CLI commands that interact with the Aruba Cloud API, you must configure your credentials:

```bash
./acloud config set --client-id <your-client-id> --client-secret <your-client-secret>
```

This will store your credentials in `~/.acloud.yaml`.

## How It Works

The CLI uses the `GetArubaClient()` function (defined in `cmd/root.go`) to initialize the SDK client with your stored credentials. This function:

1. Loads your credentials from `~/.acloud.yaml`
2. Creates an SDK client using `aruba.DefaultOptions(clientID, clientSecret)`
3. Returns the initialized client ready to make API calls

## Example: List Cloud Servers

The `cloudserver list` command demonstrates SDK integration:

```bash
./acloud compute cloudserver list --project-id <your-project-id>
```

This command:
- Initializes the SDK client using stored credentials
- Calls `client.FromCompute().CloudServers().List(ctx, projectID, nil)`
- Displays the list of cloud servers in your project

## SDK Client Methods

The Aruba Cloud SDK client provides access to different resource types:

- `client.FromCompute()` - Compute resources (cloud servers, keypairs)
- `client.FromContainer()` - Container resources (KaaS)
- `client.FromDatabase()` - Database resources (DBaaS)
- `client.FromNetwork()` - Network resources (VPC, Load Balancers, etc.)
- `client.FromStorage()` - Storage resources (block storage, snapshots)
- `client.FromSecurity()` - Security resources (KMS)
- `client.FromSchedule()` - Scheduled jobs
- `client.FromProject()` - Project management
- `client.FromAudit()` - Audit logs
- `client.FromMetric()` - Metrics

## Adding SDK Integration to Other Commands

To add SDK integration to other commands, follow this pattern:

```go
package cmd

import (
    "context"
    "fmt"
    "github.com/spf13/cobra"
)

var myResourceListCmd = &cobra.Command{
    Use:   "list",
    Short: "List resources",
    Run: func(cmd *cobra.Command, args []string) {
        // Initialize the SDK client
        client, err := GetArubaClient()
        if err != nil {
            fmt.Printf("Error initializing client: %v\n", err)
            return
        }

        // Get required parameters
        projectID, _ := cmd.Flags().GetString("project-id")
        if projectID == "" {
            fmt.Println("Error: --project-id is required")
            return
        }

        // Make API call
        ctx := context.Background()
        response, err := client.FromXXX().YYY().List(ctx, projectID, nil)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            return
        }

        // Process response
        if response != nil && response.Data != nil {
            // Display results
            fmt.Printf("Found %d resource(s)\n", len(response.Data.Values))
            for i, item := range response.Data.Values {
                fmt.Printf("  %d. Name: %s\n", i+1, item.Metadata.Name)
            }
        }
    },
}
```

## API Response Structure

Most SDK API calls return a `*types.Response[T]` where T is the response type. For example:

```go
type Response[T any] struct {
    Data *T
    // ... other fields
}
```

List responses typically have a `Values` field containing the array of items:

```go
type CloudServerList struct {
    Values []CloudServerResponse `json:"values"`
}
```

## Error Handling

Always check for errors when:
1. Initializing the client with `GetArubaClient()`
2. Making API calls to the SDK

Example:

```go
client, err := GetArubaClient()
if err != nil {
    fmt.Printf("Error initializing client: %v\n", err)
    return
}

response, err := client.FromCompute().CloudServers().List(ctx, projectID, nil)
if err != nil {
    fmt.Printf("Error listing cloud servers: %v\n", err)
    return
}
```

## Next Steps

1. Replace the stub commands with real SDK implementations
2. Add proper error handling and validation
3. Add more flags for filtering and pagination using `types.RequestParameters`
4. Consider adding output formatting options (JSON, table, etc.)
5. Add support for create, update, and delete operations using the SDK
