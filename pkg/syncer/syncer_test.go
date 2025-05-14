package syncer

import (
	"os/exec"
	"strings"
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
		{
			name:        "Target org with uppercase characters",
			sourceImage: "nginx:latest",
			targetOrg:   "MyOrg",
			expected:    "ghcr.io/myorg/nginx:latest",
		},
		{
			name:        "Target org with mixed case and ghcr.io prefix",
			sourceImage: "alpine:3.14",
			targetOrg:   "ghcr.io/MyOrg",
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

// TestImageSyncerWithMockExecutor tests the ImageSyncer with a mock executor
func TestImageSyncerWithMockExecutor(t *testing.T) {
	// Create a mock executor that records commands
	var executedCommands []string
	mockExecutor := &MockCommandExecutor{
		MockFunc: func(name string, arg ...string) *exec.Cmd {
			// Record the command
			executedCommands = append(executedCommands, name+" "+strings.Join(arg, " "))
			
			// Special case for git config to return a username
			if name == "git" && len(arg) > 0 && arg[0] == "config" {
				cmd := exec.Command("echo", "mockuser")
				return cmd
			}
			
			// Create a command that always succeeds
			cmd := exec.Command("echo", "mock command")
			return cmd
		},
	}

	// Create a syncer with the mock executor
	syncer := NewImageSyncerWithExecutor(
		"nginx:latest",
		"ghcr.io/myorg/nginx:latest",
		"fake-token",
		mockExecutor,
	)

	// Run the sync
	err := syncer.Sync()
	if err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	// Verify the commands that were executed
	expectedCommands := []string{
		"docker pull nginx:latest",
		"docker tag nginx:latest ghcr.io/myorg/nginx:latest",
		"git config user.name",
		"docker login ghcr.io -u mockuser --password-stdin",
		"docker push ghcr.io/myorg/nginx:latest",
	}

	if len(executedCommands) != len(expectedCommands) {
		t.Errorf("Expected %d commands, got %d", len(expectedCommands), len(executedCommands))
	}

	for i, cmd := range expectedCommands {
		if i < len(executedCommands) && executedCommands[i] != cmd {
			t.Errorf("Expected command %d to be '%s', got '%s'", i, cmd, executedCommands[i])
		}
	}
}

// TestMockHelperProcess is a placeholder test that does nothing
// We don't need this anymore since we simplified our mocking approach
func TestMockHelperProcess(t *testing.T) {
	t.Skip("Not a real test")
}