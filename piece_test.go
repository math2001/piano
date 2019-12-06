package main

import "testing"

func TestPiecePlaySimpleOverlap(t *testing.T) {
	p := &Piece{
		name: "simple overlap",
		notes: []Note{
			Note{
				Frequency: 440,
				Duration:  Frac{2, 1},
				Start:     Frac{0, 1},
			},
			Note{
				Frequency: 523.25,
				Duration:  Frac{2, 1},
				Start:     Frac{1, 1},
			},
			Note{
				Frequency: 440,
				Duration:  Frac{1, 2},
				Start:     Frac{3, 1},
			},
		},
	}
	_ = p
	// p.Render()
}
