package imaging

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// TestImage represents a test image with its format and data
type TestImage struct {
	data   []byte
	format string
	path   string
}

// createTestImage creates a small test image with a specific color
func createTestImage(t *testing.T, format string, c color.Color) TestImage {
	// Create a 2x2 test image
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			img.Set(x, y, c)
		}
	}

	// Encode the image to bytes
	var buf bytes.Buffer
	switch format {
	case "png":
		if err := png.Encode(&buf, img); err != nil {
			t.Fatal("Failed to encode PNG:", err)
		}
	case "jpeg":
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100}); err != nil {
			t.Fatal("Failed to encode JPEG:", err)
		}
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_*."+format)
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
	}
	if _, err := tmpFile.Write(buf.Bytes()); err != nil {
		t.Fatal("Failed to write temp file:", err)
	}
	tmpFile.Close()

	return TestImage{
		data:   buf.Bytes(),
		format: format,
		path:   tmpFile.Name(),
	}
}

func TestLoadImage(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		color       color.Color
		expectError bool
	}{
		{
			name:        "Valid PNG",
			format:      "png",
			color:       color.RGBA{R: 255, G: 0, B: 0, A: 255},
			expectError: false,
		},
		{
			name:        "Valid JPEG",
			format:      "jpeg",
			color:       color.RGBA{R: 0, G: 255, B: 0, A: 255},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testImg := createTestImage(t, tt.format, tt.color)
			defer os.Remove(testImg.path)

			img, err := LoadImage(testImg.path)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if img == nil {
				t.Fatal("Expected image, got nil")
			}

			// Check image dimensions
			bounds := img.Bounds()
			if bounds.Dx() != 2 || bounds.Dy() != 2 {
				t.Errorf("Expected 2x2 image, got %dx%d", bounds.Dx(), bounds.Dy())
			}
		})
	}
}

func TestLoadImages(t *testing.T) {
	// Create multiple test images
	testImages := []TestImage{
		createTestImage(t, "png", color.RGBA{R: 255, G: 0, B: 0, A: 255}),
		createTestImage(t, "jpeg", color.RGBA{R: 0, G: 255, B: 0, A: 255}),
		createTestImage(t, "png", color.RGBA{R: 0, G: 0, B: 255, A: 255}),
	}

	// Clean up test files
	defer func() {
		for _, img := range testImages {
			os.Remove(img.path)
		}
	}()

	// Get paths
	paths := make([]string, len(testImages))
	for i, img := range testImages {
		paths[i] = img.path
	}

	// Test loading multiple images
	images, err := LoadImages(paths...)
	if err != nil {
		t.Fatalf("Failed to load images: %v", err)
	}

	if len(images) != len(testImages) {
		t.Errorf("Expected %d images, got %d", len(testImages), len(images))
	}
}

func TestSaveImage(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		expectError bool
	}{
		{
			name:        "Save as PNG",
			format:      "png",
			expectError: false,
		},
		{
			name:        "Save as JPEG",
			format:      "jpeg",
			expectError: false,
		},
		{
			name:        "Invalid format",
			format:      "invalid",
			expectError: true,
		},
	}

	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary output path
			tmpFile := filepath.Join(os.TempDir(), "test_output."+tt.format)
			defer os.Remove(tmpFile)

			err := SaveImage(img, tmpFile)
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify file exists and is readable
			_, err = os.Stat(tmpFile)
			if err != nil {
				t.Errorf("Failed to verify output file: %v", err)
			}
		})
	}
}

func TestUnsupportedFormat(t *testing.T) {
	// Try to load a non-image file
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = LoadImage(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for unsupported format, got nil")
	}
}
