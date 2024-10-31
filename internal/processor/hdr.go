package processor

import (
	"errors"
	"fmt"
	"image"

	"github.com/mdouchement/hdr"
	"github.com/mdouchement/hdr/hdrcolor"
	"github.com/mdouchement/hdr/tmo"
)

// HDRProcessor handles HDR image processing
type HDRProcessor struct {
	toneMapper string
	params     map[string]float64
}

// NewHDRProcessor creates a new HDR processor with default settings
func NewHDRProcessor() *HDRProcessor {
	return &HDRProcessor{
		toneMapper: "reinhard05",
		params: map[string]float64{
			"gamma":     1.0,
			"intensity": 1.0,
			"light":     0.0,
		},
	}
}

// WithToneMapper sets the tone mapping operator
func (p *HDRProcessor) WithToneMapper(mapper string) *HDRProcessor {
	p.toneMapper = mapper
	return p
}

// WithParams sets tone mapping parameters
func (p *HDRProcessor) WithParams(params map[string]float64) *HDRProcessor {
	for k, v := range params {
		p.params[k] = v
	}
	return p
}

// Process creates an HDR image from multiple exposure images
func (p *HDRProcessor) Process(images []image.Image) (image.Image, error) {
	if len(images) < 2 {
		return nil, errors.New("at least two images are required")
	}

	// Convert regular images to HDR images
	hdrImages := make([]hdr.Image, len(images))
	for i, img := range images {
		if img == nil {
			return nil, fmt.Errorf("image %d is nil", i+1)
		}

		// Handle both HDR and LDR inputs
		if hdrImg, ok := img.(hdr.Image); ok {
			hdrImages[i] = hdrImg
		} else {
			hdrImages[i] = convertToHDR(img)
		}
	}

	// Validate image properties
	if err := validateImageProperties(hdrImages); err != nil {
		return nil, err
	}

	// Create merged HDR image
	bounds := hdrImages[0].Bounds()
	merged := hdr.NewRGB(bounds)

	// Simple exposure fusion (average method)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var sumR, sumG, sumB float64
			for _, img := range hdrImages {
				r, g, b, _ := img.HDRAt(x, y).HDRRGBA()
				sumR += r
				sumG += g
				sumB += b
			}
			n := float64(len(hdrImages))
			merged.Set(x, y, hdrcolor.RGB{
				R: sumR / n,
				G: sumG / n,
				B: sumB / n,
			})
		}
	}

	// Apply tone mapping
	var tm tmo.ToneMappingOperator
	switch p.toneMapper {
	case "reinhard05":
		tm = tmo.NewReinhard05(merged, p.params["intensity"], p.params["light"], p.params["gamma"])
	case "drago03":
		tm = tmo.NewDrago03(merged, p.params["gamma"])
	default:
		return nil, fmt.Errorf("unsupported tone mapper: %s", p.toneMapper)
	}

	// Return tone mapped image
	return tm.Perform(), nil
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

// validateImageProperties checks if all images have matching properties
func validateImageProperties(images []hdr.Image) error {
	if len(images) < 2 {
		return errors.New("at least two images are required")
	}

	base := images[0]
	baseBounds := base.Bounds()

	for i, img := range images[1:] {
		if img == nil {
			return fmt.Errorf("image %d is nil", i+2)
		}

		if img.Bounds() != baseBounds {
			return fmt.Errorf("image %d has different dimensions than the first image", i+2)
		}
	}

	return nil
}
