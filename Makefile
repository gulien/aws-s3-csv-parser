.PHONY: fmt lint tests todos godoc run build

# Formats the code, "optimizes" the modules' dependencies.
fmt:
	go fmt ./...
	go mod tidy

# Runs linters.
lint:
	golangci-lint run

# Runs tests.
tests:
	go test -race -covermode=atomic -coverprofile=coverage.txt ./... &&\
	go tool cover -html=coverage.txt -o coverage.html

# Shows TODOs.
todos:
	golangci-lint run \
	--no-config \
	--disable-all \
	--enable godox

# Runs a local webserver for godoc.
godoc:
	$(info http://localhost:6060/pkg/github.com/gulien/aws-s3-csv-parser)
	godoc -http=:6060

# Runs the application.
REGION=us-west-2
BUCKET=work-sample-mk
KEY=2021/04/events.csv
TIMEOUT=300

run:
	go run cmd/aws-s3-csv-parser/main.go --region=$(REGION) --bucket=$(BUCKET) --key=$(KEY) --timeout=$(TIMEOUT)


# Builds the application.
VERSION=snapshot

build:
	go build -ldflags "-X main.version=$(VERSION)" -o=aws-s3-csv-parser cmd/aws-s3-csv-parser/main.go