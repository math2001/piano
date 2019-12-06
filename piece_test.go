package main

import (
	"testing"
)

func TestPiecePlaySimultaneous(t *testing.T) {
	p := &Piece{
		notes: []Note{
			Note{
				Frequency: 440,
				Duration:  Frac{2, 1},
				Start:     Frac{0, 1},
			},
			Note{
				Frequency: 523.25,
				Duration:  Frac{2, 1},
				Start:     Frac{0, 1},
			},
		},
	}
	blocks := p.intersectionBlocks()
	expected := []block{
		// two streamers, don't really know how it's gonna be implemented
		{start: Frac{0, 1}, duration: Frac{2, 1}, frequencies: []float64{0, 0}},
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
		notes: []Note{
			Note{
				Frequency: 440,
				Duration:  Frac{3, 1},
				Start:     Frac{0, 1},
			},
			Note{
				Frequency: 523.25,
				Duration:  Frac{1, 1},
				Start:     Frac{1, 1},
			},
		},
	}
	blocks := p.intersectionBlocks()
	expected := []block{
		{start: Frac{0, 1}, duration: Frac{1, 1}, frequencies: []float64{0}},
		{start: Frac{1, 1}, duration: Frac{1, 1}, frequencies: []float64{0, 0}},
		{start: Frac{2, 1}, duration: Frac{1, 1}, frequencies: []float64{0}},
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
		notes: []Note{
			Note{
				Frequency: 440,
				Duration:  Frac{3, 1},
				Start:     Frac{0, 1},
			},
			Note{
				Frequency: 523.25,
				Duration:  Frac{3, 1},
				Start:     Frac{1, 1},
			},
		},
	}
	blocks := p.intersectionBlocks()
	expected := []block{
		{start: Frac{0, 1}, duration: Frac{1, 1}, frequencies: []float64{0}},
		{start: Frac{1, 1}, duration: Frac{2, 1}, frequencies: []float64{0, 0}},
		{start: Frac{3, 1}, duration: Frac{1, 1}, frequencies: []float64{0}},
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
