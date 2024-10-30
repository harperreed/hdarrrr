
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

# [File: internal/processor/tonemap_test.go]
package processor

import (
	"testing"
)

func TestReinhardToneMapper(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "Zero value",
			input:    0.0,
			expected: 0.0,
		},
		{
			name:     "Value of 1",
			input:    1.0,
			expected: 0.5,
		},
		{
			name:     "Large value",
			input:    10.0,
			expected: 0.9090909090909091,
		},
		{
			name:     "Negative value",
			input:    -1.0,
			expected: -0.5,
		},
	}

	toneMapper := NewReinhardToneMapper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toneMapper.ToneMap(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}
