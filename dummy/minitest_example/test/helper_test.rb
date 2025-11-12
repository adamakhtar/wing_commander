# frozen_string_literal: true

require "test_helper"

class HelperTest < Minitest::Test
  def test_helper_production_error
    helper = MinitestExample::Helper.new
    helper.process(5, 3, raise_error: true)  # Fails in Helper#process
  end

  def test_helper_production_success
    helper = MinitestExample::Helper.new
    assert_equal 8, helper.process(5, 3, raise_error: false)  # Fails in Helper#process
  end

  def test_skipped
    skip "This test is skipped"
  end
end
