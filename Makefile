# Just builds
.PHONY: build test

build:
	go test -c -i

test:
	go test -v ./...
