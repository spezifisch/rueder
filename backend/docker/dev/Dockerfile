FROM golang:1.19.12-alpine3.18
ARG S6_OVERLAY_VERSION=3.1.5.0

# set default timezone
ARG TZ=Europe/Berlin
ENV DEFAULT_TZ ${TZ}
# go config
ENV CGO_ENABLED=1
# s6-overlay config
ENV S6_KEEP_ENV=1

WORKDIR /app

ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch.tar.xz /tmp
RUN tar -C / -Jxpf /tmp/s6-overlay-noarch.tar.xz
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-x86_64.tar.xz /tmp
RUN tar -C / -Jxpf /tmp/s6-overlay-x86_64.tar.xz
ENTRYPOINT ["/init"]

RUN apk add --no-cache \
    gcc=12.2.1_git20220924-r10 \
    musl-dev=1.2.4-r1 \
    postgresql13-client=13.11-r0 \
    tzdata=2023c-r1 && \
    ln -snf /usr/share/zoneinfo/${DEFAULT_TZ} /etc/localtime && \
    # Go Watch: https://github.com/mitranim/gow
    go install github.com/mitranim/gow@v0.0.0-20230418123246-87df6e48eec6 && \
    # Soda CLI: https://gobuffalo.io/en/docs/db/toolbox/
    go install github.com/gobuffalo/pop/v6/soda@v6.0.8 && \
    # install swag
    go install github.com/swaggo/swag/cmd/swag@v1.8.12

# cache go deps
COPY go.mod go.sum ./
RUN go mod download

# s6 daemon configs
COPY ./docker/dev/s6-rc.d /etc/s6-overlay/s6-rc.d

VOLUME /app
EXPOSE 8080
