# Start from golang base image
FROM golang:alpine as builder

# ENV GO111MODULE=on

# Add Maintainer info
LABEL maintainer="Muhammad Suryono <msuryono0@gmail.com>"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
WORKDIR /app/src/main
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata fontconfig

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Copy semua font ke dalam folder font Alpine
COPY fonts/*.ttf /usr/share/fonts/truetype/

# Refresh font cache supaya dikenali
RUN fc-cache -f -v

# Atur timezone default (optional, tapi kadang diperlukan juga)
ENV TZ=Asia/Jakarta

# Expose port 8080 to the outside world
EXPOSE 8080

#Command to run the executable
CMD ["./main"]