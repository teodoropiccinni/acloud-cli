package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {

	// SecurityRule
	networkCmd.AddCommand(securityruleCmd)

	securityruleCmd.AddCommand(securityruleCreateCmd)
	securityruleCmd.AddCommand(securityruleGetCmd)
	securityruleCmd.AddCommand(securityruleUpdateCmd)
	securityruleCmd.AddCommand(securityruleDeleteCmd)
	securityruleCmd.AddCommand(securityruleListCmd)
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
