# Makefile for SMQ

# Default target
default: build

# Build the binary
build:
	go build -o bin/smq main.go

# Run the binary
run: build
	./bin/smq

# Clean up
clean:
	rm -f bin/smq

# run tests
test:
	go test -v -race ./...

lint:
	golangci-lint run --fix ./...
