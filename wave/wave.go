package wave

import (
	"time"

	"github.com/faiface/beep"
)

// N returns, based on the sample rate (typically of the speaker) and the
// wave's desired frequency, the number of samples needed for one period to
// complete
func N(sr beep.SampleRate, freq float64) int {
	// T * f = 1
	period := time.Second / time.Duration(freq)

	samplesPerPeriod := sr.N(period)

	return samplesPerPeriod
}
