# Multi-stage build for minimal final image
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o archer \
    ./cmd/archer

# Final stage - minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 archer && \
    adduser -D -u 1000 -G archer archer

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/archer /usr/local/bin/archer

# Copy templates
COPY --from=builder /build/templates /app/templates

# Use non-root user
USER archer

# Set entrypoint
ENTRYPOINT ["archer"]

# Default command
CMD ["--help"]
