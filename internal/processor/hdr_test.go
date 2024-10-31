package processor

import (
	"image"
	"image/color"
	"testing"
)

func createTestImage(width, height int, value uint8) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: value, G: value, B: value, A: 255})
		}
	}
	return img
}

func createTestGray16Image(width, height int, value uint16) image.Image {
	img := image.NewGray16(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.Gray16{Y: value})
		}
	}
	return img
}

func TestHDRProcessor_Process(t *testing.T) {
	tests := []struct {
		name        string
		images      []image.Image
		expectError bool
	}{
		{
			name: "Valid three images",
			images: []image.Image{
				createTestImage(2, 2, 50),  // Dark exposure
				createTestImage(2, 2, 128), // Mid exposure
				createTestImage(2, 2, 200), // Bright exposure
			},
			expectError: false,
		},
		{
			name: "Different sized images",
			images: []image.Image{
				createTestImage(2, 2, 128),
				createTestImage(3, 3, 128),
				createTestImage(2, 2, 128),
			},
			expectError: true,
		},
		{
			name: "Different color depth images",
			images: []image.Image{
				createTestImage(2, 2, 128),
				createTestGray16Image(2, 2, 32768),
				createTestImage(2, 2, 128),
			},
			expectError: true,
		},
		{
			name: "Wrong number of images",
			images: []image.Image{
				createTestImage(2, 2, 128),
				createTestImage(2, 2, 128),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewHDRProcessor()
			result, err := processor.Process(tt.images)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result == nil {
				t.Fatal("Expected result image, got nil")
			}

			// Check result dimensions
			bounds := result.Bounds()
			expectedBounds := tt.images[0].Bounds()
			if bounds != expectedBounds {
				t.Errorf("Expected dimensions %v, got %v", expectedBounds, bounds)
			}
		})
	}
}
