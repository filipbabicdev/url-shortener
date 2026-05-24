# ---- Build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy module files first so the dependency layer is cached
# when only source code changes (faster rebuilds).
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source.
COPY . .

# Build a static binary:
# - CGO_ENABLED=0      -> no C deps, fully static (works because pgx is pure Go)
# - -ldflags "-s -w"   -> strip debug info, smaller binary
# - GOOS=linux         -> target the alpine runtime regardless of build host
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/server \
    ./cmd/server

# ---- Run stage ----
FROM alpine:3.20

WORKDIR /app

# ca-certificates: needed for TLS connections to Neon (sslmode=require).
# Without it, the prod DB connection fails with x509 errors.
RUN apk add --no-cache ca-certificates

# Run as a non-root user (security best practice).
RUN adduser -D -u 10001 appuser
USER appuser

COPY --from=builder /app/server /app/server

EXPOSE 8090

ENTRYPOINT ["/app/server"]