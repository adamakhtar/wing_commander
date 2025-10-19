#!/bin/bash
# Run tests and output JUnit XML
cd /Users/adamakhtar/Projects/active/wing_commander/dummy/minitest
ruby -Ilib:test test/thing_test.rb
cat test/reports/TEST-ThingTest.xml
