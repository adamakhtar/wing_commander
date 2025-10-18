# Wing Commander - Project Context & Decisions

## Project Overview

A CLI/TUI tool for analyzing test failures by grouping them by backtrace similarity. Helps developers quickly identify shared root causes among multiple failing tests.

## Key Decisions Made

### 1. **Grouping Strategy**

- **Group by**: Normalized backtrace (file + method name, excluding line numbers)
- **Store**: Full 50 frames for user viewing
- **Rationale**: Line numbers change with code edits, but file+method stays stable
- **Language Support**: Works across Ruby, Python, JavaScript, Go (method names in stack traces)

### 2. **Project Structure**

- **Types Location**: `internal/types/` (not `pkg/`) - internal use only
- **Config Location**: `.wing_commander/config.yml` (not root)
- **Cache Location**: `.wing_commander/cache.json` (not root)
- **Rationale**: Keep project root clean, signal internal vs public packages

### 3. **Configuration Approach**

- **V1**: User creates config file manually
- **Future**: Tool can generate config file on request
- **Format**: YAML with exclude_patterns and test_command
- **Default Patterns**: `/gems/`, `/lib/ruby/`, `/vendor/bundle/`, etc.

### 4. **Build System**

- **Primary**: Makefile (standard, cross-platform)
- **Removed**: dev.sh (redundant wrapper)
- **Build Locations**: `bin/` (dev), `dist/` (production)
- **Commands**: `make dev`, `make test`, `make run`, `make clean`

### 5. **Dependencies**

- **TUI**: Bubbletea (not gocui) - more modern and flexible
- **Testing**: testify - comprehensive assertions
- **Config**: gopkg.in/yaml.v3 (added by go mod tidy)
- **Styling**: lipgloss (will add in Step 9)

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
- Parse method names from stack frames
- Generate hash from file+method (ignore line numbers)
- Keep both full and filtered backtraces

### UI Design (LazyGit-style)

- **3 Panes**: Groups | Tests | Backtrace
- **Navigation**: Tab/Shift+Tab between panes, arrows within panes
- **Keybindings**: f (toggle frames), o (open file), r (re-run), q (quit)
- **Highlighting**: Recently changed files marked with `[*]`

## Current State (Step 1 Complete)

- ✅ Go module initialized
- ✅ Core types defined and tested
- ✅ Basic CLI working
- ✅ Build system (Makefile) configured
- ✅ Git repository clean
- ✅ All tests passing

## Next Steps

- **Step 2**: JSON Parser (accept file path, parse test results)
- **Step 3**: Configuration System (load .wing_commander/config.yml)
- **Step 4**: Backtrace Normalizer (filter using config patterns)
- **Step 5**: Failure Grouper (group by normalized signature)

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
├── internal/types/
│   ├── types.go
│   └── types_test.go
├── Makefile
├── go.mod
├── go.sum
└── .gitignore
```

## Success Criteria

- Groups 100-1000 test failures efficiently (<1s)
- Responsive TUI navigation
- Recently changed files highlighted correctly
- Can open files in editor at correct line
- Can re-run specific test groups
- Simple workflow: run tests → see grouped failures
