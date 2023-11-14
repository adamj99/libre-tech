# Use the official Go image from the Docker Hub
FROM golang:1.20 as builder

# Enable Go modules
ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Run the tests
RUN go test -v

# Build the application
# By setting CGO_ENABLED=0, we are telling Go to not include any C libraries in the binary. This is necessary because the final scratch image does not contain these libraries.
# GOOS=linux: This environment variable sets the target operating system for the build to Linux. This is necessary because the Docker image will be run on a Linux kernel.
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Start a new stage from scratch
FROM scratch

ENV PORT_NUMBER=8080
# Copy the binary from the builder stage
COPY --from=builder /app/main /main
# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["/main"]
