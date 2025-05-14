package syncer

import (
	"os/exec"
)

// MockCmd is a helper for testing commands
func MockCmd(command string, args ...string) *exec.Cmd {
	// For simplicity, just use echo command for mocking
	return exec.Command("echo", "Mocked command execution")
}