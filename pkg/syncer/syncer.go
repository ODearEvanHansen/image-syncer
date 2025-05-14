package syncer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ImageSyncer handles the syncing of container images from a source registry to GHCR
type ImageSyncer struct {
	SourceImage string
	TargetImage string
	GHCRToken   string
}

// NewImageSyncer creates a new ImageSyncer instance
func NewImageSyncer(sourceImage, targetImage, ghcrToken string) *ImageSyncer {
	return &ImageSyncer{
		SourceImage: sourceImage,
		TargetImage: targetImage,
		GHCRToken:   ghcrToken,
	}
}

// Sync pulls the source image and pushes it to GHCR
func (s *ImageSyncer) Sync() error {
	// Pull the source image
	fmt.Printf("Pulling source image: %s\n", s.SourceImage)
	pullCmd := exec.Command("docker", "pull", s.SourceImage)
	pullCmd.Stdout = os.Stdout
	pullCmd.Stderr = os.Stderr
	if err := pullCmd.Run(); err != nil {
		return fmt.Errorf("failed to pull source image: %w", err)
	}

	// Tag the image for GHCR
	fmt.Printf("Tagging image for GHCR: %s\n", s.TargetImage)
	tagCmd := exec.Command("docker", "tag", s.SourceImage, s.TargetImage)
	tagCmd.Stdout = os.Stdout
	tagCmd.Stderr = os.Stderr
	if err := tagCmd.Run(); err != nil {
		return fmt.Errorf("failed to tag image: %w", err)
	}

	// Login to GHCR
	fmt.Println("Logging in to GHCR")
	loginCmd := exec.Command("docker", "login", "ghcr.io", "-u", "$GITHUB_ACTOR", "--password-stdin")
	loginCmd.Stdin = strings.NewReader(s.GHCRToken)
	loginCmd.Stdout = os.Stdout
	loginCmd.Stderr = os.Stderr
	if err := loginCmd.Run(); err != nil {
		return fmt.Errorf("failed to login to GHCR: %w", err)
	}

	// Push the image to GHCR
	fmt.Printf("Pushing image to GHCR: %s\n", s.TargetImage)
	pushCmd := exec.Command("docker", "push", s.TargetImage)
	pushCmd.Stdout = os.Stdout
	pushCmd.Stderr = os.Stderr
	if err := pushCmd.Run(); err != nil {
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