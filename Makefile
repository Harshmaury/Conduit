.PHONY: build test verify clean

build:
	go build -o conduit ./cmd/conduit/

test:
	go test ./...

verify:
	go vet ./...
	go build ./...

clean:
	rm -f conduit
