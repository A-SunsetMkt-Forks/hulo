<h1 align="center">欢迎使用 Hulo 👋</h1>
<center>

[![Hulo](https://img.shields.io/badge/Hulo-%238866E9.svg?logoColor=white&style=for-the-badge)](https://github.com/hulo-lang/hulo) [![Go](https://img.shields.io/badge/Go-1.24.4-%2300ADD8.svg?logo=go&logoColor=white&style=for-the-badge)](https://golang.org/) [![BashScript](https://img.shields.io/badge/Bash%20Script-%23121011.svg?logo=gnu-bash&logoColor=white&style=for-the-badge)](https://www.gnu.org/software/bash/) [![PowerShell](https://img.shields.io/badge/PowerShell-%235391FE.svg?logo=powershell&logoColor=white&style=for-the-badge)](https://learn.microsoft.com/en-us/powershell/) [![VBScript](https://img.shields.io/badge/VBScript-%234A4A4A.svg?logo=windows&logoColor=white&style=for-the-badge)](https://documentation.help/MS-Office-VBScript/VBSTOC.htm)

</center>

---

[English](README.md) | 简体中文

> Hulo 是一个现代化的、面向批处理的编程语言，可以编译为 Bash、PowerShell 和 VBScript。它旨在通过简洁一致的 DSL 来统一跨平台的脚本编写。

## 📦 安装

### **直接下载**

从 [GitHub Releases](https://github.com/hulo-lang/hulo/releases) 直接下载预构建的二进制文件：

```sh
# Linux/macOS
curl -L https://github.com/hulo-lang/hulo/releases/latest/download/install.sh | bash

# Windows (PowerShell)
irm https://github.com/hulo-lang/hulo/releases/latest/download/install.ps1 | iex
```

### **从源码构建**
```sh
# 克隆仓库
git clone https://github.com/hulo-lang/hulo.git
cd hulo

# Windows 用户
tools/scripts/setup.ps1

# Linux 用户
tools/scripts/setup.sh

# 构建
mage release:all
```

### **包管理器**

| 包管理器 | 主页 | 仓库 |
|---------|------|------|
| **npm** | [hulo-lang](https://www.npmjs.com/package/hulo-lang) | [hulo-npm](https://github.com/hulo-lang/hulo-npm) |
| **pypi** | [hulo](https://pypi.org/project/hulo) | [hulo-py](https://github.com/hulo-lang/hulo-py) |
| **scoop** |  | [scoop-hulo](https://github.com/hulo-lang/scoop-hulo) |
| **brew** |  | [homebrew-hulo](https://github.com/hulo-lang/homebrew-hulo) |


## 🚀 使用示例

```hulo
// hello.hl

echo "Hello, World!"
```

运行 `hulo hello.hl`，它将编译为：
* Unix 系统的 `hello.sh`
* Windows 的 `hello.ps1`
* 如果启用了 VBScript 输出，则为 `hello.vbs`
* 未来版本中将支持更多目标平台！

## 📖 文档

- **[官方文档](https://hulo-lang.github.io/docs)** - 完整的语言参考
- **[示例](./examples/)** - 代码示例
- **[讨论](https://github.com/hulo-lang/hulo/discussions)** - 提问和分享想法

## 🤝 贡献方式

欢迎任何形式的贡献，包括但不限于：

- 报告 Bug
- 提交问题或功能请求
- 翻译或优化脚本
- 提交 Pull Request

在提交修改之前，请先阅读我们的 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 📝 许可证

本项目基于 MIT 协议开源。详见 [LICENSE](LICENSE)。

---

_为热爱自动化的黑客、运维和开发者而生 ❤️_
