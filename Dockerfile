FROM golang:1.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# # Copy the source code
COPY . .

# Build the application
RUN make build
# CMD ["./bin/inst2any"]

# # Stage 2: Build the minimal docker image
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/bin/inst2any .
COPY --from=builder /app/configs/ ./configs/
COPY --from=builder /app/media.db .
COPY --from=builder /app/media-main.db .
COPY --from=builder /app/.creds .creds

# Expose port for the application
EXPOSE 8080

# Command to run the executable
CMD ["./inst2any", "-config", "configs/config.prod.yaml"]
