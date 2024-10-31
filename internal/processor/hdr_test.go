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
		toneMapper  string
		params      map[string]float64
		expectError bool
	}{
		{
			name: "Valid three images with Reinhard05",
			images: []image.Image{
				createTestImage(2, 2, 50),  // Dark exposure
				createTestImage(2, 2, 128), // Mid exposure
				createTestImage(2, 2, 200), // Bright exposure
			},
			toneMapper: "reinhard05",
			params: map[string]float64{
				"gamma":     1.0,
				"intensity": 1.0,
				"light":     0.0,
			},
			expectError: false,
		},
		{
			name: "Valid three images with Drago03",
			images: []image.Image{
				createTestImage(2, 2, 50),  // Dark exposure
				createTestImage(2, 2, 128), // Mid exposure
				createTestImage(2, 2, 200), // Bright exposure
			},
			toneMapper: "drago03",
			params: map[string]float64{
				"gamma": 1.0,
			},
			expectError: false,
		},
		{
			name: "Valid three images with Linear",
			images: []image.Image{
				createTestImage(2, 2, 50),  // Dark exposure
				createTestImage(2, 2, 128), // Mid exposure
				createTestImage(2, 2, 200), // Bright exposure
			},
			toneMapper: "linear",
			params:     map[string]float64{},
			expectError: false,
		},
		{
			name: "Valid three images with Logarithmic",
			images: []image.Image{
				createTestImage(2, 2, 50),  // Dark exposure
				createTestImage(2, 2, 128), // Mid exposure
				createTestImage(2, 2, 200), // Bright exposure
			},
			toneMapper: "logarithmic",
			params:     map[string]float64{},
			expectError: false,
		},
		{
			name: "Valid three images with Durand",
			images: []image.Image{
				createTestImage(2, 2, 50),  // Dark exposure
				createTestImage(2, 2, 128), // Mid exposure
				createTestImage(2, 2, 200), // Bright exposure
			},
			toneMapper: "durand",
			params:     map[string]float64{},
			expectError: false,
		},
		{
			name: "Valid three images with CustomReinhard05",
			images: []image.Image{
				createTestImage(2, 2, 50),  // Dark exposure
				createTestImage(2, 2, 128), // Mid exposure
				createTestImage(2, 2, 200), // Bright exposure
			},
			toneMapper: "custom_reinhard05",
			params:     map[string]float64{},
			expectError: false,
		},
		{
			name: "Valid three images with ICam06",
			images: []image.Image{
				createTestImage(2, 2, 50),  // Dark exposure
				createTestImage(2, 2, 128), // Mid exposure
				createTestImage(2, 2, 200), // Bright exposure
			},
			toneMapper: "icam06",
			params:     map[string]float64{},
			expectError: false,
		},
		{
			name: "Different sized images",
			images: []image.Image{
				createTestImage(2, 2, 128),
				createTestImage(3, 3, 128),
				createTestImage(2, 2, 128),
			},
			toneMapper:  "reinhard05",
			params:      map[string]float64{},
			expectError: true,
		},
		{
			name: "Single image",
			images: []image.Image{
				createTestImage(2, 2, 128),
			},
			toneMapper:  "reinhard05",
			params:      map[string]float64{},
			expectError: true,
		},
		{
			name:        "Nil images slice",
			images:      nil,
			toneMapper:  "reinhard05",
			params:      map[string]float64{},
			expectError: true,
		},
		{
			name: "One nil image",
			images: []image.Image{
				createTestImage(2, 2, 128),
				nil,
				createTestImage(2, 2, 128),
			},
			toneMapper:  "reinhard05",
			params:      map[string]float64{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor := NewHDRProcessor().WithToneMapper(tt.toneMapper).WithParams(tt.params)
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
			if len(tt.images) > 0 {
				expectedBounds := tt.images[0].Bounds()
				if result.Bounds() != expectedBounds {
					t.Errorf("Expected dimensions %v, got %v", expectedBounds, result.Bounds())
				}
			}
		})
	}
}
