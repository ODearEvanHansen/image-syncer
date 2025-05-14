package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ODearEvanHansen/image-syncer/pkg/syncer"
)

func main() {
	// Define command-line flags
	sourceImage := flag.String("source", "", "Source container image to sync (required)")
	targetOrg := flag.String("target-org", "", "Target organization in GHCR (required)")
	ghcrToken := flag.String("token", "", "GitHub token for GHCR authentication (required)")

	// Parse flags
	flag.Parse()

	// Validate required flags
	if *sourceImage == "" || *targetOrg == "" || *ghcrToken == "" {
		fmt.Println("Error: source image, target organization, and GHCR token are required")
		flag.Usage()
		os.Exit(1)
	}

	// Generate target image name
	targetImage := syncer.ParseTargetImage(*sourceImage, *targetOrg)

	// Create a new image syncer
	imageSyncer := syncer.NewImageSyncer(*sourceImage, targetImage, *ghcrToken)

	// Sync the image
	fmt.Printf("Syncing image from %s to %s\n", *sourceImage, targetImage)
	if err := imageSyncer.Sync(); err != nil {
		fmt.Printf("Error syncing image: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Image sync completed successfully")
}