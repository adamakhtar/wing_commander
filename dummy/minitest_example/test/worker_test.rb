# frozen_string_literal: true

require "test_helper"

class WorkerTest < Minitest::Test
  # Scenario 1: Production code error - Worker -> Helper backtrace
  def test_production_error
    worker = MinitestExample::Worker.new
    worker.addition(5, 3, raise_error: true)  # Fails in Helper#process
  end

  # Scenario 2: Error in test definition
  def test_test_error
    undefined_variable  # NameError in test itself
  end

  # Scenario 3: Assertion failure - Worker -> Helper backtrace
  def test_assertion_failure
    worker = MinitestExample::Worker.new
    result = worker.addition(5, 3, raise_error: false)
    assert_equal 10, result  # Fails (expected 10, got 8), backtrace shows Worker -> Helper
  end
end
