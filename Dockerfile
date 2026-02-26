# Build stage
FROM golang:1.24.2-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Optimize caching: Download dependencies before copying source code
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server cmd/api/main.go

# Final stage
FROM alpine:3.21

# Set production environment variables
ENV GIN_MODE=release \
    ENV=production


# Install runtime dependencies and create a non-root user for security
RUN apk add --no-cache ca-certificates && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/server .

# Use non-root user to run the application
USER appuser

EXPOSE 8080

CMD ["./server"]