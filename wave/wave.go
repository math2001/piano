package wave

import (
	"time"

	"github.com/faiface/beep"
)

// N returns the number of samples needed for one period of a wave of freqency
// freq
func N(sr beep.SampleRate, freq float64) int {
	// T * f = 1
	period := time.Second / time.Duration(freq)

	samplesPerPeriod := sr.N(period)

	return samplesPerPeriod
}
