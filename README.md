# Piano

Something which will eventually be able to play some music of its own.

## Convention

In tests, always print the actual result and then the expected one. I just
couldn't chose, so I went with alphabetical order.

## Piece

When using a `Piece`, `.Render` it and paste it as a comment above the `Note`
field. It makes it much easier to understand what's going on. For example:

```go
p := &Piece{
    // 440: ***        <- the .Render
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
```

**Make sure they stay updated**

(is there an easy thing to build which would automatically update those?)

## TODO

### Use `beep.Buffer` when generating waves

When I generate a wave, I just use a regular `[][2]float64` which takes up a
lot of room, as said in the wiki. I should probably use a `beep.Buffer` using a
`.wav` format

### Use a polynomial approximation instead of `math.Sin`

(Measure the performance gain before doing that)
Then, I could vary the precision and optimise it for the period of the wave


