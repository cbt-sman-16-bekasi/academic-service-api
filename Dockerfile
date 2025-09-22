# Start from golang base image
FROM golang:1.22-alpine AS builder

LABEL maintainer="Muhammad Suryono <msuryono0@gmail.com>"

# Install build dependencies
RUN apk add --no-cache git build-base

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies (cache layer)
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the Go app
WORKDIR /app/src/main
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main .

# --------- runtime stage ----------
FROM alpine:3.19

# Install runtime dependencies only
RUN apk add --no-cache ca-certificates tzdata fontconfig

WORKDIR /root/

# Copy the binary and env file
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Copy fonts (optional)
COPY fonts/*.ttf /usr/share/fonts/truetype/
RUN fc-cache -f -v

# Set timezone
ENV TZ=Asia/Jakarta

EXPOSE 8080

CMD
