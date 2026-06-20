.PHONY: run build test clean fmt tidy

# Binary name
BINARY_NAME=server

# Run the server
run:
	go run cmd/server/main.go

# Build the binary
build:
	mkdir -p bin
	go build -o bin/$(BINARY_NAME) cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Format the code
fmt:
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy

# Clean build artifacts and local sqlite db
clean:
	rm -rf bin/
	rm -f env-vault.db
