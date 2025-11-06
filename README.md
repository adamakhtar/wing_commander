# Heads up

Currently WIP and not in a workable state.

# About the code

I started this project to

1. scratch an itch (see below)
2. learn Go and TUI development
3. put LLMs through their paces and see how far they can go

Regarding 3: LLMs start great but can veer off into sloppy code. Useful for learning Go, but needs cleanup. I'm now the primary driver, using LLMs more focused. Until I reach the first working version I'll be moving fast and commits may be a little messy.

# Wing Commander

A CLI/TUI tool for both running tests and analyzing their results

## Rationale / Impetus

Test frameworks often vomit up slabs of text when tests fail. Not offensive in isolation, but since we spend so much time hopping between code and tests, making them easier to run and parse will be nicer on our already tired minds.

Some problems with tests:

1. Running a single test file is simple, but running multiple (e.g., all controllers, models, and services related to Billing) is tedious. Wouldn't it be great to fuzzy search for tests or save commonly used groups?
2. Our tests fail because of our project's code, yet backtraces often include filepaths from 3rd party libraries and frameworks. It just gets in the way.
3. They often show long unwieldy absolute paths when shorter relative paths would work
4. They don't show you the code so you can quickly ground yourself.
5. They don't show which backtrace files you've changed since your last commit - the most likely culprits.

I'm aiming to solve all these with this tool.

Currently supports Ruby and minitest. Will expand to RSpec and perhaps JavaScript.

## Quick Start

### Development Testing

```bash
# Build dev version and launch TUI against dummy minitest app
make dev-minitest
```

### Production Usage

```bash
# Run tests and analyze failures
wing_commander run --project-path /path/to/project --test-command "rails test --output .wing_commander/test_output.xml"

# Show configuration
wing_commander config

```

## Supported Test Frameworks

- Minitest (Ruby) - in development

### Minitest setup

Install this project's companion Minitest reporter gem (Coming soon)

```ruby
# test/test_helper.rb
require 'minitest/reporters'
Minitest::Reporters.use! [
  WingCommanderReporter.new(summary_output_path: "....")
]
```

This produces a test run summary at the given path which this CLI tool will read. You will likely want to gitignore the summary file.

## Development

```bash
# Build and launch the cli tool against the dummy minitest suite in this repo
make dev-minitest

# Build development version
make dev

# Run tests
make test

# Clean build artifacts
make clean
```
