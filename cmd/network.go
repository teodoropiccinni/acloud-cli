// ...existing code...
package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/types"
	"github.com/spf13/cobra"
)

// Completion functions for network resources

func completeVPCID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	response, err := client.FromNetwork().VPCs().List(ctx, projectID, nil)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string
	if response != nil && response.Data != nil {
		for _, vpc := range response.Data.Values {
			if vpc.Metadata.ID != nil && vpc.Metadata.Name != nil {
				completions = append(completions, fmt.Sprintf("%s\t%s", *vpc.Metadata.ID, *vpc.Metadata.Name))
			}
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

func completeElasticIPID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	response, err := client.FromNetwork().ElasticIPs().List(ctx, projectID, nil)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string
	if response != nil && response.Data != nil {
		for _, eip := range response.Data.Values {
			if eip.Metadata.ID != nil && eip.Metadata.Name != nil {
				completions = append(completions, fmt.Sprintf("%s\t%s", *eip.Metadata.ID, *eip.Metadata.Name))
			}
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

func completeLoadBalancerID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	response, err := client.FromNetwork().LoadBalancers().List(ctx, projectID, nil)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string
	if response != nil && response.Data != nil {
		for _, lb := range response.Data.Values {
			if lb.Metadata.ID != nil && lb.Metadata.Name != nil {
				completions = append(completions, fmt.Sprintf("%s\t%s", *lb.Metadata.ID, *lb.Metadata.Name))
			}
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

func completeVPNTunnelID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	response, err := client.FromNetwork().VPNTunnels().List(ctx, projectID, nil)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string
	if response != nil && response.Data != nil {
		for _, vpn := range response.Data.Values {
			if vpn.Metadata.ID != nil && vpn.Metadata.Name != nil {
				completions = append(completions, fmt.Sprintf("%s\t%s", *vpn.Metadata.ID, *vpn.Metadata.Name))
			}
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage networks",
	Long:  `Manage network resources in Aruba Cloud.`,
}

// ElasticIP subcommands
var elasticipCmd = &cobra.Command{
	Use:   "elasticip",
	Short: "Manage Elastic IPs",
	Long:  `Perform CRUD operations on Elastic IPs in Aruba Cloud.`,
}

var elasticipCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Elastic IP",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		name, _ := cmd.Flags().GetString("name")
		region, _ := cmd.Flags().GetString("region")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		billingPeriod, _ := cmd.Flags().GetString("billing-period")

		// Validate required fields
		if name == "" {
			fmt.Println("Error: --name is required")
			return
		}
		if region == "" {
			fmt.Println("Error: --region is required")
			return
		}

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

		// Build the create request
		createRequest := types.ElasticIPRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: name,
					Tags: tags,
				},
				Location: types.LocationRequest{
					Value: region,
				},
			},
			Properties: types.ElasticIPPropertiesRequest{
				BillingPlan: types.BillingPeriodResource{
					BillingPeriod: billingPeriod,
				},
			},
		}

		// Create the Elastic IP using the SDK
		ctx := context.Background()
		response, err := client.FromNetwork().ElasticIPs().Create(ctx, projectID, createRequest, nil)
		if err != nil {
			fmt.Printf("Error creating Elastic IP: %v\n", err)
			return
		}

		if response != nil && !response.IsSuccess() {
			fmt.Printf("Failed to create Elastic IP - Status: %d\n", response.StatusCode)
			if response.Error.Title != nil {
				fmt.Printf("Error: %s\n", *response.Error.Title)
			}
			if response.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *response.Error.Detail)
			}
			return
		}

		if response != nil && response.Data != nil {
			fmt.Println("\nElastic IP created successfully!")
			if response.Data.Metadata.ID != nil {
				fmt.Printf("ID:      %s\n", *response.Data.Metadata.ID)
			}
			if response.Data.Metadata.Name != nil {
				fmt.Printf("Name:    %s\n", *response.Data.Metadata.Name)
			}
			if response.Data.Properties.Address != nil {
				fmt.Printf("Address: %s\n", *response.Data.Properties.Address)
			}
			if len(response.Data.Metadata.Tags) > 0 {
				fmt.Printf("Tags:    %v\n", response.Data.Metadata.Tags)
			}
		} else {
			fmt.Println("Elastic IP creation initiated. Use 'list' or 'get' to check status.")
		}
	},
}

var elasticipListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Elastic IPs",
	Run: func(cmd *cobra.Command, args []string) {
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

		// List Elastic IPs using the SDK
		ctx := context.Background()
		response, err := client.FromNetwork().ElasticIPs().List(ctx, projectID, nil)
		if err != nil {
			fmt.Printf("Error listing Elastic IPs: %v\n", err)
			return
		}

		if response != nil && response.Data != nil && len(response.Data.Values) > 0 {
			// Define table columns
			headers := []TableColumn{
				{Header: "NAME", Width: 30},
				{Header: "ID", Width: 26},
				{Header: "REGION", Width: 18},
				{Header: "ADDRESS", Width: 16},
				{Header: "STATUS", Width: 15},
			}

			// Build rows
			var rows [][]string
			for _, eip := range response.Data.Values {
				name := ""
				if eip.Metadata.Name != nil && *eip.Metadata.Name != "" {
					name = *eip.Metadata.Name
				}

				id := ""
				if eip.Metadata.ID != nil {
					id = *eip.Metadata.ID
				}

				region := eip.Metadata.LocationResponse.Code

				address := ""
				if eip.Properties.Address != nil {
					address = *eip.Properties.Address
				}

				status := ""
				if eip.Status.State != nil {
					status = *eip.Status.State
				}

				rows = append(rows, []string{name, id, region, address, status})
			}

			// Print the table
			PrintTable(headers, rows)
		} else {
			fmt.Println("No Elastic IPs found")
		}
	},
}

var elasticipGetCmd = &cobra.Command{
	Use:   "get <elastic-ip-id>",
	Short: "Get Elastic IP details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		eipID := args[0]

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

		// Get Elastic IP details using the SDK
		ctx := context.Background()
		response, err := client.FromNetwork().ElasticIPs().Get(ctx, projectID, eipID, nil)
		if err != nil {
			fmt.Printf("Error getting Elastic IP details: %v\n", err)
			return
		}

		if response != nil && response.Data != nil {
			eip := response.Data

			// Display Elastic IP details
			fmt.Println("\nElastic IP Details:")
			fmt.Println("===================")

			if eip.Metadata.ID != nil {
				fmt.Printf("ID:              %s\n", *eip.Metadata.ID)
			}
			if eip.Metadata.URI != nil {
				fmt.Printf("URI:             %s\n", *eip.Metadata.URI)
			}
			if eip.Metadata.Name != nil {
				fmt.Printf("Name:            %s\n", *eip.Metadata.Name)
			}
			if eip.Metadata.LocationResponse.Code != "" {
				fmt.Printf("Region:          %s\n", eip.Metadata.LocationResponse.Code)
			}
			if eip.Properties.Address != nil {
				fmt.Printf("Address:         %s\n", *eip.Properties.Address)
			}

			fmt.Printf("Billing Period:  %s\n", eip.Properties.BillingPlan.BillingPeriod)
			fmt.Printf("Linked Resources: %d\n", len(eip.Properties.LinkedResources))

			if eip.Metadata.CreationDate != nil {
				fmt.Printf("Creation Date:   %s\n", eip.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}
			if eip.Metadata.CreatedBy != nil {
				fmt.Printf("Created By:      %s\n", *eip.Metadata.CreatedBy)
			}

			if len(eip.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", eip.Metadata.Tags)
			}

			if eip.Status.State != nil {
				fmt.Printf("Status:          %s\n", *eip.Status.State)
			}
		}
	},
}

var elasticipUpdateCmd = &cobra.Command{
	Use:   "update <elastic-ip-id>",
	Short: "Update an Elastic IP",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		eipID := args[0]

		// Get flags
		name, _ := cmd.Flags().GetString("name")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		// At least one update flag must be provided
		if name == "" && len(tags) == 0 {
			fmt.Println("Error: at least one of --name or --tags must be provided")
			return
		}

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

		// First, get the current Elastic IP to preserve existing properties
		ctx := context.Background()
		getResponse, err := client.FromNetwork().ElasticIPs().Get(ctx, projectID, eipID, nil)
		if err != nil {
			fmt.Printf("Error getting Elastic IP details: %v\n", err)
			return
		}

		if getResponse == nil || getResponse.Data == nil {
			fmt.Println("Error: Elastic IP not found")
			return
		}

		// Check if Elastic IP is in InCreation state
		if getResponse.Data.Status.State != nil && *getResponse.Data.Status.State == "InCreation" {
			fmt.Println("Error: Cannot update Elastic IP while it is in 'InCreation' state. Please wait until the Elastic IP is fully created.")
			return
		}

		// Fix region code format (IT BG -> ITBG-Bergamo)
		regionCode := getResponse.Data.Metadata.LocationResponse.Code
		if regionCode == "IT BG" {
			regionCode = "ITBG-Bergamo"
		}

		// Build the update request, preserving existing values
		updateRequest := types.ElasticIPRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: *getResponse.Data.Metadata.Name,
					Tags: getResponse.Data.Metadata.Tags,
				},
				Location: types.LocationRequest{
					Value: regionCode,
				},
			},
			Properties: types.ElasticIPPropertiesRequest{
				BillingPlan: getResponse.Data.Properties.BillingPlan,
			},
		}

		// Apply updates
		if name != "" {
			updateRequest.Metadata.Name = name
		}
		if len(tags) > 0 {
			updateRequest.Metadata.Tags = tags
		}

		// Update the Elastic IP using the SDK
		response, err := client.FromNetwork().ElasticIPs().Update(ctx, projectID, eipID, updateRequest, nil)
		if err != nil {
			fmt.Printf("Error updating Elastic IP: %v\n", err)
			return
		}

		if response != nil && response.Data != nil {
			fmt.Println("\nElastic IP updated successfully!")
			if response.Data.Metadata.ID != nil {
				fmt.Printf("ID:      %s\n", *response.Data.Metadata.ID)
			}
			if response.Data.Metadata.Name != nil {
				fmt.Printf("Name:    %s\n", *response.Data.Metadata.Name)
			}
			if len(response.Data.Metadata.Tags) > 0 {
				fmt.Printf("Tags:    %v\n", response.Data.Metadata.Tags)
			}
		} else {
			fmt.Printf("\nElastic IP %s update completed.\n", eipID)
		}
	},
}

var elasticipDeleteCmd = &cobra.Command{
	Use:   "delete <elastic-ip-id>",
	Short: "Delete an Elastic IP",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		eipID := args[0]

		// Get skip confirmation flag
		skipConfirm, _ := cmd.Flags().GetBool("yes")

		// Prompt for confirmation unless --yes flag is used
		if !skipConfirm {
			fmt.Printf("Are you sure you want to delete Elastic IP %s? This action cannot be undone.\n", eipID)
			fmt.Print("Type 'yes' to confirm: ")
			var response string
			fmt.Scanln(&response)
			if response != "yes" && response != "y" {
				fmt.Println("Delete cancelled")
				return
			}
		}

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

		// Delete the Elastic IP using the SDK
		ctx := context.Background()
		_, err = client.FromNetwork().ElasticIPs().Delete(ctx, projectID, eipID, nil)
		if err != nil {
			fmt.Printf("Error deleting Elastic IP: %v\n", err)
			return
		}

		fmt.Printf("\nElastic IP %s deleted successfully!\n", eipID)
	},
}

// LoadBalancer subcommands
var loadbalancerCmd = &cobra.Command{
	Use:   "loadbalancer",
	Short: "Manage Load Balancers",
	Long:  `View Load Balancers in Aruba Cloud. Load Balancers are read-only resources.`,
}

var loadbalancerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Load Balancers",
	Run: func(cmd *cobra.Command, args []string) {
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

		// List Load Balancers using the SDK
		ctx := context.Background()
		response, err := client.FromNetwork().LoadBalancers().List(ctx, projectID, nil)
		if err != nil {
			fmt.Printf("Error listing Load Balancers: %v\n", err)
			return
		}

		if response != nil && response.Data != nil && len(response.Data.Values) > 0 {
			// Define table columns
			headers := []TableColumn{
				{Header: "NAME", Width: 30},
				{Header: "ID", Width: 26},
				{Header: "REGION", Width: 18},
				{Header: "ADDRESS", Width: 16},
				{Header: "STATUS", Width: 15},
			}

			// Build rows
			var rows [][]string
			for _, lb := range response.Data.Values {
				name := ""
				if lb.Metadata.Name != nil && *lb.Metadata.Name != "" {
					name = *lb.Metadata.Name
				}

				id := ""
				if lb.Metadata.ID != nil {
					id = *lb.Metadata.ID
				}

				region := lb.Metadata.LocationResponse.Code

				address := ""
				if lb.Properties.Address != nil {
					address = *lb.Properties.Address
				}

				status := ""
				if lb.Status.State != nil {
					status = *lb.Status.State
				}

				rows = append(rows, []string{name, id, region, address, status})
			}

			// Print the table
			PrintTable(headers, rows)
		} else {
			fmt.Println("No Load Balancers found")
		}
	},
}

var loadbalancerGetCmd = &cobra.Command{
	Use:   "get <loadbalancer-id>",
	Short: "Get Load Balancer details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lbID := args[0]

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

		// Get Load Balancer details using the SDK
		ctx := context.Background()
		response, err := client.FromNetwork().LoadBalancers().Get(ctx, projectID, lbID, nil)
		if err != nil {
			fmt.Printf("Error getting Load Balancer details: %v\n", err)
			return
		}

		if response != nil && response.Data != nil {
			lb := response.Data

			// Display Load Balancer details
			fmt.Println("\nLoad Balancer Details:")
			fmt.Println("======================")

			if lb.Metadata.ID != nil {
				fmt.Printf("ID:              %s\n", *lb.Metadata.ID)
			}
			if lb.Metadata.URI != nil {
				fmt.Printf("URI:             %s\n", *lb.Metadata.URI)
			}
			if lb.Metadata.Name != nil {
				fmt.Printf("Name:            %s\n", *lb.Metadata.Name)
			}
			if lb.Properties.Address != nil {
				fmt.Printf("Address:         %s\n", *lb.Properties.Address)
			}
			if lb.Properties.VPC != nil && lb.Properties.VPC.URI != "" {
				fmt.Printf("VPC:             %s\n", lb.Properties.VPC.URI)
			}

			fmt.Printf("Linked Resources: %d\n", len(lb.Properties.LinkedResources))

			if lb.Metadata.CreationDate != nil {
				fmt.Printf("Creation Date:   %s\n", lb.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}
			if lb.Metadata.CreatedBy != nil {
				fmt.Printf("Created By:      %s\n", *lb.Metadata.CreatedBy)
			}

			if len(lb.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", lb.Metadata.Tags)
			}

			if lb.Status.State != nil {
				fmt.Printf("Status:          %s\n", *lb.Status.State)
			}
		}
	},
}

// VPC subcommands
var vpcCmd = &cobra.Command{
	Use:   "vpc",
	Short: "Manage VPCs",
	Long:  `Perform CRUD operations on VPCs in Aruba Cloud.`,
}

var vpcCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new VPC",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
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

		// Build the create request (default and preset are always false)
		setDefault := false
		setPreset := false
		createRequest := types.VPCRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: name,
					Tags: tags,
				},
				Location: types.LocationRequest{
					Value: region,
				},
			},
			Properties: types.VPCPropertiesRequest{
				Properties: &types.VPCProperties{
					Default: &setDefault,
					Preset:  &setPreset,
				},
			},
		}

		// Create the VPC using the SDK
		ctx := context.Background()
		response, err := client.FromNetwork().VPCs().Create(ctx, projectID, createRequest, nil)
		if err != nil {
			fmt.Printf("Error creating VPC: %v\n", err)
			return
		}

		if response != nil && response.IsError() && response.Error != nil {
			fmt.Printf("Failed to create VPC - Status: %d\n", response.StatusCode)
			if response.Error.Title != nil {
				fmt.Printf("Error: %s\n", *response.Error.Title)
			}
			if response.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *response.Error.Detail)
			}
			return
		}

		if response != nil && response.Data != nil {
			fmt.Println("\nVPC created successfully!")
			if response.Data.Metadata.ID != nil {
				fmt.Printf("ID:      %s\n", *response.Data.Metadata.ID)
			}
			if response.Data.Metadata.Name != nil {
				fmt.Printf("Name:    %s\n", *response.Data.Metadata.Name)
			}
			fmt.Printf("Default: %t\n", response.Data.Properties.Default)
			if len(response.Data.Metadata.Tags) > 0 {
				fmt.Printf("Tags:    %v\n", response.Data.Metadata.Tags)
			}
		} else {
			fmt.Println("VPC creation initiated. Use 'list' or 'get' to check status.")
		}
	},
}

var vpcUpdateCmd = &cobra.Command{
	Use:   "update <vpc-id>",
	Short: "Update a VPC",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]

		// Get flags
		name, _ := cmd.Flags().GetString("name")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		// At least one update flag must be provided
		if name == "" && len(tags) == 0 {
			fmt.Println("Error: at least one of --name or --tags must be provided")
			return
		}

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

		// First, get the current VPC to preserve existing properties
		ctx := context.Background()
		getResponse, err := client.FromNetwork().VPCs().Get(ctx, projectID, vpcID, nil)
		if err != nil {
			fmt.Printf("Error getting VPC details: %v\n", err)
			return
		}

		if getResponse == nil || getResponse.Data == nil {
			fmt.Println("Error: VPC not found")
			return
		}

		// Check if VPC is in InCreation state
		if getResponse.Data.Status.State != nil && *getResponse.Data.Status.State == "InCreation" {
			fmt.Println("Error: Cannot update VPC while it is in 'InCreation' state. Please wait until the VPC is fully created.")
			return
		}

		// Fix region code format (IT BG -> ITBG-Bergamo)
		regionCode := getResponse.Data.Metadata.LocationResponse.Code
		if regionCode == "IT BG" {
			regionCode = "ITBG-Bergamo"
		}

		// Build the update request, preserving existing values
		updateRequest := types.VPCRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: *getResponse.Data.Metadata.Name,
					Tags: getResponse.Data.Metadata.Tags,
				},
				Location: types.LocationRequest{
					Value: regionCode,
				},
			},
			Properties: types.VPCPropertiesRequest{
				Properties: &types.VPCProperties{
					Default: &getResponse.Data.Properties.Default,
				},
			},
		}

		// Apply updates
		if name != "" {
			updateRequest.Metadata.Name = name
		}
		if len(tags) > 0 {
			updateRequest.Metadata.Tags = tags
		}

		// Update the VPC using the SDK
		response, err := client.FromNetwork().VPCs().Update(ctx, projectID, vpcID, updateRequest, nil)
		if err != nil {
			fmt.Printf("Error updating VPC: %v\n", err)
			return
		}

		if response != nil && response.Data != nil {
			fmt.Println("\nVPC updated successfully!")
			if response.Data.Metadata.ID != nil {
				fmt.Printf("ID:      %s\n", *response.Data.Metadata.ID)
			}
			if response.Data.Metadata.Name != nil {
				fmt.Printf("Name:    %s\n", *response.Data.Metadata.Name)
			}
			if len(response.Data.Metadata.Tags) > 0 {
				fmt.Printf("Tags:    %v\n", response.Data.Metadata.Tags)
			}
		} else {
			fmt.Printf("\nVPC %s update completed.\n", vpcID)
		}
	},
}

var vpcDeleteCmd = &cobra.Command{
	Use:   "delete <vpc-id>",
	Short: "Delete a VPC",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]

		// Get skip confirmation flag
		skipConfirm, _ := cmd.Flags().GetBool("yes")

		// Prompt for confirmation unless --yes flag is used
		if !skipConfirm {
			fmt.Printf("Are you sure you want to delete VPC %s? This action cannot be undone.\n", vpcID)
			fmt.Print("Type 'yes' to confirm: ")
			var response string
			fmt.Scanln(&response)
			if response != "yes" && response != "y" {
				fmt.Println("Delete cancelled")
				return
			}
		}

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

		// Delete the VPC using the SDK
		ctx := context.Background()
		_, err = client.FromNetwork().VPCs().Delete(ctx, projectID, vpcID, nil)
		if err != nil {
			fmt.Printf("Error deleting VPC: %v\n", err)
			return
		}

		fmt.Printf("\nVPC %s deleted successfully!\n", vpcID)
	},
}

var vpcListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VPCs",
	Run: func(cmd *cobra.Command, args []string) {
		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// Get projectID from flag or context
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// List VPCs using the SDK
		ctx := context.Background()
		response, err := client.FromNetwork().VPCs().List(ctx, projectID, nil)
		if err != nil {
			fmt.Printf("Error listing VPCs: %v\n", err)
			return
		}

		if response != nil && response.Data != nil && len(response.Data.Values) > 0 {
			// Define table columns
			headers := []TableColumn{
				{Header: "NAME", Width: 40},
				{Header: "ID", Width: 25},
				{Header: "REGION", Width: 18},
				{Header: "SUBNETS", Width: 10},
				{Header: "STATUS", Width: 15},
			}

			// Build rows
			var rows [][]string
			for _, vpc := range response.Data.Values {
				name := ""
				if vpc.Metadata.Name != nil && *vpc.Metadata.Name != "" {
					name = *vpc.Metadata.Name
				}

				id := ""
				if vpc.Metadata.ID != nil {
					id = *vpc.Metadata.ID
				}

				region := vpc.Metadata.LocationResponse.Code

				subnets := fmt.Sprintf("%d", len(vpc.Properties.LinkedResources))

				status := ""
				if vpc.Status.State != nil {
					status = *vpc.Status.State
				}

				rows = append(rows, []string{name, id, region, subnets, status})
			}

			// Print the table
			PrintTable(headers, rows)
		} else {
			fmt.Println("No VPCs found")
		}
	},
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
				resp.Data.Metadata.LocationResponse.Code,
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

func init() {

	// VPCPeering flags
	peeringCreateCmd.Flags().String("name", "", "VPC peering name (required)")
	peeringCreateCmd.Flags().String("region", "", "Region code (required)")
	peeringCreateCmd.Flags().String("peer-vpc-id", "", "Peer VPC URI (required)")
	peeringCreateCmd.Flags().StringSlice("tags", []string{}, "Tags (comma-separated)")
	// VPCPeering flags
	// SecurityGroup flags
	securitygroupCreateCmd.Flags().String("name", "", "Security group name (required)")
	securitygroupCreateCmd.Flags().String("region", "", "Region code (required)")
	securitygroupCreateCmd.Flags().StringSlice("tags", []string{}, "Tags (comma-separated)")
	subnetCreateCmd.Flags().String("region", "", "Region for the subnet (required)")
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
			if subnet.Metadata.LocationResponse.Code != "" {
				fmt.Printf("Region:          %s\n", subnet.Metadata.LocationResponse.Code)
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
				region := subnet.Metadata.LocationResponse.Code
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

		// Normalize region code if needed (e.g., IT BG -> ITBG-Bergamo)
		regionCode := current.Metadata.LocationResponse.Code
		if regionCode == "IT BG" {
			regionCode = "ITBG-Bergamo"
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
					Value: regionCode,
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

// SecurityGroup subcommands
var securitygroupCmd = &cobra.Command{
	Use:   "securitygroup",
	Short: "Manage security groups",
	Long:  `Perform CRUD operations on security groups in Aruba Cloud.`,
}

var securitygroupCreateCmd = &cobra.Command{
	Use:   "create [vpc-id]",
	Short: "Create a new security group",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]
		name, _ := cmd.Flags().GetString("name")
		region, _ := cmd.Flags().GetString("region")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		if name == "" || region == "" {
			fmt.Println("Error: --name and --region are required")
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
		req := types.SecurityGroupRequest{
			Metadata: types.ResourceMetadataRequest{
				Name: name,
				Tags: tags,
			},
		}
		resp, err := client.FromNetwork().SecurityGroups().Create(ctx, projectID, vpcID, req, nil)
		if err != nil {
			fmt.Printf("Error creating security group: %v\n", err)
			return
		}
		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to create security group - Status: %d\n", resp.StatusCode)
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
				{Header: "STATUS", Width: 15},
			}
			row := []string{
				name,
				*resp.Data.Metadata.ID,
				resp.Data.Metadata.LocationResponse.Code,
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
			fmt.Println("Security group created, but no ID returned.")
		}
	},
}

var securitygroupGetCmd = &cobra.Command{
	Use:   "get [vpc-id] [securitygroup-id]",
	Short: "Get security group details",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]
		sgID := args[1]
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
		resp, err := client.FromNetwork().SecurityGroups().Get(ctx, projectID, vpcID, sgID, nil)
		if err != nil {
			fmt.Printf("Error getting security group: %v\n", err)
			return
		}
		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to get security group - Status: %d\n", resp.StatusCode)
			if resp.Error.Title != nil {
				fmt.Printf("Error: %s\n", *resp.Error.Title)
			}
			if resp.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *resp.Error.Detail)
			}
			return
		}
		if resp != nil && resp.Data != nil {
			sg := resp.Data
			fmt.Println("\nSecurity Group Details:")
			fmt.Println("======================")
			if sg.Metadata.ID != nil {
				fmt.Printf("ID:              %s\n", *sg.Metadata.ID)
			}
			if sg.Metadata.URI != nil {
				fmt.Printf("URI:             %s\n", *sg.Metadata.URI)
			}
			if sg.Metadata.Name != nil {
				fmt.Printf("Name:            %s\n", *sg.Metadata.Name)
			}
			if sg.Metadata.LocationResponse != nil {
				fmt.Printf("Region:          %s\n", sg.Metadata.LocationResponse.Code)
			}
			if sg.Metadata.CreationDate != nil {
				fmt.Printf("Creation Date:   %s\n", sg.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}
			if sg.Metadata.CreatedBy != nil {
				fmt.Printf("Created By:      %s\n", *sg.Metadata.CreatedBy)
			}
			if len(sg.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", sg.Metadata.Tags)
			}
			if sg.Status.State != nil {
				fmt.Printf("Status:          %s\n", *sg.Status.State)
			}
		} else {
			fmt.Println("Security group not found or no data returned.")
		}
	},
}

var securitygroupListCmd = &cobra.Command{
	Use:   "list [vpc-id]",
	Short: "List security groups for a VPC",
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
		resp, err := client.FromNetwork().SecurityGroups().List(ctx, projectID, vpcID, nil)
		if err != nil {
			fmt.Printf("Error listing security groups: %v\n", err)
			return
		}
		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to list security groups - Status: %d\n", resp.StatusCode)
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
				{Header: "STATUS", Width: 15},
			}
			var rows [][]string
			for _, sg := range resp.Data.Values {
				name := ""
				if sg.Metadata.Name != nil {
					name = *sg.Metadata.Name
				}
				id := ""
				if sg.Metadata.ID != nil {
					id = *sg.Metadata.ID
				}
				region := ""
				if sg.Metadata.LocationResponse != nil {
					region = sg.Metadata.LocationResponse.Code
				}
				status := ""
				if sg.Status.State != nil {
					status = *sg.Status.State
				}
				rows = append(rows, []string{name, id, region, status})
			}
			PrintTable(headers, rows)
		} else {
			fmt.Println("No security groups found.")
		}
	},
}

var securitygroupUpdateCmd = &cobra.Command{
	Use:   "update [vpc-id] [securitygroup-id]",
	Short: "Update a security group",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]
		sgID := args[1]
		name, _ := cmd.Flags().GetString("name")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		if name == "" && !cmd.Flags().Changed("tags") {
			fmt.Println("Error: at least one of --name or --tags must be provided")
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
		// Fetch current security group details
		getResp, err := client.FromNetwork().SecurityGroups().Get(ctx, projectID, vpcID, sgID, nil)
		if err != nil || getResp == nil || getResp.Data == nil {
			fmt.Printf("Error fetching current security group: %v\n", err)
			return
		}
		current := getResp.Data
		// Block update if security group is in 'InCreation' state
		if current.Status.State != nil && *current.Status.State == "InCreation" {
			fmt.Println("Error: Cannot update security group while it is in 'InCreation' state. Please wait until the security group is fully created.")
			return
		}
		// Normalize region code if needed
		regionCode := current.Metadata.LocationResponse.Code
		if regionCode == "IT BG" {
			regionCode = "ITBG-Bergamo"
		}
		// Build update request by merging user input with all current valid fields
		req := types.SecurityGroupRequest{
			Metadata: types.ResourceMetadataRequest{
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
		}
		resp, err := client.FromNetwork().SecurityGroups().Update(ctx, projectID, vpcID, sgID, req, nil)
		if err != nil {
			fmt.Printf("Error updating security group: %v\n", err)
			return
		}
		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to update security group - Status: %d\n", resp.StatusCode)
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
				{Header: "REGION", Width: 18},
				{Header: "STATUS", Width: 15},
			}
			row := []string{
				func() string {
					if resp.Data.Metadata.Name != nil {
						return *resp.Data.Metadata.Name
					}
					return ""
				}(),
				func() string {
					if resp.Data.Metadata.ID != nil {
						return *resp.Data.Metadata.ID
					}
					return ""
				}(),
				resp.Data.Metadata.LocationResponse.Code,
				func() string {
					if resp.Data.Status.State != nil {
						return *resp.Data.Status.State
					}
					return ""
				}(),
			}
			PrintTable(headers, [][]string{row})
		} else {
			fmt.Printf("Security group '%s' updated.\n", sgID)
		}
	},
}

var securitygroupDeleteCmd = &cobra.Command{
	Use:   "delete [vpc-id] [securitygroup-id]",
	Short: "Delete a security group",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]
		sgID := args[1]
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
		resp, err := client.FromNetwork().SecurityGroups().Delete(ctx, projectID, vpcID, sgID, nil)
		if err != nil {
			fmt.Printf("Error deleting security group: %v\n", err)
			return
		}
		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to delete security group - Status: %d\n", resp.StatusCode)
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
		PrintTable(headers, [][]string{{sgID, status}})
	},
}

// SecurityRule subcommands
var securityruleCmd = &cobra.Command{
	Use:   "securityrule",
	Short: "Manage security rules",
	Long:  `Perform CRUD operations on security rules in Aruba Cloud.`,
}

var securityruleCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new security rule",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Security rule created (stub)")
	},
}

var securityruleGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get security rule details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Security rule details (stub)")
	},
}

var securityruleUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a security rule",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Security rule updated (stub)")
	},
}

var securityruleDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a security rule",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Security rule deleted (stub)")
	},
}

var securityruleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all security rules",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Security rule list (stub)")
	},
}

// Peering subcommands
var vpcpeeringCmd = &cobra.Command{
	Use:   "vpcpeering",
	Short: "Manage VPC peering",
	Long:  `Perform CRUD operations on VPC peering in Aruba Cloud.`,
}

// vpcpeeringrouteCmd represents the VPC peering route resource commands.
//
// Provides CRUD operations for VPC peering routes in Aruba Cloud.
var vpcpeeringrouteCmd = &cobra.Command{
	Use:   "vpcpeeringroute",
	Short: "Manage VPC peering routes",
	Long:  `Perform CRUD operations on VPC peering routes in Aruba Cloud.`,
}

var vpcpeeringrouteCreateCmd = &cobra.Command{
	Use:   "create [vpc-id] [peering-id]",
	Short: "Create a new VPC peering route",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[ERROR] VPC Peering Route API is not available in the current ArubaCloud SDK. Please use the Aruba Cloud Web Console or API if available.")
	},
}

var vpcpeeringrouteGetCmd = &cobra.Command{
	Use:   "get [vpc-id] [peering-id] [route-id]",
	Short: "Get VPC peering route details",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[ERROR] VPC Peering Route API is not available in the current ArubaCloud SDK. Please use the Aruba Cloud Web Console or API if available.")
	},
}

var vpcpeeringrouteUpdateCmd = &cobra.Command{
	Use:   "update [vpc-id] [peering-id] [route-id]",
	Short: "Update a VPC peering route",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[ERROR] VPC Peering Route API is not available in the current ArubaCloud SDK. Please use the Aruba Cloud Web Console or API if available.")
	},
}

var vpcpeeringrouteDeleteCmd = &cobra.Command{
	Use:   "delete [vpc-id] [peering-id] [route-id]",
	Short: "Delete a VPC peering route",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[ERROR] VPC Peering Route API is not available in the current ArubaCloud SDK. Please use the Aruba Cloud Web Console or API if available.")
	},
}

var vpcpeeringrouteListCmd = &cobra.Command{
	Use:   "list [vpc-id] [peering-id]",
	Short: "List VPC peering routes",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[ERROR] VPC Peering Route API is not available in the current ArubaCloud SDK. Please use the Aruba Cloud Web Console or API if available.")
	},
}

var peeringCreateCmd = &cobra.Command{
	Use:   "create [vpc-id]",
	Short: "Create a new VPC peering",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]
		name, _ := cmd.Flags().GetString("name")
		peerVPCID, _ := cmd.Flags().GetString("peer-vpc-id")
		region, _ := cmd.Flags().GetString("region")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		if name == "" || peerVPCID == "" || region == "" {
			fmt.Println("Error: --name, --peer-vpc-id, and --region are required")
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
		req := types.VPCPeeringRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: name,
					Tags: tags,
				},
				Location: types.LocationRequest{Value: region},
			},
			Properties: types.VPCPeeringPropertiesRequest{
				RemoteVPC: &types.ReferenceResource{URI: peerVPCID},
			},
		}
		resp, err := client.FromNetwork().VPCPeerings().Create(ctx, projectID, vpcID, req, nil)
		if err != nil {
			fmt.Printf("Error creating VPC peering: %v\n", err)
			return
		}
		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to create VPC peering - Status: %d\n", resp.StatusCode)
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
				{Header: "PEER VPC", Width: 26},
				{Header: "REGION", Width: 18},
				{Header: "STATUS", Width: 15},
			}
			row := []string{
				name,
				*resp.Data.Metadata.ID,
				func() string {
					if resp.Data.Properties.RemoteVPC != nil {
						return resp.Data.Properties.RemoteVPC.URI
					}
					return ""
				}(),
				func() string {
					if resp.Data.Metadata.LocationResponse != nil {
						return resp.Data.Metadata.LocationResponse.Code
					}
					return ""
				}(),
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
			fmt.Println("VPC peering created, but no ID returned.")
		}
	},
}

var peeringGetCmd = &cobra.Command{
	Use:   "get [vpc-id] [peering-id]",
	Short: "Get VPC peering details",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]
		peeringID := args[1]
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
		resp, err := client.FromNetwork().VPCPeerings().Get(ctx, projectID, vpcID, peeringID, nil)
		if err != nil {
			fmt.Printf("Error getting VPC peering: %v\n", err)
			return
		}
		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to get VPC peering - Status: %d\n", resp.StatusCode)
			if resp.Error.Title != nil {
				fmt.Printf("Error: %s\n", *resp.Error.Title)
			}
			if resp.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *resp.Error.Detail)
			}
			return
		}
		if resp != nil && resp.Data != nil {
			peering := resp.Data
			fmt.Println("\nVPC Peering Details:")
			fmt.Println("====================")
			if peering.Metadata.ID != nil {
				fmt.Printf("ID:              %s\n", *peering.Metadata.ID)
			}
			if peering.Metadata.Name != nil {
				fmt.Printf("Name:            %s\n", *peering.Metadata.Name)
			}
			if peering.Properties.RemoteVPC != nil {
				fmt.Printf("Peer VPC:        %s\n", peering.Properties.RemoteVPC.URI)
			}
			if peering.Metadata.LocationResponse != nil {
				fmt.Printf("Region:          %s\n", peering.Metadata.LocationResponse.Code)
			}
			if peering.Metadata.CreationDate != nil {
				fmt.Printf("Creation Date:   %s\n", peering.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}
			if peering.Metadata.CreatedBy != nil {
				fmt.Printf("Created By:      %s\n", *peering.Metadata.CreatedBy)
			}
			if len(peering.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", peering.Metadata.Tags)
			}
			if peering.Status.State != nil {
				fmt.Printf("Status:          %s\n", *peering.Status.State)
			}
		} else {
			fmt.Println("VPC peering not found or no data returned.")
		}
	},
}

var peeringListCmd = &cobra.Command{
	Use:   "list [vpc-id]",
	Short: "List VPC peerings for a VPC",
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
		resp, err := client.FromNetwork().VPCPeerings().List(ctx, projectID, vpcID, nil)
		if err != nil {
			fmt.Printf("Error listing VPC peerings: %v\n", err)
			return
		}
		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to list VPC peerings - Status: %d\n", resp.StatusCode)
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
				{Header: "PEER VPC", Width: 26},
				{Header: "REGION", Width: 18},
				{Header: "STATUS", Width: 15},
			}
			var rows [][]string
			for _, peering := range resp.Data.Values {
				name := ""
				if peering.Metadata.Name != nil {
					name = *peering.Metadata.Name
				}
				id := ""
				if peering.Metadata.ID != nil {
					id = *peering.Metadata.ID
				}
				peerVPC := ""
				if peering.Properties.RemoteVPC != nil {
					peerVPC = peering.Properties.RemoteVPC.URI
				}
				region := ""
				if peering.Metadata.LocationResponse != nil {
					region = peering.Metadata.LocationResponse.Code
				}
				status := ""
				if peering.Status.State != nil {
					status = *peering.Status.State
				}
				rows = append(rows, []string{name, id, peerVPC, region, status})
			}
			PrintTable(headers, rows)
		} else {
			fmt.Println("No VPC peerings found.")
		}
	},
}

var peeringUpdateCmd = &cobra.Command{
	Use:   "update [vpc-id] [peering-id]",
	Short: "Update a VPC peering",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]
		peeringID := args[1]
		name, _ := cmd.Flags().GetString("name")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		if name == "" && !cmd.Flags().Changed("tags") {
			fmt.Println("Error: at least one of --name or --tags must be provided")
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
		// Fetch current peering details
		getResp, err := client.FromNetwork().VPCPeerings().Get(ctx, projectID, vpcID, peeringID, nil)
		if err != nil || getResp == nil || getResp.Data == nil {
			fmt.Printf("Error fetching current VPC peering: %v\n", err)
			return
		}
		current := getResp.Data
		// Block update if peering is in 'InCreation' state
		if current.Status.State != nil && *current.Status.State == "InCreation" {
			fmt.Println("Error: Cannot update VPC peering while it is in 'InCreation' state. Please wait until the VPC peering is fully created.")
			return
		}
		// Build update request by merging user input with all current valid fields
		req := types.VPCPeeringRequest{
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
					Value: func() string {
						if current.Metadata.LocationResponse != nil {
							return current.Metadata.LocationResponse.Code
						}
						return ""
					}(),
				},
			},
			Properties: types.VPCPeeringPropertiesRequest{
				RemoteVPC: current.Properties.RemoteVPC,
			},
		}
		resp, err := client.FromNetwork().VPCPeerings().Update(ctx, projectID, vpcID, peeringID, req, nil)
		if err != nil {
			fmt.Printf("Error updating VPC peering: %v\n", err)
			return
		}
		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to update VPC peering - Status: %d\n", resp.StatusCode)
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
				{Header: "PEER VPC", Width: 26},
				{Header: "REGION", Width: 18},
				{Header: "STATUS", Width: 15},
			}
			row := []string{
				func() string {
					if resp.Data.Metadata.Name != nil {
						return *resp.Data.Metadata.Name
					}
					return ""
				}(),
				func() string {
					if resp.Data.Metadata.ID != nil {
						return *resp.Data.Metadata.ID
					}
					return ""
				}(),
				func() string {
					if resp.Data.Properties.RemoteVPC != nil {
						return resp.Data.Properties.RemoteVPC.URI
					}
					return ""
				}(),
				func() string {
					if resp.Data.Metadata.LocationResponse != nil {
						return resp.Data.Metadata.LocationResponse.Code
					}
					return ""
				}(),
				func() string {
					if resp.Data.Status.State != nil {
						return *resp.Data.Status.State
					}
					return ""
				}(),
			}
			PrintTable(headers, [][]string{row})
		} else {
			fmt.Printf("VPC peering '%s' updated.\n", peeringID)
		}
	},
}

var peeringDeleteCmd = &cobra.Command{
	Use:   "delete [vpc-id] [peering-id]",
	Short: "Delete a VPC peering",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]
		peeringID := args[1]
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
		resp, err := client.FromNetwork().VPCPeerings().Delete(ctx, projectID, vpcID, peeringID, nil)
		if err != nil {
			fmt.Printf("Error deleting VPC peering: %v\n", err)
			return
		}
		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to delete VPC peering - Status: %d\n", resp.StatusCode)
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
		PrintTable(headers, [][]string{{peeringID, status}})
	},
}

// Route subcommands
var routeCmd = &cobra.Command{
	Use:   "route",
	Short: "Manage peering routes",
	Long:  `Perform CRUD operations on peering routes in Aruba Cloud.`,
}

var routeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new peering route",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Peering route created (stub)")
	},
}

var routeGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get peering route details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Peering route details (stub)")
	},
}

var routeUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a peering route",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Peering route updated (stub)")
	},
}

var routeDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a peering route",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Peering route deleted (stub)")
	},
}

var routeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all peering routes",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Peering route list (stub)")
	},
}

// VPNTunnel subcommands
var vpntunnelCmd = &cobra.Command{
	Use:   "vpntunnel",
	Short: "Manage VPN tunnels",
	Long:  `Perform CRUD operations on VPN tunnels in Aruba Cloud.`,
}

// vpnrouteCmd represents the VPN tunnel route resource commands.
//
// Provides CRUD operations for VPN tunnel routes in Aruba Cloud.
var vpnrouteCmd = &cobra.Command{
	Use:   "vpnroute",
	Short: "Manage VPN tunnel routes",
	Long:  `Perform CRUD operations on VPN tunnel routes in Aruba Cloud.`,
}

var vpnrouteCreateCmd = &cobra.Command{
	Use:   "create [vpn-tunnel-id]",
	Short: "Create a new VPN tunnel route",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[ERROR] VPN Tunnel Route API is not available in the current ArubaCloud SDK. Please use the Aruba Cloud Web Console or API if available.")
	},
}

var vpnrouteGetCmd = &cobra.Command{
	Use:   "get [vpn-tunnel-id] [route-id]",
	Short: "Get VPN tunnel route details",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[ERROR] VPN Tunnel Route API is not available in the current ArubaCloud SDK. Please use the Aruba Cloud Web Console or API if available.")
	},
}

var vpnrouteUpdateCmd = &cobra.Command{
	Use:   "update [vpn-tunnel-id] [route-id]",
	Short: "Update a VPN tunnel route",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[ERROR] VPN Tunnel Route API is not available in the current ArubaCloud SDK. Please use the Aruba Cloud Web Console or API if available.")
	},
}

var vpnrouteDeleteCmd = &cobra.Command{
	Use:   "delete [vpn-tunnel-id] [route-id]",
	Short: "Delete a VPN tunnel route",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[ERROR] VPN Tunnel Route API is not available in the current ArubaCloud SDK. Please use the Aruba Cloud Web Console or API if available.")
	},
}

var vpnrouteListCmd = &cobra.Command{
	Use:   "list [vpn-tunnel-id]",
	Short: "List VPN tunnel routes",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[ERROR] VPN Tunnel Route API is not available in the current ArubaCloud SDK. Please use the Aruba Cloud Web Console or API if available.")
	},
}

var vpntunnelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VPN tunnels",
	Run: func(cmd *cobra.Command, args []string) {
		// Get SDK client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// Get projectID from flag or context
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// List VPN tunnels using the SDK
		ctx := context.Background()
		response, err := client.FromNetwork().VPNTunnels().List(ctx, projectID, nil)
		if err != nil {
			fmt.Printf("Error listing VPN tunnels: %v\n", err)
			return
		}

		if response != nil && response.Data != nil && len(response.Data.Values) > 0 {
			// Define table columns
			headers := []TableColumn{
				{Header: "NAME", Width: 40},
				{Header: "ID", Width: 25},
				{Header: "REGION", Width: 18},
				{Header: "TYPE", Width: 15},
				{Header: "STATUS", Width: 15},
			}
			// Build rows
			var rows [][]string
			for _, vpn := range response.Data.Values {
				name := ""
				if vpn.Metadata.Name != nil && *vpn.Metadata.Name != "" {
					name = *vpn.Metadata.Name
				}

				id := ""
				if vpn.Metadata.ID != nil {
					id = *vpn.Metadata.ID
				}

				region := vpn.Metadata.LocationResponse.Code

				vpnType := ""
				if vpn.Properties.VPNType != nil {
					vpnType = *vpn.Properties.VPNType
				}

				status := ""
				if vpn.Status.State != nil {
					status = *vpn.Status.State
				}

				rows = append(rows, []string{name, id, region, vpnType, status})
			}

			// Print the table
			PrintTable(headers, rows)
		} else {
			fmt.Println("No VPN tunnels found")
		}
	},
}

var vpntunnelGetCmd = &cobra.Command{
	Use:   "get <vpn-tunnel-id>",
	Short: "Get VPN tunnel details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vpnID := args[0]

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

		// Get VPN tunnel details using the SDK
		ctx := context.Background()
		response, err := client.FromNetwork().VPNTunnels().Get(ctx, projectID, vpnID, nil)
		if err != nil {
			fmt.Printf("Error getting VPN tunnel details: %v\n", err)
			return
		}

		if response != nil && response.Data != nil {
			vpn := response.Data

			// Display VPN tunnel details
			fmt.Println("\nVPN Tunnel Details:")
			fmt.Println("===================")

			if vpn.Metadata.ID != nil {
				fmt.Printf("ID:              %s\n", *vpn.Metadata.ID)
			}
			if vpn.Metadata.URI != nil {
				fmt.Printf("URI:             %s\n", *vpn.Metadata.URI)
			}
			if vpn.Metadata.Name != nil {
				fmt.Printf("Name:            %s\n", *vpn.Metadata.Name)
			}
			if vpn.Metadata.LocationResponse.Code != "" {
				fmt.Printf("Region:          %s\n", vpn.Metadata.LocationResponse.Code)
			}

			if vpn.Properties.VPNType != nil {
				fmt.Printf("VPN Type:        %s\n", *vpn.Properties.VPNType)
			}
			if vpn.Properties.VPNClientProtocol != nil {
				fmt.Printf("Protocol:        %s\n", *vpn.Properties.VPNClientProtocol)
			}
			if vpn.Properties.VPNClientSettings != nil && vpn.Properties.VPNClientSettings.PeerClientPublicIP != nil {
				fmt.Printf("Peer IP:         %s\n", *vpn.Properties.VPNClientSettings.PeerClientPublicIP)
			}

			if vpn.Properties.IPConfigurations != nil {
				fmt.Println("\nIP Configuration:")
				if vpn.Properties.IPConfigurations.VPC != nil {
					fmt.Printf("  VPC:           %s\n", vpn.Properties.IPConfigurations.VPC.URI)
				}
				if vpn.Properties.IPConfigurations.Subnet != nil {
					fmt.Printf("  Subnet CIDR:   %s\n", vpn.Properties.IPConfigurations.Subnet.CIDR)
					if vpn.Properties.IPConfigurations.Subnet.Name != "" {
						fmt.Printf("  Subnet Name:   %s\n", vpn.Properties.IPConfigurations.Subnet.Name)
					}
				}
				if vpn.Properties.IPConfigurations.PublicIP != nil {
					fmt.Printf("  Public IP:     %s\n", vpn.Properties.IPConfigurations.PublicIP.URI)
				}
			}

			if vpn.Properties.BillingPlan != nil {
				fmt.Printf("\nBilling Period:  %s\n", vpn.Properties.BillingPlan.BillingPeriod)
			}

			if vpn.Metadata.CreationDate != nil {
				fmt.Printf("Creation Date:   %s\n", vpn.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}
			if vpn.Metadata.CreatedBy != nil {
				fmt.Printf("Created By:      %s\n", *vpn.Metadata.CreatedBy)
			}

			if vpn.Status.State != nil {
				fmt.Printf("Status:          %s\n", *vpn.Status.State)
			}
		}
	},
}

var vpntunnelCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new VPN tunnel",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		name, _ := cmd.Flags().GetString("name")
		region, _ := cmd.Flags().GetString("region")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		vpnType, _ := cmd.Flags().GetString("vpn-type")
		protocol, _ := cmd.Flags().GetString("protocol")
		peerIP, _ := cmd.Flags().GetString("peer-ip")
		vpcURI, _ := cmd.Flags().GetString("vpc-uri")
		subnetCIDR, _ := cmd.Flags().GetString("subnet-cidr")
		subnetName, _ := cmd.Flags().GetString("subnet-name")
		publicIPURI, _ := cmd.Flags().GetString("elastic-ip-uri")
		billingPeriod, _ := cmd.Flags().GetString("billing-period")
		psk, _ := cmd.Flags().GetString("psk")

		// IKE settings
		ikeLifetime, _ := cmd.Flags().GetInt32("ike-lifetime")
		ikeEncryption, _ := cmd.Flags().GetString("ike-encryption")
		ikeHash, _ := cmd.Flags().GetString("ike-hash")
		ikeDHGroup, _ := cmd.Flags().GetString("ike-dh-group")
		ikeDPDAction, _ := cmd.Flags().GetString("ike-dpd-action")
		ikeDPDInterval, _ := cmd.Flags().GetInt32("ike-dpd-interval")
		ikeDPDTimeout, _ := cmd.Flags().GetInt32("ike-dpd-timeout")
		// ESP settings
		espLifetime, _ := cmd.Flags().GetInt32("esp-lifetime")
		espEncryption, _ := cmd.Flags().GetString("esp-encryption")
		if espEncryption == "" {
			espEncryption = "aes256"
		}
		espHash, _ := cmd.Flags().GetString("esp-hash")
		espPFS, _ := cmd.Flags().GetString("esp-pfs")
		// PSK settings
		pskCloudSite, _ := cmd.Flags().GetString("psk-cloud-site")
		pskOnpremSite, _ := cmd.Flags().GetString("psk-onprem-site")

		// Validate required fields
		if name == "" {
			fmt.Println("Error: --name is required")
			return
		}
		if region == "" {
			fmt.Println("Error: --region is required")
			return
		}
		if peerIP == "" {
			fmt.Println("Error: --peer-ip is required")
			return
		}
		if vpcURI == "" {
			fmt.Println("Error: --vpc-uri is required (e.g., /projects/{project-id}/providers/Aruba.Network/vpcs/{vpc-id})")
			return
		}
		if subnetCIDR == "" && subnetName == "" {
			fmt.Println("Error: --subnet-cidr or --subnet-name is required")
			return
		}
		if publicIPURI == "" {
			fmt.Println("Error: --elastic-ip-uri is required (e.g., /projects/{project-id}/providers/Aruba.Network/elasticIps/{ip-id})")
			return
		}

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

		// Build subnet object with both CIDR and Name fields
		subnetRef := &types.SubnetInfo{}
		if subnetCIDR != "" {
			subnetRef.CIDR = subnetCIDR
		}
		if subnetName != "" {
			subnetRef.Name = subnetName
		}

		// Build IP configurations
		ipConfig := types.IPConfigurations{
			VPC: &types.ReferenceResource{
				URI: vpcURI,
			},
			Subnet: subnetRef,
			PublicIP: &types.ReferenceResource{
				URI: publicIPURI,
			},
		}

		// Build IKE settings
		ikeSettings := &types.IKESettings{
			Lifetime:    ikeLifetime,
			Encryption:  nil,
			Hash:        nil,
			DHGroup:     nil,
			DPDAction:   nil,
			DPDInterval: ikeDPDInterval,
			DPDTimeout:  ikeDPDTimeout,
		}
		if ikeEncryption != "" {
			ikeSettings.Encryption = &ikeEncryption
		}
		if ikeHash != "" {
			ikeSettings.Hash = &ikeHash
		}
		if ikeDHGroup != "" {
			ikeSettings.DHGroup = &ikeDHGroup
		}
		if ikeDPDAction != "" {
			ikeSettings.DPDAction = &ikeDPDAction
		}

		// Build ESP settings
		espSettings := &types.ESPSettings{
			Lifetime:   espLifetime,
			Encryption: &espEncryption, // always include, default to aes256
			Hash:       nil,
			PFS:        nil,
		}
		if espHash != "" {
			espSettings.Hash = &espHash
		}
		if espPFS != "" {
			espSettings.PFS = &espPFS
		}

		// Build PSK settings
		pskSettings := &types.PSKSettings{
			Secret:     nil,
			CloudSite:  nil,
			OnPremSite: nil,
		}
		if psk != "" {
			pskSettings.Secret = &psk
		}
		if pskCloudSite != "" {
			pskSettings.CloudSite = &pskCloudSite
		}
		if pskOnpremSite != "" {
			pskSettings.OnPremSite = &pskOnpremSite
		}

		// Build VPN client settings
		vpnClientSettings := &types.VPNClientSettings{
			IKE:                ikeSettings,
			ESP:                espSettings,
			PSK:                pskSettings,
			PeerClientPublicIP: &peerIP,
		}

		// Build the create request using custom types that match the API
		createRequest := types.VPNTunnelRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: name,
					Tags: tags,
				},
				Location: types.LocationRequest{
					Value: region,
				},
			},
			Properties: types.VPNTunnelPropertiesRequest{
				VPNType:           &vpnType,
				VPNClientProtocol: &protocol,
				IPConfigurations:  &ipConfig,
				VPNClientSettings: vpnClientSettings,
				BillingPlan: &types.BillingPeriodResource{
					BillingPeriod: billingPeriod,
				},
			},
		}

		// Debug: Print request body
		fmt.Println("\n=== DEBUG: Create Request ===")
		requestJSON, _ := json.MarshalIndent(createRequest, "", "  ")
		fmt.Println(string(requestJSON))
		fmt.Println("=== End Debug ===")

		// Create the VPN tunnel using the SDK
		// Note: SDK v0.1.2 types don't match the API, so we marshal/unmarshal to convert
		ctx := context.Background()
		requestBody, err := json.Marshal(createRequest)
		if err != nil {
			fmt.Printf("Error marshaling request: %v\n", err)
			return
		}

		// Convert our custom request to SDK request type
		var sdkRequest types.VPNTunnelRequest
		if err := json.Unmarshal(requestBody, &sdkRequest); err != nil {
			fmt.Printf("Error converting request: %v\n", err)
			return
		}

		response, err := client.FromNetwork().VPNTunnels().Create(ctx, projectID, sdkRequest, nil)
		if err != nil {
			fmt.Printf("Error creating VPN tunnel: %v\n", err)
			if response != nil && response.RawBody != nil {
				fmt.Println("--- HTTP Response Body ---")
				fmt.Println(string(response.RawBody))
				fmt.Println("--------------------------")
			}
			return
		}

		if response != nil && response.IsError() && response.Error != nil {
			fmt.Printf("Failed to create VPN tunnel - Status: %d\n", response.StatusCode)
			if response.Error.Title != nil {
				fmt.Printf("Error: %s\n", *response.Error.Title)
			}
			if response.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *response.Error.Detail)
			}
			if response.RawBody != nil {
				fmt.Println("--- HTTP Response Body ---")
				fmt.Println(string(response.RawBody))
				fmt.Println("--------------------------")
			}
			return
		}

		if response != nil && response.Data != nil {
			fmt.Println("\nVPN Tunnel created successfully!")
			if response.Data.Metadata.ID != nil {
				fmt.Printf("ID:       %s\n", *response.Data.Metadata.ID)
			}
			if response.Data.Metadata.Name != nil {
				fmt.Printf("Name:     %s\n", *response.Data.Metadata.Name)
			}
			if response.Data.Metadata.URI != nil {
				fmt.Printf("URI:      %s\n", *response.Data.Metadata.URI)
			}
			if response.Data.Properties.VPNType != nil {
				fmt.Printf("Type:     %s\n", *response.Data.Properties.VPNType)
			}
			if response.Data.Properties.VPNClientProtocol != nil {
				fmt.Printf("Protocol: %s\n", *response.Data.Properties.VPNClientProtocol)
			}
			if len(response.Data.Metadata.Tags) > 0 {
				fmt.Printf("Tags:     %v\n", response.Data.Metadata.Tags)
			}
		} else {
			fmt.Println("VPN Tunnel creation initiated. Use 'list' or 'get' to check status.")
		}
	},
}

var vpntunnelUpdateCmd = &cobra.Command{
	Use:   "update <vpn-tunnel-id>",
	Short: "Update a VPN tunnel",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vpnTunnelID := args[0]

		// Get flags
		name, _ := cmd.Flags().GetString("name")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		// At least one update flag must be provided
		if name == "" && len(tags) == 0 {
			fmt.Println("Error: at least one of --name or --tags must be provided")
			return
		}

		// Get project ID
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Get Aruba Cloud client
		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// Get current VPN tunnel configuration
		ctx := context.Background()
		getResp, err := client.FromNetwork().VPNTunnels().Get(ctx, projectID, vpnTunnelID, nil)
		if err != nil {
			fmt.Printf("Error getting VPN tunnel: %v\n", err)
			return
		}

		if getResp != nil && getResp.IsError() && getResp.Error != nil {
			fmt.Printf("Error getting VPN tunnel - Status: %d\n", getResp.StatusCode)
			if getResp.Error.Title != nil {
				fmt.Printf("Error: %s\n", *getResp.Error.Title)
			}
			if getResp.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *getResp.Error.Detail)
			}
			return
		}

		if getResp.Data == nil {
			fmt.Println("Error: VPN tunnel not found")
			return
		}

		// Check if VPN tunnel is in "InCreation" state
		if getResp.Data.Status.State != nil && *getResp.Data.Status.State == "InCreation" {
			fmt.Println("Error: Cannot update VPN tunnel while it is in 'InCreation' state. Please wait until the VPN tunnel is fully created.")
			return
		}

		// Get region code and normalize it
		regionCode := getResp.Data.Metadata.LocationResponse.Code
		if regionCode == "IT BG" {
			regionCode = "ITBG-Bergamo"
		}

		// Build update request, preserving current values
		updateReq := types.VPNTunnelRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: *getResp.Data.Metadata.Name,
					Tags: getResp.Data.Metadata.Tags,
				},
				Location: types.LocationRequest{
					Value: regionCode,
				},
			},
			Properties: types.VPNTunnelPropertiesRequest{
				VPNType:           getResp.Data.Properties.VPNType,
				VPNClientProtocol: getResp.Data.Properties.VPNClientProtocol,
				IPConfigurations:  getResp.Data.Properties.IPConfigurations,
				VPNClientSettings: getResp.Data.Properties.VPNClientSettings,
				// PeerClientPublicIP is now set via VPNClientSettings only
				BillingPlan: getResp.Data.Properties.BillingPlan,
			},
		}

		// Apply updates
		if name != "" {
			updateReq.Metadata.Name = name
		}
		if len(tags) > 0 {
			updateReq.Metadata.Tags = tags
		}

		// Update VPN tunnel
		resp, err := client.FromNetwork().VPNTunnels().Update(ctx, projectID, vpnTunnelID, updateReq, nil)
		if err != nil {
			fmt.Printf("Error updating VPN tunnel: %v\n", err)
			return
		}

		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to update VPN tunnel - Status: %d\n", resp.StatusCode)
			if resp.Error.Title != nil {
				fmt.Printf("Error: %s\n", *resp.Error.Title)
			}
			if resp.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *resp.Error.Detail)
			}
			return
		}

		if resp.Data != nil {
			fmt.Println("\nVPN Tunnel updated successfully!")
			if resp.Data.Metadata.ID != nil {
				fmt.Printf("ID:      %s\n", *resp.Data.Metadata.ID)
			}
			if resp.Data.Metadata.Name != nil {
				fmt.Printf("Name:    %s\n", *resp.Data.Metadata.Name)
			}
			if len(resp.Data.Metadata.Tags) > 0 {
				fmt.Printf("Tags:    %v\n", resp.Data.Metadata.Tags)
			}
		} else {
			fmt.Printf("\nVPN Tunnel %s update completed.\n", vpnTunnelID)
		}
	},
}

var vpntunnelDeleteCmd = &cobra.Command{
	Use:   "delete <vpn-tunnel-id>",
	Short: "Delete a VPN tunnel",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vpnTunnelID := args[0]

		// Get skip confirmation flag
		skipConfirm, _ := cmd.Flags().GetBool("yes")

		// Prompt for confirmation unless --yes flag is used
		if !skipConfirm {
			fmt.Printf("Are you sure you want to delete VPN tunnel %s? This action cannot be undone.\n", vpnTunnelID)
			fmt.Print("Type 'yes' to confirm: ")
			var response string
			fmt.Scanln(&response)
			if response != "yes" && response != "y" {
				fmt.Println("Delete cancelled")
				return
			}
		}

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

		// Debug output
		fmt.Printf("[DEBUG] Project ID: %s\n", projectID)
		fmt.Printf("[DEBUG] VPN Tunnel ID: %s\n", vpnTunnelID)

		// Delete the VPN tunnel using the SDK
		ctx := context.Background()
		response, err := client.FromNetwork().VPNTunnels().Delete(ctx, projectID, vpnTunnelID, nil)
		if err != nil {
			fmt.Printf("Error deleting VPN tunnel: %v\n", err)
			// Print full error if possible
			if response != nil {
				fmt.Printf("[DEBUG] Full response: %+v\n", response)
			}
			return
		}

		if response != nil && response.IsError() && response.Error != nil {
			fmt.Printf("Failed to delete VPN tunnel - Status: %d\n", response.StatusCode)
			if response.Error.Title != nil {
				fmt.Printf("Error: %s\n", *response.Error.Title)
			}
			if response.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *response.Error.Detail)
			}
			// Print full error for debugging
			fmt.Printf("[DEBUG] Full error: %+v\n", response.Error)
			return
		}

		fmt.Printf("VPN tunnel %s has been successfully deleted.\n", vpnTunnelID)
	},
}

// Route subcommands under VPNTunnel
var vpntunnelRouteCmd = &cobra.Command{
	Use:   "route",
	Short: "Manage VPN tunnel routes",
	Long:  `Perform CRUD operations on VPN tunnel routes in Aruba Cloud.`,
}

var vpntunnelRouteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VPN tunnel routes",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPN tunnel route list operation is not yet implemented.")
		fmt.Println("Please use the Aruba Cloud Web Console or API to manage VPN tunnel routes.")
	},
}

var vpntunnelRouteGetCmd = &cobra.Command{
	Use:   "get <route-id>",
	Short: "Get VPN tunnel route details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPN tunnel route get operation is not yet implemented.")
		fmt.Println("Please use the Aruba Cloud Web Console or API to manage VPN tunnel routes.")
	},
}

var vpntunnelRouteCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new VPN tunnel route",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPN tunnel route create operation is not yet implemented.")
		fmt.Println("Please use the Aruba Cloud Web Console or API to manage VPN tunnel routes.")
	},
}

var vpntunnelRouteUpdateCmd = &cobra.Command{
	Use:   "update <route-id>",
	Short: "Update a VPN tunnel route",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPN tunnel route update operation is not yet implemented.")
		fmt.Println("Please use the Aruba Cloud Web Console or API to manage VPN tunnel routes.")
	},
}

var vpntunnelRouteDeleteCmd = &cobra.Command{
	Use:   "delete <route-id>",
	Short: "Delete a VPN tunnel route",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPN tunnel route delete operation is not yet implemented.")
		fmt.Println("Please use the Aruba Cloud Web Console or API to manage VPN tunnel routes.")
	},
}

func init() {
	rootCmd.AddCommand(networkCmd)
	// ElasticIP
	networkCmd.AddCommand(elasticipCmd)
	elasticipCmd.AddCommand(elasticipCreateCmd)
	elasticipCmd.AddCommand(elasticipGetCmd)
	elasticipCmd.AddCommand(elasticipUpdateCmd)
	elasticipCmd.AddCommand(elasticipDeleteCmd)
	elasticipCmd.AddCommand(elasticipListCmd)
	// LoadBalancer (read-only)
	networkCmd.AddCommand(loadbalancerCmd)
	loadbalancerCmd.AddCommand(loadbalancerListCmd)
	loadbalancerCmd.AddCommand(loadbalancerGetCmd)
	// VPC
	networkCmd.AddCommand(vpcCmd)
	vpcCmd.AddCommand(vpcCreateCmd)
	vpcCmd.AddCommand(vpcUpdateCmd)
	vpcCmd.AddCommand(vpcDeleteCmd)
	vpcCmd.AddCommand(vpcListCmd)

	// Add flags for Elastic IP commands
	elasticipCreateCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	elasticipCreateCmd.Flags().String("name", "", "Name for the Elastic IP")
	elasticipCreateCmd.Flags().String("region", "", "Region code (e.g., IT-BG)")
	elasticipCreateCmd.Flags().String("billing-period", "Hour", "Billing period: Hour, Month, Year")
	elasticipCreateCmd.Flags().StringSlice("tags", []string{}, "Tags (comma-separated)")

	elasticipListCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	elasticipGetCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")

	elasticipUpdateCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	elasticipUpdateCmd.Flags().String("name", "", "New name for the Elastic IP")
	elasticipUpdateCmd.Flags().StringSlice("tags", []string{}, "New tags (comma-separated)")

	elasticipDeleteCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	elasticipDeleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	// Add flags for Load Balancer commands
	loadbalancerListCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	loadbalancerGetCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")

	// Add flags for VPC commands
	vpcCreateCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	vpcCreateCmd.Flags().String("name", "", "Name for the VPC")
	vpcCreateCmd.Flags().String("region", "", "Region code (e.g., IT-BG)")
	vpcCreateCmd.Flags().StringSlice("tags", []string{}, "Tags (comma-separated)")

	vpcUpdateCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	vpcUpdateCmd.Flags().String("name", "", "New name for the VPC")
	vpcUpdateCmd.Flags().StringSlice("tags", []string{}, "New tags (comma-separated)")

	vpcDeleteCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	vpcDeleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	// Add flags for vpc commands
	vpcListCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")

	// Add flags for VPN Tunnel commands
	vpntunnelCreateCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	vpntunnelCreateCmd.Flags().String("name", "", "Name for the VPN tunnel")
	vpntunnelCreateCmd.Flags().String("region", "", "Region code (e.g., ITBG-Bergamo)")
	vpntunnelCreateCmd.Flags().String("peer-ip", "", "Peer client public IP address")
	vpntunnelCreateCmd.Flags().String("vpc-uri", "", "VPC URI (e.g., /projects/{project-id}/providers/Aruba.Network/vpcs/{vpc-id})")
	vpntunnelCreateCmd.Flags().String("subnet-cidr", "", "Subnet CIDR (e.g., 10.0.1.0/24)")
	vpntunnelCreateCmd.Flags().String("subnet-name", "", "Subnet name (alternative to CIDR)")
	vpntunnelCreateCmd.Flags().String("elastic-ip-uri", "", "Elastic IP URI (e.g., /projects/{project-id}/providers/Aruba.Network/elasticIps/{ip-id})")
	vpntunnelCreateCmd.Flags().String("vpn-type", "Site-To-Site", "VPN type (default: Site-To-Site)")
	vpntunnelCreateCmd.Flags().String("protocol", "ikev2", "VPN protocol (default: ikev2)")
	vpntunnelCreateCmd.Flags().String("billing-period", "Hour", "Billing period: Hour, Month, Year")
	vpntunnelCreateCmd.Flags().StringSlice("tags", []string{}, "Tags (comma-separated)")

	// IKE settings
	vpntunnelCreateCmd.Flags().Int32("ike-lifetime", 0, "IKE lifetime (seconds)")
	vpntunnelCreateCmd.Flags().String("ike-encryption", "", "IKE encryption algorithm")
	vpntunnelCreateCmd.Flags().String("ike-hash", "", "IKE hash algorithm")
	vpntunnelCreateCmd.Flags().String("ike-dh-group", "", "IKE DH group")
	vpntunnelCreateCmd.Flags().String("ike-dpd-action", "", "IKE DPD action")
	vpntunnelCreateCmd.Flags().Int32("ike-dpd-interval", 0, "IKE DPD interval (seconds)")
	vpntunnelCreateCmd.Flags().Int32("ike-dpd-timeout", 0, "IKE DPD timeout (seconds)")
	// ESP settings
	vpntunnelCreateCmd.Flags().Int32("esp-lifetime", 0, "ESP lifetime (seconds)")
	vpntunnelCreateCmd.Flags().String("esp-encryption", "", "ESP encryption algorithm")
	vpntunnelCreateCmd.Flags().String("esp-hash", "", "ESP hash algorithm")
	vpntunnelCreateCmd.Flags().String("esp-pfs", "", "ESP PFS group")
	// PSK settings
	vpntunnelCreateCmd.Flags().String("psk-cloud-site", "", "PSK cloud site")
	vpntunnelCreateCmd.Flags().String("psk-onprem-site", "", "PSK on-prem site")
	vpntunnelCreateCmd.Flags().String("psk", "", "Pre-shared key for authentication (PSK secret)")

	vpntunnelListCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	vpntunnelGetCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")

	vpntunnelUpdateCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	vpntunnelUpdateCmd.Flags().String("name", "", "New name for the VPN tunnel")
	vpntunnelUpdateCmd.Flags().StringSlice("tags", []string{}, "New tags (comma-separated)")

	vpntunnelDeleteCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	vpntunnelDeleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	// Set up auto-completion for resource IDs
	elasticipGetCmd.ValidArgsFunction = completeElasticIPID
	elasticipUpdateCmd.ValidArgsFunction = completeElasticIPID
	elasticipDeleteCmd.ValidArgsFunction = completeElasticIPID

	vpcUpdateCmd.ValidArgsFunction = completeVPCID
	vpcDeleteCmd.ValidArgsFunction = completeVPCID

	loadbalancerGetCmd.ValidArgsFunction = completeLoadBalancerID

	vpntunnelGetCmd.ValidArgsFunction = completeVPNTunnelID
	vpntunnelUpdateCmd.ValidArgsFunction = completeVPNTunnelID
	vpntunnelDeleteCmd.ValidArgsFunction = completeVPNTunnelID

	// Subnet
	networkCmd.AddCommand(subnetCmd)
	subnetCmd.AddCommand(subnetCreateCmd)
	subnetCmd.AddCommand(subnetGetCmd)
	subnetCmd.AddCommand(subnetUpdateCmd)
	subnetCmd.AddCommand(subnetDeleteCmd)
	subnetCmd.AddCommand(subnetListCmd)

	subnetCreateCmd.Flags().String("vpc-id", "", "Parent VPC ID (required)")
	subnetCreateCmd.Flags().String("name", "", "Subnet name (required)")
	subnetCreateCmd.Flags().String("cidr", "", "Subnet CIDR (required)")
	subnetUpdateCmd.Flags().String("name", "", "Subnet name (optional)")
	subnetUpdateCmd.Flags().String("cidr", "", "Subnet CIDR (optional)")
	subnetUpdateCmd.Flags().StringSlice("tags", []string{}, "Subnet tags (optional)")
	subnetListCmd.Flags().String("vpc-id", "", "Parent VPC ID (required)")

	// SecurityGroup
	networkCmd.AddCommand(securitygroupCmd)
	securitygroupCmd.AddCommand(securitygroupCreateCmd)
	securitygroupCmd.AddCommand(securitygroupGetCmd)
	securitygroupCmd.AddCommand(securitygroupUpdateCmd)
	securitygroupCmd.AddCommand(securitygroupDeleteCmd)
	securitygroupCmd.AddCommand(securitygroupListCmd)
	// SecurityRule
	securitygroupCmd.AddCommand(securityruleCmd)
	securityruleCmd.AddCommand(securityruleCreateCmd)
	securityruleCmd.AddCommand(securityruleGetCmd)
	securityruleCmd.AddCommand(securityruleUpdateCmd)
	securityruleCmd.AddCommand(securityruleDeleteCmd)
	securityruleCmd.AddCommand(securityruleListCmd)
	// Peering
	networkCmd.AddCommand(vpcpeeringCmd)
	vpcpeeringCmd.AddCommand(peeringCreateCmd)
	vpcpeeringCmd.AddCommand(peeringGetCmd)
	vpcpeeringCmd.AddCommand(peeringUpdateCmd)
	vpcpeeringCmd.AddCommand(peeringDeleteCmd)
	vpcpeeringCmd.AddCommand(peeringListCmd)
	// VPC Peering Route (nested under vpcpeering)
	vpcpeeringCmd.AddCommand(vpcpeeringrouteCmd)
	vpcpeeringrouteCmd.AddCommand(vpcpeeringrouteCreateCmd)
	vpcpeeringrouteCmd.AddCommand(vpcpeeringrouteGetCmd)
	vpcpeeringrouteCmd.AddCommand(vpcpeeringrouteUpdateCmd)
	vpcpeeringrouteCmd.AddCommand(vpcpeeringrouteDeleteCmd)
	vpcpeeringrouteCmd.AddCommand(vpcpeeringrouteListCmd)

	// VPC Peering Route (nested under vpcpeering)
	vpcpeeringCmd.AddCommand(vpcpeeringrouteCmd)
	vpcpeeringrouteCmd.AddCommand(vpcpeeringrouteCreateCmd)
	vpcpeeringrouteCmd.AddCommand(vpcpeeringrouteGetCmd)
	vpcpeeringrouteCmd.AddCommand(vpcpeeringrouteUpdateCmd)
	vpcpeeringrouteCmd.AddCommand(vpcpeeringrouteDeleteCmd)
	vpcpeeringrouteCmd.AddCommand(vpcpeeringrouteListCmd)
	// VPNTunnel
	networkCmd.AddCommand(vpntunnelCmd)
	vpntunnelCmd.AddCommand(vpntunnelCreateCmd)
	vpntunnelCmd.AddCommand(vpntunnelGetCmd)
	vpntunnelCmd.AddCommand(vpntunnelUpdateCmd)
	vpntunnelCmd.AddCommand(vpntunnelDeleteCmd)
	vpntunnelCmd.AddCommand(vpntunnelListCmd)

	// VPN Route (nested under vpntunnel)
	vpntunnelCmd.AddCommand(vpnrouteCmd)
	vpnrouteCmd.AddCommand(vpnrouteCreateCmd)
	vpnrouteCmd.AddCommand(vpnrouteGetCmd)
	vpnrouteCmd.AddCommand(vpnrouteUpdateCmd)
	vpnrouteCmd.AddCommand(vpnrouteDeleteCmd)
	vpnrouteCmd.AddCommand(vpnrouteListCmd)
}
