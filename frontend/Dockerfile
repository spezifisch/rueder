## BASE
FROM node:18.17.1-alpine3.18 AS base
# TODO migrate to pnpm

# disable npm update check, it's annoying
ENV NPM_CONFIG_UPDATE_NOTIFIER=false

WORKDIR /usr/src/app

# install dependencies first for better layer caching
COPY package*.json ./
RUN npm ci --only=production --unsafe-perm --no-audit

# copy sources
COPY *.config.*js *.d.ts index.html tsconfig.json ./
COPY ./src ./src
COPY ./public ./public

## TEST STAGE
FROM base as test

# run tests to ensure we're not making a defective prod build
RUN npm ci --unsafe-perm --no-audit \
    && npm run test

## BUILD STAGE
FROM base as build

ARG VITE_BASE_URL=/
ARG VITE_RUEDER_BASE_URL_LOGIN
ARG VITE_RUEDER_BASE_URL_API
ARG VITE_RUEDER_BASE_URL_IMGPROXY

# bundle the app
RUN echo Base URLs: base=$VITE_BASE_URL login=$VITE_RUEDER_BASE_URL_LOGIN api=$VITE_RUEDER_BASE_URL_API imgproxy=$VITE_RUEDER_BASE_URL_IMGPROXY && \
    npm run build -- --base="$VITE_BASE_URL"

## FINAL STAGE
# serve with nginx on :8080
FROM nginxinc/nginx-unprivileged:1.25.1-alpine AS final

COPY --from=build /usr/src/app/dist /usr/share/nginx/html/
