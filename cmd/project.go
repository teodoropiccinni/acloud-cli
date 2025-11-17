package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var managementCmd = &cobra.Command{
	Use:   "management",
	Short: "Manage organization resources",
	Long:  `Manage organization resources in Aruba Cloud, like projects.`,
}

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
	Long:  `Perform CRUD operations on projects in Aruba Cloud.`,
}

var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Project created (stub)")
	},
}

var projectGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get project details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Project details (stub)")
	},
}

var projectUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Project updated (stub)")
	},
}

var projectDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Project deleted (stub)")
	},
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Project list (stub)")
	},
}

func init() {
	rootCmd.AddCommand(managementCmd)
	managementCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectCreateCmd)
	projectCmd.AddCommand(projectGetCmd)
	projectCmd.AddCommand(projectUpdateCmd)
	projectCmd.AddCommand(projectDeleteCmd)
	projectCmd.AddCommand(projectListCmd)
}
