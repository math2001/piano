.PHONY: run test

test: *.go wave/*.go
	go test ./...

run: piano
	./piano

piano: $(shell fd --type f --extension go)
	go build


