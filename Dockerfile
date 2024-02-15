# Copyright (c) 2020, Dobrosław Żybort
# SPDX-License-Identifier: BSD-3-Clause

FROM golang:1.21 as builder

# Set environment variables for build process
ENV \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/app

# Cache dependencies
COPY go.mod go.sum ./
RUN \
    go mod download && \
    go mod verify

COPY . .

# Run tests so we don't build app with failing tests
RUN go test ./...

RUN \
    VERSION=$(git describe --tags --dirty --always) && \
    COMMIT=$(git rev-parse --short HEAD) && \
    BUILDTIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") && \
    go build -ldflags="-s -w \
        -X main.buildType=df \
        -X main.version=${VERSION} \
        -X main.commit=${COMMIT} \
        -X main.buildTime=${BUILDTIME}" \
        -o app

# Second build step
FROM gcr.io/distroless/base-debian10

LABEL \
    org.opencontainers.image.ref.name="matrixik/sort-camera-photos" \
    org.opencontainers.image.description="No configuration camera photos sorting" \
    org.opencontainers.image.authors="Dobrosław Żybort <matrixik@gmail.com>" \
    org.opencontainers.image.documentation="https://github.com/matrixik/sort-camera-photos/blob/master/README.md" \
    org.opencontainers.image.licenses="BSD-3-Clause" \
    org.opencontainers.image.source="https://github.com/matrixik/sort-camera-photos" \
    org.opencontainers.image.url="https://hub.docker.com/r/matrixik/sort-camera-photos/"

WORKDIR /
COPY --from=builder /go/src/app/app .
ENTRYPOINT ["/app"]
