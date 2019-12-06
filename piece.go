package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/faiface/beep"
)

// Note describes how a single note is played
// right now, it's very simplistic
type Note struct {
	// Volume == 0 -> remains unchanged. < 0 decrease volume, > 0 increase volume
	Volume float64

	// Frequency is the pitch of the note
	Frequency float64

	// Duration is the coefficient of a beat
	// ie. 1 is the duration one beat
	//     2 is the duration two beats
	//     .5 is the duration of half a beat
	Duration Frac

	// Start is the starting time, as a scaling of the beat
	Start Frac
}

func (n Note) End() Frac {
	return n.Start.Add(n.Duration)
}

// Piece is a collection of notes
type Piece struct {
	name string
	// float64 is a scalar describing the start time of each note relative to
	// the start of the piece (it scales one beat)
	notes []Note
}

type block struct {
	start    Frac
	duration Frac
	// this will be become streamers when I switch to samples instead of
	// sine waves.
	frequencies []float64
}

func (b *block) end() Frac {
	return b.start.Add(b.duration)
}

func (b *block) equal(target block) bool {
	// FIXME: better comparaisson of streamers please
	return b.start.Equal(target.start) && b.duration.Equal(target.duration) && len(b.frequencies) == len(target.frequencies)
}

// Play assumes that the speaker has been initialized
func (p *Piece) Play(sr beep.SampleRate, beat time.Duration) {
}

func (p *Piece) intersectionBlocks() []block {
	// estimate the number of intersections
	intersections := make([]block, 0, len(p.notes))
	for _, a := range p.notes {
		for _, b := range p.notes {
			var splits []block
			if a.End().Equal(b.End()) && b.Start.Equal(b.Start) {
				splits = []block{
					{
						start:       a.Start,
						duration:    a.Duration,
						frequencies: []float64{a.Frequency, b.Frequency},
					},
				}
			} else {
				var hasOverlaps bool
				splits, hasOverlaps = overlapsFrom(a, b)
				if !hasOverlaps {
					// they don't overlap. Since we assume that the notes are
					// sorted, we can conclude all the remaining notes won't
					// intersect
					break
				}
			}
			intersections = append(intersections, splits...)

		}
	}
	return intersections
}

// returns the splits from two notes
// Case 1:
// |  A |  B  | C      |
// ***********
//       ***************
// Case 2:
// |A |   B  | C  |
// ***************
//     ******
// hence, it always returns 3 splits
func overlapsFrom(a Note, b Note) (splits []block, doOverlap bool) {
	if a.Start.Float() > b.Start.Float() {
		panic(fmt.Sprintf("wrong order: a.start has to be less than b.start (shouldn't happen, notes are supposed to be sorted) %v %v", a, b))
	}

	if a.Start == b.Start && a.End() == b.End() {
		panic("a and b are equivalent. This is a simple edge case, handle it yourself")
	}

	doOverlap = a.End().Float() > b.Start.Float() && a.End().Float() < b.End().Float()

	isFirstCase := a.End().Float() < b.End().Float()

	if !doOverlap {
		return splits, false
	}

	splits = make([]block, 3)

	splits[0] = block{
		start:       a.Start,
		duration:    b.Start.Minus(a.Start),
		frequencies: []float64{a.Frequency},
	}

	duration := b.Duration
	if isFirstCase {
		duration = a.End().Minus(b.Start)
	}

	splits[1] = block{
		start:       b.Start,
		duration:    duration,
		frequencies: []float64{a.Frequency, b.Frequency},
	}

	freq := a.Frequency
	start := b.End()
	if isFirstCase {
		freq = b.Frequency
		start = a.End()
	}

	splits[2] = block{
		start:       start,
		duration:    b.End().Minus(a.End()).Abs(),
		frequencies: []float64{freq},
	}

	return splits, true
}

func (p *Piece) Render() {
	// we make duration and start integers (fraction with denominator 1)
	// so that every character is the lower fraction of time in the piece

	var dens []int

	frequencies := make(map[float64][]Note)
	for _, note := range p.notes {
		frequencies[note.Frequency] = append(frequencies[note.Frequency], note)
		dens = append(dens, note.Duration.Den(), note.Start.Den())
	}

	// compute the smallest number which is a product of all the different
	// primes composing all the denominators
	sort.Ints(dens)
	scaler := 1
	for _, n := range dens {
		if scaler%n != 0 {
			scaler = scaler * n
		}
	}

	k := Frac{scaler, 1}

	for freq, notes := range frequencies {
		fmt.Printf("%3.0f: ", freq)
		// we assume the notes are sorted
		pos := 0
		for _, note := range notes {
			start := note.Start.Multiply(k).Num()
			width := note.Duration.Multiply(k).Num()
			if pos > start {
				panic(fmt.Sprintf("current position %d is later than required start position of %v: %d", pos, note, start))
			}
			fmt.Print(strings.Repeat(" ", start-pos))
			fmt.Print(strings.Repeat("*", width))
			pos = start + width
		}
		fmt.Println()
	}
}
