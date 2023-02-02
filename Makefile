$(VERBOSE).SILENT:

run-build:
	echo "Building..."
	go build ./...

run-cleanup:
	go mod tidy

run-vet:
	echo "Vetting..."
	go vet ./...

test-verifier:
	go test -v ./verifier/*.go

test-api:
	go test -v ./api-tests/*.go

test-diff:
	go test -v ./diff/*.go

test-utils:
	go test -v ./utils/*.go

test-tray:
	go test -v ./tray/*.go

run-tests: test-verifier test-api test-diff test-utils test-tray

run-integration-test: export RUN_INTEGRATION_TESTS=True
run-integration-test:
	echo "Should run integration tests: $$RUN_INTEGRATION_TESTS"
	echo "Running integration tests..."
	go test -v ./... -run Integration

all: run-build run-vet run-tests