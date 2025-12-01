# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/whoisthere main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS and whois
RUN apk add --no-cache ca-certificates whois

# Copy binary from builder
COPY --from=builder /app/bin/whoisthere .

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["./whoisthere"]

