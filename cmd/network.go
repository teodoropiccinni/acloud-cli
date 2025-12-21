package cmd

import (
	"context"
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

				address := ""
				if eip.Properties.Address != nil {
					address = *eip.Properties.Address
				}

				status := ""
				if eip.Status.State != nil {
					status = *eip.Status.State
				}

				rows = append(rows, []string{name, id, address, status})
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

				address := ""
				if lb.Properties.Address != nil {
					address = *lb.Properties.Address
				}

				status := ""
				if lb.Status.State != nil {
					status = *lb.Status.State
				}

				rows = append(rows, []string{name, id, address, status})
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

var vpcGetCmd = &cobra.Command{
	Use:   "get <vpc-id>",
	Short: "Get VPC details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vpcID := args[0]

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

		// Get VPC details using the SDK
		ctx := context.Background()
		response, err := client.FromNetwork().VPCs().Get(ctx, projectID, vpcID, nil)
		if err != nil {
			fmt.Printf("Error getting VPC details: %v\n", err)
			return
		}

		if response != nil && response.Data != nil {
			vpc := response.Data

			// Display VPC details
			fmt.Println("\nVPC Details:")
			fmt.Println("============")

			if vpc.Metadata.ID != nil {
				fmt.Printf("ID:              %s\n", *vpc.Metadata.ID)
			}
			if vpc.Metadata.URI != nil {
				fmt.Printf("URI:             %s\n", *vpc.Metadata.URI)
			}
			if vpc.Metadata.Name != nil {
				fmt.Printf("Name:            %s\n", *vpc.Metadata.Name)
			}

			fmt.Printf("Default:         %t\n", vpc.Properties.Default)
			fmt.Printf("Linked Resources: %d\n", len(vpc.Properties.LinkedResources))

			if vpc.Metadata.CreationDate != nil {
				fmt.Printf("Creation Date:   %s\n", vpc.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}
			if vpc.Metadata.CreatedBy != nil {
				fmt.Printf("Created By:      %s\n", *vpc.Metadata.CreatedBy)
			}

			if len(vpc.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", vpc.Metadata.Tags)
			}

			if vpc.Status.State != nil {
				fmt.Printf("Status:          %s\n", *vpc.Status.State)
			}
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

				subnets := fmt.Sprintf("%d", len(vpc.Properties.LinkedResources))

				status := ""
				if vpc.Status.State != nil {
					status = *vpc.Status.State
				}

				rows = append(rows, []string{name, id, subnets, status})
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
	Use:   "create",
	Short: "Create a new subnet",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Subnet created (stub)")
	},
}

var subnetGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get subnet details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Subnet details (stub)")
	},
}

var subnetUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a subnet",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Subnet updated (stub)")
	},
}

var subnetDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a subnet",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Subnet deleted (stub)")
	},
}

var subnetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all subnets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Subnet list (stub)")
	},
}

// SecurityGroup subcommands
var securitygroupCmd = &cobra.Command{
	Use:   "securitygroup",
	Short: "Manage security groups",
	Long:  `Perform CRUD operations on security groups in Aruba Cloud.`,
}

var securitygroupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new security group",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Security group created (stub)")
	},
}

var securitygroupGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get security group details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Security group details (stub)")
	},
}

var securitygroupUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a security group",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Security group updated (stub)")
	},
}

var securitygroupDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a security group",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Security group deleted (stub)")
	},
}

var securitygroupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all security groups",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Security group list (stub)")
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
var peeringCmd = &cobra.Command{
	Use:   "peering",
	Short: "Manage VPC peering",
	Long:  `Perform CRUD operations on VPC peering in Aruba Cloud.`,
}

var peeringCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new VPC peering",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPC peering created (stub)")
	},
}

var peeringGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get VPC peering details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPC peering details (stub)")
	},
}

var peeringUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a VPC peering",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPC peering updated (stub)")
	},
}

var peeringDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a VPC peering",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPC peering deleted (stub)")
	},
}

var peeringListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VPC peerings",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPC peering list (stub)")
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

var vpntunnelCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new VPN tunnel",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPN tunnel created (stub)")
	},
}

var vpntunnelGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get VPN tunnel details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPN tunnel details (stub)")
	},
}

var vpntunnelUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a VPN tunnel",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPN tunnel updated (stub)")
	},
}

var vpntunnelDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a VPN tunnel",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPN tunnel deleted (stub)")
	},
}

var vpntunnelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VPN tunnels",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPN tunnel list (stub)")
	},
}

// Route subcommands under VPNTunnel
var vpntunnelRouteCmd = &cobra.Command{
	Use:   "route",
	Short: "Manage VPN tunnel routes",
	Long:  `Perform CRUD operations on VPN tunnel routes in Aruba Cloud.`,
}

var vpntunnelRouteCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new VPN tunnel route",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPN tunnel route created (stub)")
	},
}

var vpntunnelRouteGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get VPN tunnel route details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPN tunnel route details (stub)")
	},
}

var vpntunnelRouteUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a VPN tunnel route",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPN tunnel route updated (stub)")
	},
}

var vpntunnelRouteDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a VPN tunnel route",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPN tunnel route deleted (stub)")
	},
}

var vpntunnelRouteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VPN tunnel routes",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VPN tunnel route list (stub)")
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
	vpcCmd.AddCommand(vpcGetCmd)
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

	vpcGetCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")

	vpcUpdateCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	vpcUpdateCmd.Flags().String("name", "", "New name for the VPC")
	vpcUpdateCmd.Flags().StringSlice("tags", []string{}, "New tags (comma-separated)")

	vpcDeleteCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	vpcDeleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	// Add flags for vpc commands
	vpcListCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")

	// Set up auto-completion for resource IDs
	elasticipGetCmd.ValidArgsFunction = completeElasticIPID
	elasticipUpdateCmd.ValidArgsFunction = completeElasticIPID
	elasticipDeleteCmd.ValidArgsFunction = completeElasticIPID

	vpcGetCmd.ValidArgsFunction = completeVPCID
	vpcUpdateCmd.ValidArgsFunction = completeVPCID
	vpcDeleteCmd.ValidArgsFunction = completeVPCID

	loadbalancerGetCmd.ValidArgsFunction = completeLoadBalancerID

	// Subnet
	vpcCmd.AddCommand(subnetCmd)
	subnetCmd.AddCommand(subnetCreateCmd)
	subnetCmd.AddCommand(subnetGetCmd)
	subnetCmd.AddCommand(subnetUpdateCmd)
	subnetCmd.AddCommand(subnetDeleteCmd)
	subnetCmd.AddCommand(subnetListCmd)
	// SecurityGroup
	vpcCmd.AddCommand(securitygroupCmd)
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
	vpcCmd.AddCommand(peeringCmd)
	peeringCmd.AddCommand(peeringCreateCmd)
	peeringCmd.AddCommand(peeringGetCmd)
	peeringCmd.AddCommand(peeringUpdateCmd)
	peeringCmd.AddCommand(peeringDeleteCmd)
	peeringCmd.AddCommand(peeringListCmd)
	// Route
	peeringCmd.AddCommand(routeCmd)
	routeCmd.AddCommand(routeCreateCmd)
	routeCmd.AddCommand(routeGetCmd)
	routeCmd.AddCommand(routeUpdateCmd)
	routeCmd.AddCommand(routeDeleteCmd)
	routeCmd.AddCommand(routeListCmd)
	// VPNTunnel
	networkCmd.AddCommand(vpntunnelCmd)
	vpntunnelCmd.AddCommand(vpntunnelCreateCmd)
	vpntunnelCmd.AddCommand(vpntunnelGetCmd)
	vpntunnelCmd.AddCommand(vpntunnelUpdateCmd)
	vpntunnelCmd.AddCommand(vpntunnelDeleteCmd)
	vpntunnelCmd.AddCommand(vpntunnelListCmd)
	// Route under VPNTunnel
	vpntunnelCmd.AddCommand(vpntunnelRouteCmd)
	vpntunnelRouteCmd.AddCommand(vpntunnelRouteCreateCmd)
	vpntunnelRouteCmd.AddCommand(vpntunnelRouteGetCmd)
	vpntunnelRouteCmd.AddCommand(vpntunnelRouteUpdateCmd)
	vpntunnelRouteCmd.AddCommand(vpntunnelRouteDeleteCmd)
	vpntunnelRouteCmd.AddCommand(vpntunnelRouteListCmd)
}
