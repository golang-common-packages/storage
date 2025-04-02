# Build stage
FROM golang:1.20-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# Enable CGO for SQLite support
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o storage .

# Final stage
FROM alpine:latest

# Install necessary packages
RUN apk --no-cache add ca-certificates sqlite-libs

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/storage .

# Command to run
CMD ["./storage"]
