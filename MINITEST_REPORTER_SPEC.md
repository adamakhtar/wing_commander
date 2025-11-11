# Minitest Reporter Requirements Specification

## Overview

This document specifies the requirements for the `WingCommanderReporter` - a custom Minitest reporter that provides progress markers during test execution and a detailed YAML summary of failed tests.

## Class Structure

- Extends `Minitest::Reporters::BaseReporter`
- Requires `yaml` and `fileutils` libraries
- Constructor accepts `backtrace_depth` keyword argument (default: 50)
- Constructor accepts `summary_output_path` keyword argument (default: `nil`)
  - If `nil`: summary written to stdout
  - If path string: summary written to file at specified path
- Constructor accepts `**options` hash (passed to parent)
- Uses parent's `io` accessor for progress output (stdout)

## Progress Reporter Requirements

### Timing

- Output progress markers immediately as each test completes (in `record` method)
- Output concatenated markers without newlines between them

### Markers

- `<<START>>` - Output at test suite start (`start` method)
- `P` - Output immediately when test passes
- `F` - Output immediately when test fails
- `S` - Output immediately when test is skipped
- `<<END>>` - Output after all tests complete (`report` method)

### Output Format

- Concatenated on single line: `<<START>>PPFSSP<<END>>`
- No delimiters between individual markers
- Always output to stdout via `io` accessor
- Newline after `<<END>>` before summary (if summary goes to stdout)

## Summary Reporter Requirements

### Timing

- Output only at end after `<<END>>` marker
- Includes all tests: passed, failed, and skipped
- Output empty YAML array `[]` if no tests were run

### Output Destination

- **Default behavior** (`summary_output_path` is `nil`): Summary written to stdout via `io` accessor
- **File output** (`summary_output_path` specified): Summary written to file at specified path
- **File management**:
  - If output file exists at start of test run, it is deleted before tests begin
  - Parent directory of output file is created if it doesn't exist
  - Progress markers always go to stdout regardless of summary destination

### Format

- YAML format using `YAML.dump`
- Array of hashes, one per test (passed, failed, or skipped)

### Sample Output

**Example with all test statuses (passed, failed, skipped):**

```yaml
---
- test_group_name: ThingTest
  test_file_path: /absolute/path/to/test/thing_test.rb
  test_line_number: 6
  failure_details:
  failure_file_path:
  failure_line_number:
  full_backtrace:
  test_status: passed
  duration: "0.05"
- test_group_name: ThingTest
  test_file_path: /absolute/path/to/test/thing_test.rb
  test_line_number: 10
  failure_details: NameError: uninitialized constant ThingTest::Thing
  failure_file_path: /absolute/path/to/test/thing_test.rb
  failure_line_number: 11
  full_backtrace:
    - /absolute/path/to/test/thing_test.rb:11:in `test_boom_second_case'
    - /absolute/path/to/gems/minitest-5.16.0/lib/minitest/test.rb:98:in `block in run'
    - /absolute/path/to/gems/minitest-5.16.0/lib/minitest/test.rb:164:in `run'
  test_status: failed
  duration: "0.00"
- test_group_name: ThingTest
  test_file_path: /absolute/path/to/test/thing_test.rb
  test_line_number: 14
  failure_details: RuntimeError: error in test
  failure_file_path: /absolute/path/to/test/thing_test.rb
  failure_line_number: 15
  full_backtrace:
    - /absolute/path/to/test/thing_test.rb:15:in `test_error_in_test'
    - /absolute/path/to/gems/minitest-5.16.0/lib/minitest/test.rb:98:in `block in run'
  test_status: failed
  duration: "0.00"
- test_group_name: ThingTest
  test_file_path: /absolute/path/to/test/thing_test.rb
  test_line_number: 18
  failure_details: Expected: "foo"\n  Actual: "bar"
  failure_file_path: /absolute/path/to/test/thing_test.rb
  failure_line_number: 19
  full_backtrace:
    - /absolute/path/to/test/thing_test.rb:19:in `test_expectation_not_met'
    - /absolute/path/to/gems/minitest-5.16.0/lib/minitest/test.rb:98:in `block in run'
  test_status: failed
  duration: "0.00"
- test_group_name: ThingTest
  test_file_path: /absolute/path/to/test/thing_test.rb
  test_line_number: 21
  failure_details:
  failure_file_path:
  failure_line_number:
  full_backtrace:
  test_status: skipped
  duration: "0.00"
```

**Example with no tests (empty array):**

```yaml
---
[]
```

### Required Fields (9 total)

1. **test_group_name** - Test class name (`result.klass.name`)
2. **test_file_path** - Absolute file path of test file (from `source_location`, expanded)
3. **test_line_number** - Line number of test definition (from `source_location`)
4. **failure_details** - Combined error/assertion message describing the failure (empty for passed/skipped)
5. **failure_file_path** - Absolute file path where the failure originated (blank if unknown)
6. **failure_line_number** - Line number of the failure origin (0 if unknown)
7. **full_backtrace** - Array of backtrace strings (limited to `backtrace_depth` lines; may be empty)
8. **test_status** - String: `"passed"`, `"failed"`, or `"skipped"`
9. **duration** - String format with exactly 2 decimal places (e.g., `"2.00"`, `"2.54"`)

## Data Extraction Requirements

### Test Location

- Use `get_source_location` helper method
- Pattern: Check `result.respond_to?(:klass)`, then use `result.source_location`, else use `result.method(result.name).source_location`
- Always expand to absolute path using `File.expand_path`

### Failure Details

- Extract human-readable message from `result.failure`
- Prefer `exception.backtrace_locations.first` for file path and line number when available
- Fallback to parsing first backtrace line or assertion location (`result.failure.location`)
- Parse format: `"file:line"` or `"file:line:in method"`

### Backtrace

- Extract from `result.failure.exception.backtrace`
- Limit to `@backtrace_depth` lines (default: 50)

## Edge Cases Handled

- Source location extraction: Uses DefaultReporter pattern with fallback
- File paths: All paths converted to absolute using `File.expand_path`
- Missing data: Fields included even if nil (YAML includes nil values)
- All test statuses: Passed, failed, and skipped tests are all included in summary
- Empty backtraces: Handled gracefully (nil or empty array for passed/skipped tests)
- Duration formatting: Always 2 decimal places using `sprintf('%.2f', time)`
- Failure cause detection: Distinguishes between errors and assertion failures (only present for failed tests)
- Conditional field extraction: Error fields only for errors, assertion fields only for assertion failures

## Integration

- Registerable via `Minitest::Reporters.use! [WingCommanderReporter.new]`
- Can be configured with `backtrace_depth` option
- Can be configured with `summary_output_path` option for file output
- Basic usage documentation in code comments

### Usage Examples

**Default (summary to stdout):**

```ruby
Minitest::Reporters.use! [WingCommanderReporter.new]
```

**Summary to file:**

```ruby
Minitest::Reporters.use! [
  WingCommanderReporter.new(summary_output_path: '/path/to/summary.yml')
]
```

**Both options configured:**

```ruby
Minitest::Reporters.use! [
  WingCommanderReporter.new(
    backtrace_depth: 50,
    summary_output_path: 'summary.yml'
  )
]
```
