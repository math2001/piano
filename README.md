# Piano

Something which will eventually be able to play some music of its own.

## TODO

### Use `beep.Buffer` when generating waves

When I generate a wave, I just use a regular `[][2]float64` which takes up a
lot of room, as said in the wiki. I should probably use a `beep.Buffer` using a
`.wav` format

### Use a polynomial approximation instead of `math.Sin`

(Measure the performance gain before doing that)
Then, I could vary the precision and optimise it for the length of our wave

### Play systems

#### Fill each track (streamer) and then mix

Make a track per streamer sound (eg. for sine waves that means one track per
frequency), with silence to space out each note as required by their `.Start`
and `.Duration`. Each track would be a streamer itself (of the length of the
song), and I would then mix every track together.

It's analogous to how the `Piece.Render` function works, plus a mix of
everything at the end.

Pros:

- should be easy to implement

Cons:

- depends on `beep.Mixer`. If it's quick to mix *a lot* of tracks, then it'll
  be fast
- not much can be optimized

##### Step through

Step through the `Piece` (a step is the smallest time variation within the
piece). For each step, check every `Note` and mix all of the ones
intersecting at *that* step. Then, I end up with a linear list of blocks to
play, which I can just give to `beep.Seq`.

Pros:

- a lot of possibilities for optimisation
    - Limit the number of notes we have to check:
        - statically split off the `Piece` in multiple parts
        - have boundaries moving with the current position
    - Maybe try to do some lookahead to limit the number of mixes required?


Cons:

- might be expensive to step through every smallest variation (but see first pro)
- adding a note at the end that is the smallest in the piece, or oddly placed
  will require a full re-render (that's not too bad, especially with first pro)
- might be expensive to do 20 different mix of the same streamers:

    A4: **********
    C5: **********

    Here there would be one mix per star.
