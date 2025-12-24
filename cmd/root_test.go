package cmd

import (
	"os"
	"testing"
)

func TestGetArubaClient(t *testing.T) {
	// This test requires valid credentials
	// Skip if credentials are not configured
	if os.Getenv("ACLOUD_TEST_SKIP_CLIENT") == "true" {
		t.Skip("Skipping client test (ACLOUD_TEST_SKIP_CLIENT=true)")
	}

	client, err := GetArubaClient()
	if err != nil {
		// If credentials are not configured, skip the test
		if err.Error() == "failed to load configuration: open ~/.acloud.yaml: no such file or directory" ||
			err.Error() == "client ID or client secret not configured. Please run 'acloud config set'" {
			t.Skip("Skipping test: credentials not configured")
		}
		t.Fatalf("GetArubaClient() error = %v", err)
	}

	if client == nil {
		t.Fatal("GetArubaClient() returned nil client")
	}
}

func TestGetArubaClient_Caching(t *testing.T) {
	// Test that client caching works
	if os.Getenv("ACLOUD_TEST_SKIP_CLIENT") == "true" {
		t.Skip("Skipping client test (ACLOUD_TEST_SKIP_CLIENT=true)")
	}

	// Clear cache
	clientCacheLock.Lock()
	clientCache = nil
	cachedClientID = ""
	cachedSecret = ""
	cachedDebug = false
	clientCacheLock.Unlock()

	client1, err1 := GetArubaClient()
	if err1 != nil {
		if err1.Error() == "failed to load configuration: open ~/.acloud.yaml: no such file or directory" ||
			err1.Error() == "client ID or client secret not configured. Please run 'acloud config set'" {
			t.Skip("Skipping test: credentials not configured")
		}
		t.Fatalf("GetArubaClient() error = %v", err1)
	}

	// Second call should use cached client
	client2, err2 := GetArubaClient()
	if err2 != nil {
		t.Fatalf("GetArubaClient() second call error = %v", err2)
	}

	// Should be the same instance (cached)
	if client1 != client2 {
		t.Error("GetArubaClient() should return cached client on second call")
	}
}

func TestPrintTable(t *testing.T) {
	headers := []TableColumn{
		{Header: "NAME", Width: 10},
		{Header: "ID", Width: 10},
		{Header: "STATUS", Width: 10},
	}

	rows := [][]string{
		{"test-name", "test-id-123", "Active"},
		{"another-name", "another-id-456", "Inactive"},
	}

	// Should not panic
	PrintTable(headers, rows)
}

func TestPrintTable_EmptyRows(t *testing.T) {
	headers := []TableColumn{
		{Header: "NAME", Width: 10},
		{Header: "ID", Width: 10},
	}

	rows := [][]string{}

	// Should not panic with empty rows
	PrintTable(headers, rows)
}

func TestPrintTable_LongValues(t *testing.T) {
	headers := []TableColumn{
		{Header: "NAME", Width: 5},
	}

	rows := [][]string{
		{"very-long-name-that-exceeds-width"},
	}

	// Should truncate long values
	PrintTable(headers, rows)
}

func TestPrintTable_MismatchedColumns(t *testing.T) {
	headers := []TableColumn{
		{Header: "NAME", Width: 10},
		{Header: "ID", Width: 10},
	}

	rows := [][]string{
		{"test-name"}, // Missing second column
		{"another-name", "id-123", "extra-column"}, // Extra column
	}

	// Should not panic with mismatched columns
	PrintTable(headers, rows)
}

func TestGetArubaClient_DebugFlagChange(t *testing.T) {
	// Test that client cache is invalidated when debug flag changes
	if os.Getenv("ACLOUD_TEST_SKIP_CLIENT") == "true" {
		t.Skip("Skipping client test (ACLOUD_TEST_SKIP_CLIENT=true)")
	}

	// Clear cache
	clientCacheLock.Lock()
	clientCache = nil
	cachedClientID = ""
	cachedSecret = ""
	cachedDebug = false
	clientCacheLock.Unlock()

	// First call without debug
	rootCmd.PersistentFlags().Set("debug", "false")
	client1, err1 := GetArubaClient()
	if err1 != nil {
		if err1.Error() == "failed to load configuration: open ~/.acloud.yaml: no such file or directory" ||
			err1.Error() == "client ID or client secret not configured. Please run 'acloud config set'" {
			t.Skip("Skipping test: credentials not configured")
		}
		t.Fatalf("GetArubaClient() error = %v", err1)
	}

	// Change debug flag
	rootCmd.PersistentFlags().Set("debug", "true")

	// Second call should create new client (cache invalidated)
	client2, err2 := GetArubaClient()
	if err2 != nil {
		t.Fatalf("GetArubaClient() second call error = %v", err2)
	}

	// Should be different instances (cache invalidated due to debug flag change)
	if client1 == client2 {
		t.Error("GetArubaClient() should return new client when debug flag changes")
	}
}

