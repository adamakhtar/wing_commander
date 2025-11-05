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
- test_case_name: ThingTest
  test_file_path: /absolute/path/to/test/thing_test.rb
  test_line_number: 6
  failure_cause:
  error_message:
  error_file_path:
  error_line_number:
  failed_assertion_details:
  assertion_file_path:
  assertion_line_number:
  full_backtrace:
  test_status: passed
  duration: "0.05"
- test_case_name: ThingTest
  test_file_path: /absolute/path/to/test/thing_test.rb
  test_line_number: 10
  failure_cause: error
  error_message: NameError: uninitialized constant ThingTest::Thing
  error_file_path: /absolute/path/to/test/thing_test.rb
  error_line_number: 11
  failed_assertion_details:
  assertion_file_path:
  assertion_line_number:
  full_backtrace:
    - /absolute/path/to/test/thing_test.rb:11:in `test_boom_second_case'
    - /absolute/path/to/gems/minitest-5.16.0/lib/minitest/test.rb:98:in `block in run'
    - /absolute/path/to/gems/minitest-5.16.0/lib/minitest/test.rb:164:in `run'
  test_status: failed
  duration: "0.00"
- test_case_name: ThingTest
  test_file_path: /absolute/path/to/test/thing_test.rb
  test_line_number: 14
  failure_cause: error
  error_message: RuntimeError: error in test
  error_file_path: /absolute/path/to/test/thing_test.rb
  error_line_number: 15
  failed_assertion_details:
  assertion_file_path:
  assertion_line_number:
  full_backtrace:
    - /absolute/path/to/test/thing_test.rb:15:in `test_error_in_test'
    - /absolute/path/to/gems/minitest-5.16.0/lib/minitest/test.rb:98:in `block in run'
  test_status: failed
  duration: "0.00"
- test_case_name: ThingTest
  test_file_path: /absolute/path/to/test/thing_test.rb
  test_line_number: 18
  failure_cause: failed_assertion
  error_message:
  error_file_path:
  error_line_number:
  failed_assertion_details: Expected: "foo"\n  Actual: "bar"
  assertion_file_path: /absolute/path/to/test/thing_test.rb
  assertion_line_number: 19
  full_backtrace:
    - /absolute/path/to/test/thing_test.rb:19:in `test_expectation_not_met'
    - /absolute/path/to/gems/minitest-5.16.0/lib/minitest/test.rb:98:in `block in run'
  test_status: failed
  duration: "0.00"
- test_case_name: ThingTest
  test_file_path: /absolute/path/to/test/thing_test.rb
  test_line_number: 21
  failure_cause:
  error_message:
  error_file_path:
  error_line_number:
  failed_assertion_details:
  assertion_file_path:
  assertion_line_number:
  full_backtrace:
  test_status: skipped
  duration: "0.00"
```

**Example with no tests (empty array):**

```yaml
---
[]
```

### Required Fields (13 total)

1. **test_case_name** - Test class name (`result.class.name`)
2. **test_file_path** - Absolute file path of test file (from `source_location`, expanded)
3. **test_line_number** - Line number of test definition (from `source_location`)
4. **failure_cause** - Either `"error"` or `"failed_assertion"`
   - `"error"` for `Minitest::UnexpectedError` (detected via `result.error?`)
   - `"failed_assertion"` for `Minitest::Assertion` failures
5. **error_message** - Error message (only when `failure_cause == "error"`)
6. **error_file_path** - Absolute file path where error occurred (only when `failure_cause == "error"`)
7. **error_line_number** - Line number where error occurred (only when `failure_cause == "error"`)
8. **failed_assertion_details** - Assertion failure message (only when `failure_cause == "failed_assertion"`)
9. **assertion_file_path** - Absolute file path of assertion failure (only when `failure_cause == "failed_assertion"`)
10. **assertion_line_number** - Line number of assertion failure (only when `failure_cause == "failed_assertion"`)
11. **full_backtrace** - Array of backtrace strings (limited to `backtrace_depth` lines)
12. **test_status** - String: `"passed"`, `"failed"`, or `"skipped"`
13. **duration** - String format with exactly 2 decimal places (e.g., `"2.00"`, `"2.54"`)

## Data Extraction Requirements

### Test Location

- Use `get_source_location` helper method
- Pattern: Check `result.respond_to?(:klass)`, then use `result.source_location`, else use `result.method(result.name).source_location`
- Always expand to absolute path using `File.expand_path`

### Error Details (when `failure_cause == "error"`)

- Extract from `result.failure.exception`
- Prefer `exception.backtrace_locations.first` for file path and line number
- Fallback to parsing first backtrace line if `backtrace_locations` unavailable
- Parse format: `"file:line"` or `"file:line:in method"`

### Assertion Details (when `failure_cause == "failed_assertion"`)

- Extract from `result.failure.message` for assertion details
- Parse `result.failure.location` string for file path and line number
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
