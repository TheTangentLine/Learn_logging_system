.PHONY: build vet test check

build:
	go build ./...

vet:
	go vet ./...

test:
	go test -v -race ./...

check: vet build test
