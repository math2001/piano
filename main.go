package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/math2001/piano/wave"
)

func main() {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second))

	T := time.Second / 440

	samples_per_period := sr.N(T)
	fmt.Println(samples_per_period)
	streamer := wave.Sine(samples_per_period)

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}

func debugStreamer(streamer beep.Streamer, size int) {
	buf := make([][2]float64, size)
	n, ok := streamer.Stream(buf)
	if !ok {
		log.Fatalf("reading from stream not ok")
	}
	fmt.Printf("Read %d samples\n", n)
	for _, sample := range buf {
		fmt.Print(sample[0], " ")
	}
	fmt.Println()
	for _, sample := range buf {
		fmt.Print(sample[1], " ")
	}
	fmt.Println()
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

func Noise() beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			samples[i][0] = rand.Float64()*2 - 1
			samples[i][1] = rand.Float64()*2 - 1
		}
		return len(samples), true
	})
}
