
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

// ToneMap implements the Reinhard tone mapping operator.
// It clamps negative values to 0 since HDR values represent light intensity
// which cannot be negative in reality.
func (t *ReinhardToneMapper) ToneMap(v float64) float64 {
	// Clamp negative values to 0 since light intensity cannot be negative
	v = math.Max(0, v)
	return v / (1 + v)
}

