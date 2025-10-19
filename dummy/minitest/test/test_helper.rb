require 'minitest/autorun'
require 'minitest/reporters'

# Add lib directory to load path
$LOAD_PATH.unshift File.expand_path('../lib', __dir__)

# Configure JUnit XML reporter
Minitest::Reporters.use! [Minitest::Reporters::JUnitReporter.new('test/reports')]
