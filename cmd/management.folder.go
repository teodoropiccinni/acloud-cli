package cmd

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "strings"

    "github.com/spf13/cobra"
    "golang.org/x/oauth2/clientcredentials"
    "gopkg.in/yaml.v3"
)

type folderCreateRequest struct {
    Name    string `json:"name,omitempty"`
    Default bool   `json:"default,omitempty"`
}

type folderResponse struct {
    ID           string `json:"id"`
    Name         string `json:"name"`
    ProjectCount int64  `json:"projectCount"`
    Default      bool   `json:"default"`
}

type folderListResponse struct {
    Total  int64            `json:"total"`
    Values []folderResponse `json:"values"`
}

type folderMetadataResponse struct {
    BaremetalTypologies []string `json:"baremetalTypologies"`
    CmpLegacyTypologies []string `json:"cmpLegacyTypologies"`
}

func init() {
    managementCmd.AddCommand(folderCmd)
    folderCmd.AddCommand(folderCreateCmd)
    folderCmd.AddCommand(folderGetCmd)
    folderCmd.AddCommand(folderMetadataCmd)
    folderCmd.AddCommand(folderUpdateCmd)
    folderCmd.AddCommand(folderDeleteCmd)
    folderCmd.AddCommand(folderListCmd)

    folderCreateCmd.Flags().String("name", "", "Name for the folder (required)")
    folderCreateCmd.Flags().Bool("default", false, "Set as default folder")
    folderCreateCmd.Flags().String("api-version", "1", "API version for folder endpoint")
    folderCreateCmd.Flags().String("format", "table", "Output format (table, json, yaml)")
    folderCreateCmd.MarkFlagRequired("name")

    folderGetCmd.Flags().String("api-version", "1", "API version for folder endpoint")
    folderGetCmd.Flags().String("format", "table", "Output format (table, json, yaml)")

    folderMetadataCmd.Flags().String("api-version", "1", "API version for folder endpoint")
    folderMetadataCmd.Flags().String("format", "table", "Output format (table, json, yaml)")

    folderUpdateCmd.Flags().String("name", "", "New folder name")
    folderUpdateCmd.Flags().Bool("default", false, "Set as default folder")
    folderUpdateCmd.Flags().String("api-version", "1", "API version for folder endpoint")
    folderUpdateCmd.Flags().String("format", "table", "Output format (table, json, yaml)")

    folderDeleteCmd.Flags().String("api-version", "1", "API version for folder endpoint")
    folderDeleteCmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompt")
    folderDeleteCmd.Flags().String("format", "table", "Output format (table, json, yaml)")

    folderListCmd.Flags().String("filter", "", "Filter expression")
    folderListCmd.Flags().String("sort", "", "Sort expression")
    folderListCmd.Flags().String("projection", "", "Projection expression")
    folderListCmd.Flags().Int("offset", -1, "Pagination offset")
    folderListCmd.Flags().Int("limit", -1, "Pagination limit")
    folderListCmd.Flags().String("api-version", "1", "API version for folder endpoint")
    folderListCmd.Flags().String("format", "table", "Output format (table, json, yaml)")
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

        payload := folderCreateRequest{Name: folderName, Default: setDefault}
        statusCode, respBody, err := doFolderRequest(http.MethodPost, "/folders", apiVersion, nil, payload)
        if err != nil {
            fmt.Printf("Error creating folder: %v\n", err)
            return
        }

        if statusCode < 200 || statusCode >= 300 {
            printFolderAPIError("create folder", statusCode, respBody)
            return
        }

        output := map[string]any{"status": "created", "name": folderName, "default": setDefault, "apiVersion": apiVersion}
        if len(respBody) > 0 {
            output["response"] = parseAnyJSONOrRaw(respBody)
        }

        if err := printFormattedOutput(outputFormat, output, func() {
            headers := []TableColumn{{Header: "NAME", Width: 40}, {Header: "DEFAULT", Width: 10}, {Header: "STATUS", Width: 12}}
            defaultValue := "No"
            if setDefault {
                defaultValue = "Yes"
            }
            PrintTable(headers, [][]string{{folderName, defaultValue, "created"}})
        }); err != nil {
            fmt.Println(err.Error())
        }
    },
}

var folderGetCmd = &cobra.Command{
    Use:   "get [folder-id]",
    Short: "Get folder details",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        folderID := args[0]
        apiVersion, _ := cmd.Flags().GetString("api-version")
        outputFormat, _ := cmd.Flags().GetString("format")

        statusCode, respBody, err := doFolderRequest(http.MethodGet, "/folders/"+url.PathEscape(folderID), apiVersion, nil, nil)
        if err != nil {
            fmt.Printf("Error getting folder: %v\n", err)
            return
        }
        if statusCode < 200 || statusCode >= 300 {
            printFolderAPIError("get folder", statusCode, respBody)
            return
        }

        var folder folderResponse
        _ = json.Unmarshal(respBody, &folder)
        raw := parseAnyJSONOrRaw(respBody)

        if err := printFormattedOutput(outputFormat, raw, func() {
            headers := []TableColumn{{Header: "ID", Width: 30}, {Header: "NAME", Width: 40}, {Header: "DEFAULT", Width: 10}, {Header: "PROJECTS", Width: 10}}
            defaultValue := "No"
            if folder.Default {
                defaultValue = "Yes"
            }
            PrintTable(headers, [][]string{{folder.ID, folder.Name, defaultValue, strconv.FormatInt(folder.ProjectCount, 10)}})
        }); err != nil {
            fmt.Println(err.Error())
        }
    },
}

var folderMetadataCmd = &cobra.Command{
    Use:   "metadata",
    Short: "Get folder metadata",
    Run: func(cmd *cobra.Command, args []string) {
        apiVersion, _ := cmd.Flags().GetString("api-version")
        outputFormat, _ := cmd.Flags().GetString("format")

        statusCode, respBody, err := doFolderRequest(http.MethodGet, "/folders/metadata", apiVersion, nil, nil)
        if err != nil {
            fmt.Printf("Error getting folder metadata: %v\n", err)
            return
        }
        if statusCode < 200 || statusCode >= 300 {
            printFolderAPIError("get folder metadata", statusCode, respBody)
            return
        }

        var metadata folderMetadataResponse
        _ = json.Unmarshal(respBody, &metadata)
        raw := parseAnyJSONOrRaw(respBody)

        if err := printFormattedOutput(outputFormat, raw, func() {
            headers := []TableColumn{{Header: "BAREMETAL TYPES", Width: 50}, {Header: "CMP LEGACY TYPES", Width: 50}}
            PrintTable(headers, [][]string{{strings.Join(metadata.BaremetalTypologies, ","), strings.Join(metadata.CmpLegacyTypologies, ",")}})
        }); err != nil {
            fmt.Println(err.Error())
        }
    },
}

var folderUpdateCmd = &cobra.Command{
    Use:   "update [folder-id]",
    Short: "Update folder",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        folderID := args[0]
        name, _ := cmd.Flags().GetString("name")
        setDefault, _ := cmd.Flags().GetBool("default")
        apiVersion, _ := cmd.Flags().GetString("api-version")
        outputFormat, _ := cmd.Flags().GetString("format")

        payload := map[string]any{}
        if cmd.Flags().Changed("name") {
            payload["name"] = name
        }
        if cmd.Flags().Changed("default") {
            payload["default"] = setDefault
        }

        if len(payload) == 0 {
            fmt.Println("Error: at least one flag must be provided: --name and/or --default")
            return
        }

        statusCode, respBody, err := doFolderRequest(http.MethodPut, "/folders/"+url.PathEscape(folderID), apiVersion, nil, payload)
        if err != nil {
            fmt.Printf("Error updating folder: %v\n", err)
            return
        }
        if statusCode < 200 || statusCode >= 300 {
            printFolderAPIError("update folder", statusCode, respBody)
            return
        }

        output := map[string]any{"status": "updated", "id": folderID, "changes": payload}
        if len(respBody) > 0 {
            output["response"] = parseAnyJSONOrRaw(respBody)
        }

        if err := printFormattedOutput(outputFormat, output, func() {
            headers := []TableColumn{{Header: "ID", Width: 30}, {Header: "STATUS", Width: 12}}
            PrintTable(headers, [][]string{{folderID, "updated"}})
        }); err != nil {
            fmt.Println(err.Error())
        }
    },
}

var folderDeleteCmd = &cobra.Command{
    Use:   "delete [folder-id]",
    Short: "Delete folder",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        folderID := args[0]
        skipConfirm, _ := cmd.Flags().GetBool("yes")
        apiVersion, _ := cmd.Flags().GetString("api-version")
        outputFormat, _ := cmd.Flags().GetString("format")

        if !skipConfirm {
            fmt.Printf("Are you sure you want to delete folder '%s'? (y/N): ", folderID)
            var answer string
            if _, err := fmt.Fscanln(os.Stdin, &answer); err != nil {
                answer = ""
            }
            if strings.ToLower(strings.TrimSpace(answer)) != "y" {
                fmt.Println("Delete cancelled.")
                return
            }
        }

        statusCode, respBody, err := doFolderRequest(http.MethodDelete, "/folders/"+url.PathEscape(folderID), apiVersion, nil, nil)
        if err != nil {
            fmt.Printf("Error deleting folder: %v\n", err)
            return
        }
        if statusCode < 200 || statusCode >= 300 {
            printFolderAPIError("delete folder", statusCode, respBody)
            return
        }

        output := map[string]any{"status": "deleted", "id": folderID}
        if err := printFormattedOutput(outputFormat, output, func() {
            headers := []TableColumn{{Header: "ID", Width: 30}, {Header: "STATUS", Width: 12}}
            PrintTable(headers, [][]string{{folderID, "deleted"}})
        }); err != nil {
            fmt.Println(err.Error())
        }
    },
}

var folderListCmd = &cobra.Command{
    Use:   "list",
    Short: "List folders",
    Run: func(cmd *cobra.Command, args []string) {
        apiVersion, _ := cmd.Flags().GetString("api-version")
        outputFormat, _ := cmd.Flags().GetString("format")
        filter, _ := cmd.Flags().GetString("filter")
        sortValue, _ := cmd.Flags().GetString("sort")
        projection, _ := cmd.Flags().GetString("projection")
        offset, _ := cmd.Flags().GetInt("offset")
        limit, _ := cmd.Flags().GetInt("limit")

        query := map[string]string{}
        if filter != "" {
            query["filter"] = filter
        }
        if sortValue != "" {
            query["sort"] = sortValue
        }
        if projection != "" {
            query["projection"] = projection
        }
        if offset >= 0 {
            query["offset"] = strconv.Itoa(offset)
        }
        if limit >= 0 {
            query["limit"] = strconv.Itoa(limit)
        }

        statusCode, respBody, err := doFolderRequest(http.MethodGet, "/folders", apiVersion, query, nil)
        if err != nil {
            fmt.Printf("Error listing folders: %v\n", err)
            return
        }
        if statusCode < 200 || statusCode >= 300 {
            printFolderAPIError("list folders", statusCode, respBody)
            return
        }

        var list folderListResponse
        _ = json.Unmarshal(respBody, &list)
        raw := parseAnyJSONOrRaw(respBody)

        if err := printFormattedOutput(outputFormat, raw, func() {
            if len(list.Values) == 0 {
                fmt.Println("No folders found")
                return
            }
            headers := []TableColumn{{Header: "NAME", Width: 40}, {Header: "ID", Width: 30}, {Header: "DEFAULT", Width: 10}, {Header: "PROJECTS", Width: 10}}
            rows := make([][]string, 0, len(list.Values))
            for _, f := range list.Values {
                defaultValue := "No"
                if f.Default {
                    defaultValue = "Yes"
                }
                rows = append(rows, []string{f.Name, f.ID, defaultValue, strconv.FormatInt(f.ProjectCount, 10)})
            }
            PrintTable(headers, rows)
        }); err != nil {
            fmt.Println(err.Error())
        }
    },
}

func doFolderRequest(method, path, apiVersion string, query map[string]string, payload any) (int, []byte, error) {
    config, err := LoadConfig()
    if err != nil {
        return 0, nil, fmt.Errorf("error loading configuration: %w", err)
    }

    if config.ClientID == "" || config.ClientSecret == "" {
        return 0, nil, fmt.Errorf("client ID or client secret not configured. Please run 'acloud config set --client-id YOUR_CLIENT_ID --client-secret YOUR_CLIENT_SECRET'")
    }

    baseURL := config.BaseURL
    if baseURL == "" {
        baseURL = DefaultBaseURL
    }

    tokenIssuerURL := config.TokenIssuerURL
    if tokenIssuerURL == "" {
        tokenIssuerURL = DefaultTokenIssuerURL
    }

    tokenConfig := clientcredentials.Config{ClientID: config.ClientID, ClientSecret: config.ClientSecret, TokenURL: tokenIssuerURL}
    token, err := tokenConfig.Token(context.Background())
    if err != nil {
        return 0, nil, fmt.Errorf("error getting access token: %w", err)
    }

    endpoint, err := url.Parse(strings.TrimRight(baseURL, "/") + path)
    if err != nil {
        return 0, nil, fmt.Errorf("error building endpoint: %w", err)
    }

    q := endpoint.Query()
    if apiVersion != "" {
        q.Set("api-version", apiVersion)
    }
    for key, value := range query {
        if value != "" {
            q.Set(key, value)
        }
    }
    endpoint.RawQuery = q.Encode()

    var bodyReader io.Reader
    if payload != nil {
        bodyBytes, err := json.Marshal(payload)
        if err != nil {
            return 0, nil, fmt.Errorf("error preparing request body: %w", err)
        }
        bodyReader = bytes.NewReader(bodyBytes)
    }

    req, err := http.NewRequestWithContext(context.Background(), method, endpoint.String(), bodyReader)
    if err != nil {
        return 0, nil, fmt.Errorf("error creating request: %w", err)
    }

    req.Header.Set("Authorization", "Bearer "+token.AccessToken)
    req.Header.Set("Accept", "application/json")
    if payload != nil {
        req.Header.Set("Content-Type", "application/json")
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return 0, nil, err
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    return resp.StatusCode, body, nil
}

func parseAnyJSONOrRaw(body []byte) any {
    var parsed any
    if err := json.Unmarshal(body, &parsed); err == nil {
        return parsed
    }
    return string(body)
}

func printFolderAPIError(operation string, statusCode int, body []byte) {
    if len(body) > 0 {
        fmt.Printf("Failed to %s (status %d): %s\n", operation, statusCode, string(body))
        return
    }
    fmt.Printf("Failed to %s (status %d).\n", operation, statusCode)
}

func printFormattedOutput(format string, data any, tablePrinter func()) error {
    switch strings.ToLower(strings.TrimSpace(format)) {
    case "table", "":
        tablePrinter()
        return nil
    case "json":
        b, err := json.MarshalIndent(data, "", "  ")
        if err != nil {
            return fmt.Errorf("error serializing JSON output: %w", err)
        }
        fmt.Println(string(b))
        return nil
    case "yaml":
        b, err := yaml.Marshal(data)
        if err != nil {
            return fmt.Errorf("error serializing YAML output: %w", err)
        }
        fmt.Print(string(b))
        return nil
    default:
        return fmt.Errorf("error: unsupported format '%s'. Use one of: table, json, yaml", format)
    }
}
