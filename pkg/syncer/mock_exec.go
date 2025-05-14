package syncer

import (
	"os/exec"
)

// CommandExecutor is an interface for executing commands
type CommandExecutor interface {
	Command(name string, arg ...string) *exec.Cmd
}

// RealCommandExecutor uses the real exec.Command
type RealCommandExecutor struct{}

func (e *RealCommandExecutor) Command(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

// MockCommandExecutor mocks exec.Command for testing
type MockCommandExecutor struct {
	MockFunc func(name string, arg ...string) *exec.Cmd
}

func (e *MockCommandExecutor) Command(name string, arg ...string) *exec.Cmd {
	if e.MockFunc != nil {
		return e.MockFunc(name, arg...)
	}
	// Default mock that succeeds
	cmd := exec.Command("echo", "mock command")
	return cmd
}

// MockCmd creates a mock command that always succeeds
func MockCmd(name string, arg ...string) *exec.Cmd {
	return exec.Command("echo", "mock command")
}