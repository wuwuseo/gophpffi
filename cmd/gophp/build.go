package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [source.go]",
	Short: "构建 Go 共享库",
	Long: `从源文件构建 Go 共享库 (.dll/.so/.dylib)。
	
库文件将被放置在 dist/lib/ 目录中。`,
	Args: cobra.MaximumNArgs(1),
	RunE: runBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func runBuild(cmd *cobra.Command, args []string) error {
	// Load config or use args
	config, err := loadConfig()
	if err != nil && len(args) == 0 {
		return fmt.Errorf("未指定源文件且找不到 .gophp.yaml")
	}

	var sourceFile, serviceName string
	if len(args) > 0 {
		sourceFile = args[0]
		serviceName = filepath.Base(sourceFile)
		serviceName = serviceName[:len(serviceName)-3] // Remove .go
	} else {
		sourceFile = config.Source
		serviceName = config.Service
	}

	fmt.Println("=== 正在构建 Go 共享库 ===")
	fmt.Printf("源文件：%s\n", sourceFile)
	fmt.Printf("服务名：%s\n\n", serviceName)

	// Create output directories
	distDir := "dist"
	libDir := filepath.Join(distDir, "lib")
	if err := os.MkdirAll(libDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败：%w", err)
	}

	// Get platform info
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Determine extension
	var ext string
	switch goos {
	case "windows":
		ext = "dll"
	case "darwin":
		ext = "dylib"
	default:
		ext = "so"
	}

	// Build output path
	outputName := fmt.Sprintf("%s-%s-%s.%s", serviceName, goos, goarch, ext)
	outputPath := filepath.Join(libDir, outputName)

	fmt.Printf("正在为 %s-%s 构建...\n", goos, goarch)
	fmt.Printf("输出：%s\n\n", outputPath)

	// Build command
	buildCmd := exec.Command("go", "build",
		"-buildmode=c-shared",
		"-o", outputPath,
		sourceFile,
	)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("构建失败：%w", err)
	}

	fmt.Println("\n✓ 库文件已生成到 dist/lib/")
	fmt.Println()
	fmt.Println("=== 构建完成！===")
	fmt.Printf("库文件：%s\n", outputPath)
	fmt.Printf("头文件：%s\n", outputPath[:len(outputPath)-len(ext)]+"h")

	return nil
}
