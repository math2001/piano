.PHONY: run test

test: *.go wave/*.go
	go test . ./piece ./frac ./labels ./wave

run: piano
	./piano

piano: *.go wave/*.go
	go build


