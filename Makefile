$(VERBOSE).SILENT:

run-build:
	echo "Building..."
	go build ./...

run-cleanup:
	go mod tidy

run-vet:
	echo "Vetting..."
	go vet ./...

run-tests:
	echo "Running unit tests for 'verifier'..."
	go test -v ./verifier/*.go

	echo "Running api tests for 'verifier'..."
	go test -v ./api-tests/*.go

	echo "Running unit tests for 'diff'..."
	go test -v ./diff/*.go

	echo "Running unit tests for 'utils'..."
	go test -v ./utils/*.go

	echo "Running unit tests for 'tray'..."
	go test -v ./tray/*.go

run-integration-test: export RUN_INTEGRATION_TESTS=True
run-integration-test:
	echo "Should run integration tests: $$RUN_INTEGRATION_TESTS"
	echo "Running integration tests..."
	go test -v ./... -run Integration

all: run-build run-vet run-tests