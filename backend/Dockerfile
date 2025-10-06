# Stage 1: Build the Go application
FROM golang:1.25-alpine AS builder

# Install CGO build dependencies for WebP
RUN apk add --no-cache gcc musl-dev libwebp-dev

# Set the working directory to the root of our app space
WORKDIR /app

# Copy go.mod and go.sum from the project root
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy the entire backend source code into a 'backend' sub-directory
COPY backend/ ./backend/

# Build the Go application from within the /app directory
# The main package is now at ./backend/main.go
RUN go build -a -o /app/app ./backend/main.go

# Stage 2: Create a smaller, production-ready image
FROM alpine:latest

# Set working directory for the final image
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/app .

# Copy the built frontend files from the build context
COPY frontend/dist ./frontend/dist

# Create the uploads directory
RUN mkdir -p uploads

# Expose the port
EXPOSE 8080

# Define the command to run the application
CMD ["./app"]
