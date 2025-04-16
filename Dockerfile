# Build stage
FROM golang:1.24.2-bookworm AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o ubiquiti-monitor

# Final stage
FROM debian:bookworm

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/ubiquiti-monitor .

# Expose the application port
EXPOSE 8080

# Command to run the executable
CMD ["./ubiquiti-monitor"] 