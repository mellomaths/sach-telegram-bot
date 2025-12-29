# Stage 1: Build the application
FROM golang:1.24-bookworm AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to leverage Docker's cache
COPY go.mod go.sum ./

# Download application dependencies
# This step is only re-run if go.mod or go.sum changes
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go binary
# CGO_ENABLED=0 disables CGO, making the binary static and portable
RUN CGO_ENABLED=0 go build -o sach-telegram-bot ./cmd

# Stage 2: Create a minimal runtime image
# Use a minimal base image like 'scratch' for the smallest possible image
FROM scratch

# Copy the built binary from the 'builder' stage
COPY --from=builder /sach-telegram-bot /sach-telegram-bot

# (Optional) Copy necessary SSL certificates if your app makes HTTPS calls
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose the port your application listens on at runtime
EXPOSE 8080

# Define the command to run the application when the container starts
CMD ["/sach-telegram-bot"]
