# 1. Builder Stage: Build the Go binary
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies first, for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application, creating a static binary
# Using CGO_ENABLED=0 to build a statically linked binary without C dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /auth ./cmd/server/main.go


# 2. Final Stage: Create the final, minimal image
FROM alpine:3.18

# It's good practice to run as a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Set the working directory
WORKDIR /app

# Copy the static binary from the builder stage
COPY --from=builder /auth .

# Expose the port the application will run on
EXPOSE 8080

# The command to run the application
CMD ["./auth"]
