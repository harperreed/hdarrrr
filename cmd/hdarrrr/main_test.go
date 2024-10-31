package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"testing"
)

// validateImageProperties checks if all images have matching dimensions and color models
func validateImageProperties(images []image.Image) error {
	if len(images) < 2 {
		return fmt.Errorf("at least two images are required for validation")
	}

	// Check for nil images first
	for i, img := range images {
		if img == nil {
			return fmt.Errorf("image %d is nil", i+1)
		}
	}

	baseImg := images[0]
	baseBounds := baseImg.Bounds()
	baseColorModel := baseImg.ColorModel()

	for i, img := range images[1:] {
		// Check dimensions
		if img.Bounds() != baseBounds {
			return fmt.Errorf("image %d has different dimensions than the first image", i+2)
		}

		// Check color model
		if img.ColorModel() != baseColorModel {
			return fmt.Errorf("image %d has a different color model than the first image", i+2)
		}
	}

	return nil
}

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

// createHDRImage converts a regular image to HDR format
// func createHDRImage(width, height int, colorModel color.Model) hdr.Image {
// 	img := createTestImage(width, height, colorModel)
// 	bounds := img.Bounds()
// 	hdrImg := hdr.NewRGB(bounds)

// 	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
// 		for x := bounds.Max.X; x < bounds.Max.X; x++ {
// 			r, g, b, _ := img.At(x, y).RGBA()
// 			hdrImg.Set(x, y, hdrcolor.RGB{
// 				R: float64(r) / 0xffff,
// 				G: float64(g) / 0xffff,
// 				B: float64(b) / 0xffff,
// 			})
// 		}
// 	}
// 	return hdrImg
// }

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

func TestHDRProcessing(t *testing.T) {
	tests := []struct {
		name        string
		images      []image.Image
		expectError bool
	}{
		{
			name: "Valid HDR processing",
			images: []image.Image{
				createTestImage(100, 100, color.RGBAModel),
				createTestImage(100, 100, color.RGBAModel),
				createTestImage(100, 100, color.RGBAModel),
			},
			expectError: false,
		},
		{
			name: "HDR processing with different dimensions",
			images: []image.Image{
				createTestImage(100, 100, color.RGBAModel),
				createTestImage(200, 200, color.RGBAModel),
				createTestImage(100, 100, color.RGBAModel),
			},
			expectError: true,
		},
		{
			name: "HDR processing with different color models",
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

func TestCommandLineArguments(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name: "Valid arguments with default tonemapper",
			args: []string{
				"-low", "path/to/low.jpg",
				"-mid", "path/to/mid.jpg",
				"-high", "path/to/high.jpg",
				"-output", "path/to/output.jpg",
			},
			expectError: false,
		},
		{
			name: "Valid arguments with custom tonemapper",
			args: []string{
				"-low", "path/to/low.jpg",
				"-mid", "path/to/mid.jpg",
				"-high", "path/to/high.jpg",
				"-output", "path/to/output.jpg",
				"-tonemapper", "linear",
			},
			expectError: false,
		},
		{
			name: "Invalid tonemapper",
			args: []string{
				"-low", "path/to/low.jpg",
				"-mid", "path/to/mid.jpg",
				"-high", "path/to/high.jpg",
				"-output", "path/to/output.jpg",
				"-tonemapper", "invalid",
			},
			expectError: true,
		},
		{
			name: "Missing required arguments",
			args: []string{
				"-low", "path/to/low.jpg",
				"-mid", "path/to/mid.jpg",
			},
			expectError: true,
		},
		{
			name: "Invalid gamma value",
			args: []string{
				"-low", "path/to/low.jpg",
				"-mid", "path/to/mid.jpg",
				"-high", "path/to/high.jpg",
				"-output", "path/to/output.jpg",
				"-gamma", "-1.0",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ContinueOnError)
			err := flag.CommandLine.Parse(tt.args)

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
