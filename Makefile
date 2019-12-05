.PHONY: run

run: piano
	./piano

piano: *.go wave/*.go
	go build


