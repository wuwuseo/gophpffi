package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [ServiceName]",
	Short: "初始化一个新的 Go 服务",
	Long: `使用模板代码初始化一个新的 Go 服务。
	
将会创建：
  - [ServiceName].go 包含模板函数
  - 更新构建配置`,
	Args: cobra.ExactArgs(1),
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	serviceName := args[0]
	goFile := serviceName + ".go"

	fmt.Printf("========================================\n")
	fmt.Printf("   初始化新服务：%s\n", serviceName)
	fmt.Printf("========================================\n\n")

	// Check if file exists
	if _, err := os.Stat(goFile); err == nil {
		fmt.Printf("警告：%s 已存在！\n", goFile)
		fmt.Print("是否要覆盖它？(y/N): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println("已取消。")
			return nil
		}
		fmt.Println()
	}

	// Create Go source file
	fmt.Printf("[1/3] 正在创建 %s...\n", goFile)
	goContent := fmt.Sprintf(`package main

/*
#include <stdlib.h>
*/
import "C"

//go:generate gophpffi generate %s

// Add is a simple addition function for demonstration
//
//export Add
func Add(a, b int) int {
	return a + b
}

// Echo returns the input string
//
//export Echo
func Echo(s string) string {
	return s
}

// main is required for building shared library but never called
func main() {}
`, goFile)

	if err := os.WriteFile(goFile, []byte(goContent), 0644); err != nil {
		return fmt.Errorf("创建 %s 失败：%w", goFile, err)
	}
	fmt.Printf("✓ 已创建 %s\n\n", goFile)

	// Update .gophp.yaml config file
	fmt.Println("[2/3] 正在创建 .gophp.yaml 配置文件...")
	configContent := fmt.Sprintf(`# Go-PHP FFI Service Configuration
service: %s
source: %s
output:
  dir: dist
  lib_dir: dist/lib
`, serviceName, goFile)

	if err := os.WriteFile(".gophp.yaml", []byte(configContent), 0644); err != nil {
		return fmt.Errorf("创建配置文件失败：%w", err)
	}
	fmt.Println("✓ 已创建 .gophp.yaml\n")

	fmt.Println("[3/3] 设置完成")
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("   初始化完成！")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Printf("已创建文件：\n")
	fmt.Printf("  - %s          (Go 源文件)\n", goFile)
	fmt.Printf("  - .gophp.yaml      (配置文件)\n")
	fmt.Println()
	fmt.Println("下一步：")
	fmt.Printf("  1. 编辑 %s 并添加你的导出函数\n", goFile)
	fmt.Println("  2. 运行：gophpffi make   (生成 PHP 绑定并构建)")
	fmt.Println()
	fmt.Println("导出函数示例格式：")
	fmt.Println("  //export YourFunction")
	fmt.Println("  func YourFunction(param Type) ReturnType {")
	fmt.Println("      // 你的代码")
	fmt.Println("  }")

	return nil
}
