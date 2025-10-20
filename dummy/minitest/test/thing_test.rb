require 'test_helper'
require 'thing'

class ThingTest < Minitest::Test
  def test_boom_first_case
    Thing.new.boom
  end

  def test_boom_second_case
    Thing.new.boom
  end

  def test_error_in_test
    raise "error in test"
  end

  def test_expectation_not_met
    assert_equal "foo", "bar"
  end
end
