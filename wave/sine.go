package wave

import (
	"math"
)

type Sine struct {
	pos int
	len int
	buf [][2]float64
}

func (s *Sine) Stream(target [][2]float64) (n int, ok bool) {
	for i := range target {
		target[i] = s.buf[s.pos]
		s.pos += 1
		if s.pos == s.len {
			s.pos = 0
		}
	}
	return len(target), true
}

func (s *Sine) Err() error {
	return nil
}

func (s *Sine) Len() int {
	return s.len
}

func (s *Sine) Position() int {
	return s.pos
}

func (s *Sine) Seek(pos int) error {
	// FIXME: should have a lock here?
	s.pos = pos
	return nil
}

// NewSine returns a sine wave generator which implementes the StreamSeeker
// interface. Length is the number of samples it should use to complete one
// wave
// To use it, just make it loop
func NewSine(length int) *Sine {
	// FIXME: should I manually check for negative lengths?
	buf := make([][2]float64, length)
	k := 2 * math.Pi / float64(length)

	for x := 0; x < length; x++ {
		v := math.Sin(k * float64(x))
		if math.Abs(v) < 1E-15 {
			v = 0
		}
		buf[x][0] = v
		buf[x][1] = v
	}

	return &Sine{
		pos: 0,
		len: length,
		buf: buf,
	}
}
