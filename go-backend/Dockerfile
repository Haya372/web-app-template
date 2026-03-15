FROM golang:1.25.6-alpine AS builder

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

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
