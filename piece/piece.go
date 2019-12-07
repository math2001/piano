package piece

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/math2001/piano/frac"
	"github.com/math2001/piano/wave"
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
	Duration frac.Frac

	// Start is the starting time, as a scaling of the beat
	Start frac.Frac
}

func (n Note) End() frac.Frac {
	return n.Start.Add(n.Duration)
}

// Piece is a collection of notes
type Piece struct {
	name string
	// float64 is a scalar describing the start time of each note relative to
	// the start of the piece (it scales one beat)
	Notes []Note
}

// a block is like a note but can have multiple streamer (which we mix).
// the streamers are all mixed from start during duration.
type block struct {
	duration frac.Frac
	// we don't need to know about start
	start frac.Frac
	// this will be become streamers when I switch to samples instead of
	// sine waves.
	frequencies []float64
}

func (b *block) end() frac.Frac {
	return b.start.Add(b.duration)
}

func (b *block) equal(target block) bool {

	// do the cheap comparaissons first
	equal := b.start == target.start
	equal = equal && b.duration == target.duration
	equal = equal && len(b.frequencies) == len(target.frequencies)

	for i := range b.frequencies {
		if b.frequencies[i] != target.frequencies[i] {
			return false
		}
	}

	return equal
}

// Play assumes that the speaker has been initialized
func (p *Piece) GetStreamer(sr beep.SampleRate, beat time.Duration) beep.Streamer {
	// how the algorithm works
	// get every marker
	// (a marker is the start or the end of a note. It's just a number)
	//
	// for each marker
	//     find every note that intersect (start <= prev_marker && end >= current_marker)
	//     mix all those notes together from prev_marker to current_marker

	blocks := p.intersectionBlocks()

	// a simple slice that we give to a sequencer (ie if no sound is going to
	// be played, we need to pass in silence, because otherwise the next block
	// will play straight away)

	var streamers []beep.Streamer
	for _, block := range blocks {
		nsamples := sr.N(beat * time.Duration(block.duration.Float()))

		var streamer beep.Streamer
		if len(block.frequencies) == 0 {
			streamer = beep.Silence(-1)
		} else if len(block.frequencies) == 1 {
			// FIXME: please do some caching. At least profile to check the
			// cost
			streamer = beep.Loop(-1, wave.NewSine(wave.N(sr, block.frequencies[0])))
		} else {
			mixer := &beep.Mixer{}
			for _, freq := range block.frequencies {
				mixer.Add(beep.Loop(-1, wave.NewSine(wave.N(sr, freq))))
			}
			// mixer only sums up the samples. That means if we sum up to 1s,
			// we get a two which isn't allowed. Instead, we want to take the
			// *average* of the different streamers. This is what gain does
			// here...
			streamer = &effects.Gain{
				Streamer: mixer,
				// this hacky thing is due to how Gain is implemented...
				Gain: 1.0/float64(mixer.Len()) - 1.0,
			}
		}
		streamers = append(streamers, beep.Take(nsamples, streamer))
	}

	return beep.Seq(streamers...)
}

func (p *Piece) intersectionBlocks() []block {

	markers := p.getMarkers()

	var blocks []block

	for i, currentMarker := range markers {
		if i == 0 {
			continue
		}
		prevMarker := markers[i-1]
		currentblock := block{
			start:    prevMarker,
			duration: currentMarker.Minus(prevMarker),
		}
		// FIXME: we can limit what how many notes are looping over here...
		for _, note := range p.Notes {
			intersect := note.Start.Float() <= prevMarker.Float()
			intersect = intersect && note.End().Float() >= currentMarker.Float()
			if intersect {
				currentblock.frequencies = append(currentblock.frequencies, note.Frequency)
			}
		}
		blocks = append(blocks, currentblock)
	}
	return blocks
}

// markers are where the notes start or finish
func (p *Piece) getMarkers() []frac.Frac {

	markers := make([]frac.Frac, len(p.Notes)*2+1)

	// make sure there is a marker at 0, to allow piece to start with complete
	// silence... (this is due to the fact that intersectionBlocks just bases
	// itself on markers exclusively, and hence goes straight to the first marker
	// all the time)
	markers[0] = frac.N(0)

	for i, note := range p.Notes {
		markers[i*2+1] = note.Start
		markers[i*2+2] = note.End()
	}

	sort.SliceStable(markers, func(i, j int) bool {
		return markers[i].Float() < markers[j].Float()
	})

	// remove duplicates
	j := 1
	for i := 1; i < len(markers); i++ {
		if markers[i] != markers[i-1] {
			markers[j] = markers[i]
			j++
		}
	}
	return markers[:j]
}

func (p *Piece) Render() {
	// we make duration and start integers (fraction with denominator 1)
	// so that every character is the lower fraction of time in the piece

	var dens []int

	frequencies := make(map[float64][]Note)
	for _, note := range p.Notes {
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

	k := frac.F(scaler, 1)

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

// FromBPM returns the duration of the one beat for a given bpm (beat per minute)
func FromBPM(bpm int) time.Duration {
	return time.Duration(60 * 1E9 / bpm)
}
