package cmd

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/types"
	"github.com/spf13/cobra"
)

// INIT
func init() {
	// Subnet
	networkCmd.AddCommand(subnetCmd)
	subnetCmd.AddCommand(subnetCreateCmd)
	subnetCmd.AddCommand(subnetGetCmd)
	subnetCmd.AddCommand(subnetUpdateCmd)
	subnetCmd.AddCommand(subnetDeleteCmd)
	subnetCmd.AddCommand(subnetListCmd)

	subnetCreateCmd.Flags().String("name", "", "Subnet name (required)")
	subnetCreateCmd.Flags().String("cidr", "", "Subnet CIDR (required)")
	subnetCreateCmd.Flags().String("region", "", "Region for the subnet (required)")
	subnetCreateCmd.Flags().StringSlice("tags", []string{}, "Subnet tags (optional)")
	subnetUpdateCmd.Flags().String("name", "", "Subnet name (optional)")
	subnetUpdateCmd.Flags().String("cidr", "", "Subnet CIDR (optional)")
	subnetUpdateCmd.Flags().StringSlice("tags", []string{}, "Subnet tags (optional)")
	subnetListCmd.Flags().String("vpc-id", "", "Parent VPC ID (required)")
}

// Subnet subcommands
var subnetCmd = &cobra.Command{
	Use:   "subnet",
	Short: "Manage subnets",
	Long:  `Perform CRUD operations on subnets in Aruba Cloud.`,
}

var subnetCreateCmd = &cobra.Command{
	Use:   "create [vpc-id]",
	Short: "Create a new subnet",
	Args:  cobra.ExactArgs(1),
	Long:  `Create a new subnet in a VPC. Usage: acloud network subnet create <vpc-id> --name <name> --cidr <cidr>`,
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]
		name, _ := cmd.Flags().GetString("name")
		cidr, _ := cmd.Flags().GetString("cidr")
		region, _ := cmd.Flags().GetString("region")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		if name == "" || cidr == "" || region == "" {
			fmt.Println("Error: --name, --cidr, and --region are required")
			return
		}
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
		req := types.SubnetRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: name,
					Tags: tags,
				},
				Location: types.LocationRequest{
					Value: region,
				},
			},
			Properties: types.SubnetPropertiesRequest{
				Network: &types.SubnetNetwork{
					Address: cidr,
				},
			},
		}
		resp, err := client.FromNetwork().Subnets().Create(ctx, projectID, vpcID, req, nil)
		if err != nil {
			fmt.Printf("Error creating subnet: %v\n", err)
			return
		}
		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to create subnet - Status: %d\n", resp.StatusCode)
			if resp.Error.Title != nil {
				fmt.Printf("Error: %s\n", *resp.Error.Title)
			}
			if resp.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *resp.Error.Detail)
			}
			return
		}
		if resp != nil && resp.Data != nil && resp.Data.Metadata.ID != nil {
			headers := []TableColumn{
				{Header: "NAME", Width: 30},
				{Header: "ID", Width: 26},
				{Header: "REGION", Width: 18},
				{Header: "CIDR", Width: 18},
				{Header: "STATUS", Width: 15},
			}
			row := []string{
				name,
				*resp.Data.Metadata.ID,
				resp.Data.Metadata.LocationResponse.Value,
				cidr,
				func() string {
					if resp.Data.Status.State != nil {
						return *resp.Data.Status.State
					} else {
						return ""
					}
				}(),
			}
			PrintTable(headers, [][]string{row})
		} else {
			fmt.Println("Subnet created, but no ID returned.")
		}
	},
}

var subnetGetCmd = &cobra.Command{
	Use:   "get [vpc-id] [subnet-id]",
	Short: "Get subnet details",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]
		subnetID := args[1]
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
		resp, err := client.FromNetwork().Subnets().Get(ctx, projectID, vpcID, subnetID, nil)
		if err != nil {
			fmt.Printf("Error getting subnet: %v\n", err)
			return
		}
		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to get subnet - Status: %d\n", resp.StatusCode)
			if resp.Error.Title != nil {
				fmt.Printf("Error: %s\n", *resp.Error.Title)
			}
			if resp.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *resp.Error.Detail)
			}
			return
		}
		if resp != nil && resp.Data != nil {
			subnet := resp.Data
			fmt.Println("\nSubnet Details:")
			fmt.Println("===============")
			if subnet.Metadata.ID != nil {
				fmt.Printf("ID:              %s\n", *subnet.Metadata.ID)
			}
			if subnet.Metadata.URI != nil {
				fmt.Printf("URI:             %s\n", *subnet.Metadata.URI)
			}
			if subnet.Metadata.Name != nil {
				fmt.Printf("Name:            %s\n", *subnet.Metadata.Name)
			}
			if subnet.Metadata.LocationResponse != nil && subnet.Metadata.LocationResponse.Value != "" {
				fmt.Printf("Region:          %s\n", subnet.Metadata.LocationResponse.Value)
			}
			if subnet.Properties.Network != nil {
				fmt.Printf("CIDR:            %s\n", subnet.Properties.Network.Address)
			}
			if subnet.Metadata.CreationDate != nil {
				fmt.Printf("Creation Date:   %s\n", subnet.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}
			if subnet.Metadata.CreatedBy != nil {
				fmt.Printf("Created By:      %s\n", *subnet.Metadata.CreatedBy)
			}
			if len(subnet.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", subnet.Metadata.Tags)
			} else {
				fmt.Printf("Tags:            []\n")
			}
			if subnet.Status.State != nil {
				fmt.Printf("Status:          %s\n", *subnet.Status.State)
			}
		} else {
			fmt.Println("Subnet not found or no data returned.")
		}
	},
}

var subnetListCmd = &cobra.Command{
	Use:   "list [vpc-id]",
	Short: "List subnets for a VPC",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]
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
		resp, err := client.FromNetwork().Subnets().List(ctx, projectID, vpcID, nil)
		if err != nil {
			fmt.Printf("Error listing subnets: %v\n", err)
			return
		}
		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to list subnets - Status: %d\n", resp.StatusCode)
			if resp.Error.Title != nil {
				fmt.Printf("Error: %s\n", *resp.Error.Title)
			}
			if resp.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *resp.Error.Detail)
			}
			return
		}
		if resp != nil && resp.Data != nil && len(resp.Data.Values) > 0 {
			headers := []TableColumn{
				{Header: "NAME", Width: 30},
				{Header: "ID", Width: 26},
				{Header: "REGION", Width: 18},
				{Header: "CIDR", Width: 18},
				{Header: "STATUS", Width: 15},
			}
			var rows [][]string
			for _, subnet := range resp.Data.Values {
				name := ""
				if subnet.Metadata.Name != nil {
					name = *subnet.Metadata.Name
				}
				id := ""
				if subnet.Metadata.ID != nil {
					id = *subnet.Metadata.ID
				}
				region := subnet.Metadata.LocationResponse.Value
				cidr := ""
				if subnet.Properties.Network != nil {
					cidr = subnet.Properties.Network.Address
				}
				status := ""
				if subnet.Status.State != nil {
					status = *subnet.Status.State
				}
				rows = append(rows, []string{name, id, region, cidr, status})
			}
			PrintTable(headers, rows)
		} else {
			fmt.Println("No subnets found.")
		}
	},
}

var subnetUpdateCmd = &cobra.Command{
	Use:   "update [vpc-id] [subnet-id]",
	Short: "Update a subnet",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]
		subnetID := args[1]
		name, _ := cmd.Flags().GetString("name")
		cidr, _ := cmd.Flags().GetString("cidr")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		if name == "" && cidr == "" && !cmd.Flags().Changed("tags") {
			fmt.Println("Error: at least one of --name, --cidr, or --tags must be provided")
			return
		}
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
		// Enable debug logging if supported
		if logger, ok := interface{}(client).(interface{ SetDebug(bool) }); ok {
			logger.SetDebug(true)
		}
		ctx := context.Background()
		// Fetch current subnet details
		getResp, err := client.FromNetwork().Subnets().Get(ctx, projectID, vpcID, subnetID, nil)
		if err != nil || getResp == nil || getResp.Data == nil {
			fmt.Printf("Error fetching current subnet: %v\n", err)
			return
		}
		current := getResp.Data
		// Block update if subnet is in 'InCreation' state
		if current.Status.State != nil && *current.Status.State == "InCreation" {
			fmt.Println("Error: Cannot update subnet while it is in 'InCreation' state. Please wait until the subnet is fully created.")
			return
		}

		// Get region value
		regionValue := ""
		if current.Metadata.LocationResponse != nil {
			regionValue = current.Metadata.LocationResponse.Value
		}
		if regionValue == "" {
			fmt.Println("Error: Unable to determine region value for subnet")
			return
		}

		// Build update request by merging user input with all current valid fields
		req := types.SubnetRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: func() string {
						if name != "" {
							return name
						}
						if current.Metadata.Name != nil {
							return *current.Metadata.Name
						}
						return ""
					}(),
					Tags: func() []string {
						if cmd.Flags().Changed("tags") {
							return tags
						}
						if current.Metadata.Tags != nil {
							return current.Metadata.Tags
						}
						return []string{}
					}(),
				},
				Location: types.LocationRequest{
					Value: regionValue,
				},
			},
			Properties: types.SubnetPropertiesRequest{
				Network: &types.SubnetNetwork{
					Address: func() string {
						if cidr != "" {
							return cidr
						}
						if current.Properties.Network != nil {
							return current.Properties.Network.Address
						}
						return ""
					}(),
				},
			},
		}

		resp, err := client.FromNetwork().Subnets().Update(ctx, projectID, vpcID, subnetID, req, nil)
		if err != nil {
			fmt.Printf("Error updating subnet: %v\n", err)
			return
		}
		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to update subnet - Status: %d\n", resp.StatusCode)
			if resp.Error.Title != nil {
				fmt.Printf("Error: %s\n", *resp.Error.Title)
			}
			if resp.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *resp.Error.Detail)
			}
			return
		}
		if resp != nil && resp.Data != nil {
			headers := []TableColumn{
				{Header: "NAME", Width: 30},
				{Header: "ID", Width: 26},
				{Header: "CIDR", Width: 18},
				{Header: "STATUS", Width: 15},
			}
			name := ""
			if resp.Data.Metadata.Name != nil {
				name = *resp.Data.Metadata.Name
			}
			id := ""
			if resp.Data.Metadata.ID != nil {
				id = *resp.Data.Metadata.ID
			}
			cidr := ""
			if resp.Data.Properties.Network != nil {
				cidr = resp.Data.Properties.Network.Address
			}
			status := ""
			if resp.Data.Status.State != nil {
				status = *resp.Data.Status.State
			}
			PrintTable(headers, [][]string{{name, id, cidr, status}})
		} else {
			fmt.Printf("Subnet '%s' updated.\n", subnetID)
		}
	},
}

var subnetDeleteCmd = &cobra.Command{
	Use:   "delete [vpc-id] [subnet-id]",
	Short: "Delete a subnet",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]
		subnetID := args[1]
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
		resp, err := client.FromNetwork().Subnets().Delete(ctx, projectID, vpcID, subnetID, nil)
		if err != nil {
			fmt.Printf("Error deleting subnet: %v\n", err)
			return
		}
		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to delete subnet - Status: %d\n", resp.StatusCode)
			if resp.Error.Title != nil {
				fmt.Printf("Error: %s\n", *resp.Error.Title)
			}
			if resp.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *resp.Error.Detail)
			}
			return
		}
		headers := []TableColumn{
			{Header: "ID", Width: 26},
			{Header: "STATUS", Width: 15},
		}
		status := "deleted"
		PrintTable(headers, [][]string{{subnetID, status}})
	},
}
