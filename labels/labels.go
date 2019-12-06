package labels

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"unicode"
)

// this should just generate some go code defining constants

var ErrParsingName = errors.New("parsing name")
var tones = map[string]int{
	"C":  0,
	"B#": 12, // because we go up an octave

	"C#": 1,
	"Db": 1,

	"D": 2,

	"D#": 3,
	"Eb": 3,

	"E":  4,
	"Fb": 4,

	"F":  5,
	"E#": 5,

	"F#": 6,
	"Gb": 6,

	"G": 7,

	"G#": 8,
	"Ab": 8,

	"A": 9,

	"A#": 10,
	"Bb": 10,

	"B":  11,
	"Cb": -1, // because we go down an octave
}

type Labels struct {
	// names caches names to index
	names map[string]int
	// freqs caches index to frequencies
	freqs map[int]float64
}

func NewLabels() *Labels {
	return &Labels{
		names: make(map[string]int),
		freqs: make(map[int]float64),
	}
}

func (n *Labels) FromIndex(i int) (freq float64) {
	if freq, ok := n.freqs[i]; ok {
		return freq
	}
	n.freqs[i] = math.Pow(2, (float64(i)-49)/12) * 440
	return n.freqs[i]
}

func (n *Labels) name(name string) (index int, err error) {
	if len(name) != 2 && len(name) != 3 {
		return 0, fmt.Errorf("invalid length (need 2 or 3) in %q: %w", name, ErrParsingName)
	}

	name = fmt.Sprintf("%c%s", unicode.ToUpper(rune(name[0])), name[1:])

	if index, ok := n.names[name]; ok {
		return index, nil
	}

	var octave int

	octave, err = strconv.Atoi(string(name[len(name)-1]))
	if err != nil {
		return 0, fmt.Errorf("parsing octave in %q: %s (%w)", name, err, ErrParsingName)
	}

	toneIndex, ok := tones[name[:len(name)-1]]
	if !ok {
		return 0, fmt.Errorf("parsing tone in %q: unknown tone (%w)", name, ErrParsingName)
	}

	n.names[name] = octave*12 + toneIndex - 8
	return n.names[name], nil
}

func (n *Labels) Frequency(name string) (freq float64, err error) {
	index, err := n.name(name)
	if err != nil {
		return 0, err
	}
	return n.FromIndex(index), nil
}

// F is the same as Frequency, except it panics if there is an error. Don't
// pass in user data, just static strings
func (n *Labels) F(name string) float64 {
	index, err := n.name(name)
	if err != nil {
		panic(fmt.Sprintf("%s (if you don't want panic, use .Frequency instead)", err))
	}
	return n.FromIndex(index)
}

// todo
// func (n *labels) Label(freq float64) (string, error) {

// }
