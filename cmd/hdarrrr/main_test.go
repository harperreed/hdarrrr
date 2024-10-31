package main

import (
	"image"
	"image/color"
	"testing"
)

// createTestImage creates a test image with specified dimensions and color model
func createTestImage(width, height int, colorModel color.Model) image.Image {
	if colorModel == color.RGBAModel {
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				img.Set(x, y, color.RGBA{R: 100, G: 100, B: 100, A: 255})
			}
		}
		return img
	}
	img := image.NewGray(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.Gray{Y: 100})
		}
	}
	return img
}

func TestValidateImageProperties(t *testing.T) {
	tests := []struct {
		name        string
		images      []image.Image
		expectError bool
	}{
		{
			name: "Valid aligned images - same dimensions and color model",
			images: []image.Image{
				createTestImage(100, 100, color.RGBAModel),
				createTestImage(100, 100, color.RGBAModel),
				createTestImage(100, 100, color.RGBAModel),
			},
			expectError: false,
		},
		{
			name: "Different dimensions",
			images: []image.Image{
				createTestImage(100, 100, color.RGBAModel),
				createTestImage(200, 200, color.RGBAModel),
				createTestImage(100, 100, color.RGBAModel),
			},
			expectError: true,
		},
		{
			name: "Different color models",
			images: []image.Image{
				createTestImage(100, 100, color.RGBAModel),
				createTestImage(100, 100, color.GrayModel),
				createTestImage(100, 100, color.RGBAModel),
			},
			expectError: true,
		},
		{
			name: "Single image",
			images: []image.Image{
				createTestImage(100, 100, color.RGBAModel),
			},
			expectError: true,
		},
		{
			name:        "Empty image list",
			images:      []image.Image{},
			expectError: true,
		},
		{
			name: "Nil image in list",
			images: []image.Image{
				createTestImage(100, 100, color.RGBAModel),
				nil,
				createTestImage(100, 100, color.RGBAModel),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateImageProperties(tt.images)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestHDRMethodFlag(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		expectError bool
	}{
		{
			name:        "Valid tone-mapping method",
			method:      "tone-mapping",
			expectError: false,
		},
		{
			name:        "Valid exposure-fusion method",
			method:      "exposure-fusion",
			expectError: false,
		},
		{
			name:        "Invalid method",
			method:      "invalid-method",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate command line arguments
			args := []string{"-method", tt.method}
			flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
			hdrMethod := flagSet.String("method", "tone-mapping", "HDR method: tone-mapping or exposure-fusion")
			flagSet.Parse(args)

			// Validate the method
			err := validateHDRMethod(*hdrMethod)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestProcessExposureFusion(t *testing.T) {
	tests := []struct {
		name        string
		images      []image.Image
		expectError bool
	}{
		{
			name: "Valid images",
			images: []image.Image{
				createTestImage(100, 100, color.RGBAModel),
				createTestImage(100, 100, color.RGBAModel),
				createTestImage(100, 100, color.RGBAModel),
			},
			expectError: false,
		},
		{
			name: "Different dimensions",
			images: []image.Image{
				createTestImage(100, 100, color.RGBAModel),
				createTestImage(200, 200, color.RGBAModel),
				createTestImage(100, 100, color.RGBAModel),
			},
			expectError: true,
		},
		{
			name: "Different color models",
			images: []image.Image{
				createTestImage(100, 100, color.RGBAModel),
				createTestImage(100, 100, color.GrayModel),
				createTestImage(100, 100, color.RGBAModel),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := processExposureFusion(tt.images)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}
