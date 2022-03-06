$(VERBOSE).SILENT:

verifier-build:
	echo "Building 'verifier'..."
	go build ./verifier/*.go

verifier-vet:
	echo "Vetting 'verifier'..."
	go vet ./verifier

verifier-test:
	echo "Running unit tests for 'verifier'..."
	go test -v ./verifier/*.go
	echo "Running api tests for 'verifier'..."
	go test -v ./api-tests/*.go

diff-build:
	echo "Building 'diff'..."
	go build ./diff/*.go

diff-vet:
	echo "Vetting 'diff'..."
	go vet ./diff

diff-test:
	echo "Running unit tests for 'diff'..."
	go test -v ./diff

diff-integration-test:
	echo "Running integration tests for 'diff'..."
	go build -tags=integration ./diff
	go test -v ./diff -tags=integration
	

verifier-all: verifier-build verifier-vet verifier-test
diff-all: diff-build diff-vet diff-test

build-all: diff-build verifier-build
vet-all: diff-vet verifier-vet
unit-test-all: diff-test verifier-test

all: build-all vet-all unit-test-all
