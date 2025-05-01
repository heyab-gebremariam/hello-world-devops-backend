# === Build Stage ===
# Use an official Go image as the builder environment
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# === Final Stage ===
# Use a minimal base image like alpine
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy only the built binary from the builder stage
COPY --from=builder /app/main .

# Copy migration files (needed if applying migrations from within the container)
COPY --from=builder /app/migrations ./migrations

# Copy Atlas config (needed if applying migrations from within the container)
COPY --from=builder /app/atlas.hcl .

# Expose the port the application runs on (matches Gin's r.Run)
EXPOSE 8090

CMD ["./main"]