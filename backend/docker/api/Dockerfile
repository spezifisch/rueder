FROM alpine:3.15.6 AS deploy

## build stage
FROM golang:1.18.5-alpine3.15 as build

# cache go deps
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download && \
    # install swag separately so the part above can be cached for all images
    go install github.com/swaggo/swag/cmd/swag@v1.8.1

# compile
COPY . .
RUN swag init -g ./pkg/api/http/api-v1.go --exclude ./pkg/authbackend,./pkg/events,./pkg/feedfinder && \
    go build -ldflags="-s -w" -o api ./cmd/api

## deploy stage
FROM deploy as final

WORKDIR /app
COPY --from=build /build/api .
COPY ./config/database.yml .

EXPOSE 8080
CMD ["./api"]
