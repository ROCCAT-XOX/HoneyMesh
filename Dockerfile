FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache gcc musl-dev

# Copy go.mod and go.sum and download dependencies first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o honeymesh .

# Use a smaller image for the final container
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/honeymesh /app/honeymesh
# Copy template files and assets
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/assets /app/assets

# Set environment variables
ENV GIN_MODE=release

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/honeymesh"]