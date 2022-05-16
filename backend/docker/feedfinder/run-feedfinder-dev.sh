#!/bin/sh -e

rm -f docs/docs.go docs/swagger*
swag init -g ./pkg/feedfinder/http/api.go --exclude ./pkg/api,./pkg/authbackend
exec gow run ./cmd/feedfinder --dev --log debug
