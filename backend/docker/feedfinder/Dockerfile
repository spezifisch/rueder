FROM alpine:3.15.4 AS deploy

## build stage
FROM golang:1.18.1-alpine3.15 as build

# cache go deps
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download && \
    # install swag separately so the part above can be cached for all images
    go install github.com/swaggo/swag/cmd/swag@v1.8.1

# compile
COPY . .
RUN swag init -g ./pkg/feedfinder/http/api.go --exclude ./pkg/api,./pkg/authbackend,./pkg/events && \
    go build -ldflags="-s -w" -o feedfinder ./cmd/feedfinder

## deploy stage
FROM deploy as final

WORKDIR /app
COPY --from=build /build/feedfinder .

EXPOSE 8080
CMD ["./feedfinder"]
