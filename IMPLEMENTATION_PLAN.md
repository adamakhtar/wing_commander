# Wing Commander V1 - Updated Implementation Plan

## Project Overview

A CLI/TUI tool for analyzing test failures by grouping them by backtrace similarity. Helps developers quickly identify shared root causes among multiple failing tests.

## Current Status: Steps 1-5 Complete ✅

### ✅ **Step 1: Project Foundation + Core Types** (COMPLETED)

- Go module initialized with dependencies
- Core domain types defined (`StackFrame`, `TestResult`, `FailureGroup`)
- Comprehensive unit tests
- Basic CLI with welcome message
- Build system (Makefile) configured
- Clean project structure

### ✅ **Step 2: JSON Parser** (COMPLETED)

- Parser package with JSON schema support
- RSpec and Minitest JSON format support
- Backtrace frame parsing (file:line:method)
- Comprehensive test coverage with fixtures
- CLI integration for JSON file parsing
- Framework detection removed (user specifies in config)

### ✅ **Step 3: Configuration System** (COMPLETED)

- YAML-based configuration system
- Support for multiple test frameworks (RSpec, Minitest, Pytest, Jest)
- User-configurable exclude patterns
- Default test commands per framework
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

## Remaining Implementation Steps

### 🔄 **Step 6: Git Integration** (NEXT)

**Goal**: Identify recently changed files

**Files to create**:

- `internal/git/git.go`: Execute git diff, parse changed files
- `internal/git/git_test.go`: Tests with mock exec

**CLI Update**: Mark groups/frames that touch changed files in text output

**Checkpoint**: `./wing_commander <json>` shows `[*]` next to recently changed frames

---

### 🔄 **Step 7: Test Runner**

**Goal**: Execute tests and capture JSON output

**Files to create**:

- `internal/runner/runner.go`: Execute test command from config, capture stdout
- `internal/runner/runner_test.go`: Tests with mock commands

**CLI Update**: Run tests directly (no JSON file needed), parse output

**Checkpoint**: `./wing_commander run` executes tests and displays grouped failures

---

### 🔄 **Step 8: Basic Bubbletea UI - Single Pane**

**Goal**: Replace text output with interactive TUI (groups list only)

**Files to create**:

- `internal/ui/app.go`: Bubbletea Init/Update/View
- `internal/ui/models.go`: UI state (selected group index)
- `internal/ui/styles.go`: Lipgloss styles
- Add bubbletea and lipgloss dependencies

**CLI Update**: Launch TUI showing groups list, navigate with arrows, `q` to quit

**Checkpoint**: `./wing_commander run` shows interactive groups list

---

### 🔄 **Step 9: Multi-Pane UI**

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
- Highlight recently changed files

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

### **Grouping Strategy**

- Group by ErrorLocation strategy (bottom frame only: file:line)
- Store full 50 frames for user viewing
- Use strategy pattern for future extensibility
- Groups sorted by count (most frequent failures first)

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
│   └── grouper/        # Grouping logic
│       ├── normalizer.go
│       ├── normalizer_test.go
│       ├── strategy.go
│       ├── strategy_test.go
│       ├── grouper.go
│       └── grouper_test.go
├── testdata/
│   ├── fixtures/       # Test JSON files
│   └── config/         # Sample configs
├── Makefile
├── go.mod
├── go.sum
├── .gitignore
└── CONTEXT.md
```
