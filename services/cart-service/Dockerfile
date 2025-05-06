FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o cart-service .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/cart-service .

# Copy Swagger UI files
COPY docs/swagger.html ./docs/swagger.html
COPY docs/swagger.yaml ./docs/swagger.yaml

EXPOSE 8087

CMD ["./cart-service"]
