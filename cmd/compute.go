package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var computeCmd = &cobra.Command{
	Use:   "compute",
	Short: "Manage compute resources",
	Long:  `Manage compute resources in Aruba Cloud.`,
}

// Cloudserver subcommands
var cloudserverCmd = &cobra.Command{
	Use:   "cloudserver",
	Short: "Manage cloud servers",
	Long:  `Perform CRUD operations on cloud servers in Aruba Cloud.`,
}

var cloudserverCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new cloud server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cloud server created (stub)")
	},
}

var cloudserverGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get cloud server details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cloud server details (stub)")
	},
}

var cloudserverUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a cloud server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cloud server updated (stub)")
	},
}

var cloudserverDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a cloud server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cloud server deleted (stub)")
	},
}

var cloudserverListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cloud servers",
	Run: func(cmd *cobra.Command, args []string) {
		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// Get projectID from flag or use empty string for default
		projectID, _ := cmd.Flags().GetString("project-id")
		if projectID == "" {
			fmt.Println("Error: --project-id is required")
			return
		}

		// Example: List cloud servers using the SDK
		ctx := context.Background()
		response, err := client.FromCompute().CloudServers().List(ctx, projectID, nil)
		if err != nil {
			fmt.Printf("Error listing cloud servers: %v\n", err)
			return
		}

		if response != nil && response.Data != nil && len(response.Data.Values) > 0 {
			// Define table columns
			headers := []TableColumn{
				{Header: "NAME", Width: 25},
				{Header: "LOCATION", Width: 15},
				{Header: "FLAVOR", Width: 15},
				{Header: "CPU", Width: 15},
				{Header: "RAM(GB)", Width: 15},
				{Header: "HD(GB)", Width: 15},
			}

			// Build rows
			var rows [][]string
			for _, server := range response.Data.Values {
				// Safely get name
				name := ""
				if server.Metadata.Name != "" {
					name = server.Metadata.Name
				}

				// Safely get location
				location := ""
				if server.Metadata.Location.Value != "" {
					location = server.Metadata.Location.Value
				}

				// Safely get flavor details with nil checks
				flavor := ""
				var cpu int32 = 0
				var ram int32 = 0
				var disk int32 = 0
				if server.Properties.Flavor.Name != "" {
					flavor = server.Properties.Flavor.Name
					cpu = server.Properties.Flavor.CPU
					ram = server.Properties.Flavor.RAM
					disk = server.Properties.Flavor.HD
				}

				// Add row
				rows = append(rows, []string{
					name,
					location,
					flavor,
					fmt.Sprintf("%d", cpu),
					fmt.Sprintf("%d", ram),
					fmt.Sprintf("%d", disk),
				})
			}

			// Print the table
			PrintTable(headers, rows)
		} else {
			fmt.Println("No cloud servers found")
		}
	},
}

// Keypair subcommands
var keypairCmd = &cobra.Command{
	Use:   "keypair",
	Short: "Manage keypairs",
	Long:  `Perform CRUD operations on keypairs in Aruba Cloud.`,
}

var keypairCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new keypair",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Keypair created (stub)")
	},
}

var keypairGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get keypair details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Keypair details (stub)")
	},
}

var keypairUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a keypair",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Keypair updated (stub)")
	},
}

var keypairDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a keypair",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Keypair deleted (stub)")
	},
}

var keypairListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all keypairs",
	Run: func(cmd *cobra.Command, args []string) {
		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// Get projectID from flag
		projectID, _ := cmd.Flags().GetString("project-id")
		if projectID == "" {
			fmt.Println("Error: --project-id is required")
			return
		}

		// List keypairs using the SDK
		ctx := context.Background()
		response, err := client.FromCompute().KeyPairs().List(ctx, projectID, nil)
		if err != nil {
			fmt.Printf("Error listing keypairs: %v\n", err)
			return
		}

		if response != nil && response.Data != nil && len(response.Data.Values) > 0 {
			// Define table columns
			headers := []TableColumn{
				{Header: "NAME", Width: 40},
				{Header: "PUBLIC_KEY", Width: 60},
			}

			// Build rows
			var rows [][]string
			for _, keypair := range response.Data.Values {
				name := ""
				if keypair.Metadata.Name != nil && *keypair.Metadata.Name != "" {
					name = *keypair.Metadata.Name
				}

				publicKey := ""
				if keypair.Properties.Value != "" {
					publicKey = keypair.Properties.Value
					// Show only first 50 chars of the key
					if len(publicKey) > 50 {
						publicKey = publicKey[:50] + "..."
					}
				}

				rows = append(rows, []string{name, publicKey})
			}

			// Print the table
			PrintTable(headers, rows)
		} else {
			fmt.Println("No keypairs found")
		}
	},
}

func init() {
	rootCmd.AddCommand(computeCmd)
	computeCmd.AddCommand(cloudserverCmd)
	cloudserverCmd.AddCommand(cloudserverCreateCmd)
	cloudserverCmd.AddCommand(cloudserverGetCmd)
	cloudserverCmd.AddCommand(cloudserverUpdateCmd)
	cloudserverCmd.AddCommand(cloudserverDeleteCmd)
	cloudserverCmd.AddCommand(cloudserverListCmd)

	// Add flags for cloudserver commands
	cloudserverListCmd.Flags().String("project-id", "", "Project ID (required)")
	cloudserverListCmd.MarkFlagRequired("project-id")

	computeCmd.AddCommand(keypairCmd)
	keypairCmd.AddCommand(keypairCreateCmd)
	keypairCmd.AddCommand(keypairGetCmd)
	keypairCmd.AddCommand(keypairUpdateCmd)
	keypairCmd.AddCommand(keypairDeleteCmd)
	keypairCmd.AddCommand(keypairListCmd)

	// Add flags for keypair commands
	keypairListCmd.Flags().String("project-id", "", "Project ID (required)")
	keypairListCmd.MarkFlagRequired("project-id")
}
