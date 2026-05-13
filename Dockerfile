# ============================================
# Base stage - Common dependencies
# ============================================
FROM golang:1.25-alpine AS base

WORKDIR /app

# Install git for Go modules
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .


# ============================================
# Development stage - With Air hot reload
# ============================================
FROM base AS development

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
ENV PATH="/go/bin:${PATH}"

# Expose port
EXPOSE 3000

# Run with Air
CMD ["air", "-c", ".air.toml"]


# ============================================
# Builder stage - Build the binary
# ============================================
FROM base AS builder

# Build the server application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o api \
    ./cmd/api

# Build the consumer application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o consumer \
    ./cmd/consumer

# Build the runner application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o runner \
    ./cmd/runner

# ============================================
# Production stage - Minimal runtime
# ============================================
FROM alpine:latest AS production

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /home/appuser

# Copy binaries from builder
COPY --from=builder --chown=appuser:appuser /app/api .
COPY --from=builder --chown=appuser:appuser /app/consumer .
COPY --from=builder --chown=appuser:appuser /app/runner .

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 3000

# Run the server application by default
CMD ["./api"]