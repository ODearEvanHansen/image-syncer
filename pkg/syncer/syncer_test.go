package syncer

import (
	"fmt"
	"testing"
)

func TestParseTargetImage(t *testing.T) {
	tests := []struct {
		name        string
		sourceImage string
		targetOrg   string
		expected    string
	}{
		{
			name:        "Simple image",
			sourceImage: "nginx:latest",
			targetOrg:   "myorg",
			expected:    "ghcr.io/myorg/nginx:latest",
		},
		{
			name:        "Image with registry",
			sourceImage: "docker.io/library/ubuntu:20.04",
			targetOrg:   "myorg",
			expected:    "ghcr.io/myorg/ubuntu:20.04",
		},
		{
			name:        "Target org with ghcr.io prefix",
			sourceImage: "alpine:3.14",
			targetOrg:   "ghcr.io/myorg",
			expected:    "ghcr.io/myorg/alpine:3.14",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseTargetImage(tc.sourceImage, tc.targetOrg)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// MockExecutor is a mock implementation of CommandExecutor for testing
type MockExecutor struct {
	ExpectedCommands []string
	CurrentCommand   int
	ShouldFail       bool
}

// Execute records the command and returns nil or an error based on ShouldFail
func (m *MockExecutor) Execute(name string, args ...string) error {
	if m.ShouldFail {
		return fmt.Errorf("mock execution failed")
	}
	m.CurrentCommand++
	return nil
}

// TestImageSyncerWithMockExecutor tests the ImageSyncer with a mock executor
func TestImageSyncerWithMockExecutor(t *testing.T) {
	// Create mock executors
	mockExecutor := &MockExecutor{}
	mockLoginExecutor := &MockExecutor{}
	
	// Create an ImageSyncer with mock executors
	syncer := &ImageSyncer{
		SourceImage:   "nginx:latest",
		TargetImage:   "ghcr.io/myorg/nginx:latest",
		GHCRToken:     "dummy-token",
		Executor:      mockExecutor,
		LoginExecutor: mockLoginExecutor,
	}
	
	// Test successful sync
	err := syncer.Sync()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Test failed pull
	mockExecutor.ShouldFail = true
	err = syncer.Sync()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

// TestMockHelperProcess is not a real test, it's used by MockCmd
func TestMockHelperProcess(t *testing.T) {
	// Skip this test - it's not a real test
	t.Skip("This is a helper for mocking commands and not a real test")
}