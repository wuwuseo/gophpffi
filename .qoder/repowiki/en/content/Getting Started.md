# Getting Started

<cite>
**Referenced Files in This Document**
- [QUICKSTART.md](file://QUICKSTART.md)
- [README.md](file://README.md)
- [go.mod](file://go.mod)
- [generator/main.go](file://generator/main.go)
- [build.go](file://build.go)
- [Makefile](file://Makefile)
- [build.ps1](file://build.ps1)
- [example.php](file://example.php)
- [mygo.go](file://mygo.go)
</cite>

## Table of Contents
1. [Introduction](#introduction)
2. [Prerequisites](#prerequisites)
3. [Quick Start](#quick-start)
4. [Step-by-Step Setup](#step-by-step-setup)
5. [Using the Build System](#using-the-build-system)
6. [Calling Go Functions from PHP](#calling-go-functions-from-php)
7. [Common Pitfalls and Fixes](#common-pitfalls-and-fixes)
8. [Troubleshooting Guide](#troubleshooting-guide)
9. [Conclusion](#conclusion)

## Introduction
This Getting Started guide helps you set up a development environment and quickly build a cross-platform Go shared library consumable by PHP via FFI. It covers prerequisites, environment setup, building libraries, generating PHP bindings, and running the example. It also explains CGO and cross-compilation basics and provides solutions for common setup issues.

## Prerequisites
Before you begin, ensure the following tools and configurations are in place:

- Go 1.16+ with CGO enabled
  - CGO must be enabled because the build mode is c-shared.
  - On Unix-like systems, ensure a C compiler is installed (gcc/clang).
  - On Windows, ensure a compatible C compiler is installed (MinGW-w64 or TDM-GCC).
  - Cross-compilation requires additional cross-compilers installed for target platforms.

- PHP 7.4+ with FFI extension enabled
  - The FFI extension must be enabled in php.ini.
  - The example expects FFI to be available and configured.

- Build tools
  - On Unix-like systems, GNU make is recommended.
  - On Windows, PowerShell is recommended.

These requirements are documented in the project’s documentation and build scripts.

**Section sources**
- [README.md](file://README.md#L45-L80)
- [README.md](file://README.md#L210-L237)
- [go.mod](file://go.mod#L1-L4)

## Quick Start
Follow these steps to build and test the example:

- Generate PHP bindings from Go exports
  - Run the generator to produce platform detection, FFI bindings, and loader files.

- Build the shared library for your platform
  - Use the Go build command with c-shared mode to produce the platform-specific binary.

- Copy PHP files to the distribution directory
  - Copy the generated PHP files into the dist folder alongside the library.

- Run the example
  - Execute the PHP example to load the library and call exported functions.

These steps are described in the Quick Start guide and align with the build scripts and generator.

**Section sources**
- [QUICKSTART.md](file://QUICKSTART.md#L1-L21)
- [README.md](file://README.md#L66-L109)

## Step-by-Step Setup
This section walks you through the complete setup process for first-time users.

1. Install prerequisites
   - Install Go 1.16+ and enable CGO.
   - Install a C compiler appropriate to your platform.
   - Install PHP 7.4+ and enable the FFI extension.

2. Clone the repository
   - Clone the repository to your local machine.

3. Generate PHP bindings
   - Run the generator to scan exported functions and produce PHP files for platform detection, FFI bindings, and loader.

4. Build the shared library for your platform
   - Use the Go build command with c-shared mode to produce the platform-specific binary.
   - Alternatively, use the provided build scripts to automate generation and building.

5. Copy PHP files to dist
   - Copy the generated PHP files into the dist directory.

6. Run the example
   - Execute the PHP example to load the library and call exported functions.

7. Optional: Build for all platforms
   - Use the cross-platform builder to compile for all supported platforms and architectures.

These steps are derived from the Quick Start and README documentation.

**Section sources**
- [QUICKSTART.md](file://QUICKSTART.md#L1-L21)
- [README.md](file://README.md#L66-L109)
- [README.md](file://README.md#L109-L149)

## Using the Build System
The project provides multiple ways to build the shared library and generate PHP bindings.

- Using Make (Unix-like systems)
  - make generate: Runs the generator to produce PHP files.
  - make build: Generates bindings and builds for all supported platforms.
  - make build-current: Generates bindings and builds for the current platform only.
  - make test: Builds for the current platform and runs the example.
  - make clean: Removes generated files and the dist directory.

- Using PowerShell (Windows)
  - .\build.ps1 generate: Runs the generator.
  - .\build.ps1 build: Generates bindings and builds for all platforms.
  - .\build.ps1 build-current: Generates bindings and builds for the current platform.
  - .\build.ps1 test: Builds for the current platform and runs the example.
  - .\build.ps1 clean: Removes generated files and the dist directory.

- Using the Go-based builder
  - The Go program orchestrates cross-platform builds, sets environment variables for cross-compilation, and copies generated PHP files into dist.

These commands and scripts are defined in the Makefile and build.ps1, and orchestrated by the Go builder.

**Section sources**
- [Makefile](file://Makefile#L1-L54)
- [build.ps1](file://build.ps1#L1-L152)
- [build.go](file://build.go#L1-L183)

## Calling Go Functions from PHP
The example demonstrates how to load the Go library and call exported functions through FFI. The process involves:

- Loading the library
  - Require the loader and call the loader function to initialize the library for the current platform.

- Calling functions
  - Use the returned library object to call exported functions.
  - For functions returning C strings, convert them to PHP strings and free the memory using the provided function.

- Platform information
  - Retrieve platform details to confirm the detected OS, architecture, and library filenames.

The example shows how to call integer-returning functions, boolean-returning functions, and string-returning functions, including proper memory management for strings.

**Section sources**
- [example.php](file://example.php#L1-L95)
- [README.md](file://README.md#L110-L149)

## Common Pitfalls and Fixes
This section highlights frequent setup issues and their solutions.

- FFI extension not loaded
  - Ensure the FFI extension is enabled in php.ini and restart your web server or PHP-FPM.

- Library file not found
  - Confirm that you ran the generator and built the library for your platform.
  - Ensure the dist directory contains the correct platform-specific binary and that the path passed to the loader is correct.

- CGO not enabled
  - Set the CGO enabled flag appropriately for your platform.
  - Ensure a C compiler is installed.

- DLL loading errors on Windows
  - Ensure architecture alignment between PHP and the DLL.
  - Install required runtime libraries if needed.
  - Temporarily adjust antivirus settings if necessary.

- Cross-compilation failures
  - Install the required cross-compilers for the target platform.
  - Review the cross-compilation setup in the documentation.

These pitfalls and fixes are documented in the project’s troubleshooting section.

**Section sources**
- [README.md](file://README.md#L238-L309)
- [QUICKSTART.md](file://QUICKSTART.md#L109-L134)

## Troubleshooting Guide
Use the following checklist to diagnose and resolve issues during setup and testing.

- Verify PHP FFI
  - Confirm the FFI extension is enabled and restart your server if necessary.

- Verify Go and CGO
  - Ensure Go is installed and CGO is enabled.
  - Confirm a C compiler is available for your platform.

- Verify generator output
  - Ensure the generator produced platform_detect.php, ffi_bindings.php, and loader.php.

- Verify build output
  - Confirm the dist directory contains the correct platform-specific binary and header file.
  - Ensure the binary is not empty.

- Verify loader path
  - Ensure the loader receives the correct path to the dist directory.

- Verify example execution
  - Run the example and review the console output for any errors.

If problems persist, consult the troubleshooting section in the documentation for platform-specific guidance.

**Section sources**
- [README.md](file://README.md#L238-L309)
- [example.php](file://example.php#L1-L95)

## Conclusion
You now have the essential steps to set up the environment, generate bindings, build the shared library, and call Go functions from PHP using FFI. Use the provided scripts and documentation to streamline the process and troubleshoot common issues. Extend the example by adding your own exported functions and regenerating bindings as needed.