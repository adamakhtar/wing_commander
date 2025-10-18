package parser

// This file documents the expected JSON schemas for different test frameworks

/*
RSpec JSON Output Schema (via custom formatter):

{
  "tests": [
    {
      "name": "User should be valid",
      "status": "failed",
      "message": "Expected User to be valid",
      "backtrace": [
        "app/models/user.rb:42:in `create_user'",
        "spec/models/user_spec.rb:15:in `block (2 levels) in <top (required)>'",
        "/gems/rspec-core-3.12.0/lib/rspec/core/example.rb:259:in `instance_eval'"
      ],
      "duration": 0.123,
      "file": "spec/models/user_spec.rb",
      "line": 15
    }
  ],
  "summary": {
    "total": 10,
    "passed": 8,
    "failed": 2,
    "skipped": 0
  }
}

Minitest JSON Output Schema (via custom reporter):

{
  "tests": [
    {
      "name": "test_user_creation",
      "status": "failed", 
      "message": "Expected User to be valid",
      "backtrace": [
        "app/models/user.rb:42:in `create_user'",
        "test/models/user_test.rb:15:in `test_user_creation'",
        "/gems/minitest-5.16.0/lib/minitest/test.rb:98:in `run'"
      ],
      "duration": 0.156,
      "file": "test/models/user_test.rb",
      "line": 15
    }
  ],
  "summary": {
    "total": 10,
    "passed": 8,
    "failed": 2,
    "skipped": 0
  }
}

Alternative: Array of Tests (simpler format):

[
  {
    "name": "User should be valid",
    "status": "failed",
    "message": "Expected User to be valid", 
    "backtrace": [
      "app/models/user.rb:42:in `create_user'",
      "spec/models/user_spec.rb:15:in `block (2 levels) in <top (required)>'"
    ]
  }
]

Backtrace Frame Formats Supported:

1. Ruby (RSpec/Minitest):
   "app/models/user.rb:42:in `create_user'"
   "app/models/user.rb:42"

2. Python (pytest):
   "File \"app/models/user.py\", line 42, in create_user"

3. JavaScript (Jest):
   "at createUser (app/models/user.js:42:10)"

4. Go:
   "app/models/user.go:42 +0x123 app.createUser"

Status Values Supported:

- "pass", "passed", "success" → StatusPass
- "fail", "failed", "failure" → StatusFail  
- "skip", "skipped", "pending" → StatusSkip
- Unknown values → StatusFail (default)

Framework Detection:

RSpec indicators:
- Test names containing "should", "expect", "describe", "it "

Minitest indicators:
- Test names containing "test_", "assert_"

Unknown framework:
- No clear indicators found
*/
