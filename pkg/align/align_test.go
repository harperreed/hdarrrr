package align

import (
	"image"
	"image/color"
	"testing"

	"github.com/mdouchement/hdr"
	"github.com/mdouchement/hdr/hdrcolor"
)

// createTestHDRImage creates a test HDR image with specified dimensions
func createTestHDRImage(width, height int) hdr.Image {
	img := hdr.NewRGB(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Set a test color value
			img.Set(x, y, hdrcolor.RGB{
				R: 0.5,
				G: 0.5,
				B: 0.5,
			})
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
			name: "Same size images",
			images: []image.Image{
				createTestImage(100, 100),
				createTestImage(100, 100),
			},
			expectError: false,
		},
		{
			name: "Different size images",
			images: []image.Image{
				createTestImage(100, 100),
				createTestImage(200, 200),
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

// createTestImage creates a regular test image with specified dimensions
func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 100, G: 100, B: 100, A: 255})
		}
	}
	return img
}

func TestBasicAligner(t *testing.T) {
	aligner := NewBasicAligner()

	tests := []struct {
		name        string
		images      []image.Image
		expectError bool
	}{
		{
			name: "Valid images",
			images: []image.Image{
				createTestImage(100, 100),
				createTestImage(100, 100),
			},
			expectError: false,
		},
		{
			name: "Mismatched dimensions",
			images: []image.Image{
				createTestImage(100, 100),
				createTestImage(200, 200),
			},
			expectError: true,
		},
		{
			name: "Single image",
			images: []image.Image{
				createTestImage(100, 100),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := aligner.Align(tt.images)

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

			if len(result) != len(tt.images) {
				t.Errorf("Expected %d images, got %d", len(tt.images), len(result))
			}

			if len(result) > 0 {
				baseBounds := result[0].Bounds()
				for i, img := range result[1:] {
					if img.Bounds() != baseBounds {
						t.Errorf("Image %d has different dimensions than base image", i+1)
					}
				}
			}
		})
	}
}
