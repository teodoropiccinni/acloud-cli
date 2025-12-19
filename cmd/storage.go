package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/types"
	"github.com/spf13/cobra"
)

// Completion functions for storage resources

func completeBlockStorageID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	projectID, err := GetProjectID(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := GetArubaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	ctx := context.Background()
	response, err := client.FromStorage().Volumes().List(ctx, projectID, nil)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string
	if response != nil && response.Data != nil {
		for _, volume := range response.Data.Values {
			if volume.Metadata.ID != nil && volume.Metadata.Name != nil {
				completions = append(completions, fmt.Sprintf("%s\t%s", *volume.Metadata.ID, *volume.Metadata.Name))
			}
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

func completeSnapshotID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	projectID, err := GetProjectID(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := GetArubaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	ctx := context.Background()
	response, err := client.FromStorage().Snapshots().List(ctx, projectID, nil)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string
	if response != nil && response.Data != nil {
		for _, snapshot := range response.Data.Values {
			if snapshot.Metadata.ID != nil && snapshot.Metadata.Name != nil {
				completions = append(completions, fmt.Sprintf("%s\t%s", *snapshot.Metadata.ID, *snapshot.Metadata.Name))
			}
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

func completeBackupID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	projectID, err := GetProjectID(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := GetArubaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	ctx := context.Background()
	response, err := client.FromStorage().Backups().List(ctx, projectID, nil)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string
	if response != nil && response.Data != nil {
		for _, backup := range response.Data.Values {
			if backup.Metadata.ID != nil && backup.Metadata.Name != nil {
				completions = append(completions, fmt.Sprintf("%s\t%s", *backup.Metadata.ID, *backup.Metadata.Name))
			}
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

func completeRestoreID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// For restore, we need backup-id as first arg, restore-id as second
	if len(args) == 0 {
		// First arg is backup-id, use backup completion
		return completeBackupID(cmd, args, toComplete)
	}
	if len(args) > 1 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Second arg is restore-id
	backupID := args[0]
	projectID, err := GetProjectID(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := GetArubaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	ctx := context.Background()
	response, err := client.FromStorage().Restores().List(ctx, projectID, backupID, nil)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string
	if response != nil && response.Data != nil {
		for _, restore := range response.Data.Values {
			if restore.Metadata.ID != nil && restore.Metadata.Name != nil {
				completions = append(completions, fmt.Sprintf("%s\t%s", *restore.Metadata.ID, *restore.Metadata.Name))
			}
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

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
				Type:          types.BlockStorageType(volumeType),
			},
		}

		// Add zone only if provided
		if zone != "" {
			createRequest.Properties.Zone = &zone
		}

		// Debug: print request
		fmt.Printf("\nCreating block storage with:\n")
		fmt.Printf("  Name: %s\n", name)
		fmt.Printf("  Region: %s\n", region)
		if zone != "" {
			fmt.Printf("  Zone: %s\n", zone)
		}
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

			if volume.Metadata.URI != nil {
				fmt.Printf("URI:             %s\n", *volume.Metadata.URI)
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
	Short: "Update block storage (name and/or tags)",
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

		// Check if the volume status allows updates
		if currentVolume.Status.State != nil {
			status := *currentVolume.Status.State
			if status != "Used" && status != "NotUsed" {
				fmt.Printf("Error: Cannot update block storage with status '%s'\n", status)
				fmt.Println("Block storage can only be updated when status is 'Used' or 'NotUsed'")
				return
			}
		}

		// Fix region code format (IT BG -> ITBG-Bergamo)
		regionCode := currentVolume.Metadata.LocationResponse.Code
		if regionCode == "IT BG" {
			regionCode = "ITBG-Bergamo"
		}

		// Handle zone - if empty, set to nil
		var zone *string
		if currentVolume.Properties.Zone != "" {
			zone = &currentVolume.Properties.Zone
		}

		// Build the update request with current values as defaults
		updateRequest := types.BlockStorageRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: *currentVolume.Metadata.Name,
					Tags: currentVolume.Metadata.Tags,
				},
				Location: types.LocationRequest{
					Value: regionCode,
				},
			},
			Properties: types.BlockStoragePropertiesRequest{
				SizeGB:        currentVolume.Properties.SizeGB,
				BillingPeriod: "Hour",
				Zone:          zone,
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

		// Check if the response indicates an error
		if response.StatusCode >= 400 {
			fmt.Printf("API Error (Status %d):\n", response.StatusCode)
			if response.Error != nil {
				if response.Error.Title != nil {
					fmt.Printf("  Title: %s\n", *response.Error.Title)
				}
				if response.Error.Detail != nil {
					fmt.Printf("  Detail: %s\n", *response.Error.Detail)
				}
				if response.Error.Extensions != nil {
					fmt.Printf("  Extensions: %+v\n", response.Error.Extensions)
				}
			}
			// Decode RawBody to see the actual error
			if len(response.RawBody) > 0 {
				fmt.Printf("  Raw Response: %s\n", string(response.RawBody))
			}
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
		} else {
			fmt.Println("Warning: Update may have succeeded but response is empty")
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
			fmt.Println("\n=== End Debug ===")
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

		fmt.Println("Creating snapshot with the following parameters:")
		fmt.Printf("  Name:       %s\n", name)
		fmt.Printf("  Region:     %s\n", region)
		fmt.Printf("  Volume URI: %s\n", volumeURI)
		fmt.Printf("  Tags:       %v\n", tags)
		fmt.Println()

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
			if !response.Data.Metadata.CreationDate.IsZero() {
				fmt.Printf("Creation Date:   %s\n", response.Data.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}
		} else {
			fmt.Println("Warning: Snapshot may have been created but response is empty")
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

			if snapshot.Metadata.URI != nil {
				fmt.Printf("URI:             %s\n", *snapshot.Metadata.URI)
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

			status := ""
			if snapshot.Status.State != nil {
				status = *snapshot.Status.State
			}
			fmt.Printf("Status:          %s\n", status)

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

		// Fix region code format (IT BG -> ITBG-Bergamo)
		regionCode := currentSnapshot.Metadata.LocationResponse.Code
		if regionCode == "IT BG" {
			regionCode = "ITBG-Bergamo"
		}

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
					Value: regionCode,
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

		// Check if the response indicates an error
		if response.StatusCode >= 400 {
			fmt.Printf("API Error (Status %d):\n", response.StatusCode)
			if response.Error != nil {
				if response.Error.Title != nil {
					fmt.Printf("  Title: %s\n", *response.Error.Title)
				}
				if response.Error.Detail != nil {
					fmt.Printf("  Detail: %s\n", *response.Error.Detail)
				}
				if response.Error.Extensions != nil {
					fmt.Printf("  Extensions: %+v\n", response.Error.Extensions)
				}
			}
			// Decode RawBody to see the actual error
			if len(response.RawBody) > 0 {
				fmt.Printf("  Raw Response: %s\n", string(response.RawBody))
			}
			return
		}

		if response != nil && response.Data != nil {
			fmt.Println("\nSnapshot updated successfully!")
			fmt.Printf("ID:              %s\n", *response.Data.Metadata.ID)
			fmt.Printf("Name:            %s\n", *response.Data.Metadata.Name)
			if len(response.Data.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", response.Data.Metadata.Tags)
			}
		} else {
			fmt.Println("Warning: Update may have succeeded but response is empty")
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
	Short: "List snapshots for a block storage volume",
	Run: func(cmd *cobra.Command, args []string) {
		// Get project ID from flag or context
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Get block storage URI flag
		volumeURI, _ := cmd.Flags().GetString("volume-uri")
		if volumeURI == "" {
			fmt.Println("Error: --volume-uri is required")
			fmt.Println("Use: acloud storage snapshot list --volume-uri <block-storage-uri>")
			return
		}

		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// List snapshots using the SDK (filter by volume URI on client side)
		ctx := context.Background()
		response, err := client.FromStorage().Snapshots().List(ctx, projectID, nil)
		if err != nil {
			fmt.Printf("Error listing snapshots: %v\n", err)
			return
		}

		// Check verbose flag
		verbose, _ := cmd.Flags().GetBool("verbose")
		if verbose {
			fmt.Println("=== Full API Response ===")
			fmt.Printf("%+v\n\n", response)
		}

		if response != nil && response.Data != nil && len(response.Data.Values) > 0 {
			// Filter snapshots by volume URI
			var filteredSnapshots []types.SnapshotResponse
			for _, snapshot := range response.Data.Values {
				if snapshot.Properties.Volume != nil && snapshot.Properties.Volume.URI != nil && *snapshot.Properties.Volume.URI == volumeURI {
					filteredSnapshots = append(filteredSnapshots, snapshot)
				}
			}

			if len(filteredSnapshots) == 0 {
				fmt.Printf("No snapshots found for volume: %s\n", volumeURI)
				return
			}

			// Define table columns
			headers := []TableColumn{
				{Header: "NAME", Width: 30},
				{Header: "ID", Width: 26},
				{Header: "SIZE(GB)", Width: 12},
				{Header: "STATUS", Width: 15},
			}

			// Build rows
			var rows [][]string
			for _, snapshot := range filteredSnapshots {
				name := ""
				if snapshot.Metadata.Name != nil && *snapshot.Metadata.Name != "" {
					name = *snapshot.Metadata.Name
				}

				id := ""
				if snapshot.Metadata.ID != nil && *snapshot.Metadata.ID != "" {
					id = *snapshot.Metadata.ID
				}

				size := fmt.Sprintf("%d", snapshot.Properties.SizeGB)

				status := ""
				if snapshot.Status.State != nil {
					status = *snapshot.Status.State
				}

				rows = append(rows, []string{name, id, size, status})
			}

			// Print the table
			PrintTable(headers, rows)
		} else {
			fmt.Println("No snapshots found")
		}
	},
}

// Backup command - creates an Aruba.Storage/backup resource
var storageBackupCmd = &cobra.Command{
	Use:   "backup [volume-id]",
	Short: "Create a storage backup of a block storage volume",
	Long:  `Create a storage backup resource (Aruba.Storage/backup) for a block storage volume.`,
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
		region, _ := cmd.Flags().GetString("region")
		backupType, _ := cmd.Flags().GetString("type")
		retentionDays, _ := cmd.Flags().GetInt("retention-days")
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

		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// First, get the volume details to get the full URI
		ctx := context.Background()
		volumeResponse, err := client.FromStorage().Volumes().Get(ctx, projectID, volumeID, nil)
		if err != nil {
			fmt.Printf("Error getting volume details: %v\n", err)
			return
		}

		if volumeResponse == nil || volumeResponse.Data == nil {
			fmt.Println("Volume not found")
			return
		}

		volumeURI := *volumeResponse.Data.Metadata.URI

		// Build the backup create request
		createRequest := types.StorageBackupRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: name,
					Tags: tags,
				},
				Location: types.LocationRequest{
					Value: region,
				},
			},
			Properties: types.StorageBackupPropertiesRequest{
				StorageBackupType: types.StorageBackupType(backupType),
				Origin: types.ReferenceResource{
					URI: volumeURI,
				},
			},
		}

		// Add optional fields
		if retentionDays > 0 {
			createRequest.Properties.RetentionDays = &retentionDays
		}
		if billingPeriod != "" {
			createRequest.Properties.BillingPeriod = &billingPeriod
		}

		fmt.Println("Creating storage backup with the following parameters:")
		fmt.Printf("  Name:           %s\n", name)
		fmt.Printf("  Type:           %s\n", backupType)
		fmt.Printf("  Region:         %s\n", region)
		fmt.Printf("  Volume ID:      %s\n", volumeID)
		fmt.Printf("  Volume URI:     %s\n", volumeURI)
		if retentionDays > 0 {
			fmt.Printf("  Retention Days: %d\n", retentionDays)
		}
		if billingPeriod != "" {
			fmt.Printf("  Billing Period: %s\n", billingPeriod)
		}
		if len(tags) > 0 {
			fmt.Printf("  Tags:           %v\n", tags)
		}
		fmt.Println()

		// Create the backup using the SDK
		response, err := client.FromStorage().Backups().Create(ctx, projectID, createRequest, nil)
		if err != nil {
			fmt.Printf("Error creating backup: %v\n", err)
			return
		}

		if response == nil {
			fmt.Println("Error: received nil response from API")
			return
		}

		if !response.IsSuccess() {
			fmt.Printf("Error: failed to create backup - Status: %d\n", response.StatusCode)
			if response.Error.Title != nil {
				fmt.Printf("Error Title: %s\n", *response.Error.Title)
			}
			if response.Error.Detail != nil {
				fmt.Printf("Error Detail: %s\n", *response.Error.Detail)
			}
			// Try to decode RawBody for more details
			if response.StatusCode >= 400 && response.RawBody != nil {
				var errorDetail map[string]interface{}
				if err := json.Unmarshal(response.RawBody, &errorDetail); err == nil {
					fmt.Printf("Full Error Response: %+v\n", errorDetail)
				}
			}
			return
		}

		if response.Data != nil {
			fmt.Println("Storage backup created successfully!")
			fmt.Printf("ID:              %s\n", *response.Data.Metadata.ID)
			fmt.Printf("Name:            %s\n", *response.Data.Metadata.Name)
			fmt.Printf("Type:            %s\n", response.Data.Properties.Type)
			if !response.Data.Metadata.CreationDate.IsZero() {
				fmt.Printf("Creation Date:   %s\n", response.Data.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}
		}
	},
}

// Backup List command
var storageBackupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List storage backups",
	Run: func(cmd *cobra.Command, args []string) {
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		ctx := context.Background()
		response, err := client.FromStorage().Backups().List(ctx, projectID, nil)
		if err != nil {
			fmt.Printf("Error listing backups: %v\n", err)
			return
		}

		if response != nil && response.Data != nil && len(response.Data.Values) > 0 {
			headers := []TableColumn{
				{Header: "NAME", Width: 30},
				{Header: "ID", Width: 26},
				{Header: "TYPE", Width: 12},
				{Header: "STATUS", Width: 15},
			}

			var rows [][]string
			for _, backup := range response.Data.Values {
				name := ""
				if backup.Metadata.Name != nil {
					name = *backup.Metadata.Name
				}

				id := ""
				if backup.Metadata.ID != nil {
					id = *backup.Metadata.ID
				}

				backupType := string(backup.Properties.Type)

				status := ""
				if backup.Status.State != nil {
					status = *backup.Status.State
				}

				rows = append(rows, []string{name, id, backupType, status})
			}

			PrintTable(headers, rows)
		} else {
			fmt.Println("No backups found")
		}
	},
}

// Backup Get command
var storageBackupGetCmd = &cobra.Command{
	Use:   "get [backup-id]",
	Short: "Get storage backup details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		backupID := args[0]

		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		ctx := context.Background()
		response, err := client.FromStorage().Backups().Get(ctx, projectID, backupID, nil)
		if err != nil {
			fmt.Printf("Error getting backup details: %v\n", err)
			return
		}

		if response != nil && response.Data != nil {
			backup := response.Data

			fmt.Println("\nStorage Backup Details:")
			fmt.Println("=======================")

			if backup.Metadata.ID != nil {
				fmt.Printf("ID:              %s\n", *backup.Metadata.ID)
			}
			if backup.Metadata.URI != nil {
				fmt.Printf("URI:             %s\n", *backup.Metadata.URI)
			}
			if backup.Metadata.Name != nil {
				fmt.Printf("Name:            %s\n", *backup.Metadata.Name)
			}

			fmt.Printf("Type:            %s\n", backup.Properties.Type)

			if backup.Properties.Origin.URI != "" {
				fmt.Printf("Source Volume:   %s\n", backup.Properties.Origin.URI)
			}

			if backup.Properties.RetentionDays != nil {
				fmt.Printf("Retention Days:  %d\n", *backup.Properties.RetentionDays)
			}

			if backup.Properties.BillingPeriod != nil {
				fmt.Printf("Billing Period:  %s\n", *backup.Properties.BillingPeriod)
			}

			fmt.Printf("Region:          %s\n", backup.Metadata.LocationResponse.Code)

			if backup.Status.State != nil {
				fmt.Printf("Status:          %s\n", *backup.Status.State)
			}

			if !backup.Metadata.CreationDate.IsZero() {
				fmt.Printf("Creation Date:   %s\n", backup.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}

			if backup.Metadata.CreatedBy != nil {
				fmt.Printf("Created By:      %s\n", *backup.Metadata.CreatedBy)
			}

			if len(backup.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", backup.Metadata.Tags)
			}
		} else {
			fmt.Println("Backup not found")
		}
	},
}

// Backup Update command
var storageBackupUpdateCmd = &cobra.Command{
	Use:   "update [backup-id]",
	Short: "Update a storage backup (name and/or tags)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		backupID := args[0]

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

		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// First, get the current backup details
		ctx := context.Background()
		getResponse, err := client.FromStorage().Backups().Get(ctx, projectID, backupID, nil)
		if err != nil {
			fmt.Printf("Error getting backup details: %v\n", err)
			return
		}

		if getResponse == nil || getResponse.Data == nil {
			fmt.Println("Backup not found")
			return
		}

		currentBackup := getResponse.Data

		// Fix region code format (IT BG -> ITBG-Bergamo)
		regionCode := currentBackup.Metadata.LocationResponse.Code
		if regionCode == "IT BG" {
			regionCode = "ITBG-Bergamo"
		}

		// Build the update request with current values as defaults
		updateRequest := types.StorageBackupRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: *currentBackup.Metadata.Name,
					Tags: currentBackup.Metadata.Tags,
				},
				Location: types.LocationRequest{
					Value: regionCode,
				},
			},
			Properties: types.StorageBackupPropertiesRequest{
				StorageBackupType: currentBackup.Properties.Type,
				Origin: types.ReferenceResource{
					URI: currentBackup.Properties.Origin.URI,
				},
			},
		}

		// Add optional fields if present
		if currentBackup.Properties.RetentionDays != nil {
			updateRequest.Properties.RetentionDays = currentBackup.Properties.RetentionDays
		}
		if currentBackup.Properties.BillingPeriod != nil {
			updateRequest.Properties.BillingPeriod = currentBackup.Properties.BillingPeriod
		}

		// Update only the fields that were provided
		if name != "" {
			updateRequest.Metadata.ResourceMetadataRequest.Name = name
		}

		if cmd.Flags().Changed("tags") {
			updateRequest.Metadata.ResourceMetadataRequest.Tags = tags
		}

		// Update the backup using the SDK
		response, err := client.FromStorage().Backups().Update(ctx, projectID, backupID, updateRequest, nil)
		if err != nil {
			fmt.Printf("Error updating backup: %v\n", err)
			return
		}

		if response.StatusCode >= 400 {
			fmt.Printf("API Error (Status %d):\n", response.StatusCode)
			if response.Error != nil {
				if response.Error.Title != nil {
					fmt.Printf("  Title: %s\n", *response.Error.Title)
				}
				if response.Error.Detail != nil {
					fmt.Printf("  Detail: %s\n", *response.Error.Detail)
				}
			}
			if len(response.RawBody) > 0 {
				fmt.Printf("  Raw Response: %s\n", string(response.RawBody))
			}
			return
		}

		if response != nil && response.Data != nil {
			fmt.Println("\nBackup updated successfully!")
			fmt.Printf("ID:              %s\n", *response.Data.Metadata.ID)
			fmt.Printf("Name:            %s\n", *response.Data.Metadata.Name)
			if len(response.Data.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", response.Data.Metadata.Tags)
			}
		} else {
			fmt.Println("Warning: Update may have succeeded but response is empty")
		}
	},
}

// Backup Delete command
var storageBackupDeleteCmd = &cobra.Command{
	Use:   "delete [backup-id]",
	Short: "Delete a storage backup",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		backupID := args[0]

		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		confirm, _ := cmd.Flags().GetBool("yes")
		if !confirm {
			fmt.Printf("Are you sure you want to delete backup %s? (yes/no): ", backupID)
			var response string
			fmt.Scanln(&response)
			if response != "yes" && response != "y" {
				fmt.Println("Delete cancelled")
				return
			}
		}

		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		ctx := context.Background()
		_, err = client.FromStorage().Backups().Delete(ctx, projectID, backupID, nil)
		if err != nil {
			fmt.Printf("Error deleting backup: %v\n", err)
			return
		}

		fmt.Printf("\nBackup %s deleted successfully!\n", backupID)
	},
}

// Restore command - creates an Aruba.Storage/restore resource
var storageRestoreCmd = &cobra.Command{
	Use:   "restore [backup-id] [volume-id]",
	Short: "Restore a block storage volume from a backup",
	Long:  `Create a restore operation (Aruba.Storage/restore) to restore a volume from a backup.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		backupID := args[0]
		volumeID := args[1]

		// Get project ID from flag or context
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Get flags
		name, _ := cmd.Flags().GetString("name")
		region, _ := cmd.Flags().GetString("region")
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

		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		ctx := context.Background()

		// Get the backup details
		backupResponse, err := client.FromStorage().Backups().Get(ctx, projectID, backupID, nil)
		if err != nil {
			fmt.Printf("Error getting backup details: %v\n", err)
			return
		}

		if backupResponse == nil || backupResponse.Data == nil {
			fmt.Println("Backup not found")
			return
		}

		backupURI := *backupResponse.Data.Metadata.URI

		// Get the volume details
		volumeResponse, err := client.FromStorage().Volumes().Get(ctx, projectID, volumeID, nil)
		if err != nil {
			fmt.Printf("Error getting volume details: %v\n", err)
			return
		}

		if volumeResponse == nil || volumeResponse.Data == nil {
			fmt.Println("Volume not found")
			return
		}

		volumeURI := *volumeResponse.Data.Metadata.URI

		// Build the restore request
		createRequest := types.RestoreRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: name,
					Tags: tags,
				},
				Location: types.LocationRequest{
					Value: region,
				},
			},
			Properties: types.RestorePropertiesRequest{
				Target: types.ReferenceResource{
					URI: volumeURI,
				},
			},
		}

		fmt.Println("Creating restore operation with the following parameters:")
		fmt.Printf("  Name:       %s\n", name)
		fmt.Printf("  Region:     %s\n", region)
		fmt.Printf("  Backup ID:  %s\n", backupID)
		fmt.Printf("  Backup URI: %s\n", backupURI)
		fmt.Printf("  Volume ID:  %s\n", volumeID)
		fmt.Printf("  Volume URI: %s\n", volumeURI)
		if len(tags) > 0 {
			fmt.Printf("  Tags:       %v\n", tags)
		}
		fmt.Println()

		// Create the restore using the SDK
		response, err := client.FromStorage().Restores().Create(ctx, projectID, backupID, createRequest, nil)
		if err != nil {
			fmt.Printf("Error creating restore: %v\n", err)
			return
		}

		if response == nil {
			fmt.Println("Error: received nil response from API")
			return
		}

		if !response.IsSuccess() {
			fmt.Printf("Error: failed to create restore - Status: %d\n", response.StatusCode)
			if response.Error.Title != nil {
				fmt.Printf("Error Title: %s\n", *response.Error.Title)
			}
			if response.Error.Detail != nil {
				fmt.Printf("Error Detail: %s\n", *response.Error.Detail)
			}
			// Try to decode RawBody for more details
			if response.StatusCode >= 400 && response.RawBody != nil {
				var errorDetail map[string]interface{}
				if err := json.Unmarshal(response.RawBody, &errorDetail); err == nil {
					fmt.Printf("Full Error Response: %+v\n", errorDetail)
				}
			}
			return
		}

		if response.Data != nil {
			fmt.Println("Restore operation created successfully!")
			fmt.Printf("ID:              %s\n", *response.Data.Metadata.ID)
			fmt.Printf("Name:            %s\n", *response.Data.Metadata.Name)
			if !response.Data.Metadata.CreationDate.IsZero() {
				fmt.Printf("Creation Date:   %s\n", response.Data.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}
			if response.Data.Status.State != nil {
				fmt.Printf("Status:          %s\n", *response.Data.Status.State)
			}
		}
	},
}

// Restore List command
var storageRestoreListCmd = &cobra.Command{
	Use:   "list [backup-id]",
	Short: "List restore operations for a backup",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		backupID := args[0]

		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		ctx := context.Background()
		response, err := client.FromStorage().Restores().List(ctx, projectID, backupID, nil)
		if err != nil {
			fmt.Printf("Error listing restores: %v\n", err)
			return
		}

		if response != nil && response.Data != nil && len(response.Data.Values) > 0 {
			headers := []TableColumn{
				{Header: "NAME", Width: 30},
				{Header: "ID", Width: 26},
				{Header: "STATUS", Width: 15},
			}

			var rows [][]string
			for _, restore := range response.Data.Values {
				name := ""
				if restore.Metadata.Name != nil {
					name = *restore.Metadata.Name
				}

				id := ""
				if restore.Metadata.ID != nil {
					id = *restore.Metadata.ID
				}

				status := ""
				if restore.Status.State != nil {
					status = *restore.Status.State
				}

				rows = append(rows, []string{name, id, status})
			}

			PrintTable(headers, rows)
		} else {
			fmt.Println("No restores found for this backup")
		}
	},
}

// Restore Get command
var storageRestoreGetCmd = &cobra.Command{
	Use:   "get [backup-id] [restore-id]",
	Short: "Get restore operation details",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		backupID := args[0]
		restoreID := args[1]

		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		ctx := context.Background()
		response, err := client.FromStorage().Restores().Get(ctx, projectID, backupID, restoreID, nil)
		if err != nil {
			fmt.Printf("Error getting restore details: %v\n", err)
			return
		}

		if response != nil && response.Data != nil {
			restore := response.Data

			fmt.Println("\nRestore Operation Details:")
			fmt.Println("==========================")

			if restore.Metadata.ID != nil {
				fmt.Printf("ID:              %s\n", *restore.Metadata.ID)
			}
			if restore.Metadata.URI != nil {
				fmt.Printf("URI:             %s\n", *restore.Metadata.URI)
			}
			if restore.Metadata.Name != nil {
				fmt.Printf("Name:            %s\n", *restore.Metadata.Name)
			}

			if restore.Properties.Destination.URI != "" {
				fmt.Printf("Target Volume:   %s\n", restore.Properties.Destination.URI)
			}

			fmt.Printf("Region:          %s\n", restore.Metadata.LocationResponse.Code)

			if restore.Status.State != nil {
				fmt.Printf("Status:          %s\n", *restore.Status.State)
			}

			if !restore.Metadata.CreationDate.IsZero() {
				fmt.Printf("Creation Date:   %s\n", restore.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}

			if restore.Metadata.CreatedBy != nil {
				fmt.Printf("Created By:      %s\n", *restore.Metadata.CreatedBy)
			}

			if len(restore.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", restore.Metadata.Tags)
			}
		} else {
			fmt.Println("Restore operation not found")
		}
	},
}

// Restore Update command
var storageRestoreUpdateCmd = &cobra.Command{
	Use:   "update [backup-id] [restore-id]",
	Short: "Update a restore operation (name and/or tags)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		backupID := args[0]
		restoreID := args[1]

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

		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// First, get the current restore details
		ctx := context.Background()
		getResponse, err := client.FromStorage().Restores().Get(ctx, projectID, backupID, restoreID, nil)
		if err != nil {
			fmt.Printf("Error getting restore details: %v\n", err)
			return
		}

		if getResponse == nil || getResponse.Data == nil {
			fmt.Println("Restore operation not found")
			return
		}

		currentRestore := getResponse.Data

		// Fix region code format (IT BG -> ITBG-Bergamo)
		regionCode := currentRestore.Metadata.LocationResponse.Code
		if regionCode == "IT BG" {
			regionCode = "ITBG-Bergamo"
		}

		// Build the update request with current values as defaults
		updateRequest := types.RestoreRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: *currentRestore.Metadata.Name,
					Tags: currentRestore.Metadata.Tags,
				},
				Location: types.LocationRequest{
					Value: regionCode,
				},
			},
			Properties: types.RestorePropertiesRequest{
				Target: types.ReferenceResource{
					URI: currentRestore.Properties.Destination.URI,
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

		// Update the restore using the SDK
		response, err := client.FromStorage().Restores().Update(ctx, projectID, backupID, restoreID, updateRequest, nil)
		if err != nil {
			fmt.Printf("Error updating restore: %v\n", err)
			return
		}

		if response.StatusCode >= 400 {
			fmt.Printf("API Error (Status %d):\n", response.StatusCode)
			if response.Error != nil {
				if response.Error.Title != nil {
					fmt.Printf("  Title: %s\n", *response.Error.Title)
				}
				if response.Error.Detail != nil {
					fmt.Printf("  Detail: %s\n", *response.Error.Detail)
				}
			}
			if len(response.RawBody) > 0 {
				fmt.Printf("  Raw Response: %s\n", string(response.RawBody))
			}
			return
		}

		if response != nil && response.Data != nil {
			fmt.Println("\nRestore operation updated successfully!")
			fmt.Printf("ID:              %s\n", *response.Data.Metadata.ID)
			fmt.Printf("Name:            %s\n", *response.Data.Metadata.Name)
			if len(response.Data.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", response.Data.Metadata.Tags)
			}
		} else {
			fmt.Println("Warning: Update may have succeeded but response is empty")
		}
	},
}

// Restore Delete command
var storageRestoreDeleteCmd = &cobra.Command{
	Use:   "delete [backup-id] [restore-id]",
	Short: "Delete a restore operation",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		backupID := args[0]
		restoreID := args[1]

		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		confirm, _ := cmd.Flags().GetBool("yes")
		if !confirm {
			fmt.Printf("Are you sure you want to delete restore operation %s? (yes/no): ", restoreID)
			var response string
			fmt.Scanln(&response)
			if response != "yes" && response != "y" {
				fmt.Println("Delete cancelled")
				return
			}
		}

		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		ctx := context.Background()
		_, err = client.FromStorage().Restores().Delete(ctx, projectID, backupID, restoreID, nil)
		if err != nil {
			fmt.Printf("Error deleting restore: %v\n", err)
			return
		}

		fmt.Printf("\nRestore operation %s deleted successfully!\n", restoreID)
	},
}

func init() {
	rootCmd.AddCommand(storageCmd)

	// Backup commands
	storageCmd.AddCommand(storageBackupCmd)
	storageBackupCmd.AddCommand(storageBackupListCmd)
	storageBackupCmd.AddCommand(storageBackupGetCmd)
	storageBackupCmd.AddCommand(storageBackupUpdateCmd)
	storageBackupCmd.AddCommand(storageBackupDeleteCmd)

	// Restore commands
	storageCmd.AddCommand(storageRestoreCmd)
	storageRestoreCmd.AddCommand(storageRestoreListCmd)
	storageRestoreCmd.AddCommand(storageRestoreGetCmd)
	storageRestoreCmd.AddCommand(storageRestoreUpdateCmd)
	storageRestoreCmd.AddCommand(storageRestoreDeleteCmd)

	// Block storage commands
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
	blockstorageCreateCmd.Flags().String("zone", "", "Zone/datacenter (optional)")
	blockstorageCreateCmd.Flags().Int("size", 0, "Size in GB (required)")
	blockstorageCreateCmd.Flags().String("type", "Standard", "Type: Standard or Performance")
	blockstorageCreateCmd.Flags().String("billing-period", "Hour", "Billing period: Hour, Month, Year")
	blockstorageCreateCmd.Flags().StringSlice("tags", []string{}, "Tags (comma-separated)")
	blockstorageCreateCmd.MarkFlagRequired("name")
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
	snapshotListCmd.Flags().String("volume-uri", "", "Block storage volume URI (required)")
	snapshotListCmd.Flags().BoolP("verbose", "v", false, "Show detailed debug information")
	snapshotListCmd.MarkFlagRequired("volume-uri")

	// Add flags for backup command
	storageBackupCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	storageBackupCmd.Flags().String("name", "", "Name for the backup (required)")
	storageBackupCmd.Flags().String("region", "ITBG-Bergamo", "Region code")
	storageBackupCmd.Flags().String("type", "Full", "Backup type: Full or Incremental")
	storageBackupCmd.Flags().Int("retention-days", 0, "Number of days to retain the backup")
	storageBackupCmd.Flags().String("billing-period", "", "Billing period: Hour, Month, Year")
	storageBackupCmd.Flags().StringSlice("tags", []string{}, "Tags (comma-separated)")
	storageBackupCmd.MarkFlagRequired("name")

	storageBackupListCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	storageBackupGetCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	storageBackupUpdateCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	storageBackupUpdateCmd.Flags().String("name", "", "New name for the backup")
	storageBackupUpdateCmd.Flags().StringSlice("tags", []string{}, "New tags (comma-separated)")
	storageBackupDeleteCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	storageBackupDeleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	// Add flags for restore command
	storageRestoreCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	storageRestoreCmd.Flags().String("name", "", "Name for the restore operation (required)")
	storageRestoreCmd.Flags().String("region", "ITBG-Bergamo", "Region code")
	storageRestoreCmd.Flags().StringSlice("tags", []string{}, "Tags (comma-separated)")
	storageRestoreCmd.MarkFlagRequired("name")

	storageRestoreListCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	storageRestoreGetCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	storageRestoreUpdateCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	storageRestoreUpdateCmd.Flags().String("name", "", "New name for the restore operation")
	storageRestoreUpdateCmd.Flags().StringSlice("tags", []string{}, "New tags (comma-separated)")
	storageRestoreDeleteCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	storageRestoreDeleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	// Set up auto-completion for resource IDs
	blockstorageGetCmd.ValidArgsFunction = completeBlockStorageID
	blockstorageUpdateCmd.ValidArgsFunction = completeBlockStorageID
	blockstorageDeleteCmd.ValidArgsFunction = completeBlockStorageID

	snapshotGetCmd.ValidArgsFunction = completeSnapshotID
	snapshotUpdateCmd.ValidArgsFunction = completeSnapshotID
	snapshotDeleteCmd.ValidArgsFunction = completeSnapshotID

	storageBackupGetCmd.ValidArgsFunction = completeBackupID
	storageBackupUpdateCmd.ValidArgsFunction = completeBackupID
	storageBackupDeleteCmd.ValidArgsFunction = completeBackupID

	storageRestoreGetCmd.ValidArgsFunction = completeRestoreID
	storageRestoreUpdateCmd.ValidArgsFunction = completeRestoreID
	storageRestoreDeleteCmd.ValidArgsFunction = completeRestoreID
	storageRestoreListCmd.ValidArgsFunction = completeBackupID
}
