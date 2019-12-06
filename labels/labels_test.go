package labels

import (
	"math"
	"testing"
)

// todo: test caching

func TestLabelsCases(t *testing.T) {
	names_freq := []struct {
		name string
		freq float64
	}{
		{"A0", 27.5},
		{"A4", 440},
		{"A#7", 3729.310},

		{"F6", 1396.913},
		{"C4", 261.6256},

		{"Eb1", 38.89087},
		{"D#1", 38.89087},

		// special behaviors

		// Cb == B below
		{"Cb4", 246.9417},
		{"B3", 246.9417},

		// E# == F
		{"E#5", 698.4565},
		{"F5", 698.4565},

		// B# == C above
		{"B#2", 130.8128},
		{"C3", 130.8128},

		// Fb == E
		{"Fb7", 2637.020},
		{"E7", 2637.020},
	}

	labels := NewLabels()

	for _, obj := range names_freq {
		actual, err := labels.Frequency(obj.name)
		if err != nil {
			t.Errorf("name: %q, expected: pass, got err: %s", obj.name, err)
			continue
		}

		// round because labels uses math to calculate the frequencies, and
		// hence can be arbitrarily more precise than the fixed dp we have in
		// the tests
		if roundTo(actual, 2) != roundTo(obj.freq, 2) {
			t.Errorf("name: %q, expected: %f, got: %f", obj.name, obj.freq, actual)
		}
	}
}

func roundTo(n float64, dp int) float64 {
	return math.Round(n*math.Pow(10, float64(dp))) / math.Pow(10, float64(dp))
}
