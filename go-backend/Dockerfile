FROM golang:1.26.1-alpine AS builder

RUN apk add --no-cache make

WORKDIR /app

# Copy workspace manifests first for layer caching
COPY go-backend/go.mod go-backend/go.sum go-backend/Makefile ./
RUN go mod download && make setup

# Copy source and openapi spec (oapi-codegen needs ../openapi relative to WORKDIR)
COPY go-backend/ .
COPY openapi/ /openapi/

RUN make generate
RUN go build -o server ./cmd/server/main.go

FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -S appuser && adduser -S appuser -G appuser

WORKDIR /app
COPY --from=builder --chown=appuser:appuser /app/server .

USER appuser

EXPOSE 8080
CMD ["./server"]
