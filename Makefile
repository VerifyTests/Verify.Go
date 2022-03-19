$(VERBOSE).SILENT:

run-build:
	echo "Building..."
	go build ./...

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

run-integration-test:
	echo "Running integration tests for 'diff'..."
	go build -tags=integration ./...
	go test -v ./... -tags=integration

all: run-build run-vet run-tests