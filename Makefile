.PHONY: run test

test: *.go wave/*.go
	go test ./...

run: piano
	./piano

clean:
	go clean
	go clean -testcache

piano: $(shell fd --type f --extension go)
	go build


