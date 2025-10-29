dev:
	go build -o bin/wing_commander
test:
	go test ./...
dev-minitest: dev
	cd dummy/minitest && ../../bin/wing_commander run --project-path . --test-command "./run_tests.sh"
