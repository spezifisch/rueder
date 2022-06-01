all: docs

download:
	go mod download

# run unit tests without database connection
test:
	go test ./... -parallel=8 -count=1

# run unit tests and let copyist record
test_record:
	COPYIST_RECORD=1 go test ./... -count=1

lint:
	golangci-lint run

ci: lint test

# generate api docs
docs:
	swag init -g pkg/api/http/api-v1.go

# some docker dev shortcuts
dockerargs = --force-recreate --build

run:
	docker-compose up $(dockerargs)

run_noworker:
	docker-compose up $(dockerargs) auth authbackend db dev

run_worker:
	docker-compose up $(dockerargs) worker

run_feedfinder:
	docker-compose up $(dockerargs) feedfinder

run_events:
	docker-compose up $(dockerargs) events

.PHONY: *
