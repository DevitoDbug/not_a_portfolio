# Multi-stage build for Go Web App
FROM golang:1.24-alpine AS builder

# Install necessary packages for building (including Node.js for Tailwind)
RUN apk add --no-cache git ca-certificates tzdata nodejs npm
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy package.json files for npm
COPY package.json package-lock.json ./

# Install npm dependencies
RUN npm ci

# Copy source code
COPY . .

# Generate templ files
RUN templ generate

# Build Tailwind CSS
RUN npx @tailwindcss/cli -i ./internals/web/static/css/input.css -o ./internals/web/static/css/style.css --minify

# Build the web application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o app ./cmd/web

# Final stage
FROM alpine:latest AS runner

RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user
RUN addgroup -g 1001 -S golang && \
    adduser -S golang -u 1001

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/app .

# Copy internals (templates, static assets)
COPY --from=builder /app/internals ./internals

# Change ownership to non-root user
RUN chown -R golang:golang /app

# Switch to non-root user
USER golang

# Expose port 8080
EXPOSE 8080

# Run the web app binary
CMD ["./app"]
