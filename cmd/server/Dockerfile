# Build stage
FROM --platform=$BUILDPLATFORM golang:1.24.0 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Copy the source code
COPY . .

# Build the application
ARG TARGETARCH TARGETOS
ENV GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 OS=$TARGETOS
RUN make server-build

# Final stage
FROM --platform=$TARGETPLATFORM scratch

WORKDIR /app

# Copy timezone data from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# Copy SSL certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Copy the binary from builder
COPY --from=builder /app/build/server .

# Run the server
ENTRYPOINT ["./server"]
