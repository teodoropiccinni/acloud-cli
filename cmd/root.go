package cmd

import (
	"fmt"
	"os"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "acloud",
	Short: "CLI for Aruba Cloud APIs",
	Long: `acloud is a command-line interface for interacting with Aruba Cloud APIs.
It provides a simple and intuitive way to manage your Aruba Cloud resources
directly from your terminal.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.acloud.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// GetArubaClient creates and returns an Aruba Cloud SDK client using stored credentials
func GetArubaClient() (aruba.Client, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w. Please run 'acloud config set' to configure credentials", err)
	}

	if config.ClientID == "" || config.ClientSecret == "" {
		return nil, fmt.Errorf("client ID or client secret not configured. Please run 'acloud config set'")
	}

	// Create SDK client with credentials using DefaultOptions
	options := aruba.DefaultOptions(config.ClientID, config.ClientSecret)
	client, err := aruba.NewClient(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create Aruba Cloud client: %w", err)
	}

	return client, nil
}

// TableColumn represents a column definition for the table printer
type TableColumn struct {
	Header string // Column header name
	Width  int    // Column width for formatting
}

// PrintTable prints data in a formatted table with headers
// headers: slice of TableColumn defining each column
// rows: slice of string slices, each inner slice represents a row
func PrintTable(headers []TableColumn, rows [][]string) {
	// Print header row
	formatStr := ""
	headerValues := make([]interface{}, len(headers))
	for i, col := range headers {
		formatStr += fmt.Sprintf("%%-%ds ", col.Width)
		headerValues[i] = col.Header
	}
	formatStr += "\n"
	fmt.Printf(formatStr, headerValues...)

	// Print data rows
	for _, row := range rows {
		rowValues := make([]interface{}, len(row))
		for i, val := range row {
			// Truncate if value is too long
			if len(headers) > i && len(val) > headers[i].Width {
				val = val[:headers[i].Width-3] + "..."
			}
			rowValues[i] = val
		}
		fmt.Printf(formatStr, rowValues...)
	}
}
