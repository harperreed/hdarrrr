package align

import (
	"errors"
	"fmt"
	"image"
)

// AlignImages aligns multiple images using feature matching
func AlignImages(images []image.Image) ([]image.Image, error) {
	if len(images) < 2 {
		return nil, errors.New("at least two images are required for alignment")
	}

	// Validate image dimensions
	baseBounds := images[0].Bounds()
	for i, img := range images[1:] {
		if img.Bounds() != baseBounds {
			return nil, fmt.Errorf("image %d has different dimensions than the base image", i+1)
		}
	}

	// For now, just return the original images since we haven't implemented
	// the actual alignment algorithm yet
	return images, nil
}
