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

	"github.com/mdouchement/hdr"
)

// createTestImage creates a test image with specified dimensions and color
func createTestImage(width, height int, c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

// createTestImageFile creates a temporary image file for testing
func createTestImageFile(t *testing.T, format string, c color.Color) (string, func()) {
	// Create a 2x2 test image
	img := createTestImage(2, 2, c)

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
	default:
		t.Fatalf("Unsupported format: %s", format)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_*."+format)
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
	}

	if _, err := tmpFile.Write(buf.Bytes()); err != nil {
		os.Remove(tmpFile.Name())
		t.Fatal("Failed to write temp file:", err)
	}
	tmpFile.Close()

	cleanup := func() {
		os.Remove(tmpFile.Name())
	}

	return tmpFile.Name(), cleanup
}

func TestLoadImages(t *testing.T) {
	// Create test files
	file1, cleanup1 := createTestImageFile(t, "png", color.RGBA{R: 255, A: 255})
	defer cleanup1()
	file2, cleanup2 := createTestImageFile(t, "png", color.RGBA{G: 255, A: 255})
	defer cleanup2()
	file3, cleanup3 := createTestImageFile(t, "png", color.RGBA{B: 255, A: 255})
	defer cleanup3()

	images, err := LoadImages(file1, file2, file3)
	if err != nil {
		t.Fatalf("Failed to load images: %v", err)
	}

	if len(images) != 3 {
		t.Errorf("Expected 3 images, got %d", len(images))
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
			color:       color.RGBA{R: 255, A: 255},
			expectError: false,
		},
		{
			name:        "Valid JPEG",
			format:      "jpeg",
			color:       color.RGBA{G: 255, A: 255},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, cleanup := createTestImageFile(t, tt.format, tt.color)
			defer cleanup()

			img, err := LoadImage(filePath)
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

			bounds := img.Bounds()
			if bounds.Dx() != 2 || bounds.Dy() != 2 {
				t.Errorf("Expected 2x2 image, got %dx%d", bounds.Dx(), bounds.Dy())
			}
		})
	}
}

func TestAlignImages(t *testing.T) {
	tests := []struct {
		name        string
		images      []hdr.Image
		expectError bool
	}{
		{
			name: "Aligned images",
			images: []hdr.Image{
				hdr.NewImageFromGoImage(createTestImage(2, 2, color.RGBA{R: 255, A: 255})),
				hdr.NewImageFromGoImage(createTestImage(2, 2, color.RGBA{G: 255, A: 255})),
				hdr.NewImageFromGoImage(createTestImage(2, 2, color.RGBA{B: 255, A: 255})),
			},
			expectError: false,
		},
		{
			name: "Misaligned images",
			images: []hdr.Image{
				hdr.NewImageFromGoImage(createTestImage(2, 2, color.RGBA{R: 255, A: 255})),
				hdr.NewImageFromGoImage(createTestImage(3, 3, color.RGBA{G: 255, A: 255})), // Different size
				hdr.NewImageFromGoImage(createTestImage(2, 2, color.RGBA{B: 255, A: 255})),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alignedImages, err := AlignImages(tt.images)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(alignedImages) != len(tt.images) {
				t.Errorf("Expected %d aligned images, got %d", len(tt.images), len(alignedImages))
			}

			// Check that aligned images have the same dimensions
			baseBounds := alignedImages[0].Bounds()
			for i, img := range alignedImages[1:] {
				if img.Bounds() != baseBounds {
					t.Errorf("Aligned image %d has different dimensions", i+1)
				}
			}
		})
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

	img := hdr.NewImageFromGoImage(createTestImage(2, 2, color.RGBA{R: 255, A: 255}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
