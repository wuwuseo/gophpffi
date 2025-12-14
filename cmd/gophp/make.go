package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "生成绑定并构建库（完整构建）",
	Long: `完整的构建流程：生成 PHP 绑定并构建共享库。
	
相当于依次运行 'gophpffi generate' 和 'gophpffi build'。`,
	RunE: runMake,
}

func init() {
	rootCmd.AddCommand(makeCmd)
}

func runMake(cmd *cobra.Command, args []string) error {
	fmt.Println("========================================")
	fmt.Println("   Go-PHP FFI 构建流程")
	fmt.Println("========================================")
	fmt.Println()

	// Step 1: Generate
	fmt.Println("[1/2] 正在生成 PHP 绑定...")
	if err := runGenerate(cmd, args); err != nil {
		return fmt.Errorf("生成失败：%w", err)
	}
	fmt.Println()

	// Step 2: Build
	fmt.Println("[2/2] 正在构建共享库...")
	if err := runBuild(cmd, args); err != nil {
		return fmt.Errorf("构建失败：%w", err)
	}
	fmt.Println()

	fmt.Println("========================================")
	fmt.Println("   构建完成！")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("请检查 dist/ 目录中的生成文件。")
	fmt.Println("使用方法：在你的 PHP 代码中导入 PHP 服务文件。")

	return nil
}
