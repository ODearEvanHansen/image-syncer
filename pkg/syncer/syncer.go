package syncer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CommandExecutor is an interface for executing commands
type CommandExecutor interface {
	Execute(name string, args ...string) error
}

// DefaultExecutor is the default implementation of CommandExecutor
type DefaultExecutor struct{}

// Execute runs a command with the given name and arguments
func (e *DefaultExecutor) Execute(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// DockerLoginExecutor is a special executor for Docker login that handles stdin
type DockerLoginExecutor struct {
	Token string
}

// Execute runs the Docker login command with the token provided via stdin
func (e *DockerLoginExecutor) Execute(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = strings.NewReader(e.Token)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ImageSyncer handles the syncing of container images from a source registry to GHCR
type ImageSyncer struct {
	SourceImage string
	TargetImage string
	GHCRToken   string
	Executor    CommandExecutor
	LoginExecutor CommandExecutor
}

// NewImageSyncer creates a new ImageSyncer instance with default executors
func NewImageSyncer(sourceImage, targetImage, ghcrToken string) *ImageSyncer {
	return &ImageSyncer{
		SourceImage: sourceImage,
		TargetImage: targetImage,
		GHCRToken:   ghcrToken,
		Executor:    &DefaultExecutor{},
		LoginExecutor: &DockerLoginExecutor{Token: ghcrToken},
	}
}

// Sync pulls the source image and pushes it to GHCR
func (s *ImageSyncer) Sync() error {
	// Pull the source image
	fmt.Printf("Pulling source image: %s\n", s.SourceImage)
	if err := s.Executor.Execute("docker", "pull", s.SourceImage); err != nil {
		return fmt.Errorf("failed to pull source image: %w", err)
	}

	// Tag the image for GHCR
	fmt.Printf("Tagging image for GHCR: %s\n", s.TargetImage)
	if err := s.Executor.Execute("docker", "tag", s.SourceImage, s.TargetImage); err != nil {
		return fmt.Errorf("failed to tag image: %w", err)
	}

	// Login to GHCR
	fmt.Println("Logging in to GHCR")
	githubActor := os.Getenv("GITHUB_ACTOR")
	if githubActor == "" {
		githubActor = "github-actions"
	}
	
	if err := s.LoginExecutor.Execute("docker", "login", "ghcr.io", "-u", githubActor, "--password-stdin"); err != nil {
		return fmt.Errorf("failed to login to GHCR: %w", err)
	}

	// Push the image to GHCR
	fmt.Printf("Pushing image to GHCR: %s\n", s.TargetImage)
	if err := s.Executor.Execute("docker", "push", s.TargetImage); err != nil {
		return fmt.Errorf("failed to push image to GHCR: %w", err)
	}

	fmt.Println("Image sync completed successfully")
	return nil
}

// ParseTargetImage generates the target image name based on the source image
// and the target organization
func ParseTargetImage(sourceImage, targetOrg string) string {
	// Extract the image name and tag from the source image
	parts := strings.Split(sourceImage, "/")
	var imagePart string
	if len(parts) > 1 {
		imagePart = parts[len(parts)-1]
	} else {
		imagePart = sourceImage
	}

	// If the target org doesn't include ghcr.io, add it
	if !strings.HasPrefix(targetOrg, "ghcr.io/") {
		targetOrg = "ghcr.io/" + targetOrg
	}

	return targetOrg + "/" + imagePart
}