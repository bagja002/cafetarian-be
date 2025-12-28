# Stage 1: Build binary
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Install git/bash kalo perlu dependency
RUN apk add --no-cache git bash

# Copy mod files & download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build binary statik
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
# Stage 2: Production
FROM alpine:3.18

WORKDIR /app

# Install tzdata supaya Go ngerti Asia/Jakarta
RUN apk add --no-cache tzdata

# Set timezone (opsional, tapi kadang perlu)
ENV TZ=Asia/Jakarta

# Copy semua dari project root
COPY . .

# Copy binary dari builder
COPY --from=builder /app/main .

EXPOSE 6402

CMD ["./main"]
