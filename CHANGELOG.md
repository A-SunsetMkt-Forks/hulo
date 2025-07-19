# 📜 Changelog

## [v0.2.0]() - 2025-07-19

### Added

- Bash transpiler support
- HLPM (Hulo Package Manager) commandline tool for managing Hulo packages
  - `init`, `install`, `uninstall`, `run` commands implemented
  - Package publishing and import resolution pending
- Hulo-REPL commandline tool with simple completion and theme support
  - Basic completion and theme settings implemented
  - Full integration pending

### Fixed

- Fixed string translation issue with `echo "Hello World"` [#4](https://github.com/hulo-lang/hulo/issues/4)

### Breaking Changes

- Refactored VBScript transpiler - this is a breaking change, output may not be consistent with v0.1.0
- Refactored Hulo commandline tool to read `huloc.yaml` configuration file from working directory

### Migration

- Environmental variable `HULOPATH` renamed to `HULO_PATH`

## [v0.1.0](https://github.com/hulo-lang/hulo/commit/c119ee825f518af1cb2a6fbe8012ed265970cbd2) - 2025-07-13

### Added

- First Release, only support VBScript

---

_Made with ❤️ for hackers, ops, and anyone who loves clean automation._
