# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to leverage Docker's build cache
COPY go.mod go.sum ./ 
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
# CGO_ENABLED=0 is important for static binaries, making the final image smaller
# GOOS=linux ensures the binary is built for Linux, compatible with alpine
# -a -installsuffix cgo are flags for static linking
# -o app specifies the output binary name
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app main.go

# Stage 2: Create a smaller, production-ready image
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/app .

# Copy the built frontend files from the host into the container
# Assuming frontend/dist is relative to the Dockerfile context
COPY ../frontend/dist ./frontend/dist

# Create the uploads directory. This directory should ideally be mounted as a volume
# in production to persist uploaded files.
RUN mkdir -p uploads

# Expose the port the application listens on
EXPOSE 8080

# Define the command to run the application
CMD ["./app"]