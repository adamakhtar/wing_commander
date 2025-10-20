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
  - First line: error message
  - Second line: bottom frame `file:line` and failure count

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
