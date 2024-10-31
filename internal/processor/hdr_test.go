package processor

import (
	"image"
	"image/color"
	"testing"
)

func createTestImage(width, height int, value uint8) image.Image {
	// Minimum size for ICam06 operator
	if width < 32 {
		width = 32
	}
	if height < 32 {
		height = 32
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create a gradient pattern to ensure varied pixel values
			v := uint8((x * y * int(value)) % 256)
			img.Set(x, y, color.RGBA{R: v, G: v, B: v, A: 255})
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
				createTestImage(32, 32, 50),  // Dark exposure
				createTestImage(32, 32, 128), // Mid exposure
				createTestImage(32, 32, 200), // Bright exposure
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
				createTestImage(32, 32, 50),
				createTestImage(32, 32, 128),
				createTestImage(32, 32, 200),
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
				createTestImage(32, 32, 50),
				createTestImage(32, 32, 128),
				createTestImage(32, 32, 200),
			},
			toneMapper:  "linear",
			params:      map[string]float64{},
			expectError: false,
		},
		{
			name: "Valid three images with Logarithmic",
			images: []image.Image{
				createTestImage(32, 32, 50),
				createTestImage(32, 32, 128),
				createTestImage(32, 32, 200),
			},
			toneMapper:  "logarithmic",
			params:      map[string]float64{},
			expectError: false,
		},
		{
			name: "Valid three images with Durand",
			images: []image.Image{
				createTestImage(32, 32, 50),
				createTestImage(32, 32, 128),
				createTestImage(32, 32, 200),
			},
			toneMapper: "durand",
			params: map[string]float64{
				"saturation": 0.8,
			},
			expectError: false,
		},
		{
			name: "Valid three images with CustomReinhard05",
			images: []image.Image{
				createTestImage(32, 32, 50),
				createTestImage(32, 32, 128),
				createTestImage(32, 32, 200),
			},
			toneMapper: "custom_reinhard05",
			params: map[string]float64{
				"intensity": 1.0,
				"light":     0.0,
				"gamma":     1.0,
			},
			expectError: false,
		},
		{
			name: "Valid three images with ICam06",
			images: []image.Image{
				createTestImage(32, 32, 50),
				createTestImage(32, 32, 128),
				createTestImage(32, 32, 200),
			},
			toneMapper: "icam06",
			params: map[string]float64{
				"gamma":     1.0,
				"contrast":  4.0,
				"chromatic": 0.0,
			},
			expectError: false,
		},
		{
			name: "Different sized images",
			images: []image.Image{
				createTestImage(32, 32, 128),
				createTestImage(64, 64, 128),
				createTestImage(32, 32, 128),
			},
			toneMapper:  "reinhard05",
			params:      map[string]float64{},
			expectError: true,
		},
		{
			name: "Single image",
			images: []image.Image{
				createTestImage(32, 32, 128),
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
				createTestImage(32, 32, 128),
				nil,
				createTestImage(32, 32, 128),
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
