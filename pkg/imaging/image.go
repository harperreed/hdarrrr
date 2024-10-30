package imaging

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"
)

// SupportedFormats contains the file extensions we support
var SupportedFormats = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
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
	return img, nil
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
