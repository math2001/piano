.PHONY: run test

run: piano
	./piano

test: *.go wave/*.go
	go test ./...

clean:
	go clean
	go clean -testcache

piano: $(shell fd --type f --extension go)
	go build


