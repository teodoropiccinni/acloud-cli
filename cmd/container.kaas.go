package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Arubacloud/sdk-go/pkg/types"
	"github.com/spf13/cobra"
)

func init() {
	// KaaS commands
	containerCmd.AddCommand(kaasCmd)
	kaasCmd.AddCommand(kaasCreateCmd)
	kaasCmd.AddCommand(kaasGetCmd)
	kaasCmd.AddCommand(kaasUpdateCmd)
	kaasCmd.AddCommand(kaasDeleteCmd)
	kaasCmd.AddCommand(kaasListCmd)

	// Add flags for KaaS commands
	kaasCreateCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	kaasCreateCmd.Flags().String("name", "", "Name for the KaaS cluster (required)")
	kaasCreateCmd.Flags().String("region", "", "Region code (required)")
	kaasCreateCmd.Flags().String("version", "", "Kubernetes version (optional - check SDK for correct field)")
	kaasCreateCmd.Flags().StringSlice("tags", []string{}, "Tags (comma-separated)")
	kaasCreateCmd.MarkFlagRequired("name")
	kaasCreateCmd.MarkFlagRequired("region")
	// Version may not be a direct field in KaaSPropertiesRequest - check SDK documentation

	kaasGetCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")

	kaasUpdateCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	kaasUpdateCmd.Flags().String("name", "", "New name for the KaaS cluster")
	kaasUpdateCmd.Flags().StringSlice("tags", []string{}, "New tags (comma-separated)")

	kaasDeleteCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")
	kaasDeleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")

	kaasListCmd.Flags().String("project-id", "", "Project ID (uses context if not specified)")

	// Set up auto-completion for resource IDs
	kaasGetCmd.ValidArgsFunction = completeKaaSID
	kaasUpdateCmd.ValidArgsFunction = completeKaaSID
	kaasDeleteCmd.ValidArgsFunction = completeKaaSID
}

// Completion functions for container resources
func completeKaaSID(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	projectID, err := GetProjectID(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	client, err := GetArubaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	ctx := context.Background()
	response, err := client.FromContainer().KaaS().List(ctx, projectID, nil)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string
	if response != nil && response.Data != nil {
		for _, kaas := range response.Data.Values {
			if kaas.Metadata.ID != nil && kaas.Metadata.Name != nil {
				id := *kaas.Metadata.ID
				if toComplete == "" || strings.HasPrefix(id, toComplete) {
					completions = append(completions, fmt.Sprintf("%s\t%s", id, *kaas.Metadata.Name))
				}
			}
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// KaaS subcommands
var kaasCmd = &cobra.Command{
	Use:   "kaas",
	Short: "Manage Kubernetes as a Service (KaaS)",
	Long:  `Perform CRUD operations on KaaS resources in Aruba Cloud.`,
}

var kaasCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new KaaS cluster",
	Run: func(cmd *cobra.Command, args []string) {
		projectID, err := GetProjectID(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		name, _ := cmd.Flags().GetString("name")
		region, _ := cmd.Flags().GetString("region")
		version, _ := cmd.Flags().GetString("version")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		if name == "" || region == "" || version == "" {
			fmt.Println("Error: --name, --region, and --version are required")
			return
		}

		client, err := GetArubaClient()
		if err != nil {
			fmt.Printf("Error initializing client: %v\n", err)
			return
		}

		// Build the create request
		createRequest := types.KaaSRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: name,
					Tags: tags,
				},
				Location: types.LocationRequest{
					Value: region,
				},
			},
			Properties: types.KaaSPropertiesRequest{
				// Note: Version field may not exist in KaaSPropertiesRequest
				// Check SDK documentation for correct way to specify Kubernetes version
			},
		}

		ctx := context.Background()
		response, err := client.FromContainer().KaaS().Create(ctx, projectID, createRequest, nil)
		if err != nil {
			fmt.Printf("Error creating KaaS cluster: %v\n", err)
			return
		}

		if response != nil && response.IsError() && response.Error != nil {
			fmt.Printf("Failed to create KaaS cluster - Status: %d\n", response.StatusCode)
			if response.Error.Title != nil {
				fmt.Printf("Error: %s\n", *response.Error.Title)
			}
			if response.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *response.Error.Detail)
			}
			return
		}

		if response != nil && response.Data != nil {
			headers := []TableColumn{
				{Header: "ID", Width: 30},
				{Header: "NAME", Width: 40},
				{Header: "VERSION", Width: 20},
				{Header: "REGION", Width: 20},
			}
			row := []string{
				func() string {
					if response.Data.Metadata.ID != nil {
						return *response.Data.Metadata.ID
					}
					return ""
				}(),
				func() string {
					if response.Data.Metadata.Name != nil {
						return *response.Data.Metadata.Name
					}
					return ""
				}(),
				"", // Version field may not be available in response
				func() string {
					if response.Data.Metadata.LocationResponse != nil {
						return response.Data.Metadata.LocationResponse.Value
					}
					return ""
				}(),
			}
			PrintTable(headers, [][]string{row})
		} else {
			fmt.Println("KaaS cluster created, but no data returned.")
		}
	},
}

var kaasGetCmd = &cobra.Command{
	Use:   "get [kaas-id]",
	Short: "Get KaaS cluster details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		kaasID := args[0]

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
		resp, err := client.FromContainer().KaaS().Get(ctx, projectID, kaasID, nil)
		if err != nil {
			fmt.Printf("Error getting KaaS cluster: %v\n", err)
			return
		}

		if resp != nil && resp.IsError() && resp.Error != nil {
			fmt.Printf("Failed to get KaaS cluster - Status: %d\n", resp.StatusCode)
			if resp.Error.Title != nil {
				fmt.Printf("Error: %s\n", *resp.Error.Title)
			}
			if resp.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *resp.Error.Detail)
			}
			return
		}

		if resp != nil && resp.Data != nil {
			kaas := resp.Data

			fmt.Println("\nKaaS Cluster Details:")
			fmt.Println("====================")

			if kaas.Metadata.ID != nil {
				fmt.Printf("ID:              %s\n", *kaas.Metadata.ID)
			}
			if kaas.Metadata.URI != nil {
				fmt.Printf("URI:             %s\n", *kaas.Metadata.URI)
			}
			if kaas.Metadata.Name != nil {
				fmt.Printf("Name:            %s\n", *kaas.Metadata.Name)
			}
			if kaas.Metadata.LocationResponse != nil {
				fmt.Printf("Region:          %s\n", kaas.Metadata.LocationResponse.Value)
			}
			// Version field may not be available in KaaSPropertiesResponse
			// Check SDK documentation for correct field name
			if kaas.Status.State != nil {
				fmt.Printf("Status:          %s\n", *kaas.Status.State)
			}

			if !kaas.Metadata.CreationDate.IsZero() {
				fmt.Printf("Creation Date:   %s\n", kaas.Metadata.CreationDate.Format("02-01-2006 15:04:05"))
			}
			if kaas.Metadata.CreatedBy != nil {
				fmt.Printf("Created By:      %s\n", *kaas.Metadata.CreatedBy)
			}

			if len(kaas.Metadata.Tags) > 0 {
				fmt.Printf("Tags:            %v\n", kaas.Metadata.Tags)
			} else {
				fmt.Printf("Tags:            []\n")
			}

			// Show JSON output if verbose
			verbose, _ := cmd.Flags().GetBool("verbose")
			if verbose {
				jsonData, _ := json.MarshalIndent(kaas, "", "  ")
				fmt.Println("\nFull JSON Response:")
				fmt.Println("==================")
				fmt.Println(string(jsonData))
			}
		} else {
			fmt.Println("KaaS cluster not found or no data returned.")
		}
	},
}

var kaasUpdateCmd = &cobra.Command{
	Use:   "update [kaas-id]",
	Short: "Update a KaaS cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		kaasID := args[0]

		name, _ := cmd.Flags().GetString("name")
		tags, _ := cmd.Flags().GetStringSlice("tags")

		// Note: The new SDK only supports updating properties, not metadata (name/tags)
		if name != "" || len(tags) > 0 {
			fmt.Println("Warning: Updating name and tags is not supported in the current SDK version.")
			fmt.Println("The update will only modify cluster properties (Kubernetes version, node pools, etc.).")
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
		getResponse, err := client.FromContainer().KaaS().Get(ctx, projectID, kaasID, nil)
		if err != nil {
			fmt.Printf("Error getting KaaS cluster details: %v\n", err)
			return
		}

		if getResponse == nil || getResponse.Data == nil {
			fmt.Println("Error: KaaS cluster not found")
			return
		}

		current := getResponse.Data

		// Build the update request using KaaSUpdateRequest (only supports properties)
		// Convert from response types to update request types
		updateRequest := types.KaaSUpdateRequest{
			Properties: types.KaaSPropertiesUpdateRequest{
				// KubernetesVersion is required
				KubernetesVersion: types.KubernetesVersionInfoUpdate{
					Value: func() string {
						if current.Properties.KubernetesVersion.Value != nil {
							return *current.Properties.KubernetesVersion.Value
						}
						return ""
					}(),
				},
				// NodePools is required - convert from response to request format
				NodePools: func() []types.NodePoolProperties {
					if current.Properties.NodePools == nil {
						return []types.NodePoolProperties{}
					}
					var nodePools []types.NodePoolProperties
					for _, np := range *current.Properties.NodePools {
						nodePool := types.NodePoolProperties{
							Name: func() string {
								if np.Name != nil {
									return *np.Name
								} else {
									return ""
								}
							}(),
							Nodes: func() int32 {
								if np.Nodes != nil {
									return *np.Nodes
								} else {
									return 0
								}
							}(),
							Autoscaling: np.Autoscaling,
							MinCount:    np.MinCount,
							MaxCount:    np.MaxCount,
						}
						// Instance - convert from InstanceResponse to string
						if np.Instance != nil && np.Instance.Name != nil {
							nodePool.Instance = *np.Instance.Name
						}
						// DataCenter - convert from DataCenterResponse to string
						if np.DataCenter != nil && np.DataCenter.Code != nil {
							nodePool.Zone = *np.DataCenter.Code
						}
						nodePools = append(nodePools, nodePool)
					}
					return nodePools
				}(),
				// Optional fields
				HA:      current.Properties.HA,
				Storage: current.Properties.Storage,
				BillingPlan: func() *types.BillingPeriodResource {
					if current.Properties.BillingPlan != nil && current.Properties.BillingPlan.BillingPeriod != nil {
						return &types.BillingPeriodResource{
							BillingPeriod: *current.Properties.BillingPlan.BillingPeriod,
						}
					}
					return nil
				}(),
			},
		}

		response, err := client.FromContainer().KaaS().Update(ctx, projectID, kaasID, updateRequest, nil)
		if err != nil {
			fmt.Printf("Error updating KaaS cluster: %v\n", err)
			return
		}

		if response != nil && response.IsError() && response.Error != nil {
			fmt.Printf("Failed to update KaaS cluster - Status: %d\n", response.StatusCode)
			if response.Error.Title != nil {
				fmt.Printf("Error: %s\n", *response.Error.Title)
			}
			if response.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *response.Error.Detail)
			}
			return
		}

		if response != nil && response.Data != nil {
			fmt.Println("\nKaaS cluster updated successfully!")
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
			fmt.Println("KaaS cluster update initiated. Use 'get' to check status.")
		}
	},
}

var kaasDeleteCmd = &cobra.Command{
	Use:   "delete [kaas-id]",
	Short: "Delete a KaaS cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		kaasID := args[0]

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

		// Confirmation prompt
		skipConfirm, _ := cmd.Flags().GetBool("yes")
		if !skipConfirm {
			fmt.Printf("Are you sure you want to delete KaaS cluster '%s'? (yes/no): ", kaasID)
			var confirmation string
			fmt.Scanln(&confirmation)
			if confirmation != "yes" && confirmation != "y" {
				fmt.Println("Deletion cancelled.")
				return
			}
		}

		ctx := context.Background()
		response, err := client.FromContainer().KaaS().Delete(ctx, projectID, kaasID, nil)
		if err != nil {
			fmt.Printf("Error deleting KaaS cluster: %v\n", err)
			return
		}

		if response != nil && response.IsError() && response.Error != nil {
			fmt.Printf("Failed to delete KaaS cluster - Status: %d\n", response.StatusCode)
			if response.Error.Title != nil {
				fmt.Printf("Error: %s\n", *response.Error.Title)
			}
			if response.Error.Detail != nil {
				fmt.Printf("Detail: %s\n", *response.Error.Detail)
			}
			return
		}

		fmt.Printf("KaaS cluster '%s' deleted successfully.\n", kaasID)
	},
}

var kaasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all KaaS clusters",
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
		response, err := client.FromContainer().KaaS().List(ctx, projectID, nil)
		if err != nil {
			fmt.Printf("Error listing KaaS clusters: %v\n", err)
			return
		}

		if response != nil && response.Data != nil && len(response.Data.Values) > 0 {
			headers := []TableColumn{
				{Header: "ID", Width: 30},
				{Header: "NAME", Width: 40},
				{Header: "VERSION", Width: 20},
				{Header: "REGION", Width: 20},
				{Header: "STATUS", Width: 15},
			}

			var rows [][]string
			for _, kaas := range response.Data.Values {
				id := ""
				if kaas.Metadata.ID != nil {
					id = *kaas.Metadata.ID
				}
				name := ""
				if kaas.Metadata.Name != nil {
					name = *kaas.Metadata.Name
				}
				version := "" // Version field may not be available
				region := ""
				if kaas.Metadata.LocationResponse != nil {
					region = kaas.Metadata.LocationResponse.Value
				}
				status := ""
				if kaas.Status.State != nil {
					status = *kaas.Status.State
				}

				rows = append(rows, []string{
					id,
					name,
					version,
					region,
					status,
				})
			}

			PrintTable(headers, rows)
		} else {
			fmt.Println("No KaaS clusters found")
		}
	},
}
