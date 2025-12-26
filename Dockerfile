# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.25

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates && update-ca-certificates
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/root/.cache/go-build \
    /bin/sh -c 'if [ "$TARGETARCH" = "arm" ]; then export GOARM=${TARGETVARIANT#v}; fi; \
      GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags="-s -w" -o /out/autocft ./cmd/autocft'

FROM --platform=$TARGETPLATFORM alpine:3.20 AS final
RUN adduser -D -H -u 10001 appuser
RUN mkdir -p /app/data && chown -R 10001:10001 /app/data && chmod -R 755 /app/data
WORKDIR /app
COPY --from=builder /out/autocft /app/autocft
VOLUME ["/app/data"]
USER appuser
ENTRYPOINT ["/app/autocft"]
