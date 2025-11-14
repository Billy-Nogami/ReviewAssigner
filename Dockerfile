# Build stage
FROM golang:1.25.4-alpine AS builder

# Install migrate tool
RUN apk add --no-cache git
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates postgresql-client
WORKDIR /root/

# Copy migrate binary from builder
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Copy application binary
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Health check (optional)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]