FROM golang:1.26.1-alpine AS builder

RUN apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum Makefile ./
RUN go mod download && make setup

COPY . .
RUN make generate
RUN go build -o server ./cmd/server/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
