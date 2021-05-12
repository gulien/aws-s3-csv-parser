.PHONY: fmt lint todos godoc up down run build

# Formats the code, "optimizes" the modules' dependencies.
fmt:
	go fmt ./...
	go mod tidy

# Runs linters.
lint:
	golangci-lint run

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

# Starts the MySQL and phpMyAdmin Docker containers.
up:
	docker-compose up -d

# Stops the MySQL and phpMyAdmin Docker containers.
down:
	docker-compose down

# Runs the application.
REGION=us-west-2
BUCKET=work-sample-mk
KEY=2021/04/events.csv
TIMEOUT=300
SKIP_DOWNLOAD=0

run:
	MYSQL_DATABASE_URL='agent:secret@(127.0.0.1:3306)/events' go run cmd/aws-s3-csv-parser/main.go \
	--region=$(REGION) \
	--bucket=$(BUCKET) \
	--key=$(KEY) \
	--timeout=$(TIMEOUT) \
	--skip-download=$(SKIP_DOWNLOAD)

# Builds the application.
VERSION=snapshot

build:
	go build -ldflags "-X main.version=$(VERSION)" -o=aws-s3-csv-parser cmd/aws-s3-csv-parser/main.go