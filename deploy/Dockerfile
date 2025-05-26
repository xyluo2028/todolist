# Dockerfile
# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the application source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -v -o todolist-server ./cmd/server/main.go

# Stage 2: Create the final lightweight image
FROM alpine:latest

WORKDIR /app/

# Copy the built binary from the builder stage
COPY --from=builder /app/todolist-server .

# Expose port 7071 (the port the app listens on)
EXPOSE 7071

# Command to run the application
CMD ["./todolist-server"]