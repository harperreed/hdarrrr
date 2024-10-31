// align.go
package align

import (
	"errors"
	"fmt"
	"image"
)

// Aligner defines the interface for image alignment implementations
type Aligner interface {
	Align(images []image.Image) ([]image.Image, error)
}

// BasicAligner provides simple dimension validation
type BasicAligner struct{}

// NewBasicAligner creates a new BasicAligner
func NewBasicAligner() *BasicAligner {
	return &BasicAligner{}
}

// Align validates image dimensions and returns the original images.
// This implementation ensures images are the same size but does not perform
// any pixel-level alignment.
func (a *BasicAligner) Align(images []image.Image) ([]image.Image, error) {
	if len(images) < 2 {
		return nil, errors.New("at least two images are required for alignment")
	}

	baseBounds := images[0].Bounds()
	for i, img := range images[1:] {
		if img == nil {
			return nil, fmt.Errorf("image %d is nil", i+1)
		}
		if img.Bounds() != baseBounds {
			return nil, fmt.Errorf("image %d has different dimensions than the base image", i+1)
		}
	}

	return images, nil
}

// AlignImages is a convenience function that uses the BasicAligner
func AlignImages(images []image.Image) ([]image.Image, error) {
	aligner := NewBasicAligner()
	return aligner.Align(images)
}
