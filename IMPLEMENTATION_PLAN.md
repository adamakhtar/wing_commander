# Wing Commander V1 - Updated Implementation Plan

## Project Overview

A CLI/TUI tool for analyzing test failures by grouping them by backtrace similarity. Helps developers quickly identify shared root causes among multiple failing tests.

## Current Status: Steps 1-7 Complete ✅

### ✅ **Step 1: Project Foundation + Core Types** (COMPLETED)

- Go module initialized with dependencies
- Core domain types defined (`StackFrame`, `TestResult`, `FailureGroup`)
- Comprehensive unit tests
- Basic CLI with welcome message
- Build system (Makefile) configured
- Clean project structure

### ✅ **Step 2: JUnit XML Parser** (COMPLETED)

- Parser package with JUnit XML schema support
- RSpec and Minitest JUnit XML format support
- Backtrace frame parsing (file:line:method)
- Comprehensive test coverage with fixtures
- CLI integration for XML file parsing
- Framework detection removed (user specifies in config)

### ✅ **Step 3: Configuration System - CLI-First Design** (COMPLETED)

- YAML-based configuration system with CLI-first approach
- **Priority System**: CLI options > Config file > Sensible defaults
- **CLI Options**: `--project-path`, `--test-command`, `--config` flags
- **Template Interpolation**: Go `text/template` syntax for test commands
- Support for multiple test frameworks (RSpec, Minitest, Pytest, Jest)
- User-configurable exclude patterns
- CLI config command
- Clean file organization (test files in `testdata/`)

### ✅ **Step 4: Backtrace Normalizer** (COMPLETED)

- Filter frames using config exclude patterns
- Normalize test results with filtered backtraces
- Comprehensive test coverage
- CLI integration for frame filtering statistics

### ✅ **Step 5: Failure Grouper** (COMPLETED)

- Strategy pattern implementation for grouping
- ErrorLocationStrategy groups by bottom frame (file:line)
- Groups sorted by count (most frequent first)
- Comprehensive test coverage
- Ready for CLI integration

### ✅ **Step 6: Git Integration** (COMPLETED)

- Line-level change detection with 3 intensity levels
- Uncommitted changes (intensity 3), last commit (intensity 2), previous commit (intensity 1)
- Unified diff parsing for precise line number detection
- Integration with grouper workflow
- Comprehensive test coverage

### ✅ **Step 7: Test Runner Service** (COMPLETED)

- TestRunner service for GUI-driven test execution
- Execute test commands from config and parse JUnit XML output
- Complete workflow integration (parse → normalize → group → detect changes)
- CLI `run` command implementation with `--config` flag support
- Config file path customization via command line flags
- Comprehensive test coverage
- Ready for GUI integration

### ✅ **Dummy Projects Setup** (COMPLETED)

- **Minitest Dummy Project**: Complete Ruby project for testing minitest support
  - `dummy/minitest/lib/thing.rb`: Simple class with `boom` method that raises error
  - `dummy/minitest/test/thing_test.rb`: Two failing test cases
  - `dummy/minitest/Gemfile`: Dependencies (minitest, ci_reporter_minitest)
  - `dummy/minitest/test/test_helper.rb`: JUnit XML reporting configuration
  - `dummy/minitest/Rakefile`: Test execution with XML output
  - Generates JUnit XML reports in `test/reports/` directory
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
wing_commander run --config custom-config.yml --test-command "pytest {{.Paths}} --junit-xml=results.xml"

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

**CLI Update**: Added `demo` command to launch TUI with XML fixture data

**Features implemented**:

- 3-pane layout (Groups | Tests | Backtrace)
- Arrow key navigation within panes
- Tab/Shift+Tab to switch between panes
- 'q' to quit
- Real data processing through existing pipeline

**Checkpoint**: `./bin/wing_commander demo` shows interactive 3-pane TUI

---

### 🔄 **Step 9: Advanced UI Features** (NEXT)

**Goal**: Add tests pane and backtrace pane

**Files to create**:

- `internal/ui/views.go`: 3-pane layout rendering
- Update `models.go`: Track active pane, selections per pane

**UI Updates**: `Tab` switches panes, each pane navigable

**Checkpoint**: Full 3-pane navigation working

---

### 🔄 **Step 10: Advanced UI Features**

**Goal**: Add keybindings for actions (toggle, open file, re-run)

**Files to create**:

- Update `app.go`: Handle `f`, `o`, `r` keybindings
- `internal/editor/editor.go`: Open file in $EDITOR at line

**UI Updates**:

- `f`: Toggle full/filtered frames
- `o`: Open selected file in editor
- `r`: Re-run tests in selected group
- Highlight recently changed files with intensity levels

**Checkpoint**: All keybindings functional

---

### 🔄 **Step 11: Polish & Documentation**

**Goal**: Production-ready V1

**Files to create**:

- `README.md`: Installation, usage, configuration guide
- Example `.wing_commander/config.yml` in docs
- Error messages polish
- Help screen in UI

**Checkpoint**: Ready for release

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

### **Grouping Strategy**

- Group by ErrorLocation strategy (bottom frame only: file:line)
- Store full 50 frames for user viewing
- Use strategy pattern for future extensibility
- Groups sorted by count (most frequent failures first)
- Line-level change detection with 3 intensity levels

## Success Criteria

- Groups 100-1000 test failures efficiently (<1s)
- Responsive TUI navigation
- Recently changed files highlighted correctly
- Can open files in editor at correct line
- Can re-run specific test groups
- Simple workflow: run tests → see grouped failures

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
│   ├── grouper/        # Grouping logic
│   │   ├── normalizer.go
│   │   ├── normalizer_test.go
│   │   ├── strategy.go
│   │   ├── strategy_test.go
│   │   ├── grouper.go
│   │   └── grouper_test.go
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
