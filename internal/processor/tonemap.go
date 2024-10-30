package processor

import "math"

// ToneMapper interface defines the tone mapping operation
type ToneMapper interface {
	ToneMap(value float64) float64
}

// ReinhardToneMapper implements Reinhard tone mapping
type ReinhardToneMapper struct{}

func NewReinhardToneMapper() *ReinhardToneMapper {
	return &ReinhardToneMapper{}
}

// ToneMap implements the Reinhard tone mapping operator
func (t *ReinhardToneMapper) ToneMap(v float64) float64 {
	// Clamp negative values to 0 since light intensity cannot be negative
	v = math.Max(0, v)
	return v / (1 + v)
}

// DragoToneMapper implements Drago tone mapping
type DragoToneMapper struct {
	LdMax float64 // Maximum display luminance (cd/m²)
	B     float64 // Bias parameter
}

func NewDragoToneMapper(ldMax, b float64) *DragoToneMapper {
	if ldMax <= 0 {
		ldMax = 100.0 // Default display luminance
	}
	if b <= 0 {
		b = 0.85 // Default bias
	}
	return &DragoToneMapper{
		LdMax: ldMax,
		B:     b,
	}
}

// ToneMap implements the Drago tone mapping operator
func (t *DragoToneMapper) ToneMap(v float64) float64 {
	// Clamp negative values to 0
	v = math.Max(0, v)

	// Avoid log(0)
	const eps = 1e-6
	v = math.Max(v, eps)

	// Calculate bias parameter
	biasParam := math.Log(t.B) / math.Log(0.5)

	// Scale factor adjustment
	const scaleFactor = 0.01

	// Apply Drago operator
	numerator := math.Log1p(v * scaleFactor)
	denominator := math.Log(2 + 8*math.Pow((v/t.LdMax), biasParam))

	// Calculate result with proper scaling
	result := numerator / denominator

	// Scale to display range and ensure output is in [0,1]
	return math.Max(0, math.Min(1, result))
}
