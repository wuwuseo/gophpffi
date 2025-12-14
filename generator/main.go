package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ExportedFunc 表示一个导出的 Go 函数
type ExportedFunc struct {
	Name       string
	Comment    string
	Signature  string
	ReturnType string
	Params     []Param
}

// Param 表示一个函数参数
type Param struct {
	Name string
	Type string
}

func main() {
	fmt.Println("=== Go-PHP FFI Code Generator ===")

	// 从命令行参数获取 Go 源文件，默认为 mygo.go
	sourceFile := "mygo.go"
	if len(os.Args) > 1 {
		sourceFile = os.Args[1]
	}

	// 获取源文件的绝对路径
	absSourceFile, err := filepath.Abs(sourceFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting absolute path: %v\n", err)
		os.Exit(1)
	}
	sourceFile = absSourceFile

	fmt.Printf("Parsing Go source file: %s\n", sourceFile)

	// 提取库基本名称（不含 .go 扩展名）
	baseName := strings.TrimSuffix(filepath.Base(sourceFile), ".go")
	fmt.Printf("Library base name: %s\n", baseName)

	// 从源文件解析导出的函数
	exports, err := parseExports(sourceFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing exports: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d exported functions\n", len(exports))
	for _, exp := range exports {
		fmt.Printf("  - %s\n", exp.Name)
	}

	// 生成 PHP 文件（输出到 dist 目录）
	sourceDir := filepath.Dir(sourceFile)
	distDir := filepath.Join(sourceDir, "dist")
	libDir := filepath.Join(distDir, "lib")

	// 创建 dist 和 dist/lib 目录
	if err := os.MkdirAll(libDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directories: %v\n", err)
		os.Exit(1)
	}

	if err := generateFFIBindings(exports, sourceFile, distDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating Service.php: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Generated Service.php in dist/")
	fmt.Println("✓ Created dist/lib/ directory for library files")

	fmt.Println("\n=== Code generation complete! ===")
	fmt.Println("Next step: Run 'go run build.go' to build shared libraries for all platforms")
}

// parseExports 读取 mygo.go 并提取导出的函数
func parseExports(filename string) ([]ExportedFunc, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var exports []ExportedFunc
	scanner := bufio.NewScanner(file)

	exportRegex := regexp.MustCompile(`^//export\s+(\w+)`)
	funcRegex := regexp.MustCompile(`^func\s+(\w+)\s*\((.*?)\)\s*(.*)`)

	var currentComment strings.Builder
	var isExported bool

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// 检查 //export 指令
		if matches := exportRegex.FindStringSubmatch(trimmed); matches != nil {
			isExported = true
			continue
		}

		// 收集注释
		if strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "//export") && !strings.HasPrefix(trimmed, "//go:") {
			comment := strings.TrimPrefix(trimmed, "//")
			comment = strings.TrimSpace(comment)
			if comment != "" {
				currentComment.WriteString(comment)
				currentComment.WriteString(" ")
			}
			continue
		}

		// 解析函数声明
		if isExported && strings.HasPrefix(trimmed, "func ") {
			if matches := funcRegex.FindStringSubmatch(trimmed); matches != nil {
				funcName := matches[1]
				paramsStr := matches[2]
				returnStr := strings.TrimSpace(matches[3])

				params := parseParams(paramsStr)
				returnType := parseReturnType(returnStr)

				exports = append(exports, ExportedFunc{
					Name:       funcName,
					Comment:    strings.TrimSpace(currentComment.String()),
					Signature:  trimmed,
					ReturnType: returnType,
					Params:     params,
				})

				currentComment.Reset()
				isExported = false
			}
		}

		// 如果遇到非注释行且不是函数，则重置注释
		if !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "func ") && trimmed != "" {
			currentComment.Reset()
		}
	}

	return exports, scanner.Err()
}

// parseParams 从函数签名中提取参数
func parseParams(paramsStr string) []Param {
	if paramsStr == "" {
		return []Param{}
	}

	var params []Param
	parts := strings.Split(paramsStr, ",")

	// 先收集所有参数，然后处理类型
	var tempParams []struct {
		Name string
		Type string
	}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		fields := strings.Fields(part)
		if len(fields) >= 2 {
			// 有类型声明的参数，如 "a int" 或 "arr []int"
			paramName := fields[0]
			paramType := strings.Join(fields[1:], " ")
			tempParams = append(tempParams, struct {
				Name string
				Type string
			}{Name: paramName, Type: paramType})
		} else if len(fields) == 1 {
			// 只有参数名
			tempParams = append(tempParams, struct {
				Name string
				Type string
			}{Name: fields[0], Type: ""})
		}
	}

	// 反向处理，将类型向前传播（处理 a, b int 这种情况）
	var currentType string
	for i := len(tempParams) - 1; i >= 0; i-- {
		if tempParams[i].Type != "" {
			currentType = tempParams[i].Type
		} else if currentType != "" {
			tempParams[i].Type = currentType
		}
	}

	// 转换为 Param 切片
	for _, tp := range tempParams {
		if tp.Type != "" {
			params = append(params, Param{
				Name: tp.Name,
				Type: tp.Type,
			})
		}
	}

	return params
}

// parseReturnType 从函数签名中提取返回类型
func parseReturnType(returnStr string) string {
	returnStr = strings.TrimSpace(returnStr)
	if returnStr == "" {
		return "void"
	}
	// 移除 { 符号（函数体开始）
	if idx := strings.Index(returnStr, "{"); idx != -1 {
		returnStr = returnStr[:idx]
		returnStr = strings.TrimSpace(returnStr)
	}
	// 如果存在括号则移除
	returnStr = strings.TrimPrefix(returnStr, "(")
	returnStr = strings.TrimSuffix(returnStr, ")")
	returnStr = strings.TrimSpace(returnStr)
	if returnStr == "" {
		return "void"
	}
	return returnStr
}

// toSnakeCase 将 PascalCase 转换为 snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// toPascalCase 将 snake_case 转换为 PascalCase
func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

// generateFFIBindings 生成 service
func generateFFIBindings(exports []ExportedFunc, filename string, outputDir string) error {
	var sb strings.Builder

	// 将 filename 转换为首字母大写驼峰格式
	baseName := strings.TrimSuffix(filepath.Base(filename), ".go")
	className := toPascalCase(baseName)
	snakeName := toSnakeCase(baseName)

	// 生成 PHP 类头部
	sb.WriteString(fmt.Sprintf(`<?php
/**
 * FFI Bindings Service
 * Auto-generated by Go-PHP FFI Code Generator
 * 
 * Provides PHP interface to Go shared library functions
 */

namespace app\%s\service;

use Wuwuseo\PhpffiGoLibrary\GoLibraryBase;

class %sService extends GoLibraryBase {

	/**
     * 获取基础目录
     * @return string
     */
    protected function getBaseDir(): string
    {
        return dirname(__DIR__);
    }
`, snakeName, className))

	// 为每个导出的函数生成包装方法
	for _, exp := range exports {
		sb.WriteString(fmt.Sprintf("    /**\n"))
		if exp.Comment != "" {
			sb.WriteString(fmt.Sprintf("     * %s\n", exp.Comment))
		}
		sb.WriteString(generatePHPMethodSignatureDoc(exp))
		sb.WriteString(fmt.Sprintf("     */\n"))
		sb.WriteString(generatePHPMethod(exp))
		sb.WriteString("\n")
	}

	sb.WriteString(`}
`)

	// 使用动态生成的文件名，输出到指定目录
	outputFile := filepath.Join(outputDir, fmt.Sprintf("%sService.php", className))
	return os.WriteFile(outputFile, []byte(sb.String()), 0644)
}

// generatePHPMethodSignatureDoc 为方法生成 PHPDoc
func generatePHPMethodSignatureDoc(exp ExportedFunc) string {
	var sb strings.Builder

	for _, param := range exp.Params {
		phpType := cTypeToPHPType(param.Type)
		sb.WriteString(fmt.Sprintf("     * @param %s $%s\n", phpType, param.Name))
	}

	returnPHPType := cTypeToPHPType(exp.ReturnType)
	sb.WriteString(fmt.Sprintf("     * @return %s\n", returnPHPType))

	return sb.String()
}

// generatePHPMethod 生成 PHP 方法包装器
func generatePHPMethod(exp ExportedFunc) string {
	var sb strings.Builder

	// 方法签名
	sb.WriteString(fmt.Sprintf("    public function %s(", exp.Name))

	paramStrs := []string{}
	for _, param := range exp.Params {
		phpType := cTypeToPHPDoc(param.Type)
		if phpType != "" {
			paramStrs = append(paramStrs, fmt.Sprintf("%s $%s", phpType, param.Name))
		} else {
			paramStrs = append(paramStrs, fmt.Sprintf("$%s", param.Name))
		}
	}
	sb.WriteString(strings.Join(paramStrs, ", "))
	sb.WriteString(")")

	// 返回类型
	returnType := cTypeToPHPDoc(exp.ReturnType)
	if returnType != "" {
		sb.WriteString(fmt.Sprintf(": %s", returnType))
	}

	sb.WriteString(" {\n")

	// 方法体 - 调用 FFI 函数
	if exp.ReturnType == "void" || exp.ReturnType == "" {
		sb.WriteString(fmt.Sprintf("        $this->ffi->%s(", exp.Name))
	} else {
		sb.WriteString(fmt.Sprintf("        return $this->ffi->%s(", exp.Name))
	}

	callParams := []string{}
	for _, param := range exp.Params {
		callParams = append(callParams, fmt.Sprintf("$%s", param.Name))
	}
	sb.WriteString(strings.Join(callParams, ", "))
	sb.WriteString(");\n")

	sb.WriteString("    }\n")

	return sb.String()
}

// cTypeToPHPType 将 C/Go 类型转换为 PHP 类型用于文档
func cTypeToPHPType(cType string) string {
	cType = strings.TrimSpace(cType)

	// 字符串类型（char* 和 string 类型）
	if strings.Contains(cType, "*C.char") || cType == "*C.char" || cType == "char*" {
		return "string"
	}
	if cType == "string" || cType == "GoString" {
		return "string"
	}

	if strings.Contains(cType, "GoMap") {
		return "array"
	}
	if strings.Contains(cType, "GoSlice") {
		return "array"
	}
	if strings.Contains(cType, "[]") {
		return "array"
	}

	// 整数类型 - C
	if cType == "int" || cType == "int8" || cType == "int16" || cType == "int32" || cType == "int64" {
		return "int"
	}
	if cType == "uint" || cType == "uint8" || cType == "uint16" || cType == "uint32" || cType == "uint64" {
		return "int"
	}
	if cType == "char" || cType == "short" || cType == "long" || cType == "long long" {
		return "int"
	}
	if cType == "unsigned char" || cType == "unsigned short" || cType == "unsigned int" || cType == "unsigned long" {
		return "int"
	}
	if cType == "size_t" || cType == "ssize_t" {
		return "int"
	}

	// 整数类型 - Go/CGO
	if cType == "GoInt" || cType == "GoInt8" || cType == "GoInt16" || cType == "GoInt32" || cType == "GoInt64" {
		return "int"
	}
	if cType == "GoUint" || cType == "GoUint8" || cType == "GoUint16" || cType == "GoUint32" || cType == "GoUint64" {
		return "int"
	}
	if cType == "C.int" || cType == "C.long" || cType == "C.short" || cType == "C.char" {
		return "int"
	}

	// 浮点数类型 - C
	if cType == "float" || cType == "float32" || cType == "float64" || cType == "double" {
		return "float"
	}

	// 浮点数类型 - Go/CGO
	if cType == "GoFloat32" || cType == "GoFloat64" {
		return "float"
	}
	if cType == "C.float" || cType == "C.double" {
		return "float"
	}

	// 布尔类型
	if cType == "bool" || cType == "GoBool" || cType == "_Bool" {
		return "bool"
	}

	// void 类型
	if cType == "void" || cType == "" {
		return "void"
	}

	return "mixed"
}

// cTypeToPHPDoc 将 C/Go 类型转换为 PHP 类型提示
func cTypeToPHPDoc(cType string) string {
	cType = strings.TrimSpace(cType)

	// 字符串类型（char* 和 string 类型）
	if strings.Contains(cType, "*C.char") || cType == "*C.char" || cType == "char*" {
		return "string"
	}
	if cType == "string" || cType == "GoString" {
		return "string"
	}

	// 数组/指针类型（排除 char*）
	// Go Map、Slice 和 Array 类型
	if strings.Contains(cType, "map[") || strings.HasPrefix(cType, "map[") {
		return "array"
	}
	if strings.Contains(cType, "GoMap") {
		return "array"
	}
	if strings.Contains(cType, "GoSlice") {
		return "array"
	}
	if strings.Contains(cType, "[]") {
		return "array"
	}
	if strings.HasPrefix(cType, "*") && !strings.Contains(cType, "char") {
		return "array"
	}

	// 整数类型 - C
	if cType == "int" || cType == "int8" || cType == "int16" || cType == "int32" || cType == "int64" {
		return "int"
	}
	if cType == "uint" || cType == "uint8" || cType == "uint16" || cType == "uint32" || cType == "uint64" {
		return "int"
	}
	if cType == "char" || cType == "short" || cType == "long" || cType == "long long" {
		return "int"
	}
	if cType == "unsigned char" || cType == "unsigned short" || cType == "unsigned int" || cType == "unsigned long" {
		return "int"
	}
	if cType == "size_t" || cType == "ssize_t" {
		return "int"
	}

	// 整数类型 - Go/CGO
	if cType == "GoInt" || cType == "GoInt8" || cType == "GoInt16" || cType == "GoInt32" || cType == "GoInt64" {
		return "int"
	}
	if cType == "GoUint" || cType == "GoUint8" || cType == "GoUint16" || cType == "GoUint32" || cType == "GoUint64" {
		return "int"
	}
	if cType == "C.int" || cType == "C.long" || cType == "C.short" || cType == "C.char" {
		return "int"
	}

	// 浮点数类型 - C
	if cType == "float" || cType == "float32" || cType == "float64" || cType == "double" {
		return "float"
	}

	// 浮点数类型 - Go/CGO
	if cType == "GoFloat32" || cType == "GoFloat64" {
		return "float"
	}
	if cType == "C.float" || cType == "C.double" {
		return "float"
	}

	// 布尔类型
	if cType == "bool" || cType == "GoBool" || cType == "_Bool" {
		return "bool"
	}

	// void 类型（不使用类型提示）
	if cType == "void" || cType == "" {
		return ""
	}

	// 未知类型：返回空（不使用类型提示）
	return ""
}
