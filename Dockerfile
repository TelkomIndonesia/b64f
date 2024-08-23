# syntax = docker/dockerfile:1.2

FROM golang:1.23 AS builder

WORKDIR /src
COPY ./ ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o b64f



FROM alpine:3.20

COPY --from=builder /src/b64f /usr/local/bin/b64f
ENTRYPOINT ["/usr/local/bin/b64f"]