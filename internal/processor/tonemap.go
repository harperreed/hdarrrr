package processor

// ToneMapper interface defines the tone mapping operation
type ToneMapper interface {
	ToneMap(value float64) float64
}

// ReinhardToneMapper implements Reinhard tone mapping
type ReinhardToneMapper struct{}

func NewReinhardToneMapper() *ReinhardToneMapper {
	return &ReinhardToneMapper{}
}

func (t *ReinhardToneMapper) ToneMap(v float64) float64 {
	return v / (1 + v)
}
