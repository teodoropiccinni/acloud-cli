package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var storageCmd = &cobra.Command{
	Use:   "storage",
	Short: "Manage storage resources",
	Long:  `Manage storage resources in Aruba Cloud.`,
}

// Blockstorage subcommands
var blockstorageCmd = &cobra.Command{
	Use:   "blockstorage",
	Short: "Manage block storage",
	Long:  `Perform CRUD operations on block storage in Aruba Cloud.`,
}

var blockstorageCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new block storage",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Block storage created (stub)")
	},
}

var blockstorageGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get block storage details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Block storage details (stub)")
	},
}

var blockstorageUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update block storage",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Block storage updated (stub)")
	},
}

var blockstorageDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete block storage",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Block storage deleted (stub)")
	},
}

var blockstorageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all block storage",
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

		// List block storage using the SDK
		ctx := context.Background()
		response, err := client.FromStorage().Volumes().List(ctx, projectID, nil)
		if err != nil {
			fmt.Printf("Error listing block storage: %v\n", err)
			return
		}

		if response != nil && response.Data != nil && len(response.Data.Values) > 0 {
			// Define table columns
			headers := []TableColumn{
				{Header: "NAME", Width: 30},
				{Header: "SIZE(GB)", Width: 12},
				{Header: "REGION", Width: 15},
				{Header: "ZONE", Width: 15},
				{Header: "TYPE", Width: 15},
				{Header: "STATUS", Width: 15},
			}

			// Build rows
			var rows [][]string
			for _, volume := range response.Data.Values {
				name := ""
				if volume.Metadata.Name != nil && *volume.Metadata.Name != "" {
					name = *volume.Metadata.Name
				}

				size := fmt.Sprintf("%d", volume.Properties.SizeGB)

				region := volume.Metadata.LocationResponse.Code
				zone := volume.Properties.Zone

				volumeType := fmt.Sprintf("%v", volume.Properties.Type)

				status := fmt.Sprintf("%v", volume.Status)

				rows = append(rows, []string{name, size, region, zone, volumeType, status})
			}

			// Print the table
			PrintTable(headers, rows)
		} else {
			fmt.Println("No block storage found")
		}
	},
}

// Snapshot subcommands
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage snapshots",
	Long:  `Perform CRUD operations on snapshots in Aruba Cloud.`,
}

var snapshotCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Snapshot created (stub)")
	},
}

var snapshotGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get snapshot details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Snapshot details (stub)")
	},
}

var snapshotUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Snapshot updated (stub)")
	},
}

var snapshotDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Snapshot deleted (stub)")
	},
}

var snapshotListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all snapshots",
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

		// List snapshots using the SDK
		ctx := context.Background()
		response, err := client.FromStorage().Snapshots().List(ctx, projectID, nil)
		if err != nil {
			fmt.Printf("Error listing snapshots: %v\n", err)
			return
		}

		if response != nil && response.Data != nil && len(response.Data.Values) > 0 {
			// Define table columns
			headers := []TableColumn{
				{Header: "NAME", Width: 30},
				{Header: "SIZE(GB)", Width: 12},
				{Header: "SOURCE", Width: 30},
				{Header: "STATUS", Width: 15},
			}

			// Build rows
			var rows [][]string
			for _, snapshot := range response.Data.Values {
				name := ""
				if snapshot.Metadata.Name != nil && *snapshot.Metadata.Name != "" {
					name = *snapshot.Metadata.Name
				}

				size := fmt.Sprintf("%d", snapshot.Properties.SizeGB)

				source := ""
				if snapshot.Properties.Volume != nil && snapshot.Properties.Volume.URI != nil && *snapshot.Properties.Volume.URI != "" {
					source = *snapshot.Properties.Volume.URI
				}

				status := fmt.Sprintf("%v", snapshot.Status)

				rows = append(rows, []string{name, size, source, status})
			}

			// Print the table
			PrintTable(headers, rows)
		} else {
			fmt.Println("No snapshots found")
		}
	},
}

func init() {
	rootCmd.AddCommand(storageCmd)
	storageCmd.AddCommand(blockstorageCmd)
	blockstorageCmd.AddCommand(blockstorageCreateCmd)
	blockstorageCmd.AddCommand(blockstorageGetCmd)
	blockstorageCmd.AddCommand(blockstorageUpdateCmd)
	blockstorageCmd.AddCommand(blockstorageDeleteCmd)
	blockstorageCmd.AddCommand(blockstorageListCmd)

	// Add flags for blockstorage commands
	blockstorageListCmd.Flags().String("project-id", "", "Project ID (required)")
	blockstorageListCmd.MarkFlagRequired("project-id")

	storageCmd.AddCommand(snapshotCmd)
	snapshotCmd.AddCommand(snapshotCreateCmd)
	snapshotCmd.AddCommand(snapshotGetCmd)
	snapshotCmd.AddCommand(snapshotUpdateCmd)
	snapshotCmd.AddCommand(snapshotDeleteCmd)
	snapshotCmd.AddCommand(snapshotListCmd)

	// Add flags for snapshot commands
	snapshotListCmd.Flags().String("project-id", "", "Project ID (required)")
	snapshotListCmd.MarkFlagRequired("project-id")
}
