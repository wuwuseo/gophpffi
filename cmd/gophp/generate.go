package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate [source.go]",
	Short: "生成 PHP FFI 绑定",
	Long: `从 Go 源文件生成 PHP FFI 绑定。
	
将会解析 Go 源文件并在 dist/ 目录中创建 PHP 服务类。`,
	Args: cobra.MaximumNArgs(1),
	RunE: runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Load config or use args
	config, err := loadConfig()
	if err != nil && len(args) == 0 {
		return fmt.Errorf("未指定源文件且找不到 .gophp.yaml")
	}

	var sourceFile string
	if len(args) > 0 {
		sourceFile = args[0]
	} else {
		sourceFile = config.Source
	}

	fmt.Println("=== Go-PHP FFI 代码生成器 ===")
	fmt.Printf("正在为以下文件生成 PHP 绑定：%s\n\n", sourceFile)

	// Get absolute path
	absPath, err := filepath.Abs(sourceFile)
	if err != nil {
		return fmt.Errorf("获取绝对路径失败：%w", err)
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败：%w", err)
	}

	// Run the generator
	generatorDir := filepath.Join(cwd, "generator")
	genCmd := exec.Command("go", "run", ".", absPath)
	genCmd.Dir = generatorDir
	genCmd.Stdout = os.Stdout
	genCmd.Stderr = os.Stderr

	if err := genCmd.Run(); err != nil {
		return fmt.Errorf("生成失败：%w", err)
	}

	return nil
}
