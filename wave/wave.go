package wave

import (
	"math"

	"github.com/faiface/beep"
)

// Sine generates a sine wave using n sample per period
//
// Hence, the actual pitch of the wave will depend on your sample rate
// FIXME: cache one period and then read from that
func Sine(n int) beep.Streamer {
	var k float64 = 2 * math.Pi / float64(n)
	buf := make([][2]float64, n)
	for x := 0; x < n; x++ {
		v := math.Sin(k*float64(x)) * 0.1
		if math.Abs(v) < 1E-15 {
			v = 0
		}
		buf[x][0] = v
		buf[x][1] = v
	}
	x := 0

	return beep.StreamerFunc(func(samples [][2]float64) (int, bool) {
		for i := range samples {
			samples[i] = buf[x]
			x += 1
			if x == len(buf) {
				x = 0
			}
		}
		return len(samples), true
	})
}
