.PHONY: all
all: docs

.PHONY: download
download:
	go mod download

.PHONY: test
test:
	go test ./... -parallel=8 -count=1

.PHONY: test_record
test_record:
	COPYIST_RECORD=1 go test ./... -count=1

.PHONY: lint
lint:
	golangci-lint run

.PHONY: ci
ci: lint test

.PHONY: docs
docs:
	swag init -g pkg/api/http/api-v1.go
