dev:
	go build -o bin/wing_commander
test:
	go test ./...
dev-minitest: dev
	./bin/wing_commander start dummy/minitest_example --run-command "bundle exec rake test" --test-file-pattern "test/*_test.rb" --test-results-path ".wing_commander/test_results/summary.yml" --debug
