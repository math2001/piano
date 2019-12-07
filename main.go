package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/math2001/piano/frac"
	"github.com/math2001/piano/labels"
	"github.com/math2001/piano/piece"
	"github.com/math2001/piano/wave"
)

func main() {
	lb := labels.NewLabels()
	p := &piece.Piece{
		Notes: []piece.Note{
			piece.Note{
				Frequency: lb.F("A4"),
				Duration:  frac.F(1, 1),
				Start:     frac.F(0, 1),
			},
			piece.Note{
				Frequency: lb.F("C5"),
				Duration:  frac.F(1, 1),
				Start:     frac.F(1, 1),
			},
			piece.Note{
				Frequency: lb.F("A4"),
				Duration:  frac.F(1, 1),
				Start:     frac.F(2, 1),
			},
			piece.Note{
				Frequency: lb.F("F4"),
				Duration:  frac.F(1, 1),
				Start:     frac.F(3, 1),
			},
			piece.Note{
				Frequency: lb.F("F4"),
				Duration:  frac.F(3, 1),
				Start:     frac.F(4, 1),
			},
			piece.Note{
				Frequency: lb.F("A4"),
				Duration:  frac.F(3, 1),
				Start:     frac.F(4, 1),
			},
			piece.Note{
				Frequency: lb.F("C5"),
				Duration:  frac.F(3, 1),
				Start:     frac.F(4, 1),
			},
		},
	}
	p = &piece.Piece{
		// 440: *
		// 512: *
		Notes: []piece.Note{
			piece.Note{
				Frequency: 440,
				Duration:  frac.F(3, 2),
				Start:     frac.N(0),
			},
			piece.Note{
				Frequency: 512,
				Duration:  frac.F(3, 4),
				Start:     frac.F(1, 2),
			},
		},
	}
	p.Render()
	return

	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/6))

	// put it in a gain, just so we don't play full throttle
	streamer := &effects.Gain{
		Streamer: p.GetStreamer(sr, piece.FromBPM(60)),
		Gain:     -0.1,
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}

func notePlayer() {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/6))

	lb := labels.NewLabels()

	sine := wave.NewSine(wave.N(sr, 440))
	loop := beep.Loop(-1, sine)
	ctrl := &beep.Ctrl{Streamer: loop, Paused: false}

	speaker.Play(ctrl)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("reading line: %s", err)
		}
		line = line[:len(line)-1]
		if line == "pause" {
			speaker.Lock()
			ctrl.Paused = !ctrl.Paused
			speaker.Unlock()
		} else {
			freq, err := lb.Frequency(line)
			fmt.Println("frequency:", freq, "Hz")
			if err != nil {
				log.Fatalf("converting label %q to frequency: %s", line, err)
			}
			speaker.Lock()
			ctrl.Streamer = beep.Loop(-1, wave.NewSine(sr.N(time.Second/time.Duration(freq))))
			speaker.Unlock()
		}
	}
}

func debugStreamerDesmos(streamer beep.Streamer, size int) {
	buf := make([][2]float64, size)
	n, ok := streamer.Stream(buf)
	if !ok {
		log.Fatalf("reading from stream not ok")
	}
	fmt.Printf("Read %d samples\n", n)
	for x, sample := range buf {
		fmt.Printf("(%d, %f)\n", x, sample[0])
	}
	fmt.Println()
	for x, sample := range buf {
		fmt.Printf("(%d, %f)\n", x, sample[0])
	}
	fmt.Println()
}
