# Go-PHP FFI 服务库生成器

一个用于为 Go 共享库生成 PHP FFI 绑定的 CLI 工具。

## 快速开始

### 1. 安装 CLI 工具

```bash
go install ./cmd/gophp
# 或者
make install
```

### 2. 初始化新服务

```bash
gophpffi init [ServiceName]
```

**示例：**
```bash
gophpffi init Product
```

这将创建：
- `Product.go` - 包含模板函数的 Go 源文件
- `.gophp.yaml` - 配置文件

### 3. 编辑你的 Go 文件

打开生成的 `.go` 文件并添加你的导出函数：

```go
//export YourFunction
func YourFunction(param Type) ReturnType {
    // 你的实现代码
    return result
}
```

**重要提示：**
- 每个导出函数必须有 `//export FunctionName` 注释
- 注释必须直接位于函数上方（不能有空行）
- 函数必须在 `package main` 中
- 保留 `main()` 函数（构建共享库所需）

### 4. 构建所有内容

```bash
gophpffi make
```

这将：
1. 在 `dist/` 中生成 PHP 服务类
2. 在 `dist/lib/` 中构建共享库

## 目录结构

构建后，你将拥有：

```
dist/
├── [ServiceName]Service.php          (PHP 服务类)
└── lib/
    ├── [ServiceName]-windows-amd64.dll  (共享库)
    └── [ServiceName]-windows-amd64.h    (C 头文件)
```

## CLI 命令

### 仅生成 PHP 绑定
```bash
gophpffi generate
# 或者
go generate ./...
```

### 仅构建共享库
```bash
gophpffi build
```

### 完整构建
```bash
gophpffi make
```

## 配置文件

项目使用 `.gophp.yaml` 配置文件：

```yaml
service: ServiceName        # 服务名称
source: ServiceName.go      # Go 源文件路径
output:
  dir: dist                 # 输出目录
  lib_dir: dist/lib         # 库文件目录
```

你可以手动编辑此文件来自定义构建设置。

## 示例

### 步骤 1：初始化
```bash
gophpffi init User
```

### 步骤 2：编辑 User.go
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

### 步骤 3：构建
```bash
gophpffi make
```

### 步骤 4：在 PHP 中使用
```php
<?php
require_once 'dist/UserService.php';

use app\user\service\UserService;

$service = new UserService();
$name = $service->GetUserName(123);
echo $name; // Output: User_123
```

## 类型映射

| Go 类型 | PHP 类型 |
|---------|----------|
| int, int8, int16, int32, int64 | int |
| uint, uint8, uint16, uint32, uint64 | int |
| float32, float64 | float |
| string | string |
| bool | bool |
| []T (slice) | array |
| map[K]V | array |

## 故障排除

### 构建失败
- 确保已安装 Go：`go version`
- 检查所有导出函数都有 `//export` 注释
- 确保没有未使用的导入

### 找不到生成的 PHP 文件
- 首先运行 `gophpffi generate` 或 `go generate ./...`
- 检查你的 `.go` 文件是否有 `//go:generate` 指令

### 找不到 DLL
- 运行 `gophpffi build` 来编译共享库
- 确保你在 Windows 上（对于 .dll 文件）

### 找不到 CLI 命令
- 如果使用 `go install` 安装，确保 `$GOPATH/bin` 在你的 PATH 中
- 否则，从项目目录使用 `./gophpffi.exe`

## 文件说明

- `gophpffi` - CLI 工具（命令行界面）
- `cmd/gophp/` - CLI 工具源代码
- `generator/main.go` - 代码生成器源码
- `.gophp.yaml` - 项目配置文件
- `AGENTS.MD` - 开发指南（中文）
- `AGENTS_EN.MD` - 开发指南（英文）
- `README_EN.MD` - 用户指南（英文）
- `Makefile` - 构建自动化

## 高级用法

### 构建 CLI 工具

```bash
# 构建 CLI 工具
go build -o gophpffi.exe ./cmd/gophp

# 或使用 Makefile
make build
```

### 安装到系统

```bash
# 安装到 GOPATH/bin
go install ./cmd/gophp

# 或使用 Makefile
make install
```

### 多服务支持

你可以在同一项目中创建多个服务：

```bash
gophpffi init UserService
gophpffi init ProductService
gophpffi init OrderService
```

为不同服务构建时，需要指定源文件：

```bash
gophpffi generate UserService.go
gophpffi build UserService.go
```

## 系统要求

- Go 1.16 或更高版本
- Windows（用于构建 .dll 文件）、Linux（.so）或 macOS（.dylib）
- PHP 7.4 或更高版本，需启用 FFI 扩展

### 启用 PHP FFI

在 `php.ini` 中：
```ini
extension=ffi
ffi.enable=true
```

检查 FFI 是否已启用：
```bash
php -m | grep ffi
```

## 最佳实践

### 1. 函数设计
- 保持函数签名简单
- 尽可能使用基本类型（int、string、bool）
- 避免复杂的嵌套结构
- 返回错误码而非 Go 的 error 类型

### 2. 错误处理

**在 Go 中：**
```go
//export ProcessData
func ProcessData(data string) int {
    if data == "" {
        return -1  // 错误码
    }
    // 处理数据...
    return 0  // 成功
}
```

**在 PHP 中：**
```php
$result = $service->ProcessData($data);
if ($result < 0) {
    throw new Exception("处理失败，错误码: " . $result);
}
```

### 3. 性能优化
- 减少跨语言调用次数
- 尽可能批量处理
- 使用持久化 PHP 进程（PHP-FPM）
- 缓存 FFI 库实例

## 许可证

MIT License - 详见 LICENSE 文件

## 贡献

欢迎贡献！请查看 [AGENTS.MD](AGENTS.MD) 了解开发指南。

## 支持

- 开发文档：[AGENTS.MD](AGENTS.MD)（中文）、[AGENTS_EN.MD](AGENTS_EN.MD)（英文）
- 英文文档：[README_EN.MD](README_EN.MD)
- 问题反馈：[GitHub Issues](https://github.com/wuwuseo/gophpffi/issues)
