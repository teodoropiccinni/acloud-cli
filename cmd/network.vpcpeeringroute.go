package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {

	networkCmd.AddCommand(vpcpeeringrouteCmd)

	vpcpeeringrouteCmd.AddCommand(vpcpeeringrouteCreateCmd)
	vpcpeeringrouteCmd.AddCommand(vpcpeeringrouteGetCmd)
	vpcpeeringrouteCmd.AddCommand(vpcpeeringrouteUpdateCmd)
	vpcpeeringrouteCmd.AddCommand(vpcpeeringrouteDeleteCmd)
	vpcpeeringrouteCmd.AddCommand(vpcpeeringrouteListCmd)

}

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
