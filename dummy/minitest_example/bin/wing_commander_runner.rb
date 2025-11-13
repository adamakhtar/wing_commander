# frozen_string_literal: true

require "shellwords"


# Builds and executes minitest commands based on command line arguments.
#
# Processes test file patterns, directories, and specific test cases, then constructs
# and executes the appropriate ruby command to run minitest.
#
# @example Basic usage
#   runner = WingCommanderRunner.new
#   runner.call(["test/worker_test.rb"])
#
# @example With custom patterns
#   runner = WingCommanderRunner.new(
#     test_file_patterns: ["spec/**/*_spec.rb"],
#     lib_dirs: %w[lib spec]
#   )
#   runner.call(["spec/models/user_spec.rb"])
#
# @example Command line usage patterns
#   # Run all tests matching default glob pattern
#   bin/test
#
#   # Run specific test file
#   bin/test test/post_test.rb
#
#   # Run specific test cases
#   bin/test test/post_test.rb:test_users_validation test/post_test.rb:test_create_post
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
  #   - Specific test cases: "test/worker_test.rb:test_success"
  #   If empty, runs all tests matching test_file_patterns.
  #
  # @return [void] Executes the command (does not return)
  def call(argv = ARGV)
    builder = CommandBuilder.new(test_file_patterns: @test_file_patterns, lib_dirs: @lib_dirs)
    command = builder.call(argv)
    exec command
  end

  private

  class CommandBuilder
    def initialize(test_file_patterns: ["test/**/*_test.rb", "test/**/test_*.rb"], lib_dirs: %w[lib test .])
      @test_file_patterns = test_file_patterns
      @lib_dirs = lib_dirs
    end

    def call(argv = ARGV)
      test_file_paths = []
      specific_test_cases = []

      argv.each do |arg|
        next if process_specific_test_case(arg, specific_test_cases)
        next if process_directory_path(arg, test_file_paths)
        next if process_file_path(arg, test_file_paths)

        raise "Not a valid test pattern: #{arg}. Valid patterns are: tests/, test/some_test.rb, test/some_test.rb:test_name"
      end

      if test_file_paths.any? && specific_test_cases.any?
        raise "Cannot specify both test file paths and specific test cases"
      end

      if specific_test_cases.any?
        generate_test_command(
          test_files: specific_test_cases.map { |test_case| test_case[:file_path] },
          test_case_names: specific_test_cases.map { |test_case| test_case[:test_case_name] }
        )

      else
        if test_file_paths.empty?
          generate_test_command(test_files: Dir[*@test_file_patterns].sort.uniq, test_case_names: [])
        else
          generate_test_command(test_files: test_file_paths.map { |test_file| test_file[:path] })
        end
      end
    end

    def generate_test_command(test_files:, test_case_names: [])
      # Construct runner code: require minitest/autorun and all test files
      requires = test_files.map { |f| %(require "#{f}") }
      runner = ["require \"minitest/autorun\"", *requires].join("; ")

      # Build command arguments
      cmd_args = []
      cmd_args << "-I#{@lib_dirs.join(File::PATH_SEPARATOR)}" unless @lib_dirs.empty?
      cmd_args << "-w" # Enable warnings
      cmd_args << "-e"

      cmd_args << "'#{runner}'" # Single quotes like minitest/test_task.rb

      if test_case_names.any?
        # TODO not sure why this -- is required. If removed we get an error. hashbangs
        # between the test class name and case name may be problematic since we
        # place these into a regex below. Perhpas this somehow corretly escapes them.
        cmd_args << "--"
        cmd_args << "-n"
        pattern = "/#{test_case_names.join("|")}/"
        cmd_args << "'#{pattern}'"
      end

      "ruby #{cmd_args.join(' ')}"
    end

    private

    def valid_test_file?(file_path)
      return false unless File.exist?(file_path)

      @test_file_patterns.any? do |pattern|
        File.fnmatch(pattern, file_path, File::FNM_PATHNAME)
      end
    end

    def process_specific_test_case(arg, specific_test_cases)
      if arg.include?(":")
        file_path, test_case_name = arg.split(":")

        raise "Specific test case not found for #{arg}. #{file_path} does not exist" if !valid_test_file?(file_path)

        specific_test_cases << {file_path:, test_case_name:}
        true
      else
        false
      end
    end

    def process_directory_path(arg, test_file_paths)
      if Dir.exist?(arg)
        expanded = Dir[arg]

        matching_test_files = expanded.select { valid_test_file?(_1) }
        matching_test_files.each { test_file_paths << {path: _1} }
        true
      else
        false
      end
    end

    def process_file_path(arg, test_file_paths)
      if File.exist?(arg)
        test_file_paths << {path: arg}
        true
      else
        false
      end
    end
  end
end
