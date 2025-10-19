package parser

// This file documents the expected JUnit XML schemas for different test frameworks

/*
JUnit XML Output Schema (standard format):

<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
  <testsuite name="TestSuite" tests="3" failures="2" skipped="0" time="0.257">
    <testcase classname="ClassName" name="test_name" time="0.123">
      <failure message="Error message">
        app/models/user.rb:42:in `create_user'
        spec/models/user_spec.rb:15:in `block (2 levels)'
        /gems/rspec-core-3.12.0/lib/rspec/core/example.rb:259:in `instance_eval'
      </failure>
    </testcase>
    <testcase classname="ClassName" name="test_name_2" time="0.089">
      <skipped message="Skipped reason"/>
    </testcase>
    <testcase classname="ClassName" name="test_name_3" time="0.045">
    </testcase>
  </testsuite>
</testsuites>

RSpec JUnit XML Output (via rspec_junit_formatter gem):

<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
  <testsuite name="RSpec" tests="3" failures="2" skipped="0" time="0.257">
    <testcase classname="User" name="should be valid" time="0.123">
      <failure message="Expected User to be valid">
        app/models/user.rb:42:in `create_user'
        spec/models/user_spec.rb:15:in `block (2 levels) in &lt;top (required)&gt;'
      </failure>
    </testcase>
  </testsuite>
</testsuites>

Minitest JUnit XML Output (via ci_reporter_minitest gem):

<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
  <testsuite name="Minitest" tests="3" failures="2" skipped="0" time="0.268">
    <testcase classname="UserTest" name="test_user_creation" time="0.156">
      <failure message="Expected User to be valid">
        app/models/user.rb:42:in `create_user'
        test/models/user_test.rb:15:in `test_user_creation'
        /gems/minitest-5.16.0/lib/minitest/test.rb:98:in `run'
      </failure>
    </testcase>
  </testsuite>
</testsuites>

Pytest JUnit XML Output (via --junit-xml flag):

<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
  <testsuite name="pytest" tests="3" failures="2" skipped="0" time="0.234">
    <testcase classname="test_user" name="test_user_creation" time="0.123">
      <failure message="AssertionError: Expected user to be valid">
        File "app/models/user.py", line 42, in create_user
        File "test/test_user.py", line 15, in test_user_creation
      </failure>
    </testcase>
  </testsuite>
</testsuites>

Jest JUnit XML Output (via jest-junit reporter):

<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
  <testsuite name="jest" tests="3" failures="2" skipped="0" time="0.345">
    <testcase classname="User" name="should be valid" time="0.123">
      <failure message="Expected user to be valid">
        at createUser (app/models/user.js:42:10)
        at Object.test (test/user.test.js:15:5)
      </failure>
    </testcase>
  </testsuite>
</testsuites>

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

Test Status Mapping:

- Presence of <failure> element → StatusFail
- Presence of <skipped> element → StatusSkip
- Neither failure nor skipped → StatusPass

Framework Detection:

RSpec indicators:
- Test names containing "should", "expect", "describe", "it "
- Classname often matches model/class names

Minitest indicators:
- Test names containing "test_", "assert_"
- Classname often ends with "Test"

Unknown framework:
- No clear indicators found
*/
