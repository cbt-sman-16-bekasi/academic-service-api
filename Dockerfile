# build stage
FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go build -o app .

# final image
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["./app"]
