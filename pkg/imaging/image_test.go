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
	"github.com/mdouchement/hdr/hdrcolor"
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

// createHDRImage creates an HDR image from a regular image
func createHDRImage(width, height int, c color.Color) hdr.Image {
	img := createTestImage(width, height, c)
	bounds := img.Bounds()
	hdrImg := hdr.NewRGB(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			hdrImg.Set(x, y, hdrcolor.RGB{
				R: float64(r) / 0xffff,
				G: float64(g) / 0xffff,
				B: float64(b) / 0xffff,
			})
		}
	}
	return hdrImg
}

// createTestImageFile creates a temporary image file for testing
func createTestImageFile(t *testing.T, format string, c color.Color) (string, func()) {
	img := createTestImage(2, 2, c)

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

	img := createTestImage(2, 2, color.RGBA{R: 255, A: 255})

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

func TestLoadImagesWithDifferentProperties(t *testing.T) {
	file1, cleanup1 := createTestImageFile(t, "png", color.RGBA{R: 255, A: 255})
	defer cleanup1()
	file2, cleanup2 := createTestImageFile(t, "jpeg", color.RGBA{G: 255, A: 255})
	defer cleanup2()

	_, err := LoadImages(file1, file2)
	if err == nil {
		t.Error("Expected error due to different image properties, got nil")
	}
}

func TestLoadImageUnsupportedFormat(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_*.bmp")
	if err != nil {
		t.Fatal("Failed to create temp file:", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = LoadImage(tmpFile.Name())
	if err == nil {
		t.Error("Expected error due to unsupported format, got nil")
	}
}

func TestSaveImageUnsupportedFormat(t *testing.T) {
	img := createTestImage(2, 2, color.RGBA{R: 255, A: 255})

	tmpFile := filepath.Join(os.TempDir(), "test_output.bmp")
	defer os.Remove(tmpFile)

	err := SaveImage(img, tmpFile)
	if err == nil {
		t.Error("Expected error due to unsupported format, got nil")
	}
}
