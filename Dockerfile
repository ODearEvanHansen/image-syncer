FROM golang:1.19-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./

# Copy the source code
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/

# Build the application
RUN go build -o image-syncer ./cmd/image-syncer

# Use a smaller base image for the final image
FROM alpine:3.17

# Install Docker client
RUN apk add --no-cache docker-cli

# Copy the binary from the builder stage
COPY --from=builder /app/image-syncer /usr/local/bin/image-syncer

ENTRYPOINT ["image-syncer"]