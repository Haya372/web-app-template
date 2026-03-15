FROM golang:1.26.1-alpine AS builder

RUN apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum Makefile ./
RUN go mod download && make setup

COPY . .
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
