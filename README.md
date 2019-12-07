# Piano

Something which will eventually be able to play some music of its own.

## TODO

### Use `beep.Buffer` when generating waves

When I generate a wave, I just use a regular `[][2]float64` which takes up a
lot of room, as said in the wiki. I should probably use a `beep.Buffer` using a
`.wav` format

### Use a polynomial approximation instead of `math.Sin`

(Measure the performance gain before doing that)
Then, I could vary the precision and optimise it for the period of the wave
