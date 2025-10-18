# Wing Commander Makefile

.PHONY: build dev test clean install

# Default target
all: dev

# Development build (fast, no optimizations)
dev:
	@echo "🔨 Building development version..."
	@mkdir -p bin
	go build -o bin/wing_commander ./cmd/wing_commander
	@echo "✅ Development build complete: bin/wing_commander"

# Production build (optimized)
build:
	@echo "🚀 Building production version..."
	@mkdir -p dist
	go build -ldflags="-s -w" -o dist/wing_commander ./cmd/wing_commander
	@echo "✅ Production build complete: dist/wing_commander"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report: coverage.html"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/ dist/ build/ coverage.out coverage.html
	@echo "✅ Clean complete"

# Install to GOPATH/bin
install:
	@echo "📦 Installing to GOPATH/bin..."
	go install ./cmd/wing_commander
	@echo "✅ Installed: wing_commander"

# Run development build
run: dev
	@echo "🏃 Running development build..."
	./bin/wing_commander

# Run with specific command
run-cmd: dev
	@echo "🏃 Running: wing_commander $(CMD)"
	./bin/wing_commander $(CMD)

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "🔍 Linting code..."
	golangci-lint run

# Help
help:
	@echo "Wing Commander Build Commands:"
	@echo "  make dev          - Build development version (fast)"
	@echo "  make build        - Build production version (optimized)"
	@echo "  make test         - Run all tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make run          - Build and run development version"
	@echo "  make run-cmd CMD=version - Run with specific command"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make install      - Install to GOPATH/bin"
	@echo "  make fmt          - Format code"
	@echo "  make lint         - Lint code"
	@echo "  make help         - Show this help"
