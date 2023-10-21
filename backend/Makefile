COMPOSE=podman-compose
dockerargs = --force-recreate --build

all: docs

download:
	go mod download

# run unit tests without database connection
test:
	go test ./... -parallel=4 -count=1

# run unit tests and let copyist record
test_record:
	COPYIST_RECORD=1 go test ./... -count=1

lint:
	golangci-lint run

ci: lint test

# generate api docs
docs:
	swag init -g pkg/api/http/api-v1.go

run:
	$(COMPOSE) up $(dockerargs)

.PHONY: *
