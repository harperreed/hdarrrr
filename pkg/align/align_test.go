package align

import (
	"image"
	"testing"

	"github.com/mdouchement/hdr"
)

func createTestImage(width, height int) hdr.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	return hdr.NewImageFromGoImage(img)
}

func TestAlignImages(t *testing.T) {
	tests := []struct {
		name        string
		images      []hdr.Image
		expectError bool
	}{
		{
			name: "Same size images",
			images: []hdr.Image{
				createTestImage(100, 100),
				createTestImage(100, 100),
			},
			expectError: false,
		},
		{
			name: "Different size images",
			images: []hdr.Image{
				createTestImage(100, 100),
				createTestImage(200, 200),
			},
			expectError: true,
		},
		{
			name:        "Empty image list",
			images:      []hdr.Image{},
			expectError: true,
		},
		{
			name: "Single image",
			images: []hdr.Image{
				createTestImage(100, 100),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aligned, err := AlignImages(tt.images)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(aligned) != len(tt.images) {
				t.Errorf("Expected %d aligned images, got %d", len(tt.images), len(aligned))
			}

			// Check all aligned images have the same dimensions
			if len(aligned) > 0 {
				bounds := aligned[0].Bounds()
				for i, img := range aligned[1:] {
					if img.Bounds() != bounds {
						t.Errorf("Image %d has different bounds than first image", i+1)
					}
				}
			}
		})
	}
}
