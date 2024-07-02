# Start from the base golang image
FROM golang:1.18-alpine

# Install necessary tools
RUN apk add --no-cache gcc musl-dev

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Copy test data if building for testing
ARG ENV=production
COPY testdata /app/testdata

# Ensure test directory is copied
COPY test /app/test

# Ensure .env file is copied
COPY .env /app/.env

# Set the working directory to cmd/server for the build
WORKDIR /app/cmd/server

# Run tests only if in testing environment
RUN if [ "$ENV" = "test" ]; then go test -v /app/test/...; fi

# Build the Go app from the cmd/server directory
RUN go build -o /app/main .

# Set the working directory back to /app
WORKDIR /app

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
