# frozen_string_literal: true

$LOAD_PATH.unshift File.expand_path("../lib", __dir__)
require "minitest_example"
require "minitest/autorun"
require "minitest/reporters"
require_relative "wing_commander_reporter"

Minitest::Reporters.use! [
  # Minitest::Reporters::JUnitReporter.new('.wing_commander/test_results/')
  WingCommanderReporter.new(backtrace_depth: 50, summary_output_path: '.wing_commander/test_results/summary.yml')
]
