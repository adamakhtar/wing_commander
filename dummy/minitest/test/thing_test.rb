require 'test_helper'
require 'thing'

class ThingTest < Minitest::Test
  def test_boom_first_case
    Thing.new.boom
  end

  def test_boom_second_case
    Thing.new.boom
  end
end
