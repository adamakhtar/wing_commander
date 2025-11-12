# Wing Commander V1 - Updated Implementation Plan

## Project Overview

A CLI/TUI tool for running tests and reviewing failure details with normalized backtraces. Helps developers quickly inspect failing tests without the noise of third-party frames.

## Current Status: Steps 1-10 Complete ✅

### ✅ **Step 1: Project Foundation + Core Types** (COMPLETED)

- Go module initialized with dependencies
- Core domain types defined (`StackFrame`, `TestResult`)
- Comprehensive unit tests
- Basic CLI with welcome message
- Build system (Makefile) configured
- Clean project structure

### ✅ **Step 5: Failure List Presentation** (COMPLETED)

- Present flattened failure list without grouping logic
- Ensure stable ordering and data structures based solely on `TestResult`
- Comprehensive test coverage
- Ready for CLI/TUI integration

### ✅ **Step 6: Git Integration** (COMPLETED)

- Line-level change detection with 3 intensity levels
- Uncommitted changes (intensity 3), last commit (intensity 2), previous commit (intensity 1)
- Unified diff parsing for precise line number detection
- Integrated with normalized failure workflow

### ✅ **Step 7: Test Runner Service** (COMPLETED)

- TestRunner service for GUI-driven test execution
- Execute test commands from config and parse WingCommanderReporter YAML summaries
- Complete workflow integration (parse → normalize → detect changes)
- CLI `run` command implementation with `--config` flag support
- Config file path customization via command line flags
- Comprehensive test coverage
- Ready for GUI integration

### ✅ **Dummy Projects Setup** (COMPLETED)

- **Minitest Dummy Project**: Complete Ruby project for testing minitest support
  - `dummy/minitest/lib/thing.rb`: Simple class with `boom` method that raises error
  - `dummy/minitest/test/thing_test.rb`: Two failing test cases
  - `dummy/minitest/Gemfile`: Dependencies (minitest, ci_reporter_minitest)
- `dummy/minitest/test/test_helper.rb`: WingCommanderReporter configuration
- `dummy/minitest/Rakefile`: Test execution producing YAML summary
- Generates `.wing_commander/test_results/summary.yml`
  - Runnable with `bundle install` and `bundle exec rake test`

## ✅ **CLI-First Configuration System** (NEWLY COMPLETED)

### **Configuration Priority System**

1. **CLI Options** (Highest Priority)

   - `--project-path PATH`: Project directory for test execution
   - `--test-command CMD`: Test runner command with template interpolation
   - `--config PATH`: Custom config file location

2. **Config File** (Medium Priority)

   - `.wing_commander/config.yml` (default location)
   - YAML format with `project_path`, `test_framework`, `test_command`, `exclude_patterns`

3. **Sensible Defaults** (Lowest Priority)
   - Project path: Current working directory
   - Test framework: Minitest (supported framework)
   - Test command: Must be specified (no hardcoded default)

### **Template Interpolation System**

- **Engine**: Go's `text/template` package
- **Syntax**: `{{.Paths}}` for test paths (empty by default)
- **Future Ready**: Extensible for specific test file selection
- **Example**: `rails test {{.Paths}} --output junit`

### **Usage Examples**

```bash
# CLI-first approach - override everything via command line
wing_commander run --project-path /path/to/project --test-command "rails test {{.Paths}} --output junit"

# Mix CLI and config - project path from CLI, test command from config
wing_commander run --project-path /path/to/project

# Custom config file with CLI overrides
wing_commander start /path/to/project --config custom-config.yml --run-command "bundle exec rake test {{.Paths}}"

# Traditional approach - everything from config file
wing_commander run
```

## Remaining Implementation Steps

### ✅ **Step 8: Basic Bubbletea UI - Multi-Pane** (COMPLETED)

**Goal**: Replace text output with interactive TUI (3-pane layout)

**Files implemented**:

- `internal/ui/model.go`: Bubbletea Init/Update/View with 3-pane layout
- `internal/ui/styles.go`: Lipgloss styles for panes and text
- Added bubbletea and lipgloss dependencies

**CLI Update**: Added TUI with XML fixture data

**Features implemented**:

- 3-pane layout (Test Runs | Tests | Backtrace)
- Arrow key navigation within panes
- Tab/Shift+Tab to switch between panes
- 'q' to quit
- Real data processing through existing pipeline

**Checkpoint**: `./bin/wing_commander run` shows interactive 3-pane TUI

---

### ✅ **Step 9: Advanced UI Features** (COMPLETED)

**Goal**: Add keybindings for actions (toggle frames, open file, re-run tests)

**Files implemented**:

- `internal/editor/editor.go`: File opening functionality with editor detection
- `internal/editor/editor_test.go`: Comprehensive tests for editor functionality
- Updated `internal/ui/model.go`: Enhanced UI model with new keybindings
- Updated `cmd/wing_commander/main.go`: Pass TestRunner to UI model
- Updated `Makefile`: Added `dev-minitest` command for development testing

**UI Features implemented**:

- `f`: Toggle full/filtered frames display
- `o`: Open selected file in external editor at specific line
- `r`: Re-run tests for selected entry
- Async message handling for file opening and test re-running
- Updated status bar with all available keybindings

**Development workflow**:

- `make dev-minitest`: Build dev version and launch TUI against dummy minitest app
- Real test execution with WingCommanderReporter YAML summaries
- Interactive TUI with actual test failures

**Checkpoint**: Complete interactive TUI with file opening and test re-running capabilities

---

### ✅ **Step 10: Advanced UI Features** (COMPLETED)

**Goal**: Add keybindings for actions (toggle, open file, re-run)

**Files implemented**:

- `internal/editor/editor.go`: File opening functionality with editor detection
- `internal/editor/editor_test.go`: Comprehensive tests for editor functionality
- Updated `internal/ui/model.go`: Enhanced UI model with new keybindings
- Updated `cmd/wing_commander/main.go`: Pass TestRunner to UI model
- Updated `Makefile`: Added `dev-minitest` command for development testing

**UI Features implemented**:

- `f`: Toggle full/filtered frames display
- `o`: Open selected file in external editor at specific line
- `r`: Re-run tests for selected entry
- Async message handling for file opening and test re-running
- Updated status bar with all available keybindings

**Development workflow**:

- `make dev-minitest`: Build dev version and launch TUI against dummy minitest app
- Real test execution with WingCommanderReporter YAML summaries
- Interactive TUI with actual test failures

**Checkpoint**: Complete interactive TUI with file opening and test re-running capabilities

---

### 🔄 **Step 11: Polish & Documentation** (NEXT)

**Goal**: Production-ready V1

**Files to create**:

- `README.md`: Installation, usage, configuration guide
- `CHANGELOG.md`: Version history and features
- `LICENSE`: Open source license

**Polish items**:

- Error handling improvements
- Performance optimizations
- Cross-platform testing
- Final documentation review

**Checkpoint**: Production-ready V1 release

---

## Key Design Decisions Made

### **Simplified V1 Approach**

- **No Caching**: Fresh test run every time (keeps it simple)
- **User-Specified Framework**: No auto-detection (more reliable)
- **Clean File Organization**: Test files in `testdata/`, user configs ignored
- **Makefile Only**: Removed redundant dev.sh script

### **Configuration Format**

```yaml
test_framework: rspec
test_command: "bundle exec rspec --format json"
exclude_patterns:
  - "/gems/"
  - "/lib/ruby/"
  - "/vendor/bundle/"
```

### **CLI Configuration Support - CLI-First Design**

- **Priority System**: CLI options > Config file > Sensible defaults
- **CLI Options**:
  - `--project-path PATH`: Project directory (default: current working directory)
  - `--test-command CMD`: Test runner with interpolation (e.g., `rails test {{.Paths}} --output junit`)
  - `--config PATH`: Config file path (default: `.wing_commander/config.yml`)
- **Template Interpolation**: Uses Go `text/template` syntax for command interpolation
- **Backward Compatibility**: Existing usage without flags continues to work
- **Help Documentation**: `wing_commander run --help` shows all flag usage

### **Failure Presentation**

- Display failures as a flat list ordered by recent execution
- Store full 50 frames for user viewing while showing filtered project frames in the UI
- Keep the codebase ready for future enhancements without relying on grouping strategies
- Line-level change detection with 3 intensity levels highlights relevant frames

## Success Criteria

- Displays 100-1000 test failures efficiently (<1s)
- Responsive TUI navigation
- Recently changed files highlighted correctly
- Can open files in editor at correct line
- Can re-run selected tests
- Simple workflow: run tests → review failures

## Development Workflow

```bash
# Build and test
make dev
make test

# Run CLI
make run

# Run with CLI options (CLI-first approach)
./bin/wing_commander run --project-path /path/to/project --test-command "rails test {{.Paths}} --output junit"

# Run with custom config
./bin/wing_commander run --config /path/to/config.yml

# Clean up
make clean
```

## Current File Structure

```
wing_commander/
├── cmd/wing_commander/main.go
├── internal/
│   ├── types/          # Core domain types
│   ├── parser/         # JSON parsing
│   ├── config/         # Configuration system
│   ├── backtrace/      # Backtrace normalization helpers
│   │   ├── normalizer.go
│   │   └── normalizer_test.go
│   ├── git/            # Git change detection
│   │   ├── changes.go
│   │   └── changes_test.go
│   └── runner/         # Test execution service
│       ├── runner.go
│       └── runner_test.go
├── testdata/
│   ├── fixtures/       # Test JSON files
│   └── config/         # Sample configs
├── dummy/
│   └── minitest/       # Test framework dummy projects
│       ├── lib/
│       │   └── thing.rb
│       ├── test/
│       │   ├── test_helper.rb
│       │   └── thing_test.rb
│       ├── Gemfile
│       └── Rakefile
├── Makefile
├── go.mod
├── go.sum
├── .gitignore
└── CONTEXT.md
```
