# Copyright (c) 2020, Dobrosław Żybort
# SPDX-License-Identifier: BSD-3-Clause

FROM golang:1.14 as builder

# Set environmet variables for build process
ENV \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/app
COPY . .
RUN go get -d -t -v ./...
RUN go mod verify
# Run tests so we don't build app with failing tests
RUN go test ./...

RUN go build -o app

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=builder /go/src/app/app .
ENTRYPOINT ["/app"]
