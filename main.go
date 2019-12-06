package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/math2001/piano/wave"
)

func main() {
	return
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/6))

	T := time.Second / 440

	labels := NewLabels()

	samples_per_period := sr.N(T)
	sine := wave.NewSine(samples_per_period)
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
			freq, err := labels.Frequency(line)
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
