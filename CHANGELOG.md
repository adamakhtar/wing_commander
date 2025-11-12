# Changelog

All notable changes to Wing Commander will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Failure cause classification (test definition error, production code error, assertion failure) with simple parser heuristics
- UI grouping by failure cause with section headers and Unicode dividers
- Failure cause icons: 🚀 (production), 🔧 (test definition), ❌ (assertion)
- Yellow styling for counts and error messages in groups pane
- Advanced UI keybindings (`f`, `o`, `r`)
- Editor integration for opening files at specific lines
- Test re-running functionality
- Development workflow with `make dev-minitest`
- Comprehensive test suite
- CLI-first configuration system

### Changed

- Improved backtrace parsing accuracy
- Enhanced test grouping by error location
- Updated UI with better navigation and status display
- TUI layout updated: Panel 1 shows error + bottom frame, Panel 2 shows test + tail frames, Panel 3 shows full test backtrace with highlighting
- Parser now exclusively supports WingCommanderReporter YAML summaries (removed JUnit XML handling)

### Fixed

- Corrected test grouping to use first frame (error origin) instead of last frame
- Fixed parser to only capture properly indented stack frames
- Prevented parsing of error message lines as stack frames

## [0.1.0] - 2025-10-19

### Added

- Initial release
- Core test failure analysis functionality
- WingCommanderReporter YAML parsing support
- Interactive TUI with Bubbletea
- Support for RSpec, Minitest, Pytest, Jest
- Backtrace filtering and normalization
- Git change detection with intensity levels
- Error location grouping strategy
- CLI interface with configuration support
- Build system with Makefile
- Comprehensive documentation

### Features

- **Test Parsing**: Parse WingCommanderReporter YAML summaries
- **Failure Grouping**: Group tests by backtrace similarity using error location strategy
- **Interactive TUI**: Navigate through failure groups, tests, and backtraces
- **Git Integration**: Highlight recently changed files with 3 intensity levels
- **Configuration**: CLI-first approach with config file fallback
- **Cross-Platform**: Support for multiple operating systems and editors

### Supported Test Frameworks

- Minitest (Ruby via WingCommanderReporter)

### Supported Editors

- VS Code (`code`)
- Sublime Text (`subl`)
- Atom (`atom`)
- Vim/Neovim (`vim`, `nvim`)
- Emacs (`emacs`)
- Any editor set in `$EDITOR` or `$VISUAL` environment variables
