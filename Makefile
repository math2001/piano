.PHONY: run test

test: *.go wave/*.go
	go test ./...

run: piano
	./piano

piano: *.go wave/*.go
	go build


