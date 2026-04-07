package cmd

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"

    "github.com/spf13/cobra"
    "golang.org/x/oauth2/clientcredentials"
    "gopkg.in/yaml.v3"
)

type folderCreateRequest struct {
    Name    string `json:"name"`
    Default bool   `json:"default"`
}

func init() {
    managementCmd.AddCommand(folderCmd)
    folderCmd.AddCommand(folderCreateCmd)

    folderCreateCmd.Flags().String("name", "", "Name for the folder (required)")
    folderCreateCmd.Flags().Bool("default", false, "Set as default folder")
    folderCreateCmd.Flags().String("api-version", "1", "API version for folder endpoint")
    folderCreateCmd.Flags().String("format", "table", "Output format (table, json, yaml)")
    folderCreateCmd.MarkFlagRequired("name")
}

var folderCmd = &cobra.Command{
    Use:   "folder",
    Short: "Manage folders",
    Long:  `Manage folders in Aruba Cloud governance APIs.`,
}

var folderCreateCmd = &cobra.Command{
    Use:   "create",
    Short: "Create a new folder",
    Run: func(cmd *cobra.Command, args []string) {
        folderName, _ := cmd.Flags().GetString("name")
        setDefault, _ := cmd.Flags().GetBool("default")
        apiVersion, _ := cmd.Flags().GetString("api-version")
        outputFormat, _ := cmd.Flags().GetString("format")
        outputFormat = strings.ToLower(strings.TrimSpace(outputFormat))

        if folderName == "" {
            fmt.Println("Error: --name is required")
            return
        }

        config, err := LoadConfig()
        if err != nil {
            fmt.Printf("Error loading configuration: %v\n", err)
            return
        }

        if config.ClientID == "" || config.ClientSecret == "" {
            fmt.Println("Error: client ID or client secret not configured. Please run 'acloud config set --client-id YOUR_CLIENT_ID --client-secret YOUR_CLIENT_SECRET'")
            return
        }

        baseURL := config.BaseURL
        if baseURL == "" {
            baseURL = DefaultBaseURL
        }

        tokenIssuerURL := config.TokenIssuerURL
        if tokenIssuerURL == "" {
            tokenIssuerURL = DefaultTokenIssuerURL
        }

        tokenConfig := clientcredentials.Config{
            ClientID:     config.ClientID,
            ClientSecret: config.ClientSecret,
            TokenURL:     tokenIssuerURL,
        }

        token, err := tokenConfig.Token(context.Background())
        if err != nil {
            fmt.Printf("Error getting access token: %v\n", err)
            return
        }

        payload := folderCreateRequest{
            Name:    folderName,
            Default: setDefault,
        }

        bodyBytes, err := json.Marshal(payload)
        if err != nil {
            fmt.Printf("Error preparing request body: %v\n", err)
            return
        }

        endpoint, err := url.Parse(strings.TrimRight(baseURL, "/") + "/folders")
        if err != nil {
            fmt.Printf("Error building folder endpoint: %v\n", err)
            return
        }

        q := endpoint.Query()
        q.Set("api-version", apiVersion)
        endpoint.RawQuery = q.Encode()

        req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint.String(), bytes.NewReader(bodyBytes))
        if err != nil {
            fmt.Printf("Error creating request: %v\n", err)
            return
        }

        req.Header.Set("Authorization", "Bearer "+token.AccessToken)
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Accept", "application/json")

        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            fmt.Printf("Error creating folder: %v\n", err)
            return
        }
        defer resp.Body.Close()

        respBody, _ := io.ReadAll(resp.Body)

        if resp.StatusCode >= 200 && resp.StatusCode < 300 {
            output := map[string]any{
                "status":     "created",
                "name":       folderName,
                "default":    setDefault,
                "apiVersion": apiVersion,
            }

            if len(respBody) > 0 {
                var responseBody any
                if err := json.Unmarshal(respBody, &responseBody); err == nil {
                    output["response"] = responseBody
                } else {
                    output["response"] = string(respBody)
                }
            }

            switch outputFormat {
            case "json":
                b, err := json.MarshalIndent(output, "", "  ")
                if err != nil {
                    fmt.Printf("Error serializing JSON output: %v\n", err)
                    return
                }
                fmt.Println(string(b))
            case "yaml":
                b, err := yaml.Marshal(output)
                if err != nil {
                    fmt.Printf("Error serializing YAML output: %v\n", err)
                    return
                }
                fmt.Print(string(b))
            case "table":
                headers := []TableColumn{
                    {Header: "NAME", Width: 40},
                    {Header: "DEFAULT", Width: 10},
                    {Header: "STATUS", Width: 12},
                }
                defaultValue := "No"
                if setDefault {
                    defaultValue = "Yes"
                }
                PrintTable(headers, [][]string{{folderName, defaultValue, "created"}})
            default:
                fmt.Printf("Error: unsupported format '%s'. Use one of: table, json, yaml\n", outputFormat)
                return
            }

            return
        }

        if len(respBody) > 0 {
            fmt.Printf("Failed to create folder (status %d): %s\n", resp.StatusCode, string(respBody))
            return
        }

        fmt.Printf("Failed to create folder (status %d).\n", resp.StatusCode)
    },
}
