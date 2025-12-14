# Go-PHP FFI Service Library Generator

A CLI tool to generate PHP FFI bindings for Go shared libraries.

> **Note:** This project now uses the `gophpffi` CLI tool. See [CLI-README.md](CLI-README.md) for complete documentation.

## Quick Start

### 1. Install the CLI Tool

```bash
go install ./cmd/gophp
# or
make install
```

### 2. Initialize a New Service

```bash
gophpffi init [ServiceName]
```

**Example:**
```bash
gophpffi init Product
```

This will create:
- `Product.go` - Go source file with template functions
- `.gophp.yaml` - Configuration file

### 3. Edit Your Go File

Open the generated `.go` file and add your exported functions:

```go
//export YourFunction
func YourFunction(param Type) ReturnType {
    // Your implementation
    return result
}
```

**Important Notes:**
- Every exported function must have the `//export FunctionName` comment
- The comment must be directly above the function (no blank lines)
- The function must be in `package main`
- Keep the `main()` function (required for shared library)

### 4. Build Everything

```bash
gophpffi make
```

This will:
1. Generate PHP service class in `dist/`
2. Build shared library in `dist/lib/`

## Directory Structure

After building, you'll have:

```
dist/
├── [ServiceName]Service.php          (PHP Service Class)
└── lib/
    ├── [ServiceName]-windows-amd64.dll  (Shared Library)
    └── [ServiceName]-windows-amd64.h    (C Header)
```

## CLI Commands

### Generate PHP Bindings Only
```bash
gophpffi generate
# or
go generate ./...
```

### Build Shared Library Only
```bash
gophpffi build
```

### Complete Build
```bash
gophpffi make
```

For detailed CLI documentation, see [CLI-README.md](CLI-README.md).

## Example

### Step 1: Initialize
```bash
gophpffi init User
```

### Step 2: Edit User.go
```go
package main

/*
#include <stdlib.h>
*/
import "C"

//go:generate gophpffi generate User.go

// GetUserName returns a user name
//
//export GetUserName
func GetUserName(id int) string {
    return "User_" + strconv.Itoa(id)
}

// main is required for building shared library but never called
func main() {}
```

### Step 3: Build
```bash
gophpffi make
```

### Step 4: Use in PHP
```php
<?php
require_once 'dist/UserService.php';

use app\user\service\UserService;

$service = new UserService();
$name = $service->GetUserName(123);
echo $name; // Output: User_123
```

## Type Mapping

| Go Type | PHP Type |
|---------|----------|
| int, int8, int16, int32, int64 | int |
| uint, uint8, uint16, uint32, uint64 | int |
| float32, float64 | float |
| string | string |
| bool | bool |
| []T (slice) | array |
| map[K]V | array |

## Troubleshooting

### Build Fails
- Ensure you have Go installed: `go version`
- Check that all exported functions have the `//export` comment
- Make sure there are no unused imports

### Generated PHP File Not Found
- Run `gophpffi generate` or `go generate ./...` first
- Check that your `.go` file has the `//go:generate` directive

### DLL Not Found
- Run `gophpffi build` to compile the shared library
- Make sure you're on Windows (for .dll files)

### CLI Command Not Found
- If installed with `go install`, ensure `$GOPATH/bin` is in your PATH
- Otherwise, use `./gophpffi.exe` from the project directory

## Advanced Usage

### Custom Configuration

Edit `.gophp.yaml` to customize build settings:

```yaml
service: MyService
source: MyService.go
output:
  dir: dist
  lib_dir: dist/lib
```

### Multiple Services

You can create multiple services in the same project:

```bash
gophpffi init UserService
gophpffi init ProductService
gophpffi init OrderService
```

Each service will have its own:
- Go source file (`[ServiceName].go`)
- Configuration file (`.gophp.yaml`)
- PHP service class in `dist/`
- Shared library in `dist/lib/`

**Note:** For multiple services, you'll need to manage separate `.gophp.yaml` files or specify source files explicitly:

```bash
gophpffi generate UserService.go
gophpffi build UserService.go
```

### Cross-Platform Builds

While the CLI builds for the current platform by default, you can cross-compile for other platforms:

**Windows (from Linux/macOS):**
```bash
GOOS=windows GOARCH=amd64 go build -buildmode=c-shared -o dist/lib/Service-windows-amd64.dll Service.go
```

**Linux (from Windows/macOS):**
```bash
GOOS=linux GOARCH=amd64 go build -buildmode=c-shared -o dist/lib/Service-linux-amd64.so Service.go
```

**macOS (from Windows/Linux):**
```bash
GOOS=darwin GOARCH=amd64 go build -buildmode=c-shared -o dist/lib/Service-darwin-amd64.dylib Service.go
```

### Using with Composer

You can integrate the generated PHP classes with Composer:

1. Add to your `composer.json`:
```json
{
    "autoload": {
        "files": ["dist/UserService.php"]
    }
}
```

2. Run composer update:
```bash
composer dump-autoload
```

3. Use in your PHP code:
```php
<?php
require 'vendor/autoload.php';

use app\user\service\UserService;

$service = new UserService();
```

## Best Practices

### 1. Function Design
- Keep function signatures simple
- Use basic types when possible (int, string, bool)
- Avoid complex nested structures
- Return error codes instead of Go's error type

### 2. Error Handling

**In Go:**
```go
//export ProcessData
func ProcessData(data string) int {
    if data == "" {
        return -1  // Error code
    }
    // Process data...
    return 0  // Success
}
```

**In PHP:**
```php
$result = $service->ProcessData($data);
if ($result < 0) {
    throw new Exception("Processing failed with code: " . $result);
}
```

### 3. Memory Management

- Go manages memory automatically
- FFI memory is managed by PHP
- Large data should be processed in chunks
- Clean up resources in PHP when done

### 4. Performance Optimization

- Minimize cross-language calls
- Batch operations when possible
- Use persistent PHP processes (PHP-FPM)
- Cache FFI library instances

### 5. Testing

Create test files for your services:

```php
<?php
// test_user_service.php
require_once 'dist/UserService.php';

use app\user\service\UserService;

$service = new UserService();

// Test basic functionality
assert($service->GetUserName(1) === "User_1");
assert($service->GetUserName(100) === "User_100");

echo "All tests passed!\n";
```

## Project Files

- `gophpffi` - CLI tool (command-line interface)
- `cmd/gophp/` - CLI tool source code
- `generator/main.go` - Code generator source
- `.gophp.yaml` - Project configuration file
- `AGENTS.MD` - Development guide (Chinese)
- `AGENTS_EN.MD` - Development guide (English)
- `README.md` - User guide (Chinese)
- `README_EN.MD` - User guide (English)

## Requirements

- Go 1.16 or later
- Windows (for building .dll files), Linux (for .so), or macOS (for .dylib)
- PHP 7.4 or later with FFI extension enabled

### Enabling PHP FFI

**In php.ini:**
```ini
extension=ffi
ffi.enable=true
```

**Check if FFI is enabled:**
```bash
php -m | grep ffi
```

**Or in PHP code:**
```php
<?php
if (!extension_loaded('ffi')) {
    die('FFI extension is not loaded');
}
echo "FFI is enabled!\n";
```

## Platform-Specific Notes

### Windows
- Requires MinGW-w64 or Visual Studio with C++ tools
- DLL files are created in `dist/lib/`
- May need to add `dist/lib/` to system PATH
- Antivirus software may flag DLL files

### Linux
- Requires GCC compiler
- SO files are created in `dist/lib/`
- May need to set `LD_LIBRARY_PATH`
- Check library dependencies with `ldd`

### macOS
- Requires Xcode Command Line Tools
- DYLIB files are created in `dist/lib/`
- May require code signing for distribution
- Check with `otool -L`

## Security Considerations

### Input Validation
Always validate inputs before passing to Go functions:

```php
<?php
function safeGetUserName($id) {
    if (!is_int($id) || $id < 0) {
        throw new InvalidArgumentException("Invalid user ID");
    }
    return $service->GetUserName($id);
}
```

### Error Handling
Never expose internal paths or sensitive information in error messages.

### Production Deployment
- Strip debug symbols from binaries
- Use read-only permissions for library files
- Monitor for unusual FFI calls
- Keep PHP and Go versions updated

## Migration Guide

### From Batch Scripts to CLI

If you're migrating from the old batch script system:

| Old Command | New Command |
|-------------|-------------|
| `init.bat Product` | `gophpffi init Product` |
| `generate.bat` | `gophpffi generate` |
| `build.bat` | `gophpffi build` |
| `make.bat` | `gophpffi make` |

### From gophp to gophpffi

If you're upgrading from the old `gophp` command:

1. Rebuild the CLI tool:
```bash
go build -o gophpffi.exe ./cmd/gophp
```

2. Update your scripts:
```bash
# Old
gophp init Product

# New
gophpffi init Product
```

3. Update `//go:generate` directives in your Go files:
```go
// Old
//go:generate gophp generate Product.go

// New
//go:generate gophpffi generate Product.go
```

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

### Development Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/serviceLib.git
cd serviceLib

# Build the CLI tool
go build -o gophpffi.exe ./cmd/gophp

# Run tests
make test
```

## License

[Your License Here]

## Support

- Documentation: [AGENTS.MD](AGENTS.MD) (Chinese), [AGENTS_EN.MD](AGENTS_EN.MD) (English)
- Issues: [GitHub Issues](https://github.com/yourusername/serviceLib/issues)
- Discussions: [GitHub Discussions](https://github.com/yourusername/serviceLib/discussions)

## Changelog

### Version 1.0.0
- ✅ Initial release with `gophpffi` CLI tool
- ✅ Full Chinese localization
- ✅ Configuration file support
- ✅ Cross-platform build support
- ✅ Automatic type mapping
- ✅ Template-based initialization

---

**Made with ❤️ using Go and PHP FFI**
