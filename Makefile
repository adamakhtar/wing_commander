dev:
	go build -o bin/wing_commander
test:
	go test ./...
dev-minitest: dev
	./bin/wing_commander start dummy/minitest_example --run-command "rake test && cat test/reports/TEST-ThingTest.xml" --test-file-pattern "test/*_test.rb" --debug
