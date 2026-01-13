/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	ClientID       string `yaml:"clientId"`
	ClientSecret   string `yaml:"clientSecret"`
	BaseURL        string `yaml:"baseUrl,omitempty"`
	TokenIssuerURL string `yaml:"tokenIssuerUrl,omitempty"`
}

const (
	DefaultBaseURL        = "https://api.arubacloud.com"
	DefaultTokenIssuerURL = "https://login.aruba.it/auth/realms/cmp-new-apikey/protocol/openid-connect/token"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage acloud configuration",
	Long:  `Configure acloud with your Aruba Cloud API credentials (clientId and clientSecret).`,
}

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Long:  `Set configuration values for acloud, such as clientId and clientSecret.`,
	Run: func(cmd *cobra.Command, args []string) {
		clientID, _ := cmd.Flags().GetString("client-id")
		clientSecret, _ := cmd.Flags().GetString("client-secret")
		baseURL, _ := cmd.Flags().GetString("base-url")
		tokenIssuerURL, _ := cmd.Flags().GetString("token-issuer-url")

		// Load existing config or create new one
		config, err := LoadConfig()
		if err != nil {
			// If config doesn't exist, create a new one
			config = &Config{}
		}

		// Validate required fields
		// If setting up for the first time, both client-id and client-secret are required
		// If updating, at least one must be provided, but final config must have both
		if config.ClientID == "" && clientID == "" {
			fmt.Println("Error: --client-id is required")
			fmt.Println("Please run: acloud config set --client-id YOUR_CLIENT_ID --client-secret YOUR_CLIENT_SECRET")
			os.Exit(1)
		}
		if config.ClientSecret == "" && clientSecret == "" {
			fmt.Println("Error: --client-secret is required")
			fmt.Println("Please run: acloud config set --client-id YOUR_CLIENT_ID --client-secret YOUR_CLIENT_SECRET")
			os.Exit(1)
		}

		// Update only provided values
		if clientID != "" {
			config.ClientID = clientID
		}
		if clientSecret != "" {
			config.ClientSecret = clientSecret
		}
		if baseURL != "" {
			config.BaseURL = baseURL
		}
		if tokenIssuerURL != "" {
			config.TokenIssuerURL = tokenIssuerURL
		}

		// Final validation: both clientID and clientSecret must be set
		if config.ClientID == "" || config.ClientSecret == "" {
			fmt.Println("Error: Both --client-id and --client-secret are required")
			fmt.Println("Please run: acloud config set --client-id YOUR_CLIENT_ID --client-secret YOUR_CLIENT_SECRET")
			os.Exit(1)
		}

		// Save config
		if err := SaveConfig(config); err != nil {
			fmt.Printf("Error saving configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Configuration updated successfully")
		if clientID != "" {
			fmt.Printf("  Client ID: %s\n", clientID)
		}
		if clientSecret != "" {
			fmt.Println("  Client Secret: ********")
		}
		if baseURL != "" {
			fmt.Printf("  Base URL: %s\n", baseURL)
		}
		if tokenIssuerURL != "" {
			fmt.Printf("  Token Issuer URL: %s\n", tokenIssuerURL)
		}
	},
}

// configShowCmd represents the config show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current acloud configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := LoadConfig()
		if err != nil {
			fmt.Println("No configuration found. Please run 'acloud config set' to create one.")
			return
		}

		fmt.Println("Current configuration:")
		fmt.Printf("  Client ID: %s\n", config.ClientID)
		if config.ClientSecret != "" {
			fmt.Println("  Client Secret: ********")
		} else {
			fmt.Println("  Client Secret: (not set)")
		}
		baseURL := config.BaseURL
		if baseURL == "" {
			baseURL = DefaultBaseURL + " (default)"
		}
		fmt.Printf("  Base URL: %s\n", baseURL)
		tokenIssuerURL := config.TokenIssuerURL
		if tokenIssuerURL == "" {
			tokenIssuerURL = DefaultTokenIssuerURL + " (default)"
		}
		fmt.Printf("  Token Issuer URL: %s\n", tokenIssuerURL)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)

	// Flags for config set command
	configSetCmd.Flags().String("client-id", "", "Aruba Cloud API client ID (required)")
	configSetCmd.Flags().String("client-secret", "", "Aruba Cloud API client secret (required)")
	configSetCmd.Flags().String("base-url", "", "Base URL for Aruba Cloud API (optional, default: https://api.arubacloud.com)")
	configSetCmd.Flags().String("token-issuer-url", "", "Token issuer URL for authentication (optional, default: https://login.aruba.it/auth/realms/cmp-new-apikey/protocol/openid-connect/token)")
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".acloud.yaml"), nil
}

// LoadConfig loads the configuration from the config file
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves the configuration to the config file
func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}
