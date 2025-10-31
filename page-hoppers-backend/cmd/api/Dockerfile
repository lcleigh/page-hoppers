# Use Go image directly (no multi-stage build)
FROM golang:1.25-alpine

# Install git (needed for go modules)
RUN apk add --no-cache git

# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for caching deps)
COPY ../../go.mod ../../go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire backend source code
COPY ../../ ./

# Expose port
EXPOSE 8080

# Run the app using Air for live reload
CMD ["air", "-c", ".air.toml"]