package cmd

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/types"
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
		// Get project ID from flag or context
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Get flags
		name, _ := cmd.Flags().GetString("name")
		region, _ := cmd.Flags().GetString("region")
		zone, _ := cmd.Flags().GetString("zone")
		size, _ := cmd.Flags().GetInt("size")
		volumeType, _ := cmd.Flags().GetString("type")
		billingPeriod, _ := cmd.Flags().GetString("billing-period")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		// Validate required fields
		if name == "" {
			fmt.Println("Error: --name is required")
			return
		}
		if region == "" {
			fmt.Println("Error: --region is required")
			return
		}
		if zone == "" {
			fmt.Println("Error: --zone is required")
			return
		}
		if size <= 0 {
			fmt.Println("Error: --size must be greater than 0")
			return
		}

		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// Build the create request
		createRequest := types.BlockStorageRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: name,
					Tags: tags,
				},
				Location: types.LocationRequest{
					Value: region,
				},
			},
			Properties: types.BlockStoragePropertiesRequest{
				SizeGB:        size,
				BillingPeriod: billingPeriod,
				Zone:          zone,
				Type:          types.BlockStorageType(volumeType),
			},
		}

		// Debug: print request
		fmt.Printf("\nCreating block storage with:\n")
		fmt.Printf("  Name: %s\n", name)
		fmt.Printf("  Region: %s\n", region)
		fmt.Printf("  Zone: %s\n", zone)
		fmt.Printf("  Size: %d GB\n", size)
		fmt.Printf("  Type: %s\n", volumeType)
		fmt.Printf("  Billing Period: %s\n", billingPeriod)
		fmt.Printf("  Project ID: %s\n", projectID)

		// Create the block storage using the SDK
		ctx := context.Background()
		response, err := client.FromStorage().Volumes().Create(ctx, projectID, createRequest, nil)
		if err != nil {
			fmt.Printf("Error creating block storage: %v\n", err)
			return
		}

		if response == nil {
			fmt.Println("Error: received nil response from API")
			return
		}

		if !response.IsSuccess() {
			fmt.Printf("Error: failed to create block storage - Status: %d\n", response.StatusCode)
			if response.Error.Title != nil {
				fmt.Printf("Error Title: %s\n", *response.Error.Title)
			}
			if response.Error.Detail != nil {
				fmt.Printf("Error Detail: %s\n", *response.Error.Detail)
			}
			// Print full error response for debugging
			fmt.Printf("Full Error Response: %+v\n", response.Error)
			return
		}

		if response.Data != nil {
			fmt.Println("\nBlock storage created successfully!")
			fmt.Printf("ID:              %s\n", *response.Data.Metadata.ID)
			fmt.Printf("Name:            %s\n", *response.Data.Metadata.Name)
			fmt.Printf("Size (GB):       %d\n", response.Data.Properties.SizeGB)
			fmt.Printf("Type:            %s\n", response.Data.Properties.Type)
			fmt.Printf("Zone:            %s\n", response.Data.Properties.Zone)
			fmt.Printf("Region:          %s\n", response.Data.Metadata.LocationResponse.Code)
			if response.Data.Status.State != nil {
				fmt.Printf("Status:          %s\n", *response.Data.Status.State)
			}
			if !response.Data.Metadata.CreationDate.IsZero() {
				fmt.Printf("Creation Date:   %s\n", response.Data.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}
		} else {
			fmt.Println("Block storage created but no details returned")
		}
	},
}

var blockstorageGetCmd = &cobra.Command{
	Use:   "get [volume-id]",
	Short: "Get block storage details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		volumeID := args[0]

		// Get project ID from flag or context
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// Get block storage details using the SDK
		ctx := context.Background()
		response, err := client.FromStorage().Volumes().Get(ctx, projectID, volumeID, nil)
		if err != nil {
			fmt.Printf("Error getting block storage details: %v\n", err)
			return
		}

		if response != nil && response.Data != nil {
			volume := response.Data

			// Display volume details
			fmt.Println("\nBlock Storage Details:")
			fmt.Println("======================")

			if volume.Metadata.ID != nil {
				fmt.Printf("ID:              %s\n", *volume.Metadata.ID)
			}

			if volume.Metadata.Name != nil {
				fmt.Printf("Name:            %s\n", *volume.Metadata.Name)
			}

			fmt.Printf("Size (GB):       %d\n", volume.Properties.SizeGB)
			fmt.Printf("Type:            %s\n", volume.Properties.Type)
			fmt.Printf("Zone:            %s\n", volume.Properties.Zone)

			if volume.Metadata.LocationResponse != nil {
				fmt.Printf("Region:          %s\n", volume.Metadata.LocationResponse.Code)
			}

			if volume.Properties.Bootable != nil {
				fmt.Printf("Bootable:        %t\n", *volume.Properties.Bootable)
			}

			if volume.Status.State != nil {
				fmt.Printf("Status:          %s\n", *volume.Status.State)
			}

			if !volume.Metadata.CreationDate.IsZero() {
				fmt.Printf("Creation Date:   %s\n", volume.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}

			if volume.Metadata.CreatedBy != nil {
				fmt.Printf("Created By:      %s\n", *volume.Metadata.CreatedBy)
			}

			if len(volume.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", volume.Metadata.Tags)
			}

			fmt.Println()
		} else {
			fmt.Println("Block storage not found")
		}
	},
}

var blockstorageUpdateCmd = &cobra.Command{
	Use:   "update [volume-id]",
	Short: "Update block storage (name and/or tags only)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		volumeID := args[0]

		// Get project ID from flag or context
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Get flags
		name, _ := cmd.Flags().GetString("name")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		// At least one field must be provided
		if name == "" && !cmd.Flags().Changed("tags") {
			fmt.Println("Error: at least one of --name or --tags must be provided")
			fmt.Println("Note: Size update is not supported by the API yet")
			return
		}

		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// First, get the current volume details to preserve existing values
		ctx := context.Background()
		getResponse, err := client.FromStorage().Volumes().Get(ctx, projectID, volumeID, nil)
		if err != nil {
			fmt.Printf("Error getting block storage details: %v\n", err)
			return
		}

		if getResponse == nil || getResponse.Data == nil {
			fmt.Println("Block storage not found")
			return
		}

		currentVolume := getResponse.Data

		// Build the update request with current values as defaults
		updateRequest := types.BlockStorageRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: *currentVolume.Metadata.Name,
					Tags: currentVolume.Metadata.Tags,
				},
				Location: types.LocationRequest{
					Value: currentVolume.Metadata.LocationResponse.Code,
				},
			},
			Properties: types.BlockStoragePropertiesRequest{
				SizeGB:        currentVolume.Properties.SizeGB,
				BillingPeriod: "Hour",
				Zone:          currentVolume.Properties.Zone,
				Type:          currentVolume.Properties.Type,
			},
		}

		// Update only the fields that were provided
		if name != "" {
			updateRequest.Metadata.ResourceMetadataRequest.Name = name
		}

		if cmd.Flags().Changed("tags") {
			updateRequest.Metadata.ResourceMetadataRequest.Tags = tags
		}

		// Update the block storage using the SDK
		response, err := client.FromStorage().Volumes().Update(ctx, projectID, volumeID, updateRequest, nil)
		if err != nil {
			fmt.Printf("Error updating block storage: %v\n", err)
			return
		}

		if response != nil && response.Data != nil {
			fmt.Println("\nBlock storage updated successfully!")
			fmt.Printf("ID:              %s\n", *response.Data.Metadata.ID)
			fmt.Printf("Name:            %s\n", *response.Data.Metadata.Name)
			if len(response.Data.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", response.Data.Metadata.Tags)
			}
			fmt.Printf("Size (GB):       %d\n", response.Data.Properties.SizeGB)
			fmt.Printf("Type:            %s\n", response.Data.Properties.Type)
		}
	},
}

var blockstorageDeleteCmd = &cobra.Command{
	Use:   "delete [volume-id]",
	Short: "Delete block storage",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		volumeID := args[0]

		// Get project ID from flag or context
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Get flags
		confirm, _ := cmd.Flags().GetBool("yes")

		// If not confirmed, ask for confirmation
		if !confirm {
			fmt.Printf("Are you sure you want to delete block storage %s? (yes/no): ", volumeID)
			var response string
			fmt.Scanln(&response)
			if response != "yes" && response != "y" {
				fmt.Println("Delete cancelled")
				return
			}
		}

		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// Delete the block storage using the SDK
		ctx := context.Background()
		_, err = client.FromStorage().Volumes().Delete(ctx, projectID, volumeID, nil)
		if err != nil {
			fmt.Printf("Error deleting block storage: %v\n", err)
			return
		}

		fmt.Printf("\nBlock storage %s deleted successfully!\n", volumeID)
	},
}

var blockstorageListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all block storage",
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		verbose, _ := cmd.Flags().GetBool("verbose")

		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// Get project ID from flag or context
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// List block storage using the SDK
		ctx := context.Background()
		response, err := client.FromStorage().Volumes().List(ctx, projectID, nil)
		if err != nil {
			fmt.Printf("Error listing block storage: %v\n", err)
			return
		}

		// Debug output
		if verbose {
			fmt.Println("\n=== DEBUG: Raw Response ===")
			fmt.Printf("Status Code: %d\n", response.StatusCode)
			fmt.Printf("Success: %v\n", response.IsSuccess())
			if response.Data != nil {
				fmt.Printf("Number of volumes: %d\n", len(response.Data.Values))
				if len(response.Data.Values) > 0 {
					fmt.Println("\n=== Volumes Detail ===")
					for i, vol := range response.Data.Values {
						fmt.Printf("\n--- Volume %d ---\n", i+1)
						if vol.Metadata.ID != nil {
							fmt.Printf("  ID: %s\n", *vol.Metadata.ID)
						}
						if vol.Metadata.Name != nil {
							fmt.Printf("  Name: %s\n", *vol.Metadata.Name)
						}
						fmt.Printf("  Size: %d GB\n", vol.Properties.SizeGB)
						fmt.Printf("  Type: %v\n", vol.Properties.Type)
						fmt.Printf("  Zone: %s\n", vol.Properties.Zone)
						fmt.Printf("  Region: %s\n", vol.Metadata.LocationResponse.Code)
						if vol.Status.State != nil {
							fmt.Printf("  Status State: %s\n", *vol.Status.State)
						} else {
							fmt.Printf("  Status State: <nil>\n")
						}
						fmt.Printf("  Full Status: %+v\n", vol.Status)
					}
				}
			}
			fmt.Println("\n=== End Debug ===\n")
		}

		if response != nil && response.Data != nil && len(response.Data.Values) > 0 {
			// Define table columns
			headers := []TableColumn{
				{Header: "NAME", Width: 30},
				{Header: "ID", Width: 26},
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

				id := ""
				if volume.Metadata.ID != nil {
					id = *volume.Metadata.ID
				}

				size := fmt.Sprintf("%d", volume.Properties.SizeGB)

				region := volume.Metadata.LocationResponse.Code
				zone := volume.Properties.Zone

				volumeType := fmt.Sprintf("%v", volume.Properties.Type)

				status := ""
				if volume.Status.State != nil {
					status = *volume.Status.State
				}

				rows = append(rows, []string{name, id, size, region, zone, volumeType, status})
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
		// Get project ID from flag or context
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Get flags
		name, _ := cmd.Flags().GetString("name")
		region, _ := cmd.Flags().GetString("region")
		volumeURI, _ := cmd.Flags().GetString("volume-uri")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		// Validate required fields
		if name == "" {
			fmt.Println("Error: --name is required")
			return
		}
		if region == "" {
			fmt.Println("Error: --region is required")
			return
		}
		if volumeURI == "" {
			fmt.Println("Error: --volume-uri is required")
			return
		}

		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// Build the create request
		createRequest := types.SnapshotRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: name,
					Tags: tags,
				},
				Location: types.LocationRequest{
					Value: region,
				},
			},
			Properties: types.SnapshotPropertiesRequest{
				Volume: types.ReferenceResource{
					URI: volumeURI,
				},
			},
		}

		// Create the snapshot using the SDK
		ctx := context.Background()
		response, err := client.FromStorage().Snapshots().Create(ctx, projectID, createRequest, nil)
		if err != nil {
			fmt.Printf("Error creating snapshot: %v\n", err)
			return
		}

		if response != nil && response.Data != nil {
			fmt.Println("\nSnapshot created successfully!")
			fmt.Printf("ID:              %s\n", *response.Data.Metadata.ID)
			fmt.Printf("Name:            %s\n", *response.Data.Metadata.Name)
			fmt.Printf("Size (GB):       %d\n", response.Data.Properties.SizeGB)
			if !response.Data.Metadata.CreationDate.IsZero() {
				fmt.Printf("Creation Date:   %s\n", response.Data.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}
		}
	},
}

var snapshotGetCmd = &cobra.Command{
	Use:   "get [snapshot-id]",
	Short: "Get snapshot details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		snapshotID := args[0]

		// Get project ID from flag or context
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// Get snapshot details using the SDK
		ctx := context.Background()
		response, err := client.FromStorage().Snapshots().Get(ctx, projectID, snapshotID, nil)
		if err != nil {
			fmt.Printf("Error getting snapshot details: %v\n", err)
			return
		}

		if response != nil && response.Data != nil {
			snapshot := response.Data

			// Display snapshot details
			fmt.Println("\nSnapshot Details:")
			fmt.Println("=================")

			if snapshot.Metadata.ID != nil {
				fmt.Printf("ID:              %s\n", *snapshot.Metadata.ID)
			}

			if snapshot.Metadata.Name != nil {
				fmt.Printf("Name:            %s\n", *snapshot.Metadata.Name)
			}

			fmt.Printf("Size (GB):       %d\n", snapshot.Properties.SizeGB)

			if snapshot.Properties.Volume != nil && snapshot.Properties.Volume.URI != nil {
				fmt.Printf("Source Volume:   %s\n", *snapshot.Properties.Volume.URI)
			}

			if snapshot.Metadata.LocationResponse != nil {
				fmt.Printf("Region:          %s\n", snapshot.Metadata.LocationResponse.Code)
			}

			fmt.Printf("Status:          %v\n", snapshot.Status)

			if !snapshot.Metadata.CreationDate.IsZero() {
				fmt.Printf("Creation Date:   %s\n", snapshot.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}

			if snapshot.Metadata.CreatedBy != nil {
				fmt.Printf("Created By:      %s\n", *snapshot.Metadata.CreatedBy)
			}

			if len(snapshot.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", snapshot.Metadata.Tags)
			}

			fmt.Println()
		} else {
			fmt.Println("Snapshot not found")
		}
	},
}

var snapshotUpdateCmd = &cobra.Command{
	Use:   "update [snapshot-id]",
	Short: "Update a snapshot (name and/or tags only)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		snapshotID := args[0]

		// Get project ID from flag or context
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Get flags
		name, _ := cmd.Flags().GetString("name")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		// At least one field must be provided
		if name == "" && !cmd.Flags().Changed("tags") {
			fmt.Println("Error: at least one of --name or --tags must be provided")
			return
		}

		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// First, get the current snapshot details to preserve existing values
		ctx := context.Background()
		getResponse, err := client.FromStorage().Snapshots().Get(ctx, projectID, snapshotID, nil)
		if err != nil {
			fmt.Printf("Error getting snapshot details: %v\n", err)
			return
		}

		if getResponse == nil || getResponse.Data == nil {
			fmt.Println("Snapshot not found")
			return
		}

		currentSnapshot := getResponse.Data

		// Build the update request with current values as defaults
		volumeURI := ""
		if currentSnapshot.Properties.Volume != nil && currentSnapshot.Properties.Volume.URI != nil {
			volumeURI = *currentSnapshot.Properties.Volume.URI
		}

		updateRequest := types.SnapshotRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: *currentSnapshot.Metadata.Name,
					Tags: currentSnapshot.Metadata.Tags,
				},
				Location: types.LocationRequest{
					Value: currentSnapshot.Metadata.LocationResponse.Code,
				},
			},
			Properties: types.SnapshotPropertiesRequest{
				Volume: types.ReferenceResource{
					URI: volumeURI,
				},
			},
		}

		// Update only the fields that were provided
		if name != "" {
			updateRequest.Metadata.ResourceMetadataRequest.Name = name
		}

		if cmd.Flags().Changed("tags") {
			updateRequest.Metadata.ResourceMetadataRequest.Tags = tags
		}

		// Update the snapshot using the SDK
		response, err := client.FromStorage().Snapshots().Update(ctx, projectID, snapshotID, updateRequest, nil)
		if err != nil {
			fmt.Printf("Error updating snapshot: %v\n", err)
			return
		}

		if response != nil && response.Data != nil {
			fmt.Println("\nSnapshot updated successfully!")
			fmt.Printf("ID:              %s\n", *response.Data.Metadata.ID)
			fmt.Printf("Name:            %s\n", *response.Data.Metadata.Name)
			if len(response.Data.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", response.Data.Metadata.Tags)
			}
		}
	},
}

var snapshotDeleteCmd = &cobra.Command{
	Use:   "delete [snapshot-id]",
	Short: "Delete a snapshot",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		snapshotID := args[0]

		// Get project ID from flag or context
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Get flags
		confirm, _ := cmd.Flags().GetBool("yes")

		// If not confirmed, ask for confirmation
		if !confirm {
			fmt.Printf("Are you sure you want to delete snapshot %s? (yes/no): ", snapshotID)
			var response string
			fmt.Scanln(&response)
			if response != "yes" && response != "y" {
				fmt.Println("Delete cancelled")
				return
			}
		}

		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// Delete the snapshot using the SDK
		ctx := context.Background()
		_, err = client.FromStorage().Snapshots().Delete(ctx, projectID, snapshotID, nil)
		if err != nil {
			fmt.Printf("Error deleting snapshot: %v\n", err)
			return
		}

		fmt.Printf("\nSnapshot %s deleted successfully!\n", snapshotID)
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

		// Get project ID from flag or context
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
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
	blockstorageCreateCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	blockstorageCreateCmd.Flags().String("name", "", "Name for the block storage (required)")
	blockstorageCreateCmd.Flags().String("region", "ITBG-Bergamo", "Region code")
	blockstorageCreateCmd.Flags().String("zone", "", "Zone/datacenter (required)")
	blockstorageCreateCmd.Flags().Int("size", 0, "Size in GB (required)")
	blockstorageCreateCmd.Flags().String("type", "Standard", "Type: Standard or Performance")
	blockstorageCreateCmd.Flags().String("billing-period", "Hour", "Billing period: Hour, Month, Year")
	blockstorageCreateCmd.Flags().StringSlice("tags", []string{}, "Tags (comma-separated)")
	blockstorageCreateCmd.MarkFlagRequired("name")
	blockstorageCreateCmd.MarkFlagRequired("zone")
	blockstorageCreateCmd.MarkFlagRequired("size")

	blockstorageGetCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")

	blockstorageUpdateCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	blockstorageUpdateCmd.Flags().String("name", "", "New name for the block storage")
	blockstorageUpdateCmd.Flags().StringSlice("tags", []string{}, "New tags (comma-separated)")

	blockstorageDeleteCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	blockstorageDeleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	blockstorageListCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	blockstorageListCmd.Flags().BoolP("verbose", "v", false, "Show detailed debug information")

	storageCmd.AddCommand(snapshotCmd)
	snapshotCmd.AddCommand(snapshotCreateCmd)
	snapshotCmd.AddCommand(snapshotGetCmd)
	snapshotCmd.AddCommand(snapshotUpdateCmd)
	snapshotCmd.AddCommand(snapshotDeleteCmd)
	snapshotCmd.AddCommand(snapshotListCmd)

	// Add flags for snapshot commands
	snapshotCreateCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	snapshotCreateCmd.Flags().String("name", "", "Name for the snapshot (required)")
	snapshotCreateCmd.Flags().String("region", "", "Region code (required)")
	snapshotCreateCmd.Flags().String("volume-uri", "", "Source volume URI (required)")
	snapshotCreateCmd.Flags().StringSlice("tags", []string{}, "Tags (comma-separated)")
	snapshotCreateCmd.MarkFlagRequired("name")
	snapshotCreateCmd.MarkFlagRequired("region")
	snapshotCreateCmd.MarkFlagRequired("volume-uri")

	snapshotGetCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")

	snapshotUpdateCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	snapshotUpdateCmd.Flags().String("name", "", "New name for the snapshot")
	snapshotUpdateCmd.Flags().StringSlice("tags", []string{}, "New tags (comma-separated)")

	snapshotDeleteCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	snapshotDeleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	snapshotListCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
}
