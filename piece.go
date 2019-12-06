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

// Piece is a collection of notes
type Piece struct {
	name string
	// float64 is a scalar describing the start time of each note relative to
	// the start of the piece (it scales one beat)
	notes []Note
}

type block struct {
	start     float64
	duration  float64
	streamers []beep.Streamer
}

func (b *block) end() float64 {
	return b.start + b.duration
}

// Play assumes that the speaker has been initialized
func (p *Piece) Play(sr beep.SampleRate, beat time.Duration) {
	// intersections := make([]block, len(p.notes))
}

func (p *Piece) Render() {
	// we make duration and start integers (fraction with denominator 1)
	// so that every character is the lower fraction of time in the piece

	var dens []int

	frequencies := make(map[float64][]Note)
	for _, note := range p.notes {
		frequencies[note.Frequency] = append(frequencies[note.Frequency], note)
		dens = append(dens, note.Duration.Den, note.Start.Den)
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
			start := note.Start.Multiply(k).Num
			width := note.Duration.Multiply(k).Num
			if pos > start {
				panic(fmt.Sprintf("current position %d is later than required start position of %v", pos, note))
			}
			fmt.Print(strings.Repeat(" ", start-pos))
			fmt.Print(strings.Repeat("*", width))
			pos = start + width
		}
		fmt.Println()
	}
}
