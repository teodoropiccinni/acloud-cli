package formatter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	FormatTable = "table"
	FormatJSON  = "json"
	FormatYAML  = "yaml"
)

// AddFormatFlag attaches a reusable --format flag to a command.
func AddFormatFlag(cmd *cobra.Command, defaultFormat string) {
	if strings.TrimSpace(defaultFormat) == "" {
		defaultFormat = FormatTable
	}
	cmd.Flags().String("format", defaultFormat, "Output format (table, json, yaml)")
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
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling to JSON: %v", err)
		}
		fmt.Println(string(jsonData))
		return nil
	case FormatYAML:
		yamlData, err := yaml.Marshal(data)
		if err != nil {
			return fmt.Errorf("error marshaling to YAML: %v", err)
		}
		fmt.Print(string(yamlData))
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
