package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainCLI(t *testing.T) {
	// Skip if not in test mode
	if os.Getenv("TEST_CLI") != "true" {
		t.Skip("Skipping CLI test; set TEST_CLI=true to run")
	}

	// Test with missing arguments
	cmd := exec.Command("go", "run", "main.go")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("Expected error when running without arguments, but got none")
	}
	
	// Check if the output contains the usage information
	outputStr := string(output)
	if !strings.Contains(outputStr, "source image, target organization, and GHCR token are required") {
		t.Errorf("Expected usage information in output, but got: %s", outputStr)
	}

	// Test with valid arguments (mock mode)
	// This would require setting up environment variables and mocks
	// For a real test, we would need to set up Docker and GHCR credentials
}