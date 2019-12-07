package piece

import (
	"testing"
	"time"

	"github.com/math2001/piano/frac"
)

func TestIntersectionSimultaneous(t *testing.T) {
	p := &Piece{
		// 440: **
		// 523: **
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.N(2),
				Start:     frac.N(0),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.N(2),
				Start:     frac.N(0),
			},
		},
	}

	actual := p.intersectionBlocks()
	expected := []block{
		// two streamers, don't really know how it's gonna be implemented
		{start: frac.N(0), duration: frac.N(2), frequencies: []float64{440, 523.25}},
	}
	CompareBlocks(t, actual, expected)
}

func TestIntersectionContainingOverlap(t *testing.T) {
	p := &Piece{
		// 440: ***
		// 523:  *
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.N(3),
				Start:     frac.N(0),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.N(1),
				Start:     frac.N(1),
			},
		},
	}
	actual := p.intersectionBlocks()
	expected := []block{
		{start: frac.N(0), duration: frac.N(1), frequencies: []float64{440}},
		{start: frac.N(1), duration: frac.N(1), frequencies: []float64{440, 523.25}},
		{start: frac.N(2), duration: frac.N(1), frequencies: []float64{440}},
	}
	CompareBlocks(t, actual, expected)
}

func TestIntersectionSilence(t *testing.T) {
	p := &Piece{
		// 440:   *
		// 523:     ***
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.N(1),
				Start:     frac.N(2),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.N(3),
				Start:     frac.N(4),
			},
		},
	}
	actual := p.intersectionBlocks()
	expected := []block{
		{start: frac.N(0), duration: frac.N(2), frequencies: []float64{}},
		{start: frac.N(2), duration: frac.N(1), frequencies: []float64{440}},
		{start: frac.N(3), duration: frac.N(1), frequencies: []float64{}},
		{start: frac.N(4), duration: frac.N(3), frequencies: []float64{523.25}},
	}
	CompareBlocks(t, actual, expected)
}
func TestIntersectionOverlap(t *testing.T) {
	p := &Piece{
		// 440: ***
		// 523:  ***
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.N(3),
				Start:     frac.N(0),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.N(3),
				Start:     frac.N(1),
			},
		},
	}
	actual := p.intersectionBlocks()
	expected := []block{
		{start: frac.N(0), duration: frac.N(1), frequencies: []float64{440}},
		{start: frac.N(1), duration: frac.N(2), frequencies: []float64{440, 523.25}},
		{start: frac.N(3), duration: frac.N(1), frequencies: []float64{523.25}},
	}
	CompareBlocks(t, actual, expected)
}

func TestGetMarkersSimple(t *testing.T) {
	p := &Piece{
		// 440: ***
		// 523:  ***
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.N(3),
				Start:     frac.N(0),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.N(3),
				Start:     frac.N(1),
			},
		},
	}

	actual := p.getMarkers()
	expected := []frac.Frac{
		frac.N(0),
		frac.N(1),
		frac.N(3),
		frac.N(4),
	}

	if len(actual) != len(expected) {
		t.Fatalf("markers length don't match: \n(%d) %v\n(%d) %v", len(actual), actual, len(expected), expected)
	}
	for i := range actual {
		if actual[i] != expected[i] {
			t.Fatalf("markers #%d doesn't match\n%v\n%v", i, actual, expected)
		}
	}
}

func TestGetMarkersDuplicates(t *testing.T) {
	p := &Piece{
		Notes: []Note{
			// 349:   ***
			// 440: ***
			// 523:  ***
			Note{
				Frequency: 440,
				Duration:  frac.N(3),
				Start:     frac.N(0),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.N(3),
				Start:     frac.N(1),
			},
			Note{
				Frequency: 349.23,
				Duration:  frac.N(3),
				Start:     frac.N(2),
			},
		},
	}

	actual := p.getMarkers()
	expected := []frac.Frac{
		frac.N(0),
		frac.N(1),
		frac.N(2),
		frac.N(3),
		frac.N(4),
		frac.N(5),
	}

	if len(actual) != len(expected) {
		t.Fatalf("markers length don't match: \n(%d) %v\n(%d) %v", len(actual), actual, len(expected), expected)
	}
	for i := range actual {
		if actual[i] != expected[i] {
			t.Fatalf("markers #%d doesn't match\n%v\n%v", i, actual, expected)
		}
	}
}

func TestFromBPM(t *testing.T) {
	var bpmDuration = []struct {
		bpm      int
		duration time.Duration
	}{
		{120, 500 * time.Millisecond},
		{60, time.Second},
		{100, 600 * time.Millisecond},
	}

	for _, row := range bpmDuration {
		actual := FromBPM(row.bpm)
		if actual != row.duration {
			t.Errorf("bpm: %d, actual: %v, expected: %v", row.bpm, actual, row.duration)
		}
	}
}

func CompareBlocks(t *testing.T, actual, expected []block) {
	t.Helper()
	if len(actual) != len(expected) {
		t.Fatalf("intersection blocks length don't match: \n(%3d) %v\n(%3d) %v", len(actual), actual, len(expected), expected)
	}
	for i, block := range actual {
		if !block.equal(expected[i]) {
			t.Fatalf("intersection block #%d doesn't match: \n%v\n%v", i, block, expected[i])
		}
	}
}
