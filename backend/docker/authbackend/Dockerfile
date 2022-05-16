FROM alpine:3.15.4 AS deploy

## build stage
FROM golang:1.18.1-alpine3.15 as build

# cache go deps
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download && \
    # Soda CLI: https://gobuffalo.io/en/docs/db/toolbox/
    go install github.com/gobuffalo/pop/soda@v4.13.1+incompatible

# compile
COPY . .
RUN go build -ldflags="-s -w" -o authbackend ./cmd/authbackend

## deploy stage
FROM deploy as final

# add pg_dump for db migrations
RUN apk add --no-cache postgresql13-client=13.6-r0

WORKDIR /go/bin
COPY --from=build /go/bin/soda .

WORKDIR /app
COPY --from=build /build/authbackend .
COPY ./config/database.yml .

# soda migration tools
COPY ./utils/initdb-production.sh ./initdb.sh
COPY ./migrations ./migrations

EXPOSE 8080
CMD ["./authbackend"]
