# Wing Commander - Project Context & Decisions

## Project Overview

A CLI/TUI tool for analyzing test failures by grouping them by backtrace similarity. Helps developers quickly identify shared root causes among multiple failing tests.

## Key Decisions Made

### 1. **Grouping Strategy**

- **Group by**: ErrorLocation strategy - bottom frame only (file + line number)
- **Store**: Full 50 frames for user viewing
- **Rationale**: Groups failures by where the error actually surfaced, making it easy to identify root causes
- **Language Support**: Works across Ruby, Python, JavaScript, Go (file:line format)
- **Architecture**: Strategy pattern allows easy addition of new grouping strategies (CallPath, ErrorPattern, etc.)
- **Change Detection**: Line-level change detection with 3 intensity levels:
  - **Intensity 3**: Uncommitted changes to frame's line number
  - **Intensity 2**: Line changed in last commit
  - **Intensity 1**: Line changed in commit before last

### 2. **Project Structure**

- **Types Location**: `internal/types/` (not `pkg/`) - internal use only
- **Config Location**: `.wing_commander/config.yml` (not root)
- **Test Files**: `testdata/` directory for fixtures and sample configs
- **Rationale**: Keep project root clean, signal internal vs public packages

### 3. **Configuration Approach**

- **V1**: User creates config file manually
- **Future**: Tool can generate config file on request
- **Format**: YAML with test_framework, test_command, and exclude_patterns
- **Framework Support**: RSpec, Minitest, Pytest, Jest
- **Default Patterns**: `/gems/`, `/lib/ruby/`, `/vendor/bundle/`, etc.

### 4. **Build System**

- **Primary**: Makefile (standard, cross-platform)
- **Removed**: dev.sh (redundant wrapper)
- **Build Locations**: `bin/` (dev), `dist/` (production)
- **Commands**: `make dev`, `make test`, `make run`, `make clean`
- **Clean Development**: No development files cluttering production code

### 5. **Dependencies**

- **TUI**: Bubbletea (not gocui) - more modern and flexible
- **Testing**: testify - comprehensive assertions
- **Config**: gopkg.in/yaml.v3 - YAML parsing
- **Styling**: lipgloss (will add in Step 9)
- **No Framework Detection**: User specifies framework in config (simpler)

### 6. **Data Flow**

- **V1**: User runs all tests, tool executes test command and groups results
- **Future**: File watching, selective test runs, incremental results
- **Input**: JSON from test frameworks (RSpec, Minitest)
- **Output**: Grouped failures in TUI
- **No Caching**: V1 keeps it simple - fresh run every time

### 7. **Error Handling**

- **JSON Parsing**: Graceful handling of missing fields
- **Git Integration**: Degrade gracefully if not in git repo
- **Test Execution**: Clear error messages for command failures
- **No Cache**: Simple approach - no cache files to handle

## Technical Implementation Details

### Core Types (internal/types/)

```go
type StackFrame struct {
    File     string
    Line     int
    Function string
}

type TestResult struct {
    Name              string
    Status            TestStatus
    ErrorMessage      string
    FullBacktrace     []StackFrame  // 50 frames max
    FilteredBacktrace []StackFrame  // project frames only
}

type FailureGroup struct {
    Hash                string
    ErrorMessage        string
    NormalizedBacktrace []StackFrame
    Tests               []TestResult
    Count               int
}
```

### Backtrace Normalization Rules

- Remove frames matching exclude patterns from config
- Generate grouping key from bottom frame only (file:line format)
- Keep both full and filtered backtraces for display
- Use ErrorLocationStrategy for V1 grouping
- Detect line-level changes using git diff commands:
  - `git diff --unified=0` for uncommitted changes
  - `git diff HEAD~1 --unified=0` for last commit changes
  - `git diff HEAD~2 HEAD~1 --unified=0` for previous commit changes

### UI Design (LazyGit-style)

- **3 Panes**: Groups | Tests | Backtrace
- **Navigation**: Tab/Shift+Tab between panes, arrows within panes
- **Keybindings**: f (toggle frames), o (open file), r (re-run), q (quit)
- **Highlighting**: Line-level change detection with 3 intensity levels:
  - **Intensity 3**: Bright highlight for uncommitted changes
  - **Intensity 2**: Medium highlight for last commit changes
  - **Intensity 1**: Weak highlight for previous commit changes
  - **Intensity 0**: No highlight for unchanged lines

## Current Implementation Status

### ✅ **Completed (Steps 1-7)**

- **Step 1**: Go module initialized, core types defined and tested
- **Step 2**: JSON parser with RSpec/Minitest support, comprehensive tests
- **Step 3**: Configuration system with YAML support, framework specification
- **Step 4**: Backtrace normalizer (filter using config exclude patterns)
- **Step 5**: Failure grouper with ErrorLocation strategy (group by bottom frame)
- **Step 6**: Git integration with line-level change detection (3 intensity levels)
- **Step 7**: Test runner service for GUI-driven test execution
- **Build System**: Makefile configured, clean development workflow
- **CLI**: Basic commands working (version, config, JSON parsing, run)
- **Testing**: All unit tests passing, comprehensive test coverage
- **Project Structure**: Clean organization, proper gitignore

### 🔄 **Next Steps (Steps 8-9)**

- **Step 8**: Basic Bubbletea UI (single pane, then multi-pane)
- **Step 9**: Multi-pane UI with test runner integration

### 🎯 **Future Steps (Steps 9-12)**

- **Step 9**: Advanced UI features (keybindings, file opening, re-run)
- **Step 10**: Multi-pane UI (groups, tests, backtrace panes)
- **Step 11**: Polish & documentation (README, help screens)
- **Step 12**: Production release (error handling, final testing)

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

## File Structure

```
wing_commander/
├── cmd/wing_commander/main.go
├── internal/
│   ├── types/
│   │   ├── types.go
│   │   └── types_test.go
│   ├── parser/
│   │   ├── parser.go
│   │   ├── parser_test.go
│   │   └── schema.go
│   └── config/
│       ├── config.go
│       └── config_test.go
├── testdata/
│   ├── fixtures/
│   │   ├── rspec_failures.json
│   │   └── minitest_failures.json
│   └── config/
│       └── sample_config.yml
├── Makefile
├── go.mod
├── go.sum
├── .gitignore
└── CONTEXT.md
```

## Success Criteria

- Groups 100-1000 test failures efficiently (<1s)
- Responsive TUI navigation
- Recently changed files highlighted correctly
- Can open files in editor at correct line
- Can re-run specific test groups
- Simple workflow: run tests → see grouped failures
