FROM alpine:3.15.4 AS deploy

## build stage
FROM golang:1.18.1-alpine3.15 as build

# cache go deps
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

# compile
COPY . .
RUN go build -ldflags="-s -w" -o rueder-events ./cmd/events

## deploy stage
FROM deploy as final

WORKDIR /app
COPY --from=build /build/rueder-events .
COPY ./config/database.yml .

EXPOSE 8080
CMD ["./rueder-events"]
