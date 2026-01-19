# --- Build Stage ---
FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY . .
# Build the Go application
ARG TARGET
RUN CGO_ENABLED=0 go build -o /main ./${TARGET}-src/main.go

# --- Run Stage ---
FROM nicolaka/netshoot:latest
# Copy the built binary from the builder stage
COPY --from=builder /main /usr/local/bin/app-node
RUN chmod +x /usr/local/bin/app-node


CMD ["app-node"]