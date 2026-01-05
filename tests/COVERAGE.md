# 测试执行与代码覆盖率指南

本文档介绍如何执行测试用例、生成代码覆盖率报告以及查看覆盖率详情。

## 目录

- [快速开始](#快速开始)
- [执行测试用例](#执行测试用例)
- [生成代码覆盖率报告](#生成代码覆盖率报告)
- [查看覆盖率报告](#查看覆盖率报告)
- [测试命令详解](#测试命令详解)
- [覆盖率目标](#覆盖率目标)

## 快速开始

### Windows 系统

```bash
# 运行所有测试
tests\run_tests.bat

# 运行测试并生成覆盖率报告
tests\run_tests.bat -cover -coverprofile=coverage.out

# 运行测试并生成 HTML 覆盖率报告
tests\run_tests.bat -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Linux/macOS 系统

```bash
# 运行所有测试
./tests/run_tests.sh

# 运行测试并生成覆盖率报告
./tests/run_tests.sh -cover -coverprofile=coverage.out

# 运行测试并生成 HTML 覆盖率报告
./tests/run_tests.sh -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 使用 Makefile

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make coverage

# 运行测试并生成 HTML 覆盖率报告
make coverage-html
```

## 执行测试用例

### 1. 运行所有测试

```bash
go test ./...
```

### 2. 运行特定包的测试

```bash
# 运行 pkg 包的测试
go test ./pkg/...

# 运行 internal 包的测试
go test ./internal/...

# 运行 scripts 包的测试
go test ./scripts/...
```

### 3. 运行特定文件的测试

```bash
# 运行 image_test.go
go test ./pkg -run TestLoadImageWithDecode
```

### 4. 详细输出模式

```bash
# 显示详细的测试输出
go test -v ./...

# 显示详细的测试输出并运行 pkg 包
go test -v ./pkg/...
```

### 5. 启用竞态检测

```bash
# 启用竞态检测器运行测试
go test -race ./...
```

## 生成代码覆盖率报告

### 1. 生成覆盖率报告

```bash
# 生成覆盖率报告（输出到控制台）
go test -cover ./...

# 生成覆盖率报告并保存到文件
go test -coverprofile=coverage.out ./...

# 生成特定包的覆盖率报告
go test -coverprofile=coverage.out ./pkg/...
```

### 2. 覆盖率报告格式

覆盖率报告会显示每个包的覆盖率百分比：

```
ok      WaterMark/pkg   0.511s  coverage: 90.4% of statements
```

### 3. 生成覆盖率报告（详细模式）

```bash
# 生成覆盖率报告并显示详细输出
go test -cover -v ./pkg/...
```

## 查看覆盖率报告

### 1. 查看函数级覆盖率

```bash
# 显示每个函数的覆盖率
go tool cover -func=coverage.out
```

输出示例：

```
WaterMark/pkg/any.go:11:                AnyToString                     93.3%
WaterMark/pkg/compress.go:10:           ZlibCompress                    77.8%
WaterMark/pkg/image.go:15:              LoadImageWithDecode             87.0%
total:                                  (statements)                    90.4%
```

### 2. 生成 HTML 覆盖率报告

```bash
# 在浏览器中打开 HTML 覆盖率报告
go tool cover -html=coverage.out

# 生成 HTML 覆盖率报告到文件
go tool cover -html=coverage.out -o coverage.html
```

**注意：在 PowerShell 中使用 `go tool cover` 命令时，需要用引号包裹参数：**

```bash
# PowerShell 中正确的命令格式
go tool cover '-html=coverage.out'

# 或者使用绝对路径
go tool cover '-html=E:/WaterMark/coverage.out'

# 生成 HTML 文件到指定位置
go tool cover '-html=E:/WaterMark/coverage.out' -o E:/WaterMark/coverage.html
```

**原因：** PowerShell 会将 `-html=coverage.out` 解析为多个参数，导致 "too many arguments" 错误。使用引号包裹参数可以避免这个问题。

### 3. 查看覆盖率统计

```bash
# 显示覆盖率统计信息
go tool cover -func=coverage.out | grep total
```

## 测试命令详解

### run_tests.go 参数

测试入口程序 `tests/run_tests.go` 支持以下参数：

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-v` | 详细输出模式 | false |
| `-cover` | 启用覆盖率 | false |
| `-coverprofile` | 覆盖率报告文件路径 | "" |
| `-race` | 启用竞态检测 | false |
| `-pkgs` | 要测试的包路径 | "./..." |

### 使用示例

```bash
# 详细输出模式
go run tests/run_tests.go -v

# 启用覆盖率
go run tests/run_tests.go -cover

# 生成覆盖率报告
go run tests/run_tests.go -cover -coverprofile=coverage.out

# 启用竞态检测
go run tests/run_tests.go -race

# 测试特定包
go run tests/run_tests.go -pkgs=./pkg/...
```

### Go test 参数

| 参数 | 说明 |
|------|------|
| `-v` | 详细输出模式，显示每个测试的名称和结果 |
| `-cover` | 启用覆盖率分析 |
| `-coverprofile` | 指定覆盖率报告文件路径 |
| `-covermode` | 覆盖率模式：set（默认）、count、atomic |
| `-race` | 启用竞态检测器 |
| `-run` | 运行匹配正则表达式的测试 |
| `-timeout` | 设置测试超时时间 |

## 覆盖率目标

### 当前覆盖率状态

| 包名 | 覆盖率 | 状态 |
|------|--------|------|
| pkg | 90.4% | ✅ 优秀 |
| internal | 13.3% | ⚠️ 需要改进 |
| scripts | 78.2% | ✅ 良好 |

### 覆盖率目标

- **优秀**: ≥ 90%
- **良好**: 80% - 89%
- **需要改进**: < 80%

### 提高覆盖率的建议

1. **为覆盖率较低的函数添加测试用例**
   - [getLine](../pkg/csv.go#L28) - 66.7%
   - [ZlibCompress](../pkg/compress.go#L10) - 77.8%
   - [Generate](../pkg/csv.go#L47) - 77.3%

2. **添加错误处理的测试用例**
   - 文件读取错误
   - 网络请求错误
   - 无效输入处理

3. **添加边界情况的测试用例**
   - 空值处理
   - 极大值处理
   - 特殊字符处理

## 常见问题

### 1. 测试失败怎么办？

```bash
# 运行测试并显示详细输出
go test -v ./...

# 查看具体哪个测试失败
go test -v -run TestFunctionName ./...
```

### 2. 如何调试测试？

```bash
# 使用 dlv 调试器
dlv test ./pkg/...

# 在测试中添加日志
t.Logf("Debug information: %v", variable)
```

### 3. 覆盖率报告文件在哪里？

覆盖率报告文件默认保存在项目根目录，文件名为 `coverage.out`。可以使用 `-coverprofile` 参数指定其他路径。

**注意：** 当使用 `go test ./...` 测试多个包时，Go 会忽略 `-coverprofile` 参数中的扩展名，直接使用基础文件名（例如 `coverage` 而不是 `coverage.out`）。`run_tests.bat` 脚本会自动将生成的 `coverage` 文件重命名为 `coverage.out`。

### 4. 为什么生成的覆盖率文件名是 `coverage` 而不是 `coverage.out`？

这是 Go 的一个已知行为。当测试多个包时（使用 `./...`），`go test` 命令会忽略 `-coverprofile` 参数中的文件扩展名，只使用基础文件名。

**解决方案：**
- 使用 `run_tests.bat` 脚本，它会自动重命名文件
- 或者手动重命名：`ren coverage coverage.out`
- 或者只测试单个包：`go test -coverprofile=coverage.out ./pkg/...`

### 5. 如何只运行失败的测试？

```bash
# 记录失败的测试
go test -v ./... > test_output.txt

# 手动运行失败的测试
go test -v -run TestFailedFunctionName ./...
```

### 6. 如何测试并发代码？

```bash
# 启用竞态检测器
go test -race ./...

# 设置并发数
go test -parallel=4 ./...
```

## 最佳实践

1. **编写测试时遵循 AAA 模式**
   - Arrange（准备）
   - Act（执行）
   - Assert（断言）

2. **使用表驱动测试**
   - 便于添加新的测试用例
   - 提高代码可读性

3. **测试命名规范**
   - 使用 `Test` 前缀
   - 描述性命名：`TestFunctionName_Scenario_ExpectedResult`

4. **保持测试独立性**
   - 每个测试应该独立运行
   - 不依赖其他测试的执行顺序

5. **使用临时目录**
   - 使用 `t.TempDir()` 创建临时目录
   - 避免污染项目文件系统

## 相关资源

- [Go Testing 官方文档](https://golang.org/pkg/testing/)
- [Go Coverage 工具](https://golang.org/pkg/cmd/go/internal/test/)
- [Table-Driven Tests in Go](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)

## 更新日志

- 2025-12-29: 初始版本，添加测试执行和覆盖率生成指南
