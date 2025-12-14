package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gophpffi",
	Short: "Go-PHP FFI 服务库生成器",
	Long: `用于为 Go 共享库生成 PHP FFI 绑定的 CLI 工具。
	
完整文档请访问：https://github.com/wuwuseo/gophpffi`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
