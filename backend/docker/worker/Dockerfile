FROM alpine:3.15.6 AS deploy

## build stage
FROM golang:1.18.5-alpine3.15 as build

# cache go deps
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

# compile
COPY . .
RUN go build -ldflags="-s -w" -o worker ./cmd/worker

## deploy stage
FROM deploy as final

WORKDIR /app
COPY --from=build /build/worker .
COPY ./config/database.yml .

CMD ["./worker"]
