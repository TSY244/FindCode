# Use Go 1.21 as the base image for building
FROM golang:1.23-alpine AS builder

# Install necessary C libraries and tools
RUN apk add --no-cache gcc musl-dev sqlite

# Set the working directory
WORKDIR /build

# Set Go proxy
RUN go env -w GOPROXY=https://goproxy.cn,direct

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o FindCodeServer  cmd/default/main.go

# Use a minimal base image for the final container
FROM alpine:latest

# Install sqlite3 for runtime
RUN apk add --no-cache sqlite

WORKDIR /app
COPY --from=builder /build/FindCodeServer .
COPY --from=builder /build/etc etc/
COPY --from=builder /build/rule rule/
COPY --from=builder /build/web web/
COPY --from=builder /build/script script/
COPY --from=builder /build/prompt prompt/

# Expose the necessary ports
EXPOSE 8080

# Run the application
CMD ["./FindCodeServer","-mode","server"]