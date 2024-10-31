package main

import (
	"image"
	"image/color"
	"testing"

	"github.com/mdouchement/hdr"
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
		images      []hdr.Image
		expectError bool
	}{
		{
			name: "Valid aligned images - same dimensions and color model",
			images: []hdr.Image{
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
			},
			expectError: false,
		},
		{
			name: "Different dimensions",
			images: []hdr.Image{
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
				hdr.NewImageFromGoImage(createTestImage(200, 200, color.RGBAModel)),
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
			},
			expectError: true,
		},
		{
			name: "Different color models",
			images: []hdr.Image{
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.GrayModel)),
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
			},
			expectError: true,
		},
		{
			name: "Single image",
			images: []hdr.Image{
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
			},
			expectError: true,
		},
		{
			name:        "Empty image list",
			images:      []hdr.Image{},
			expectError: true,
		},
		{
			name: "Nil image in list",
			images: []hdr.Image{
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
				nil,
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
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

func TestHDRProcessing(t *testing.T) {
	tests := []struct {
		name        string
		images      []hdr.Image
		expectError bool
	}{
		{
			name: "Valid HDR processing",
			images: []hdr.Image{
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
			},
			expectError: false,
		},
		{
			name: "HDR processing with different dimensions",
			images: []hdr.Image{
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
				hdr.NewImageFromGoImage(createTestImage(200, 200, color.RGBAModel)),
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
			},
			expectError: true,
		},
		{
			name: "HDR processing with different color models",
			images: []hdr.Image{
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.GrayModel)),
				hdr.NewImageFromGoImage(createTestImage(100, 100, color.RGBAModel)),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := hdr.Merge(tt.images)

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
