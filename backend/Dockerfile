# Use the latest official Golang image
FROM golang:1.22

# Set working directory inside the container
WORKDIR /app

# Download and install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2

# Copy go.mod and go.sum first to leverage Docker caching
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application and place the executable in /app directory
RUN go build -o /app/main cmd/main.go

# Ensure the executable has the correct permissions
RUN chmod +x /app/main

# Expose the port the application will run on
EXPOSE 8080

# Command to run the Go application
CMD ["/app/main"]
