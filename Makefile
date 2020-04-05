GOFILES = $(shell find . -name '*.go')
GOPACKAGES = $(shell go list ./...)

default: build

build:
	go build -o bin/rps_financial main.go

run:
	go run main.go

test: test-all

test-all:
	export RPS_DB_USERNAME='root'
	export RPS_DB_PASSWORD='my-secret-pw'
	export RPS_DB_HOST_PORT='localhost:3306'
	export RPS_DB_NAME='rps'
	@go test -v $(GOPACKAGES)

build-docker-image:
	docker build -f deploy/Dockerfile -t rps_financial .
