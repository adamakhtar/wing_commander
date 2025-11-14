# frozen_string_literal: true

require "shellwords"


# Builds and executes minitest commands based on command line arguments.
#
# Processes test file patterns, directories, and specific test cases, then constructs
# and executes the appropriate ruby command to run minitest.
#
# @example Basic usage when executing specific test file paths / directories
#   runner = WingCommanderRunner.new
#   runner.call(["test/worker_test.rb"])
#
# @example Basic usage when executing specific test cases
#   runner = WingCommanderRunner.new
#   runner.call(test_cases: ["TestModel#test_something", "TestModel#test_something_else"])
#
# @example With custom test file globs. When no test files are provided the runner will run all tests matching the test file patterns.
# With test files provided the runner will ensure they match the test file globs.
#   runner = WingCommanderRunner.new(
#     test_file_patterns: ["spec/**/*_spec.rb"],
#     lib_dirs: %w[lib spec]
#   )
#   runner.call(["spec/models/user_spec.rb"]) # OK
#   runner.call(["spec/models/user_test.rb"]) # Not OK, test file glob pattern does not match
#
# @example Command line usage patterns
#   # Run all tests matching default glob pattern
#   bin/test
#
#   # Run specific test file
#   bin/test test/post_test.rb
#
#   # Run specific test cases
#   bin/test --test-cases 'TestModel#test_users_validation,TestModel#test_create_post'
#
#   # Run all tests in directories
#   bin/test test/subscriptions test/orders/coupons
class WingCommanderRunner
  # @param test_file_patterns [Array<String>] Glob patterns for matching test files.
  #   Default: ["test/**/*_test.rb", "test/**/test_*.rb"]
  # @param lib_dirs [Array<String>] Directories to add to Ruby's load path (-I flag).
  #   Default: %w[lib test .]
  def initialize(test_file_patterns: ["test/**/*_test.rb", "test/**/test_*.rb"], lib_dirs: %w[lib test .])
    @test_file_patterns = test_file_patterns
    @lib_dirs = lib_dirs
  end

  # Processes command line arguments and executes the minitest command.
  #
  # @param argv [Array<String>] Command line arguments. Default: ARGV
  #   Accepts:
  #   - Test file paths: "test/worker_test.rb"
  #   - Directory paths: "test/models"
  #   - Specific test cases: "TestModel#test_users_validation,TestModel#test_create_post"
  #   If empty, runs all tests matching test_file_patterns.
  #
  # @return [void] Executes the command (does not return)
  def call(test_files = ARGV, test_cases: [])
    builder = CommandBuilder.new(test_file_patterns: @test_file_patterns, lib_dirs: @lib_dirs)
    command = builder.call(test_files, test_cases: test_cases)
    exec command
  end

  private

  class CommandBuilder
    def initialize(test_file_patterns: ["test/**/*_test.rb", "test/**/test_*.rb"], lib_dirs: %w[lib test .])
      @test_file_patterns = test_file_patterns
      @lib_dirs = lib_dirs
    end

    def call(raw_test_paths = ARGV, test_cases: [])
      if raw_test_paths.any? && test_cases.any?
        raise "Cannot specify both test file paths and specific test cases"
      end

      if test_cases.any?
        return generate_test_command(test_cases:)
      end

      processed_test_paths = []

      raw_test_paths.each do |raw_test_path|
        next if process_directory_path(raw_test_path, processed_test_paths)
        next if process_file_path(raw_test_path, processed_test_paths)

        raise "Not a valid test path: #{raw_test_path}. Valid patterns are: tests/, test/some_test.rb"
      end

      return generate_test_command(test_files: processed_test_paths)
    end

    def generate_test_command(test_files: [], test_cases: [])
      if test_files.any? && test_cases.any?
        raise "Cannot specify both test files and test cases"
      end

      # Build command arguments
      cmd_args = []
      cmd_args << "-I#{@lib_dirs.join(File::PATH_SEPARATOR)}" unless @lib_dirs.empty?
      cmd_args << "-w" # Enable warnings

      if test_files.any?
        # Construct runner code: require minitest/autorun and all test files
        requires = test_files.map { |f| %(require "#{f}") }
        runner = ["require \"minitest/autorun\"", *requires].join("; ")

        cmd_args << "-e" # flag to execute code
        cmd_args << "'#{runner}'" # Single quotes like minitest/test_task.rb

        shell_escaped_test_files = test_files.map{ "'#{_1}'" }
        space_delimited_test_files = shell_escaped_test_files.join(' ')

        return "ruby #{cmd_args.join(' ')} #{space_delimited_test_files}"
      elsif test_cases.any?
        all_test_files = @test_file_patterns.flat_map { |pattern| Dir[pattern] }.select { |f| File.file?(f) }

        if all_test_files.empty?
          raise "No test files found matching patterns: #{@test_file_patterns.join(', ')}"
        end

        requires = all_test_files.map { |f| %(require "#{f}") }
        runner = ["require \"minitest/autorun\"", *requires].join("; ")

        cmd_args << "-e"
        cmd_args << "'#{runner}'"
        # TODO not sure why this -- is required. If removed we get an error. hashbangs
        # between the test class name and case name may be problematic since we
        # place these into a regex below for minitest. Perhpas this somehow corretly escapes them.
        cmd_args << "--"
        cmd_args << "-n"
        pattern = "/#{test_cases.join("|")}/"
        cmd_args << "'#{pattern}'"

        return "ruby #{cmd_args.join(' ')}"
      else # run all tests in the project using the test file globs
        all_test_files = @test_file_patterns.flat_map { |pattern| Dir[pattern] }.select { |f| File.file?(f) }

        if all_test_files.empty?
          raise "No test files found matching patterns: #{@test_file_patterns.join(', ')}"
        end

        requires = all_test_files.map { |f| %(require "#{f}") }
        runner = ["require \"minitest/autorun\"", *requires].join("; ")

        cmd_args << "-e"
        cmd_args << "'#{runner}'"

        return "ruby #{cmd_args.join(' ')}"
      end
    end

    private

    def valid_test_file?(file_path)
      return false unless File.file?(file_path)

      @test_file_patterns.any? do |pattern|
        File.fnmatch(pattern, file_path, File::FNM_PATHNAME)
      end
    end

    def process_directory_path(arg, test_file_paths)
      if Dir.exist?(arg)
        expanded = Dir[File.join(arg, "**", "*")]

        matching_test_files = expanded.select { valid_test_file?(_1) }
        test_file_paths.concat(matching_test_files)
        true
      else
        false
      end
    end

    def process_file_path(arg, test_file_paths)
      if valid_test_file?(arg)
        test_file_paths << arg
        return true
      end

      return false
    end
  end
end
