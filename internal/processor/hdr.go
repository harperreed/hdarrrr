package processor

import (
	"errors"
	"fmt"
	"image"
	"image/color"

	"github.com/harperreed/hdarrrr/pkg/align"
)

// HDRPixel represents a pixel with high dynamic range values
type HDRPixel struct {
	R, G, B float64
}

// HDRProcessor handles the HDR image processing
type HDRProcessor struct {
	toneMapper ToneMapper
}

// NewHDRProcessor creates a new HDR processor with default settings
func NewHDRProcessor() *HDRProcessor {
	return &HDRProcessor{
		toneMapper: NewReinhardToneMapper(),
	}
}

// Process creates an HDR image from three exposure images
func (p *HDRProcessor) Process(images []image.Image) (image.Image, error) {
	if len(images) != 3 {
		return nil, errors.New("exactly three images are required")
	}

	// Validate image properties
	if err := validateImageProperties(images); err != nil {
		return nil, err
	}

	// Align images
	alignedImages, err := align.AlignImages(images)
	if err != nil {
		return nil, fmt.Errorf("image alignment failed: %v", err)
	}

	width := alignedImages[0].Bounds().Max.X
	height := alignedImages[0].Bounds().Max.Y

	// Create HDR image
	hdrImage := make([][]HDRPixel, height)
	for i := range hdrImage {
		hdrImage[i] = make([]HDRPixel, width)
	}

	// Merge exposures into HDR image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r1, g1, b1, _ := alignedImages[0].At(x, y).RGBA()
			r2, g2, b2, _ := alignedImages[1].At(x, y).RGBA()
			r3, g3, b3, _ := alignedImages[2].At(x, y).RGBA()

			hdrImage[y][x] = HDRPixel{
				R: p.mergeExposures(float64(r1)/65535, float64(r2)/65535, float64(r3)/65535),
				G: p.mergeExposures(float64(g1)/65535, float64(g2)/65535, float64(g3)/65535),
				B: p.mergeExposures(float64(b1)/65535, float64(b2)/65535, float64(b3)/65535),
			}
		}
	}

	// Tone map HDR image to LDR
	output := image.NewRGBA(alignedImages[0].Bounds())
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := hdrImage[y][x]
			r := p.toneMapper.ToneMap(pixel.R)
			g := p.toneMapper.ToneMap(pixel.G)
			b := p.toneMapper.ToneMap(pixel.B)

			output.Set(x, y, color.RGBA{
				R: uint8(r * 255),
				G: uint8(g * 255),
				B: uint8(b * 255),
				A: 255,
			})
		}
	}

	return output, nil
}

func (p *HDRProcessor) mergeExposures(v1, v2, v3 float64) float64 {
	w1 := p.weight(v1)
	w2 := p.weight(v2)
	w3 := p.weight(v3)

	sumWeights := w1 + w2 + w3
	if sumWeights == 0 {
		return 0
	}

	return (v1*w1 + v2*w2 + v3*w3) / sumWeights
}

func (p *HDRProcessor) weight(v float64) float64 {
	if v <= 0 || v >= 1 {
		return 0
	}
	return 1 - (2*v-1)*(2*v-1)
}

// Only showing the validation function since the rest of the file remains the same
func validateImageProperties(images []image.Image) error {
	if len(images) < 2 {
		return errors.New("at least two images are required for validation")
	}

	baseImg := images[0]
	baseColorModel := baseImg.ColorModel()
	baseBounds := baseImg.Bounds()

	for i, img := range images[1:] {
		// Check dimensions
		if img.Bounds() != baseBounds {
			return fmt.Errorf("image %d has different dimensions than the first image", i+2)
		}

		// Check color model
		if img.ColorModel() != baseColorModel {
			return fmt.Errorf("image %d has a different color model than the first image", i+2)
		}
	}

	return nil
}
