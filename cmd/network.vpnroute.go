package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {

	networkCmd.AddCommand(vpnrouteCmd)

	vpnrouteCmd.AddCommand(vpnrouteCreateCmd)
	vpnrouteCmd.AddCommand(vpnrouteGetCmd)
	vpnrouteCmd.AddCommand(vpnrouteUpdateCmd)
	vpnrouteCmd.AddCommand(vpnrouteDeleteCmd)
	vpnrouteCmd.AddCommand(vpnrouteListCmd)
}

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
