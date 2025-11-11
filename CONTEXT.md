# Wing Commander - Project Context & Decisions

## Project Overview

A CLI/TUI tool for running tests and reviewing failure details with normalized backtraces. Helps developers quickly inspect failing tests without the noise of third-party frames.

## Key Decisions Made

### 1. **Failure Presentation**

- **Display**: Flat list of failed tests
- **Store**: Full 50 frames for user viewing
- **Rationale**: Keeps the UI simple while still surfacing relevant backtrace context
- **Language Support**: Works across Ruby, Python, JavaScript, Go (file:line format)
- **Architecture**: Backtrace normalization is handled via a dedicated package
- **Change Detection**: Line-level change detection with 3 intensity levels:
  - **Intensity 3**: Uncommitted changes to frame's line number
  - **Intensity 2**: Line changed in last commit
  - **Intensity 1**: Line changed in commit before last

### 2. **Project Structure**

- **Types Location**: `internal/types/` (not `pkg/`) - internal use only
- **Config Location**: `.wing_commander/config.yml` (not root)
- **Test Files**: `testdata/` directory for fixtures and sample configs
- **Rationale**: Keep project root clean, signal internal vs public packages

### 3. **Configuration Approach - CLI-First Design**

- **Priority System**: CLI options > Config file > Sensible defaults
- **CLI Options**:
  - `--project-path PATH`: Project directory (default: current working directory)
  - `--test-command CMD`: Test runner with interpolation support (e.g., `rails test {{.Paths}} --output junit`)
  - `--config PATH`: Config file path (default: `.wing_commander/config.yml`)
- **Template Interpolation**: Uses Go `text/template` syntax (`{{.Paths}}` for test paths)
- **Format**: YAML with `project_path`, `test_framework`, `test_command`, and `exclude_patterns`
- **Framework Support**: RSpec, Minitest, Pytest, Jest
- **Default Patterns**: `/gems/`, `/lib/ruby/`, `/vendor/bundle/`, etc.
- **CLI Support**: `wing_commander run --project-path /path/to/project --test-command "rails test {{.Paths}} --output junit"`

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

- **V1**: User runs all tests, tool executes test command and aggregates results
- **Future**: File watching, selective test runs, incremental results
- **Input**: JUnit XML from test frameworks (RSpec, Minitest)
- **Output**: Flat list of failures in TUI
- **No Caching**: V1 keeps it simple - fresh run every time

### 7. **Test Framework Support**

- **Dummy Projects**: `dummy/` directory contains test projects for each supported framework
- **Minitest**: Complete Ruby project with failing tests and JUnit XML reporting
  - Uses `ci_reporter_minitest` gem for JUnit XML format output
  - Two failing test cases calling `Thing.new.boom` method
  - Reports generated in `test/reports/` directory
  - Runnable with `bundle install` and `bundle exec rake test`

### 8. **Error Handling**

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
    FailureDetails    string
    FullBacktrace     []StackFrame  // 50 frames max
    FilteredBacktrace []StackFrame  // project frames only
}
```

### CLI-First Configuration Implementation

```go
// Config struct with CLI-first design
type Config struct {
    ProjectPath     string        `yaml:"project_path"`
    TestFramework   TestFramework `yaml:"test_framework"`
    TestCommand     string        `yaml:"test_command"`
    ExcludePatterns []string      `yaml:"exclude_patterns"`
}

// CLI options loading with priority system
func loadConfigWithCLIOptions(configPath, projectPath, testCommand string) (*Config, error) {
    // 1. Load base config from file
    cfg, err := config.LoadConfig(configPath)

    // 2. Override with CLI options (highest priority)
    if projectPath != "" {
        cfg.ProjectPath = resolveAbsolutePath(projectPath)
    } else {
        cfg.ProjectPath = getCurrentWorkingDirectory()
    }

    if testCommand != "" {
        cfg.TestCommand = testCommand
    }

    return cfg, nil
}

// Template interpolation in test runner
func (r *TestRunner) executeTestCommand() (string, error) {
    cmdTemplate, err := template.New("testCommand").Parse(r.config.TestCommand)
    templateData := struct{ Paths string }{Paths: ""} // Empty by default
    var cmdBuilder strings.Builder
    cmdTemplate.Execute(&cmdBuilder, templateData)
    finalCommand := cmdBuilder.String()
    // Execute command in project directory...
}
```

### Backtrace Normalization Rules

- Remove frames matching exclude patterns from config
- Keep both full and filtered backtraces for display (no grouping)
- Detect line-level changes using git diff commands:
  - `git diff --unified=0` for uncommitted changes
  - `git diff HEAD~1 --unified=0` for last commit changes
  - `git diff HEAD~2 HEAD~1 --unified=0` for previous commit changes

### UI Design (LazyGit-style)

- **3 Panes**: Test Runs | Tests | Backtrace
- **Navigation**: Tab/Shift+Tab between panes, arrows within panes
- **Keybindings**: f (toggle frames), o (open file), r (re-run), q (quit)
- **Highlighting**: Line-level change detection with 3 intensity levels:
  - **Intensity 3**: Bright highlight for uncommitted changes
  - **Intensity 2**: Medium highlight for last commit changes
  - **Intensity 1**: Weak highlight for previous commit changes
  - **Intensity 0**: No highlight for unchanged lines

## Current Implementation Status

### ✅ **Completed (Steps 1-8)**

- **Step 1**: Go module initialized, core types defined and tested
- **Step 2**: JSON parser with RSpec/Minitest support, comprehensive tests
- **Step 3**: Configuration system with YAML support, framework specification
- **Step 4**: Backtrace normalizer extracted for reuse across the app
- **Step 5**: Flat failure list replaces grouping logic
- **Step 6**: Git integration with line-level change detection (3 intensity levels)
- **Step 7**: Test runner service for GUI-driven test execution
- **Step 8**: Basic Bubbletea TUI with 3-pane layout
- **Build System**: Makefile configured, clean development workflow
- **CLI**: Basic commands working (version, config, JSON parsing, run)
- **CLI Flags**: `--config`, `--project-path`, `--test-command` flags for run command
- **Testing**: All unit tests passing, comprehensive test coverage
- **Project Structure**: Clean organization, proper gitignore

### 🔄 **Next Steps (Step 9)**

- **Step 9**: Advanced UI features (keybindings, file opening, re-run)

### 🎯 **Future Steps (Steps 10-12)**

- **Step 10**: Multi-pane UI refinements
- **Step 11**: Polish & documentation (README, help screens)
- **Step 12**: Production release (error handling, final testing)

## Development Workflow

```bash
# Build and test
make dev
make test

# Run CLI
make run

# Run with CLI options
./bin/wing_commander run --project-path /path/to/project --test-command "rails test {{.Paths}} --output junit"

# Run with custom config
./bin/wing_commander run --config /path/to/config.yml

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
├── dummy/
│   └── minitest/          # Test framework dummy projects
│       ├── lib/
│       │   └── thing.rb   # Simple class with boom method
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

## Success Criteria

- Displays 100-1000 test failures efficiently (<1s)
- Responsive TUI navigation
- Recently changed files highlighted correctly
- Can open files in editor at correct line
- Can re-run selected tests
- Simple workflow: run tests → review failures
