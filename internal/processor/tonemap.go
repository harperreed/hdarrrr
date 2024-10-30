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

// DragoToneMapper implements Drago tone mapping
type DragoToneMapper struct {
	LdMax float64
	B     float64
}

func NewDragoToneMapper(ldMax, b float64) *DragoToneMapper {
	return &DragoToneMapper{
		LdMax: ldMax,
		B:     b,
	}
}

// ToneMap implements the Drago tone mapping operator.
func (t *DragoToneMapper) ToneMap(v float64) float64 {
	// Clamp negative values to 0 since light intensity cannot be negative
	v = math.Max(0, v)
	return (t.LdMax * (v * (1 + (v / (t.B * t.B))))) / (1 + v)
}
