# frozen_string_literal: true

# WingCommanderReporter - Custom Minitest reporter for Wing Commander CLI
#
# Usage:
#   require 'custom_reporter'
#   Minitest::Reporters.use! [WingCommanderReporter.new]
#
#   # Or with custom options:
#   Minitest::Reporters.use! [
#     WingCommanderReporter.new(
#       backtrace_depth: 50,
#       summary_output_path: '/path/to/summary.yml'  # Optional: save summary to file
#     )
#   ]
#
# Output format:
#   Progress markers: <<START>>PPFSSP<<END>> (P=pass, F=fail, S=skip) - always to stdout
#   Summary: YAML array of all test details - to stdout or file if specified

require 'yaml'
require 'minitest/reporters'
require 'fileutils'


class WingCommanderReporter < Minitest::Reporters::BaseReporter
  def initialize(backtrace_depth: 50, summary_output_path: nil, **options)
    super(options)
    @backtrace_depth = backtrace_depth
    @summary_output_path = summary_output_path
    @all_tests = []
  end

  def start
    super
    # Delete existing summary file if output path is configured
    if @summary_output_path && File.exist?(@summary_output_path)
      File.delete(@summary_output_path)
    end
    io.puts '<<START>>'
  end

  def record(result)
    super

    # Output progress marker immediately
    if result.passed?
      io.print 'P'
    elsif result.skipped?
      io.print 'S'
    elsif result.failure
      io.print 'F'
    end

    # Store all tests for summary
    @all_tests << result
  end

  def report
    super
    io.puts
    io.puts '<<END>>'

    # Output YAML summary of all tests
    summary = @all_tests.map { |test| build_test_summary(test) }
    summary_yaml = YAML.dump(summary)

    if @summary_output_path
      # Write summary to file
      summary_dir = File.dirname(@summary_output_path)
      FileUtils.mkdir_p(summary_dir) unless summary_dir == '.' || summary_dir.empty?
      File.write(@summary_output_path, summary_yaml)
    else
      # Write summary to stdout
      io.puts summary_yaml
    end
  end

  private

  def build_test_summary(result)
    test_class_name = if result.respond_to?(:klass)
      klass = result.klass
      klass.is_a?(String) ? klass : klass.name
    elsif result.respond_to?(:test)
      result.test.class.name
    else
      result.class.name
    end

    summary = {
      'test_group_name' => test_class_name,
      'test_case_name' => result.name,
      'test_status' => determine_status(result),
      'duration' => format_duration(result.time)
    }

    # Test file path and line number
    source_location = get_source_location(result)
    if source_location
      summary['test_file_path'] = File.expand_path(source_location[0])
      summary['test_line_number'] = source_location[1]
    end

    # Failure cause and details
    if result.failure
      failure_details = extract_failure_details(result)
      failure_details.each do |key, value|
        next if value.nil? || value == ''
        summary[key] = value
      end

      # Full backtrace
      if result.failure.exception
        backtrace = result.failure.exception.backtrace
        if backtrace
          summary['full_backtrace'] = backtrace.first(@backtrace_depth)
        end
      end
    end

    summary
  end

  def determine_status(result)
    if result.passed?
      'passed'
    elsif result.skipped?
      'skipped'
    else
      'failed'
    end
  end

  def extract_failure_details(result)
    failure = result.failure
    return {} unless failure

    if failure.is_a?(Minitest::Assertion) && !result.error?
      extract_assertion_details(failure)
    else
      extract_error_details(result)
    end
  end

  def extract_error_details(result)
    exception = result.failure.exception
    return {} unless exception

    details = {
      'failure_details' => exception.message
    }

    # Try to get error location from backtrace_locations first
    if exception.respond_to?(:backtrace_locations) && exception.backtrace_locations&.first
      location = exception.backtrace_locations.first
      details['failure_file_path'] = File.expand_path(location.path)
      details['failure_line_number'] = location.lineno
    elsif exception.backtrace&.first
      # Fallback to parsing first backtrace line
      file_path, line_number = parse_backtrace_line(exception.backtrace.first)
      if file_path
        details['failure_file_path'] = File.expand_path(file_path)
        details['failure_line_number'] = line_number
      end
    end

    details
  end

  def extract_assertion_details(failure)
    return {} unless failure

    details = {
      'failure_details' => failure.message
    }

    # Parse location string (format: "file:line" or "file:line:in method")
    if failure.location
      file_path, line_number = parse_location_string(failure.location)
      if file_path
        details['failure_file_path'] = File.expand_path(file_path)
        details['failure_line_number'] = line_number
      end
    end

    details
  end

  def parse_location_string(location_str)
    return [nil, nil] unless location_str

    # Format: "file:line" or "file:line:in method"
    match = location_str.match(/^(.+?):(\d+)/)
    return [nil, nil] unless match

    file_path = match[1]
    line_number = match[2].to_i
    [file_path, line_number]
  end

  def parse_backtrace_line(backtrace_line)
    return [nil, nil] unless backtrace_line

    # Format: "file:line:in method" or "file:line"
    match = backtrace_line.match(/^(.+?):(\d+)/)
    return [nil, nil] unless match

    file_path = match[1]
    line_number = match[2].to_i
    [file_path, line_number]
  end

  def get_source_location(result)
    if result.respond_to?(:klass)
      result.source_location
    else
      result.method(result.name).source_location
    end
  end

  def format_duration(time)
    sprintf('%.2f', time)
  end
end