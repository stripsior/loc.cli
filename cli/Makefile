.PHONY: build build-all test clean install

# Build for current platform
build:
	go build -o bin/codeloc ./cmd/codeloc

# Build for all platforms
build-all:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/codeloc-linux-amd64 ./cmd/codeloc
	GOOS=linux GOARCH=arm64 go build -o dist/codeloc-linux-arm64 ./cmd/codeloc
	GOOS=darwin GOARCH=amd64 go build -o dist/codeloc-darwin-amd64 ./cmd/codeloc
	GOOS=darwin GOARCH=arm64 go build -o dist/codeloc-darwin-arm64 ./cmd/codeloc
	GOOS=windows GOARCH=amd64 go build -o dist/codeloc-windows-amd64.exe ./cmd/codeloc

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/ dist/

# Install locally
install:
	go install ./cmd/codeloc

# Run the tool
run:
	go run ./cmd/codeloc .

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build for current platform"
	@echo "  build-all  - Build for all platforms"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build artifacts"
	@echo "  install    - Install locally"
	@echo "  run        - Run the tool on current directory"
	@echo "  fmt        - Format code"
	@echo "  lint       - Lint code"
