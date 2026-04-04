FROM golang:1.26.1-alpine AS builder

RUN apk add --no-cache make wget

# Use /repo/go-backend to mirror repository structure so that
# relative paths (../openapi, ../proto, ../buf.gen.yaml) resolve correctly.
WORKDIR /repo/go-backend

# Copy workspace manifests first for layer caching
COPY go-backend/go.mod go-backend/go.sum go-backend/Makefile ./
RUN go mod download && make setup

# Install buf (required by generate-proto)
RUN BUF_VERSION=1.66.1 && \
    wget -q -O /usr/local/bin/buf \
      "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-Linux-x86_64" && \
    chmod +x /usr/local/bin/buf

# Copy source and generation inputs
COPY go-backend/ .
COPY openapi/ /repo/openapi/
COPY proto/ /repo/proto/
COPY buf.gen.yaml /repo/

RUN make generate
RUN go build -o server ./cmd/server/main.go

FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -S appuser && adduser -S appuser -G appuser

WORKDIR /app
COPY --from=builder --chown=appuser:appuser /repo/go-backend/server .

USER appuser

EXPOSE 8080
CMD ["./server"]
