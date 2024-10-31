package align

import (
	"image"
	"image/color"
	"testing"
)

func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 100, G: 100, B: 100, A: 255})
		}
	}
	return img
}

func TestAlignImages(t *testing.T) {
	tests := []struct {
		name        string
		images      []image.Image
		expectError bool
	}{
		{
			name: "Valid images same size",
			images: []image.Image{
				createTestImage(100, 100),
				createTestImage(100, 100),
				createTestImage(100, 100),
			},
			expectError: false,
		},
		{
			name: "Different sized images",
			images: []image.Image{
				createTestImage(100, 100),
				createTestImage(200, 200),
				createTestImage(100, 100),
			},
			expectError: true,
		},
		{
			name:        "Empty image list",
			images:      []image.Image{},
			expectError: true,
		},
		{
			name: "Single image",
			images: []image.Image{
				createTestImage(100, 100),
			},
			expectError: true,
		},
		{
			name: "Nil image in list",
			images: []image.Image{
				createTestImage(100, 100),
				nil,
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

			// Check that all aligned images have the same dimensions
			if len(aligned) > 0 {
				baseBounds := aligned[0].Bounds()
				for i, img := range aligned[1:] {
					if img.Bounds() != baseBounds {
						t.Errorf("Aligned image %d has different dimensions than first image", i+1)
					}
				}
			}
		})
	}
}
