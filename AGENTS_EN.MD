# AGENTS_EN.MD - Go-PHP FFI Project Development Guide

## Project Overview

Go-PHP FFI Service Library Generator is an automation tool for generating PHP FFI bindings for Go shared libraries. The project uses the Cobra framework to build CLI tools, providing a complete code generation and build process.

### Core Features
- Automatic PHP FFI binding code generation
- Cross-platform shared library building (.dll/.so/.dylib)
- Type mapping and safe memory management
- Template-based service initialization

## Technology Stack

### Backend Technologies
- **Language**: Go 1.16+
- **CLI Framework**: Cobra
- **Configuration Management**: YAML (gopkg.in/yaml.v3)
- **Build Mode**: c-shared
- **Code Generation**: Go AST parsing

### PHP Integration
- **PHP Version**: 7.4+
- **Extension Requirements**: FFI extension
- **Type System**: C type-based FFI bindings

## Project Structure

```
serviceLib/
├── cmd/
│   └── gophp/              # CLI tool source code
│       ├── main.go         # Main entry and root command
│       ├── init.go         # init command implementation
│       ├── generate.go     # generate command implementation
│       ├── build.go        # build command implementation
│       ├── make.go         # make command implementation
│       └── config.go       # Configuration file loading
├── generator/
│   └── main.go             # Code generator core logic
├── dist/                   # Build output directory
│   ├── *Service.php        # Generated PHP service class
│   └── lib/                # Shared library files
├── .gophp.yaml             # Project configuration file
├── Makefile                # Build automation
├── README.md               # Project documentation
└── go.mod                  # Go module definition
```

## CLI Tool Architecture

### Command Structure

```
gophpffi
├── init [ServiceName]      # Initialize new service
├── generate [source.go]    # Generate PHP bindings
├── build [source.go]       # Build shared library
└── make                    # Complete build process
```

### Command Implementation Details

#### 1. `gophpffi init`
- **Function**: Create template files for new service
- **Output**:
  - `[ServiceName].go` - Go source file with example functions
  - `.gophp.yaml` - Project configuration file
- **File**: `cmd/gophp/init.go`

#### 2. `gophpffi generate`
- **Function**: Parse Go source file and generate PHP FFI bindings
- **Process**:
  1. Read configuration file or command-line arguments
  2. Call `generator/main.go` to parse Go AST
  3. Generate PHP service class to `dist/` directory
- **File**: `cmd/gophp/generate.go`

#### 3. `gophpffi build`
- **Function**: Compile Go source file to shared library
- **Build Parameters**: `-buildmode=c-shared`
- **Output**: 
  - `[ServiceName]-[os]-[arch].dll/.so/.dylib`
  - `[ServiceName]-[os]-[arch].h`
- **File**: `cmd/gophp/build.go`

#### 4. `gophpffi make`
- **Function**: Execute complete build process
- **Steps**:
  1. Generate PHP bindings (generate)
  2. Build shared library (build)
- **File**: `cmd/gophp/make.go`

## Configuration File

### `.gophp.yaml` Structure

```yaml
service: ServiceName        # Service name
source: ServiceName.go      # Go source file path
output:
  dir: dist                 # Output directory
  lib_dir: dist/lib         # Library file directory
```

### Configuration Loading Logic
- Location: `cmd/gophp/config.go`
- Functions: 
  - Read and parse `.gophp.yaml`
  - Provide default values
  - Validate configuration

## Code Generator

### Core File: `generator/main.go`

#### Functional Modules
1. **Go AST Parsing**
   - Parse exported functions (with `//export` comments)
   - Extract function signatures, parameters, return values
   - Identify type information

2. **Type Mapping**
   ```
   Go Type          →  PHP Type
   ----------------------------------------
   int/int64        →  int
   string           →  string
   bool             →  bool
   float32/float64  →  float
   []T (slice)      →  array
   map[K]V          →  array
   ```

3. **Code Generation**
   - Generate PHP class definitions
   - Create FFI binding methods
   - Add type hints and documentation comments

## Development Workflow

### Standard Development Process

```bash
# 1. Initialize new service
gophpffi init MyService

# 2. Edit MyService.go, add exported functions
# //export MyFunction
# func MyFunction(param Type) ReturnType { ... }

# 3. Generate and build
gophpffi make

# 4. Use in PHP
# require_once 'dist/MyServiceService.php';
```

### Debugging Tips

1. **View Generated Code**
   ```bash
   gophpffi generate
   cat dist/*Service.php
   ```

2. **Verify Build Output**
   ```bash
   gophpffi build
   ls -la dist/lib/
   ```

3. **Test Individual Commands**
   ```bash
   gophpffi generate MyService.go
   gophpffi build MyService.go
   ```

## Code Standards

### Go Code Standards

1. **Exported Function Requirements**
   ```go
   // ✅ Correct
   //export MyFunction
   func MyFunction(a int) int {
       return a * 2
   }
   
   // ❌ Wrong - Missing export comment
   func MyFunction(a int) int {
       return a * 2
   }
   
   // ❌ Wrong - Blank line between comment and function
   //export MyFunction
   
   func MyFunction(a int) int {
       return a * 2
   }
   ```

2. **Package Requirements**
   - Must use `package main`
   - Must include `main()` function (even if empty)
   - Must import `"C"`

3. **Comment Standards**
   - Use Chinese comments for code
   - Add functional descriptions for exported functions
   - Document parameter and return value meanings

### CLI Code Standards

1. **Command Implementation**
   - One file per command
   - Use `cobra.Command` structure
   - Provide Chinese help messages

2. **Error Handling**
   - Use `fmt.Errorf` to wrap errors
   - Provide Chinese error messages
   - Return meaningful error information

3. **Output Format**
   - Use consistent separator line style
   - Step numbering: `[1/3]`, `[2/3]`
   - Success marker: `✓`

## Build and Release

### Local Build

```bash
# Build CLI tool
go build -o gophpffi.exe ./cmd/gophp

# Or use Makefile
make build
```

### System Installation

```bash
# Install to GOPATH/bin
go install ./cmd/gophp

# Or use Makefile
make install
```

### Test Build

```bash
# Run tests
make test

# Manual testing
./gophpffi.exe --help
./gophpffi.exe init TestService
./gophpffi.exe make
```

## Troubleshooting

### Build Issues

1. **CGO Not Enabled**
   - Symptom: `CGO_ENABLED` related errors
   - Solution: Set environment variable `CGO_ENABLED=1`

2. **Compiler Not Found**
   - Symptom: `gcc` or `cl.exe` not found
   - Solution: Install MinGW-w64 (Windows) or GCC (Linux)

3. **Exported Functions Not Generated**
   - Symptom: Methods missing in PHP class
   - Check: Ensure `//export` comment exists with no blank lines

### Runtime Issues

1. **DLL Loading Failed**
   - Check: Platform architecture match (64-bit/32-bit)
   - Check: Required runtime libraries installed

2. **Configuration File Read Failed**
   - Check: `.gophp.yaml` format correct
   - Check: YAML indentation uses spaces not tabs

## Extension Development

### Adding New Commands

1. Create new file under `cmd/gophp/`
2. Define `cobra.Command` structure
3. Register to `rootCmd` in `init()`
4. Implement command logic

Example:
```go
var newCmd = &cobra.Command{
    Use:   "new",
    Short: "New command description",
    Long:  `Detailed explanation`,
    RunE:  runNew,
}

func init() {
    rootCmd.AddCommand(newCmd)
}

func runNew(cmd *cobra.Command, args []string) error {
    // Implementation logic
    return nil
}
```

### Modifying Code Generator

1. Edit `generator/main.go`
2. Modify AST parsing logic
3. Update PHP template generation
4. Test generated code

## Version History

### Current Version Features
- ✅ CLI tool fully localized in Chinese
- ✅ Command name migrated from `gophp` to `gophpffi`
- ✅ Configuration file-driven support
- ✅ Automatic type mapping
- ✅ Cross-platform build support
- ✅ Bilingual documentation system (Chinese/English)
- ✅ MIT open source license

### Migrating from Old Versions

| Old Command | New Command |
|-------------|-------------|
| `init.bat Name` | `gophpffi init Name` |
| `generate.bat` | `gophpffi generate` |
| `build.bat` | `gophpffi build` |
| `make.bat` | `gophpffi make` |
| `gophp *` | `gophpffi *` |

## Best Practices

### 1. Project Organization
- One Go file per service
- Use meaningful service names
- Keep function signatures simple

### 2. Function Design
- Avoid complex nested types
- Prefer basic types
- Provide clear documentation for slices and maps

### 3. Error Handling
- Go functions return error codes instead of error type
- Check return values on PHP side
- Log detailed error information

### 4. Performance Optimization
- Reduce cross-language call frequency
- Batch process data
- Use persistent PHP processes (PHP-FPM)

## Contributing Guidelines

### Code Commit Standards
- Use Chinese for commit messages
- Follow existing code style
- Add necessary comments
- Update related documentation

### Testing Requirements
- Test all CLI commands
- Verify generated code is runnable
- Check cross-platform compatibility

## Technical Support

### Documentation Resources
- `README.md` - User guide (Chinese)
- `README_EN.MD` - User guide (English)
- `AGENTS.MD` - Development guide (Chinese)
- `AGENTS_EN.MD` - Development guide (English)
- `LICENSE` - MIT License

### Issue Reporting
- GitHub Issues: https://github.com/wuwuseo/gophpffi/issues
- Check existing documentation
- Review common issues
- Provide detailed error information and environment description

## License

MIT License

Copyright (c) 2025 wuwuseo

See [LICENSE](LICENSE) file for details.

## API Reference

### CLI Commands API

#### `gophpffi init [ServiceName]`
Initialize a new service with template code.

**Arguments:**
- `ServiceName` (required): Name of the service to create

**Output:**
- `[ServiceName].go`: Go source file
- `.gophp.yaml`: Configuration file

**Example:**
```bash
gophpffi init Product
```

#### `gophpffi generate [source.go]`
Generate PHP FFI bindings from Go source file.

**Arguments:**
- `source.go` (optional): Go source file path. If omitted, uses config file.

**Output:**
- `dist/[ServiceName]Service.php`: PHP service class

**Example:**
```bash
gophpffi generate Product.go
# Or with config file
gophpffi generate
```

#### `gophpffi build [source.go]`
Build Go shared library.

**Arguments:**
- `source.go` (optional): Go source file path. If omitted, uses config file.

**Output:**
- `dist/lib/[ServiceName]-[os]-[arch].dll/.so/.dylib`: Shared library
- `dist/lib/[ServiceName]-[os]-[arch].h`: C header file

**Example:**
```bash
gophpffi build Product.go
# Or with config file
gophpffi build
```

#### `gophpffi make`
Execute complete build process (generate + build).

**Arguments:** None

**Output:**
- All outputs from generate and build commands

**Example:**
```bash
gophpffi make
```

## Advanced Topics

### Custom Type Handling

When dealing with complex types, consider these approaches:

1. **Flatten Structures**: Convert complex Go structs to basic types
2. **JSON Serialization**: Pass complex data as JSON strings
3. **Opaque Pointers**: Use integer handles for Go objects

### Memory Management

- Go manages memory automatically via garbage collector
- PHP FFI requires manual memory management for C types
- Use `defer` in Go for cleanup operations
- Free allocated memory in PHP when done

### Cross-Platform Considerations

**Windows:**
- Requires MinGW-w64 or Visual C++
- DLL files in `dist/lib/`
- Path separators: backslash `\`

**Linux:**
- Requires GCC
- SO files in `dist/lib/`
- Path separators: forward slash `/`

**macOS:**
- Requires Xcode Command Line Tools
- DYLIB files in `dist/lib/`
- May require code signing for distribution

## Security Considerations

### Input Validation
- Always validate inputs on both Go and PHP sides
- Sanitize string inputs to prevent injection
- Check array bounds before access

### Error Information
- Avoid exposing internal paths in error messages
- Log sensitive operations securely
- Use structured logging for debugging

### Distribution
- Strip debug symbols from production builds
- Sign binaries for distribution
- Document security requirements for deployment

## Performance Tuning

### Optimization Strategies

1. **Reduce FFI Overhead**
   - Batch operations when possible
   - Cache frequently used values
   - Minimize string conversions

2. **Go Side Optimizations**
   - Use sync.Pool for object reuse
   - Preallocate slices when size is known
   - Profile with pprof for bottlenecks

3. **PHP Side Optimizations**
   - Use persistent connections (PHP-FPM)
   - Cache FFI library instances
   - Avoid repeated library loads

### Benchmarking

```bash
# Go benchmarks
go test -bench=. -benchmem ./...

# PHP benchmarks
php -d opcache.enable_cli=1 benchmark.php
```

## Continuous Integration

### GitHub Actions Example

```yaml
name: Build and Test

on: [push, pull_request]

jobs:
  build:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.16'
      - name: Build CLI
        run: go build -o gophpffi.exe ./cmd/gophp
      - name: Test
        run: |
          ./gophpffi.exe init TestService
          ./gophpffi.exe make
```

## Glossary

- **FFI**: Foreign Function Interface - allows calling C code from PHP
- **CGO**: Go's mechanism for calling C code
- **Shared Library**: Dynamic library (.dll, .so, .dylib)
- **Export Comment**: `//export FunctionName` - marks function for FFI export
- **AST**: Abstract Syntax Tree - code structure representation
- **Cobra**: Popular Go CLI framework
- **YAML**: Configuration file format

## FAQ

**Q: Can I use Go's error type in exported functions?**
A: No, FFI doesn't support Go's error interface. Return error codes (int) instead.

**Q: How do I pass arrays between Go and PHP?**
A: Use slices in Go, which map to arrays in PHP via FFI bindings.

**Q: Can I export methods on structs?**
A: No, only package-level functions can be exported. Use wrapper functions.

**Q: How do I debug FFI calls?**
A: Enable FFI logging in PHP, use Go's log package, and check generated header files.

**Q: Is it thread-safe?**
A: Go functions are thread-safe if properly designed. Ensure proper synchronization.

---

**Last Updated**: December 2025  
**Project Version**: 1.0  
**Minimum Go Version**: 1.16  
**Minimum PHP Version**: 7.4
