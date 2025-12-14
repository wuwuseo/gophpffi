# Makefile for Go-PHP FFI CLI tool

.PHONY: build install clean test help

# Build the CLI tool
build:
	@echo "Building gophpffi CLI tool..."
	go build -o gophpffi.exe ./cmd/gophp
	@echo "✓ Built gophpffi.exe"

# Install the CLI tool to GOPATH/bin
install:
	@echo "Installing gophpffi to GOPATH/bin..."
	go install ./cmd/gophp
	@echo "✓ Installed gophpffi"
	@echo ""
	@echo "Make sure GOPATH/bin is in your PATH"
	@echo "You can now run: gophpffi --help"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@if exist gophpffi.exe del gophpffi.exe
	@if exist dist rd /s /q dist
	@echo "✓ Cleaned"

# Test the CLI tool
test: build
	@echo "Testing gophpffi CLI..."
	.\gophpffi.exe --help
	@echo ""
	@echo "✓ CLI tool working"

# Show help
help:
	@echo Available targets:
	@echo   build    - Build the gophpffi CLI tool
	@echo   install  - Install gophpffi to GOPATH/bin
	@echo   clean    - Remove build artifacts
	@echo   test     - Build and test the CLI tool
	@echo   help     - Show this help message
