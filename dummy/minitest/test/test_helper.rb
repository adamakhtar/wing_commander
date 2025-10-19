require 'minitest/autorun'
require 'ci/reporter/rake/minitest_loader'

# Add lib directory to load path
$LOAD_PATH.unshift File.expand_path('../lib', __dir__)

# Configure ci_reporter for JUnit XML output
ENV['CI_REPORTS'] = File.expand_path('../test/reports', __dir__)
