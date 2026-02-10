# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 go build -a -installsuffix cgo \
    -ldflags "-w -s \
    -X github.com/fgiudici/headertrace/cmd.version=${VERSION:-v0.0.1} \
    -X github.com/fgiudici/headertrace/cmd.gitCommit=${COMMIT}" \
    -o headertrace .

# Runtime stage
FROM scratch

# Copy binary from builder
COPY --from=builder /app/headertrace /headertrace

# Expose default port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/headertrace"]
