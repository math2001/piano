package piece

import (
	"testing"
	"time"

	"github.com/math2001/piano/frac"
)

func TestPiecePlaySimultaneous(t *testing.T) {
	p := &Piece{
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.F(2, 1),
				Start:     frac.F(0, 1),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.F(2, 1),
				Start:     frac.F(0, 1),
			},
		},
	}

	blocks := p.intersectionBlocks()
	expected := []block{
		// two streamers, don't really know how it's gonna be implemented
		{start: frac.F(0, 1), duration: frac.F(2, 1), frequencies: []float64{0, 0}},
	}
	if len(blocks) != len(expected) {
		t.Fatalf("intersection blocks length don't match: \n%v\n%v", blocks, expected)
	}
	for i, block := range blocks {
		if !block.equal(expected[i]) {
			t.Fatalf("intersection block %d doesn't match: \n%v\n%v", i, block, expected[i])
		}
	}
}

func TestPiecePlayContainingOverlap(t *testing.T) {
	p := &Piece{
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.F(3, 1),
				Start:     frac.F(0, 1),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.F(1, 1),
				Start:     frac.F(1, 1),
			},
		},
	}
	blocks := p.intersectionBlocks()
	expected := []block{
		{start: frac.F(0, 1), duration: frac.F(1, 1), frequencies: []float64{0}},
		{start: frac.F(1, 1), duration: frac.F(1, 1), frequencies: []float64{0, 0}},
		{start: frac.F(2, 1), duration: frac.F(1, 1), frequencies: []float64{0}},
	}
	if len(blocks) != len(expected) {
		t.Fatalf("intersection blocks length don't match: \n%v\n%v", blocks, expected)
	}
	for i, block := range blocks {
		if !block.equal(expected[i]) {
			t.Fatalf("intersection block %d doesn't match: \n%v\n%v", i, block, expected[i])
		}
	}
}

func TestPiecePlaySilence(t *testing.T) {
	p := &Piece{
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.F(1, 1),
				Start:     frac.F(2, 1),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.F(3, 1),
				Start:     frac.F(4, 1),
			},
		},
	}
	blocks := p.intersectionBlocks()
	expected := []block{
		{start: frac.F(0, 1), duration: frac.F(2, 1), frequencies: []float64{}},
		{start: frac.F(2, 1), duration: frac.F(1, 1), frequencies: []float64{0}},
		{start: frac.F(3, 1), duration: frac.F(1, 1), frequencies: []float64{}},
		{start: frac.F(4, 1), duration: frac.F(3, 1), frequencies: []float64{0}},
	}
	if len(blocks) != len(expected) {
		t.Fatalf("intersection blocks length don't match: \n%v\n%v", blocks, expected)
	}
	for i, block := range blocks {
		if !block.equal(expected[i]) {
			t.Fatalf("intersection block %d doesn't match: \n%v\n%v", i, block, expected[i])
		}
	}
}
func TestPiecePlayIntersectingOverlap(t *testing.T) {
	p := &Piece{
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.F(3, 1),
				Start:     frac.F(0, 1),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.F(3, 1),
				Start:     frac.F(1, 1),
			},
		},
	}
	blocks := p.intersectionBlocks()
	expected := []block{
		{start: frac.F(0, 1), duration: frac.F(1, 1), frequencies: []float64{0}},
		{start: frac.F(1, 1), duration: frac.F(2, 1), frequencies: []float64{0, 0}},
		{start: frac.F(3, 1), duration: frac.F(1, 1), frequencies: []float64{0}},
	}
	if len(blocks) != len(expected) {
		t.Fatalf("intersection blocks length don't match: \n%v\n%v", blocks, expected)
	}
	for i, block := range blocks {
		if !block.equal(expected[i]) {
			t.Fatalf("intersection block %d doesn't match: \n%v\n%v", i, block, expected[i])
		}
	}
}

func TestGetMarkersSimple(t *testing.T) {
	p := &Piece{
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.F(3, 1),
				Start:     frac.F(0, 1),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.F(3, 1),
				Start:     frac.F(1, 1),
			},
		},
	}

	actual := p.getMarkers()
	expected := []frac.Frac{
		frac.F(0, 1),
		frac.F(1, 1),
		frac.F(3, 1),
		frac.F(4, 1),
	}

	if len(actual) != len(expected) {
		t.Fatalf("markers length don't match: \n(%d) %v\n(%d) %v", len(actual), actual, len(expected), expected)
	}
}

func TestGetMarkersDuplicates(t *testing.T) {
	p := &Piece{
		Notes: []Note{
			Note{
				Frequency: 440,
				Duration:  frac.F(3, 1),
				Start:     frac.F(0, 1),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.F(3, 1),
				Start:     frac.F(1, 1),
			},
			Note{
				Frequency: 523.25,
				Duration:  frac.F(3, 1),
				Start:     frac.F(2, 1),
			},
		},
	}

	actual := p.getMarkers()
	expected := []frac.Frac{
		frac.F(0, 1),
		frac.F(1, 1),
		frac.F(2, 1),
		frac.F(3, 1),
		frac.F(4, 1),
		frac.F(5, 1),
	}

	if len(actual) != len(expected) {
		t.Fatalf("markers length don't match: \n(%d) %v\n(%d) %v", len(actual), actual, len(expected), expected)
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
