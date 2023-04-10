# First stage: build the Go application
FROM golang:1.20 AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source files
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Second stage: create a minimal image
FROM alpine:latest

LABEL component="api-compiler"

# Add ca-certificates to trust SSL certificates
# Add AVR compilation dependencies
RUN apk --no-cache add ca-certificates gcc-avr avr-libc binutils-avr

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8080

# Run the Go application
CMD ["./main"]
