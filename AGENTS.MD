# AGENTS.MD - Go-PHP FFI 项目开发指南

## 项目概述

Go-PHP FFI Service Library Generator 是一个自动化工具，用于为 Go 共享库生成 PHP FFI 绑定。该项目使用 Cobra 框架构建 CLI 工具，提供完整的代码生成和构建流程。

### 核心功能
- 自动生成 PHP FFI 绑定代码
- 构建跨平台共享库（.dll/.so/.dylib）
- 提供类型映射和安全的内存管理
- 基于模板的服务初始化

## 技术栈

### 后端技术
- **语言**: Go 1.16+
- **CLI 框架**: Cobra
- **配置管理**: YAML (gopkg.in/yaml.v3)
- **构建模式**: c-shared
- **代码生成**: Go AST 解析

### PHP 集成
- **PHP 版本**: 7.4+
- **扩展要求**: FFI extension
- **类型系统**: 基于 C 类型的 FFI 绑定

## 项目结构

```
serviceLib/
├── cmd/
│   └── gophp/              # CLI 工具源代码
│       ├── main.go         # 主入口和根命令
│       ├── init.go         # init 命令实现
│       ├── generate.go     # generate 命令实现
│       ├── build.go        # build 命令实现
│       ├── make.go         # make 命令实现
│       └── config.go       # 配置文件加载
├── generator/
│   └── main.go             # 代码生成器核心逻辑
├── dist/                   # 构建输出目录
│   ├── *Service.php        # 生成的 PHP 服务类
│   └── lib/                # 共享库文件
├── .gophp.yaml             # 项目配置文件
├── Makefile                # 构建自动化
├── README.md               # 项目文档
└── go.mod                  # Go 模块定义
```

## CLI 工具架构

### 命令结构

```
gophpffi
├── init [ServiceName]      # 初始化新服务
├── generate [source.go]    # 生成 PHP 绑定
├── build [source.go]       # 构建共享库
└── make                    # 完整构建流程
```

### 命令实现细节

#### 1. `gophpffi init`
- **功能**: 创建新服务的模板文件
- **输出**:
  - `[ServiceName].go` - 包含示例函数的 Go 源文件
  - `.gophp.yaml` - 项目配置文件
- **文件**: `cmd/gophp/init.go`

#### 2. `gophpffi generate`
- **功能**: 解析 Go 源文件，生成 PHP FFI 绑定
- **流程**:
  1. 读取配置文件或命令行参数
  2. 调用 `generator/main.go` 解析 Go AST
  3. 生成 PHP 服务类到 `dist/` 目录
- **文件**: `cmd/gophp/generate.go`

#### 3. `gophpffi build`
- **功能**: 编译 Go 源文件为共享库
- **构建参数**: `-buildmode=c-shared`
- **输出**: 
  - `[ServiceName]-[os]-[arch].dll/.so/.dylib`
  - `[ServiceName]-[os]-[arch].h`
- **文件**: `cmd/gophp/build.go`

#### 4. `gophpffi make`
- **功能**: 执行完整构建流程
- **步骤**:
  1. 生成 PHP 绑定 (generate)
  2. 构建共享库 (build)
- **文件**: `cmd/gophp/make.go`

## 配置文件

### `.gophp.yaml` 结构

```yaml
service: ServiceName        # 服务名称
source: ServiceName.go      # Go 源文件路径
output:
  dir: dist                 # 输出目录
  lib_dir: dist/lib         # 库文件目录
```

### 配置加载逻辑
- 位置: `cmd/gophp/config.go`
- 功能: 
  - 读取并解析 `.gophp.yaml`
  - 提供默认值
  - 验证配置有效性

## 代码生成器

### 核心文件: `generator/main.go`

#### 功能模块
1. **Go AST 解析**
   - 解析导出函数（带 `//export` 注释）
   - 提取函数签名、参数、返回值
   - 识别类型信息

2. **类型映射**
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

3. **代码生成**
   - 生成 PHP 类定义
   - 创建 FFI 绑定方法
   - 添加类型提示和文档注释

## 开发工作流

### 标准开发流程

```bash
# 1. 初始化新服务
gophpffi init MyService

# 2. 编辑 MyService.go，添加导出函数
# //export MyFunction
# func MyFunction(param Type) ReturnType { ... }

# 3. 生成和构建
gophpffi make

# 4. 在 PHP 中使用
# require_once 'dist/MyServiceService.php';
```

### 调试技巧

1. **查看生成的代码**
   ```bash
   gophpffi generate
   cat dist/*Service.php
   ```

2. **验证构建输出**
   ```bash
   gophpffi build
   ls -la dist/lib/
   ```

3. **测试单个命令**
   ```bash
   gophpffi generate MyService.go
   gophpffi build MyService.go
   ```

## 代码规范

### Go 代码规范

1. **导出函数要求**
   ```go
   // ✅ 正确
   //export MyFunction
   func MyFunction(a int) int {
       return a * 2
   }
   
   // ❌ 错误 - 缺少 export 注释
   func MyFunction(a int) int {
       return a * 2
   }
   
   // ❌ 错误 - 注释和函数之间有空行
   //export MyFunction
   
   func MyFunction(a int) int {
       return a * 2
   }
   ```

2. **包要求**
   - 必须使用 `package main`
   - 必须包含 `main()` 函数（即使为空）
   - 必须导入 `"C"`

3. **注释规范**
   - 使用中文注释
   - 为导出函数添加功能说明
   - 记录参数和返回值含义

### CLI 代码规范

1. **命令实现**
   - 每个命令一个文件
   - 使用 `cobra.Command` 结构
   - 提供中文帮助信息

2. **错误处理**
   - 使用 `fmt.Errorf` 包装错误
   - 提供中文错误消息
   - 返回有意义的错误信息

3. **输出格式**
   - 使用统一的分隔线样式
   - 步骤编号: `[1/3]`, `[2/3]`
   - 成功标记: `✓`

## 构建和发布

### 本地构建

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

### 测试构建

```bash
# 运行测试
make test

# 手动测试
./gophpffi.exe --help
./gophpffi.exe init TestService
./gophpffi.exe make
```

## 常见问题排查

### 构建问题

1. **CGO 未启用**
   - 症状: `CGO_ENABLED` 相关错误
   - 解决: 设置环境变量 `CGO_ENABLED=1`

2. **找不到编译器**
   - 症状: `gcc` 或 `cl.exe` 未找到
   - 解决: 安装 MinGW-w64 (Windows) 或 GCC (Linux)

3. **导出函数未生成**
   - 症状: PHP 类中缺少方法
   - 检查: 确保有 `//export` 注释且无空行

### 运行时问题

1. **DLL 加载失败**
   - 检查: 平台架构匹配（64位/32位）
   - 检查: 依赖的运行库是否安装

2. **配置文件读取失败**
   - 检查: `.gophp.yaml` 格式正确
   - 检查: YAML 缩进使用空格而非 Tab

## 扩展开发

### 添加新命令

1. 在 `cmd/gophp/` 下创建新文件
2. 定义 `cobra.Command` 结构
3. 在 `init()` 中注册到 `rootCmd`
4. 实现命令逻辑

示例：
```go
var newCmd = &cobra.Command{
    Use:   "new",
    Short: "新命令描述",
    Long:  `详细说明`,
    RunE:  runNew,
}

func init() {
    rootCmd.AddCommand(newCmd)
}

func runNew(cmd *cobra.Command, args []string) error {
    // 实现逻辑
    return nil
}
```

### 修改代码生成器

1. 编辑 `generator/main.go`
2. 修改 AST 解析逻辑
3. 更新 PHP 模板生成
4. 测试生成的代码

## 版本历史

### 当前版本特性
- ✅ CLI 工具完整中文化
- ✅ 命令名从 `gophp` 迁移到 `gophpffi`
- ✅ 支持配置文件驱动
- ✅ 自动类型映射
- ✅ 跨平台构建支持
- ✅ 双语文档体系（中文/英文）
- ✅ MIT 开源许可证

### 从旧版本迁移

| 旧命令 | 新命令 |
|--------|--------|
| `init.bat Name` | `gophpffi init Name` |
| `generate.bat` | `gophpffi generate` |
| `build.bat` | `gophpffi build` |
| `make.bat` | `gophpffi make` |
| `gophp *` | `gophpffi *` |

## 最佳实践

### 1. 项目组织
- 每个服务一个 Go 文件
- 使用有意义的服务名
- 保持函数签名简单

### 2. 函数设计
- 避免使用复杂的嵌套类型
- 优先使用基本类型
- 为切片和映射提供清晰的文档

### 3. 错误处理
- Go 函数返回错误码而非 error 类型
- 在 PHP 端检查返回值
- 记录详细的错误信息

### 4. 性能优化
- 减少跨语言调用次数
- 批量处理数据
- 使用持久化 PHP 进程（PHP-FPM）

## 贡献指南

### 代码提交规范
- 提交信息使用中文
- 遵循现有代码风格
- 添加必要的注释
- 更新相关文档

### 测试要求
- 测试所有 CLI 命令
- 验证生成的代码可运行
- 检查跨平台兼容性

## 技术支持

### 文档资源
- `README.md` - 用户指南（中文）
- `README_EN.MD` - 用户指南（英文）
- `AGENTS.MD` - 开发指南（中文）
- `AGENTS_EN.MD` - 开发指南（英文）
- `LICENSE` - MIT 许可证

## 项目文件说明

- `gophpffi` - CLI 工具（命令行界面）
- `cmd/gophp/` - CLI 工具源代码
- `generator/main.go` - 代码生成器源码
- `.gophp.yaml` - 项目配置文件
- `README.md` - 用户指南（中文）
- `README_EN.MD` - 用户指南（英文）
- `AGENTS.MD` - 开发指南（中文）
- `AGENTS_EN.MD` - 开发指南（英文）
- `LICENSE` - MIT 许可证
- `Makefile` - 构建自动化

### 问题反馈
- GitHub Issues: https://github.com/wuwuseo/gophpffi/issues
- 检查现有文档
- 查看常见问题
- 提供详细的错误信息和环境描述

## 许可证

MIT License

Copyright (c) 2025 wuwuseo

详见 [LICENSE](LICENSE) 文件。
