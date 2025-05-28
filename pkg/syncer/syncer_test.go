package syncer

import "testing"

func TestNewImageSyncer(t *testing.T) {
	tests := []struct {
		name        string
		sourceImage string
		targetImage string
		ghcrToken   string
		shouldError bool
		expectNil   bool
	}{
		{
			name:        "Valid inputs",
			sourceImage: "nginx:latest",
			targetImage: "ghcr.io/myorg/nginx:latest",
			ghcrToken:   "token",
			shouldError: false,
		},
		{
			name:        "Empty source image",
			sourceImage: "",
			targetImage: "ghcr.io/myorg/nginx:latest",
			ghcrToken:   "token",
			shouldError: true,
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			syncer := NewImageSyncer(tt.sourceImage, tt.targetImage, tt.ghcrToken)
			if tt.shouldError {
				if syncer != nil {
					t.Errorf("Expected nil syncer for error case")
				}
				return
			}
			if syncer == nil {
				t.Errorf("Expected non-nil syncer")
			}
		})
	}
}