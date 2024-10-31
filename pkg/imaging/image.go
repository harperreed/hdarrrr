package imaging

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"

	"github.com/mdouchement/hdr"
	"github.com/mdouchement/hdr/hdrcolor"
)

// SupportedFormats contains the file extensions we support
var SupportedFormats = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
}

// ImageProperties holds the properties of an image
type ImageProperties struct {
	Width      int
	Height     int
	ColorDepth int
	Format     string
}

// GetImageProperties retrieves the properties of an image
func GetImageProperties(img image.Image, format string) ImageProperties {
	bounds := img.Bounds()
	colorDepth := 8 // Assuming 8 bits per channel for simplicity
	return ImageProperties{
		Width:      bounds.Dx(),
		Height:     bounds.Dy(),
		ColorDepth: colorDepth,
		Format:     format,
	}
}

// ValidateImageProperties checks if two images have the same properties
func ValidateImageProperties(baseProps, props ImageProperties) bool {
	return baseProps.Width == props.Width &&
		baseProps.Height == props.Height &&
		baseProps.ColorDepth == props.ColorDepth &&
		baseProps.Format == props.Format
}

// LoadImages loads multiple images from file paths
func LoadImages(paths ...string) ([]image.Image, error) {
	images := make([]image.Image, len(paths))

	for i, path := range paths {
		img, err := LoadImage(path)
		if err != nil {
			return nil, err
		}
		images[i] = img
	}

	// Validate image properties
	if len(images) > 1 {
		baseProps := GetImageProperties(images[0], strings.ToLower(path.Ext(paths[0])))
		for i, img := range images[1:] {
			props := GetImageProperties(img, strings.ToLower(path.Ext(paths[i+1])))
			if !ValidateImageProperties(baseProps, props) {
				return nil, errors.New("image properties do not match")
			}
		}
	}

	return images, nil
}

// LoadImage loads a single image from a file path
func LoadImage(filepath string) (image.Image, error) {
	ext := strings.ToLower(path.Ext(filepath))
	if !SupportedFormats[ext] {
		return nil, errors.New("unsupported image format: " + ext + ". Supported formats: PNG, JPEG")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var img image.Image
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	}

	if err != nil {
		return nil, err
	}

	// Convert to HDR format
	return convertToHDR(img), nil
}

// SaveImage saves an image to a file path
func SaveImage(img image.Image, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	ext := strings.ToLower(path.Ext(outputPath))
	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
	case ".png":
		return png.Encode(file, img)
	default:
		return errors.New("unsupported output format: " + ext + ". Supported formats: PNG, JPEG")
	}
}

// convertToHDR converts a regular image to HDR format
func convertToHDR(img image.Image) hdr.Image {
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
