package piece

import (
	"testing"
	"time"

	"github.com/math2001/piano/frac"
)

func TestIntersectionSimultaneous(t *testing.T) {
	p := &Piece{
		// 440: *
		// 523: *
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.F(1, 2),
				Start:     frac.N(0),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.F(1, 2),
				Start:     frac.N(0),
			},
		},
	}

	actual := p.intersectionBlocks()
	expected := []block{
		{start: frac.N(0), duration: frac.F(1, 2), frequencies: []float64{440, 523.25}},
	}
	CompareBlocks(t, actual, expected)
}

func TestIntersectionContainingOverlap(t *testing.T) {
	p := &Piece{
		//  * : 1/2 beat
		// 440: *****
		// 523:   **
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.F(5, 2),
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
		{start: frac.N(2), duration: frac.F(1, 2), frequencies: []float64{440}},
	}
	CompareBlocks(t, actual, expected)
}

func TestIntersectionSilence(t *testing.T) {
	p := &Piece{
		//  * : 1/3 beat
		// 440:   *
		// 523:     *****
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.F(1, 3),
				Start:     frac.F(2, 3),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.F(5, 3),
				Start:     frac.F(4, 3),
			},
		},
	}
	actual := p.intersectionBlocks()
	expected := []block{
		{start: frac.N(0), duration: frac.F(2, 3), frequencies: []float64{}},
		{start: frac.F(2, 3), duration: frac.F(1, 3), frequencies: []float64{440}},
		{start: frac.F(3, 3), duration: frac.F(1, 3), frequencies: []float64{}},
		{start: frac.F(4, 3), duration: frac.F(5, 3), frequencies: []float64{523.25}},
	}
	CompareBlocks(t, actual, expected)
}
func TestIntersectionOverlap(t *testing.T) {
	p := &Piece{
		//  * : 1/6 beat
		// 440: *******
		// 523:     *****
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.F(7, 6),
				Start:     frac.N(0),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.F(5, 6),
				Start:     frac.F(4, 6),
			},
		},
	}
	actual := p.intersectionBlocks()
	expected := []block{
		{start: frac.N(0), duration: frac.F(4, 6), frequencies: []float64{440}},
		{start: frac.F(4, 6), duration: frac.F(3, 6), frequencies: []float64{440, 523.25}},
		{start: frac.F(7, 6), duration: frac.F(2, 6), frequencies: []float64{523.25}},
	}
	CompareBlocks(t, actual, expected)
}

func TestGetMarkersSimple(t *testing.T) {
	p := &Piece{
		//  * : 1 / 2 beat
		// 440: ******
		// 523:  ***
		// 722:   *****
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.F(6, 2),
				Start:     frac.N(0),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.F(3, 2),
				Start:     frac.F(1, 2),
			},
			Note{
				Frequency: 722,
				Duration:  frac.F(5, 2),
				Start:     frac.F(2, 2),
			},
		},
	}

	actual := p.getMarkers()
	expected := []frac.Frac{
		frac.N(0),
		frac.F(1, 2),
		frac.F(2, 2),
		frac.F(4, 2),
		frac.F(6, 2),
		frac.F(7, 2),
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
			//  * : 1/3 beat
			// 200:    *
			// 349:   ***
			// 440:  **
			// 523:  ***
			// notice how we don't start at 0
			Note{
				Frequency: 200,
				Duration:  frac.F(1, 3),
				Start:     frac.F(3, 3),
			},
			Note{
				Frequency: 349.23,
				Duration:  frac.F(3, 3),
				Start:     frac.F(2, 3),
			},
			Note{
				Frequency: 440,
				Duration:  frac.F(2, 3),
				Start:     frac.F(1, 3),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.F(3, 3),
				Start:     frac.F(1, 3),
			},
		},
	}

	actual := p.getMarkers()
	expected := []frac.Frac{
		frac.N(0),
		frac.F(1, 3),
		frac.F(2, 3),
		frac.F(3, 3),
		frac.F(4, 3),
		frac.F(5, 3),
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
			t.Errorf("intersection block #%d doesn't match: \n%v\n%v", i, block, expected[i])
		}
	}
}
