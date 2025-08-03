<h1 align="center">Welcome to Hulo 👋</h1>
<center>

[![Hulo](https://img.shields.io/badge/Hulo-%238866E9.svg?logoColor=white&style=for-the-badge)](https://github.com/hulo-lang/hulo) [![Go](https://img.shields.io/badge/Go-1.24.4-%2300ADD8.svg?logo=go&logoColor=white&style=for-the-badge)](https://golang.org/) [![BashScript](https://img.shields.io/badge/Bash%20Script-%23121011.svg?logo=gnu-bash&logoColor=white&style=for-the-badge)](https://www.gnu.org/software/bash/) [![PowerShell](https://img.shields.io/badge/PowerShell-%235391FE.svg?logo=powershell&logoColor=white&style=for-the-badge)](https://learn.microsoft.com/en-us/powershell/) [![VBScript](https://img.shields.io/badge/VBScript-%234A4A4A.svg?logo=windows&logoColor=white&style=for-the-badge)](https://documentation.help/MS-Office-VBScript/VBSTOC.htm)

</center>

---

English | [简体中文](README.zh-CN.md)

> Hulo /ˈhjuːloʊ/ is a modern, batch-oriented programming language that compiles to Bash, PowerShell, and VBScript. It is designed to unify scripting across platforms with a clean and consistent DSL.

## 📦 Install

### **Direct Downloads**

Download pre-built binaries directly from [GitHub Releases](https://github.com/hulo-lang/hulo/releases):

```sh
# Linux/macOS
curl -L https://github.com/hulo-lang/hulo/releases/latest/download/install.sh | bash

# Windows (PowerShell)
irm https://github.com/hulo-lang/hulo/releases/latest/download/install.ps1 | iex
```

### **From Source**
```sh
# Clone repository
git clone https://github.com/hulo-lang/hulo.git
cd hulo

# for Windows
tools/scripts/setup.ps1

# for linux
tools/scripts/setup.sh

mage r
```

### **Package Managers**

| Package Manager | HomePage | Repository |
|----------------|------------|--------|
| **npm** | [hulo-lang](https://www.npmjs.com/package/hulo-lang) | [hulo-npm](https://github.com/hulo-lang/hulo-npm) |
| **pypi** | [hulo](https://pypi.org/project/hulo) | [hulo-py](https://github.com/hulo-lang/hulo-py) |
| **scoop** |  | [scoop-hulo](https://github.com/hulo-lang/scoop-hulo) |
| **brew** |  | [homebrew-hulo](https://github.com/hulo-lang/homebrew-hulo) |


## 🚀 Usage

```hulo
// hello.hl

echo "Hello, World!"
```

Run `hulo hello.hl`, and it will compile into:
* `hello.sh` for Unix-like systems
* `hello.ps1` for Windows PowerShell
* `hello.bat` for Windows Batch
* `hello.vbs` for Windows VBScript
* And more targets in future releases!

## 📖 Documentation

- **[Official Docs](https://hulo-lang.github.io/docs)** - Complete language reference
- **[Examples](./examples/)** - Code examples
- **[Discussions](https://github.com/hulo-lang/hulo/discussions)** - Ask questions and share ideas
- **[Hulo Dev](https://github.com/hulo-lang/hulo-dev)** - Learn compiler construction by building Hulo language step by step - from lexer to code generation

## 🤝 Contributing

All contributions are welcome, including:

- Reporting bugs
- Submitting issues or feature requests
- Translating or rewriting scripts
- Sending pull requests

Please see our [CONTRIBUTING.md](CONTRIBUTING.md) before submitting changes.

## 📝 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for more details.

---

_Made with ❤️ for hackers, ops, and anyone who loves clean automation._
