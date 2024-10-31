package processor

import (
	"errors"
	"fmt"
	"image"
	"math"

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
			"gamma":      1.0,
			"intensity":  1.0,
			"light":      0.0,
			"saturation": 0.8,
			"contrast":   4.0,
			"chromatic":  0.0,
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

func (p *HDRProcessor) validateImageSize(img image.Image) error {
	bounds := img.Bounds()
	if p.toneMapper == "icam06" && (bounds.Dx() < 32 || bounds.Dy() < 32) {
		return fmt.Errorf("ICam06 operator requires images of at least 32x32 pixels, got %dx%d",
			bounds.Dx(), bounds.Dy())
	}
	return nil
}

// Process creates an HDR image from multiple exposure images
func (p *HDRProcessor) Process(images []image.Image) (image.Image, error) {
	if len(images) < 2 {
		return nil, errors.New("at least two images are required")
	}

	// Validate image size requirements
	if err := p.validateImageSize(images[0]); err != nil {
		return nil, err
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

	// Weighted exposure fusion
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var sumR, sumG, sumB, sumWeight float64
			for _, img := range hdrImages {
				r, g, b, _ := img.HDRAt(x, y).HDRRGBA()
				// Calculate weight based on pixel brightness
				brightness := (r + g + b) / 3.0
				weight := p.calculateWeight(brightness)
				sumR += r * weight
				sumG += g * weight
				sumB += b * weight
				sumWeight += weight
			}

			if sumWeight > 0 {
				merged.Set(x, y, hdrcolor.RGB{
					R: sumR / sumWeight,
					G: sumG / sumWeight,
					B: sumB / sumWeight,
				})
			}
		}
	}

	// Apply tone mapping
	var tm tmo.ToneMappingOperator
	switch p.toneMapper {
	case "reinhard05":
		tm = tmo.NewReinhard05(merged, p.params["intensity"], p.params["light"], p.params["gamma"])
	case "drago03":
		tm = tmo.NewDrago03(merged, p.params["gamma"])
	case "linear":
		tm = tmo.NewLinear(merged)
	case "logarithmic":
		tm = tmo.NewLogarithmic(merged)
	case "durand":
		tm = tmo.NewDurand(merged, p.params["saturation"])
	case "custom_reinhard05":
		tm = tmo.NewCustomReinhard05(merged, p.params["intensity"], p.params["light"], p.params["gamma"])
	case "icam06":
		tm = tmo.NewICam06(merged, p.params["gamma"], p.params["contrast"], p.params["chromatic"])
	default:
		return nil, fmt.Errorf("unsupported tone mapper: %s", p.toneMapper)
	}

	return tm.Perform(), nil
}

// calculateWeight returns a weight for exposure fusion based on pixel brightness
func (p *HDRProcessor) calculateWeight(v float64) float64 {
	// Weight function favoring mid-range values
	const mid = 0.5
	diff := v - mid
	return math.Exp(-diff * diff * 4)
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
