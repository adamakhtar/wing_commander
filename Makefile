dev:
	go build -o bin/wing_commander ./cmd/wing_commander
dev-minitest: dev
	cd dummy/minitest && ../../bin/wing_commander run --project-path . --test-command "./run_tests.sh"
