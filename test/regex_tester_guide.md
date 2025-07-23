# 正则测试器指南 (`regex_tester.go`)

本文档旨在说明如何使用 `regex_tester.go` 脚本，以便在将其添加到主程序 URLFinder 之前，能够安全地测试和验证正则表达式。

## 1. 目的

该脚本的主要目的是为开发者提供一个沙盒环境，用于：
- 针对样本文本测试新的正则表达式。
- 验证对现有正则表达式的更改是否能如期工作。
- 无需运行完整的 URLFinder 应用程序，即可确保正则表达式能够正确捕获预期的数据。

## 2. 如何使用

您可以在终端中以多种模式运行此脚本。

### 模式一：测试所有已定义的正则

要对脚本中 `Infofind` map 内定义的所有正则分类进行全面测试，只需运行不带任何参数的命令：

```bash
go run regex_tester.go
```

### 模式二：测试特定的正则

如果您只想测试一个或多个特定的分类，可以将其名称作为命令行参数传递。这对于集中调试非常有用。

脚本支持使用**空格**、**逗号**或**两者混合**的方式来分隔参数。

**示例：**

- **单个分类：**
  ```bash
  go run regex_tester.go AKSK
  ```

- **多个分类（使用空格分隔）：**
  ```bash
  go run regex_tester.go AKSK Phone Jdbc
  ```

- **多个分类（使用逗号分隔）：**
  ```bash
  go run regex_tester.go AKSK,Phone,Jdbc
  ```

- **多个分类（混合分隔）：**
  ```bash
  go run regex_tester.go AKSK,Phone Jdbc
  ```

## 3. 如何修改以进行测试

您可以轻松地自定义脚本以满足您的特定测试需求。

### a. 修改测试文本

在 `main` 函数中找到 `text` 变量。您可以将这个多行字符串的内容替换为您想用来测试正则的任何数据。

```go
// 1. 在这里放入你的测试文本
text := `
      // 在这里放入你的测试数据...
      some_key: your_value
`
```

### b. 添加或修改正则

脚本中包含一个 `Infofind` map，其结构与 `config/config.go` 中的结构完全一致。您可以在此处添加新分类或编辑现有分类中的正则表达式以进行测试。

```go
// 2. 在这里定义所有可能用到的正则表达式
	var (
		AKSK      = []string{`在这里放入你新的aksk正则`}
		test2    = []string{`在这里放入你新的aksk正则`}
	)

	// 将独立的正则变量组装成 map
	Infofind := map[string][]string{
		"AKSK":     AKSK,
		"test2":     test2,
	}

```

通过遵循这些步骤，您可以高效、安全地验证您的正则表达式。