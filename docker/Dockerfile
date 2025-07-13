# Use the official Golang image as the base image
FROM golang:1.23.2-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o mqqt_go .

# Start a new, smaller image for the final executable
FROM alpine:latest

WORKDIR /root/

# Copy the compiled application from the builder stage
COPY --from=builder /app/mqqt_go .
COPY --from=builder /app/docker/.env .
COPY --from=builder /app/latest_meter_data.txt .

# Expose the port your application listens on (if any, assuming default for now)
# EXPOSE 8080

# Command to run the executable
CMD ["./mqqt_go"]
