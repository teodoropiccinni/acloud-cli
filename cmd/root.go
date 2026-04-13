package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Arubacloud/sdk-go/pkg/aruba"
	"github.com/drewstinnett/gout/v2"
	"github.com/drewstinnett/gout/v2/formats"
	"github.com/spf13/cobra"
)

var (
	// Cached client and its configuration
	clientCache       aruba.Client
	clientCacheLock   sync.Mutex
	cachedClientID    string
	cachedSecret      string
	cachedDebug       bool
	cachedBaseURL     string
	cachedTokenIssuer string
	outputFormat      string
)

const (
	FormatTable = "table"
	FormatJSON  = "json"
	FormatYAML  = "yaml"
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

	// Add global debug flag
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug logging (shows HTTP requests/responses)")
	AddPersistentFormatFlag(rootCmd, FormatTable)
}

// AddPersistentFormatFlag attaches a reusable persistent --format flag to a command.
func AddPersistentFormatFlag(cmd *cobra.Command, defaultFormat string) {
	if strings.TrimSpace(defaultFormat) == "" {
		defaultFormat = FormatTable
	}
	cmd.PersistentFlags().StringVar(&outputFormat, "format", defaultFormat, "Output format (table, json, yaml)")
}

// GetOutputFormat reads, normalizes, and validates the --format flag.
func GetOutputFormat(cmd *cobra.Command) (string, error) {
	format, _ := cmd.Flags().GetString("format")
	format = strings.ToLower(strings.TrimSpace(format))
	if format == "" {
		format = FormatTable
	}

	switch format {
	case FormatTable, FormatJSON, FormatYAML:
		return format, nil
	default:
		return "", fmt.Errorf("unknown format: %s. Supported formats: table, json, yaml", format)
	}
}

// RenderOutput prints the provided data according to the selected format.
// tablePrinter is called only when format=table.
func RenderOutput(format string, data any, tablePrinter func()) error {
	switch format {
	case FormatJSON:
		creator, ok := formats.Formats[FormatJSON]
		if !ok {
			return fmt.Errorf("json formatter not registered in gout")
		}
		var out bytes.Buffer
		g := gout.New(gout.WithWriter(&out), gout.WithFormatter(creator()))
		err := g.Print(data)
		if err != nil {
			return fmt.Errorf("error formatting JSON output: %v", err)
		}
		fmt.Print(out.String())
		return nil
	case FormatYAML:
		creator, ok := formats.Formats[FormatYAML]
		if !ok {
			return fmt.Errorf("yaml formatter not registered in gout")
		}
		var out bytes.Buffer
		g := gout.New(gout.WithWriter(&out), gout.WithFormatter(creator()))
		err := g.Print(data)
		if err != nil {
			return fmt.Errorf("error formatting YAML output: %v", err)
		}
		fmt.Print(out.String())
		return nil
	case FormatTable:
		if tablePrinter == nil {
			return fmt.Errorf("table output requires a table printer")
		}
		tablePrinter()
		return nil
	default:
		return fmt.Errorf("unknown format: %s. Supported formats: table, json, yaml", format)
	}
}

// GetArubaClient creates and returns an Aruba Cloud SDK client using stored credentials
// It automatically checks for the --debug flag from the root command to enable verbose logging
// The client is cached to avoid recreating it on every call, but is invalidated if credentials or debug flag change
func GetArubaClient() (aruba.Client, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w. Please run 'acloud config set' to configure credentials", err)
	}

	if config.ClientID == "" || config.ClientSecret == "" {
		return nil, fmt.Errorf("client ID or client secret not configured. Please run 'acloud config set --client-id YOUR_CLIENT_ID --client-secret YOUR_CLIENT_SECRET'")
	}

	// Get base URL and token issuer URL from config, or use defaults
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	tokenIssuerURL := config.TokenIssuerURL
	if tokenIssuerURL == "" {
		tokenIssuerURL = DefaultTokenIssuerURL
	}

	// Check if debug flag is set on root command
	debugEnabled := false
	if rootCmd != nil {
		debugEnabled, _ = rootCmd.PersistentFlags().GetBool("debug")
	}

	// Check if we can reuse the cached client
	clientCacheLock.Lock()
	defer clientCacheLock.Unlock()

	// Reuse cached client if credentials, URLs, and debug flag haven't changed
	if clientCache != nil &&
		cachedClientID == config.ClientID &&
		cachedSecret == config.ClientSecret &&
		cachedDebug == debugEnabled &&
		cachedBaseURL == baseURL &&
		cachedTokenIssuer == tokenIssuerURL {
		return clientCache, nil
	}

	// Create SDK client with credentials using DefaultOptions
	// Only enable native logger when debug is enabled
	options := aruba.DefaultOptions(config.ClientID, config.ClientSecret)

	// Apply base URL if provided
	if baseURL != "" {
		options = options.WithBaseURL(baseURL)
	}

	// Apply token issuer URL if provided
	if tokenIssuerURL != "" {
		options = options.WithTokenIssuerURL(tokenIssuerURL)
	}

	if debugEnabled {
		options = options.WithNativeLogger()
		// Configure logger output for debug mode
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
		log.SetPrefix("[ArubaSDK] ")
	} else {
		// Ensure logger is disabled when debug is off
		options = options.WithDefaultLogger()
	}

	client, err := aruba.NewClient(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create Aruba Cloud client: %w", err)
	}

	// Cache the client and its configuration
	clientCache = client
	cachedClientID = config.ClientID
	cachedSecret = config.ClientSecret
	cachedDebug = debugEnabled
	cachedBaseURL = baseURL
	cachedTokenIssuer = tokenIssuerURL

	return client, nil
}

// GetProjectID returns the project ID from the flag or current context
func GetProjectID(cmd *cobra.Command) (string, error) {
	// Try to get from flag first
	projectID, _ := cmd.Flags().GetString("project-id")
	if projectID != "" {
		return projectID, nil
	}

	// Try to get from context
	projectID, err := GetCurrentProjectID()
	if err != nil {
		return "", fmt.Errorf("project ID not specified. Use --project-id flag or set a context with 'acloud context use <name>'")
	}

	return projectID, nil
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
	format, err := GetOutputFormat(rootCmd)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data := make([]map[string]string, 0, len(rows))
	for _, row := range rows {
		entry := map[string]string{}
		for i, col := range headers {
			key := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(col.Header), " ", "_"))
			value := ""
			if i < len(row) {
				value = row[i]
			}
			entry[key] = value
		}
		data = append(data, entry)
	}

	err = RenderOutput(format, data, func() {
		formatStr := ""
		headerValues := make([]interface{}, len(headers))
		for i, col := range headers {
			formatStr += fmt.Sprintf("%%-%ds ", col.Width)
			headerValues[i] = col.Header
		}
		formatStr += "\n"
		fmt.Printf(formatStr, headerValues...)

		for _, row := range rows {
			rowValues := make([]interface{}, len(row))
			for i, val := range row {
				if len(headers) > i && len(val) > headers[i].Width {
					val = val[:headers[i].Width-3] + "..."
				}
				rowValues[i] = val
			}
			fmt.Printf(formatStr, rowValues...)
		}
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}
