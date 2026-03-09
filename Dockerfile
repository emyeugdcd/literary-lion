# Start from the official Golang image for building
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git and build-base for cgo (gcc, musl-dev, etc.)
RUN apk add --no-cache git build-base

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Enable CGO for go-sqlite3
ENV CGO_ENABLED=1

# Build the Go app (CGO enabled for sqlite3)
RUN go build -o app .

# Final image
FROM alpine:latest

WORKDIR /app

# Install sqlite3 dependencies for runtime
RUN apk add --no-cache libstdc++ sqlite-libs

# Copy the built binary and necessary folders
COPY --from=builder /app/app .
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/images ./images
COPY --from=builder /app/database ./database

EXPOSE 8080

CMD ["./app"]