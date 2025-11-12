package parser

// This file documents the expected WingCommanderReporter YAML summary schema.
/*
WingCommanderReporter Summary Schema:

---
- test_group_name: WorkerTest
  test_case_name: test_assertion_failure
  test_status: failed
  duration: "0.00"
  test_file_path: "/abs/path/to/test/worker_test.rb"
  test_line_number: 18
  failure_details: "Expected: 10\n  Actual: 8"
  failure_file_path: "/abs/path/to/test/worker_test.rb"
  failure_line_number: 21
  full_backtrace:
    - "/abs/path/to/test/worker_test.rb:21:in `test_assertion_failure'"

Field Definitions:

- test_group_name: Class or group name for the test case.
- test_case_name: Individual test name (Minitest method name).
- test_status: "passed", "failed", or "skipped".
- duration: String formatted to two decimal places.
- test_file_path: Absolute path to the test definition file.
- test_line_number: Definition line number.
- failure_details: Human readable failure message (empty for pass/skip).
- failure_file_path: Absolute path where the failure originated.
- failure_line_number: Line number associated with the failure.
- full_backtrace: Array of strings representing the captured backtrace.

Additional Notes:

- WingCommanderReporter always emits an array (possibly empty).
- Paths are expanded to absolute paths before serialization.
- Backtrace entries are limited to the configured `backtrace_depth`.
*/
