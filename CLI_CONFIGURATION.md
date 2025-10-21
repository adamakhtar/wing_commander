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

### `--test-command CMD`

- **Purpose**: Specify the test runner command with template interpolation support
- **Type**: String (command with optional template variables)
- **Default**: Must be provided (no hardcoded default)
- **Template Syntax**: Uses Go `text/template` syntax
- **Available Variables**: `{{.Paths}}` (for test paths, empty by default)
- **Example**: `--test-command "rails test {{.Paths}} --output .wing_commander/test_output.xml"`

For Minitest using `minitest-reporters` JUnit reporter, point the command to generate JUnit XML. For example, when your `test_helper.rb` configures:

```ruby
require 'minitest/reporters'
Minitest::Reporters.use! [
  Minitest::Reporters::JUnitReporter.new('.wing_commander')
]
```

run Wing Commander with the project path and a command that triggers your test suite (the reporter writes XML files to `.wing_commander/` which Wing Commander reads from combined output):

```bash
wing_commander run \
  --project-path /path/to/project \
  --test-command "bundle exec rake test"
```

### `--config PATH`

- **Purpose**: Specify custom configuration file location
- **Type**: String (file path)
- **Default**: `.wing_commander/config.yml`
- **Example**: `--config /path/to/custom-config.yml`

## Configuration File Format

```yaml
# Project path (can be overridden with --project-path CLI option)
project_path: ""

# Test framework (rspec, minitest, pytest, jest)
test_framework: rspec

# Test command with template interpolation
test_command: "bundle exec rspec {{.Paths}} --format RspecJunitFormatter --out .wing_commander/test_output.xml"

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
- **Example**: `rails test {{.Paths}} --output .wing_commander/test_output.xml`

## Usage Examples

```bash
# CLI-first approach - override everything via command line
wing_commander run --project-path /path/to/project --test-command "rails test {{.Paths}} --output .wing_commander/test_output.xml"

# Mix CLI and config - project path from CLI, test command from config
wing_commander run --project-path /path/to/project

# Custom config file with CLI overrides
wing_commander run --config custom-config.yml --test-command "pytest {{.Paths}} --junit-xml=.wing_commander/test_output.xml"

# Traditional approach - everything from config file
wing_commander run

# Show current configuration (includes CLI-applied settings)
wing_commander config
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
