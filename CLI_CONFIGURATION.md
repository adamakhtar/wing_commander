# CLI-First Configuration System

## Overview

Wing Commander now uses a CLI-first configuration approach where command-line options take precedence over configuration files, which in turn take precedence over sensible defaults.

## Configuration Priority

1. **CLI Options** (Highest Priority)
2. **Config File** (Medium Priority)
3. **Sensible Defaults** (Lowest Priority)

## CLI Options

### `--project-path PATH`

- **Purpose**: Specify the project directory whose tests are being observed
- **Type**: String (absolute or relative path)
- **Default**: Current working directory
- **Behavior**: Relative paths are automatically converted to absolute paths
- **Example**: `--project-path /path/to/my/project`

### `--run-command CMD`

- **Purpose**: Command that executes your test suite (WingCommanderReporter must be registered in `test_helper.rb`)
- **Type**: String (command with optional template variables)
- **Template Syntax**: Uses Go `text/template` syntax
- **Available Variables**: `{{.Paths}}` (for test paths, empty by default)
- **Example**: `--run-command "bundle exec rake test {{.Paths}}"`

### `--test-results-path FILE`

- **Purpose**: Absolute or relative path to the YAML summary produced by `WingCommanderReporter`
- **Type**: String (file path)
- **Validation**: Must point to an existing file when the CLI starts
- **Example**: `--test-results-path ".wing_commander/test_results/summary.yml"`

## Configuration File Format

```yaml
# Project path (can be overridden with --project-path CLI option)
project_path: ""

# Test framework (rspec, minitest, pytest, jest)
test_framework: rspec

# Test command with template interpolation (WingCommanderReporter is responsible for writing the summary)
test_command: "bundle exec rake test {{.Paths}}"

# File written by WingCommanderReporter
test_results_path: ".wing_commander/test_results/summary.yml"

# Patterns to exclude from backtrace analysis
exclude_patterns:
  - "/gems/"
  - "/lib/ruby/"
  - "/vendor/bundle/"
```

## Template Interpolation

- **Engine**: Go's `text/template` package
- **Syntax**: `{{.VariableName}}`
- **Available Variables**:
  - `{{.Paths}}`: Test paths (empty by default, ready for future expansion)
- **Example**: `--run-command "bundle exec rake test {{.Paths}}"`

## Usage Examples

```bash
# CLI-first approach - override everything via command line
wing_commander start /path/to/project \
  --run-command "bundle exec rake test" \
  --test-file-pattern "test/**/*_test.rb" \
  --test-results-path ".wing_commander/test_results/summary.yml"

# Mix CLI and config - project path from CLI, rest from config
wing_commander start /path/to/project

# Custom config file with CLI overrides
wing_commander start /path/to/project \
  --config custom-config.yml \
  --run-command "bundle exec rake test {{.Paths}}"
```

## Implementation Details

### Key Files Modified

- `cmd/wing_commander/main.go`: Added CLI option parsing and priority system
- `internal/config/config.go`: Added `ProjectPath` field and updated defaults
- `internal/runner/runner.go`: Added template interpolation and project path support
- `testdata/config/sample_config.yml`: Updated example configuration

### Technical Architecture

- **CLI Framework**: Go's standard `flag` package
- **Template Engine**: Go's `text/template` package
- **Path Resolution**: Automatic conversion of relative to absolute paths
- **Error Handling**: Clear messages guiding users to proper configuration

### Backward Compatibility

- Existing configuration files continue to work unchanged
- CLI options are additive - they don't break existing workflows
- Default behavior remains the same when no CLI options are provided

## Future Extensions

The template interpolation system is designed to be extensible:

- **Test Path Selection**: Future versions can populate `{{.Paths}}` with specific test files
- **Additional Variables**: Easy to add more template variables as needed
- **Framework-Specific Variables**: Can add framework-specific interpolation variables

## Testing

All existing tests have been updated and pass. The CLI-first configuration system has been tested with:

- Various CLI option combinations
- Template interpolation functionality
- Project path resolution from different directories
- Configuration priority system validation
