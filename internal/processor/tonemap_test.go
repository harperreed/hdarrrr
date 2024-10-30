package processor

import (
	"math"
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
			name:     "Very large value",
			input:    1000.0,
			expected: 0.999,
		},
		{
			name:     "Negative value",
			input:    -1.0,
			expected: 0.0, // Negative values should be clamped to 0
		},
		{
			name:     "Small positive value",
			input:    0.1,
			expected: 0.0909090909090909,
		},
	}

	toneMapper := NewReinhardToneMapper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toneMapper.ToneMap(tt.input)

			// Use approximate equality for floating point comparisons
			if !approximatelyEqual(result, tt.expected, 0.0001) {
				t.Errorf("Expected approximately %f, got %f", tt.expected, result)
			}
		})
	}
}

// approximatelyEqual checks if two float64 values are equal within a small epsilon
func approximatelyEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

// Add benchmark tests
func BenchmarkReinhardToneMapper(b *testing.B) {
	toneMapper := NewReinhardToneMapper()
	values := []float64{0.0, 0.5, 1.0, 2.0, 10.0, 100.0, -1.0}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			toneMapper.ToneMap(v)
		}
	}
}
