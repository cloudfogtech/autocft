# syntax=docker/dockerfile:1.7
# Multi-stage build for autocft
# Supports building from source (default) and suitable for multi-arch buildx usage.
# The release workflow will use buildx to produce multi-arch images.

ARG GO_VERSION=1.25

############################
# Builder stage
############################
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

# Install build dependencies (git for go modules if needed)
RUN apk add --no-cache git ca-certificates && update-ca-certificates

# Leverage caching of go mod download steps
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source
COPY . .

# Build arguments for cross compilation via buildx
ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG TARGETVARIANT

# Enable static build (modernc sqlite is pure go, so CGO can remain disabled)
ENV CGO_ENABLED=0

# Output binary path
RUN --mount=type=cache,target=/root/.cache/go-build \
    /bin/sh -c ' \
      if [ "$TARGETARCH" = "arm" ]; then export GOARM=${TARGETVARIANT#v}; fi; \
      GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -trimpath -ldflags="-s -w" -o /out/autocft ./cmd/autocft \
    '

############################
# Final minimal image (Alpine to support 386 + CA certs)
############################
FROM alpine:3.20 AS final
# Add non-root user
RUN adduser -D -H -u 10001 appuser
WORKDIR /app
COPY --from=builder /out/autocft /app/autocft
VOLUME ["/app/pb_data"]
USER appuser
ENTRYPOINT ["/app/autocft"]

# Example multi-arch build:
# docker buildx build --platform linux/amd64,linux/386,linux/arm64,linux/arm/v7 \
#   -t your/repo:latest --push .
