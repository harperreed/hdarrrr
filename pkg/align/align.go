package align

import (
	"fmt"
	"image"
)

// AlignImages aligns multiple images using feature matching
func AlignImages(images []image.Image) ([]image.Image, error) {
	if len(images) < 2 {
		return nil, fmt.Errorf("at least two images are required for alignment")
	}

	// Validate image dimensions
	baseBounds := images[0].Bounds()
	for i, img := range images[1:] {
		if img == nil {
			return nil, fmt.Errorf("image %d is nil", i+2)
		}
		if img.Bounds() != baseBounds {
			return nil, fmt.Errorf("image %d has different dimensions than the first image", i+2)
		}
	}

	// For now, return the original images since actual alignment is not implemented yet
	return images, nil
}
