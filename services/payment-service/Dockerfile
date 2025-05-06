# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY src/go.mod src/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY src/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o payment-service .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/payment-service .

# Copy Swagger UI files
COPY docs/swagger.html ./docs/swagger.html
COPY docs/swagger.yaml ./docs/swagger.yaml

# Expose port
EXPOSE 8083

# Run the application
CMD ["./payment-service"]
