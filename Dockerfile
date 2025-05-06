# Build stage
FROM golang:latest AS builder

WORKDIR /build
COPY . . 

# Ensure go.mod exists
RUN go mod tidy

# Build for Linux with a static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goaira

# Final minimal image
FROM scratch

# Copy CA certificates (needed for HTTPS)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the built binary from builder stage
COPY --from=builder /build/goaira /goaira

# Copy everything inside resouces folder to docker
COPY --from=builder /build/resources/ /resources/

# Set the binary to run
CMD ["/goaira"]
