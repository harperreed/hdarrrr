package processor

import (
	"testing"
)

func TestReinhardToneMapper(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "Zero value",
			input:    0.0,
			expected: 0.0,
		},
		{
			name:     "Value of 1",
			input:    1.0,
			expected: 0.5,
		},
		{
			name:     "Large value",
			input:    10.0,
			expected: 0.9090909090909091,
		},
		{
			name:     "Negative value",
			input:    -1.0,
			expected: -0.5,
		},
	}

	toneMapper := NewReinhardToneMapper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toneMapper.ToneMap(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

