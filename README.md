# Wing Commander

A CLI/TUI tool for analyzing test failures by grouping them by backtrace similarity. Helps developers quickly identify shared root causes among multiple failing tests.

## Quick Start

### Development Testing

```bash
# Build dev version and launch TUI against dummy minitest app
make dev-minitest
```

### Production Usage

```bash
# Run tests and analyze failures
wing_commander run --project-path /path/to/project --test-command "rails test --output junit"

# Show configuration
wing_commander config

# Demo mode with sample data
wing_commander demo
```

## Keybindings

- `↑↓` - Navigate between items
- `Tab` - Switch between panes (Groups/Tests/Backtrace)
- `f` - Toggle full/filtered frames display
- `o` - Open selected file in editor at specific line
- `r` - Re-run tests for selected group
- `q` - Quit

## TUI Panels

- **Panel 1 – Failure Groups**

  - Groups failures by cause (Production Code Error, Test Definition Error, Failed Assertion)
  - Each section shows: `{icon} {count} - {error message}` and `{bottom frame}`
  - Failure cause icons: 🚀 (production), 🔧 (test definition), ❌ (assertion)
  - Count and error message displayed in yellow
  - Bottom frame shows relative file path and line number

- **Panel 2 – Tests in Selected Group**

  - First line: test name
  - Second line: tail frames (all frames except the shared bottom frame) shown as a chain: `file1:line → file2:line → ...`
  - If there are no additional frames, shows `(no additional frames)`

- **Panel 3 – Backtrace**
  - Full backtrace of the selected test
  - Frames are highlighted by change intensity as per Git integration
  - `f` toggles between filtered (project-only) and full frames

## Supported Test Frameworks

- RSpec (Ruby)
- Minitest (Ruby)
- Pytest (Python)
- Jest (JavaScript)

### Minitest setup

If you use Minitest, ensure a JUnit XML reporter is enabled so Wing Commander can parse failures:

```ruby
# test/test_helper.rb
require 'minitest/reporters'
Minitest::Reporters.use! [
  Minitest::Reporters::JUnitReporter.new('test/reports')
]
```

This produces XML files under `test/reports` with embedded failure/error details. Wing Commander extracts stack frames from both the formatted failure body and any `SystemErr` output, including embedded `file.rb:LINE` tokens like `[test/thing_test.rb:18]`.

## Development

```bash
# Build development version
make dev

# Run tests
make test

# Clean build artifacts
make clean
```

## Architecture

- **CLI-first configuration**: CLI options > Config file > Sensible defaults
- **JUnit XML parsing**: Supports standard test output format
- **Backtrace grouping**: Groups failures by error location similarity
- **Git integration**: Highlights recently changed files with intensity levels
- **Interactive TUI**: Built with Bubbletea for smooth navigation

## Failure Cause Classification

Wing Commander assigns each failed test one of three broad causes:

- Test definition error: Failure originates in test code, framework, setup, or teardown.
- Production code error: Exception raised by the application under test (stack points to app code).
- Assertion failure: Test body completed; an expectation/matcher reported a mismatch.

Heuristics are intentionally simple: assertion-like messages imply assertion failure; frames that clearly reference test paths/framework internals imply test definition error; otherwise the failure is attributed to production code. Hangs that produce no report are out of scope. Currently, path indicators include RSpec and Minitest conventions and can be extended as more frameworks are supported.
