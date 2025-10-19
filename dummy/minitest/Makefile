dev:
	go build -o bin/wing_commander ./cmd/wing_commander
dev-minitest: dev
	cd dummy/minitest && ../../bin/wing_commander run --project-path . --test-command "ruby -Ilib:test test/thing_test.rb && cat test/reports/TEST-ThingTest.xml"
