# Start with a base image that includes Go runtime
FROM golang:1.17-alpine AS builder

# Set necessary environment variables
ENV GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o web-crawler .

# Start a new stage from scratch
FROM alpine:latest  

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/web-crawler .

# Copy the .env file from the host to the Docker image
COPY .env ./.env

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./web-crawler"]
