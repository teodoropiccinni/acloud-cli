package cmd

import (
	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage schedules",
	Long:  `Manage schedules in Aruba Cloud.`,
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
	// job subcommands
	var jobCmd = &cobra.Command{
		Use:   "job",
		Short: "Manage scheduled jobs",
		Long:  `Perform CRUD operations on scheduled jobs in Aruba Cloud.`,
	}

	var jobCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new scheduled job",
		Run: func(cmd *cobra.Command, args []string) {
			println("Scheduled job created (stub)")
		},
	}

	var jobGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get scheduled job details",
		Run: func(cmd *cobra.Command, args []string) {
			println("Scheduled job details (stub)")
		},
	}

	var jobUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update a scheduled job",
		Run: func(cmd *cobra.Command, args []string) {
			println("Scheduled job updated (stub)")
		},
	}

	var jobDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete a scheduled job",
		Run: func(cmd *cobra.Command, args []string) {
			println("Scheduled job deleted (stub)")
		},
	}

	var jobListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all scheduled jobs",
		Run: func(cmd *cobra.Command, args []string) {
			println("Scheduled job list (stub)")
		},
	}

	scheduleCmd.AddCommand(jobCmd)
	jobCmd.AddCommand(jobCreateCmd)
	jobCmd.AddCommand(jobGetCmd)
	jobCmd.AddCommand(jobUpdateCmd)
	jobCmd.AddCommand(jobDeleteCmd)
	jobCmd.AddCommand(jobListCmd)
}
