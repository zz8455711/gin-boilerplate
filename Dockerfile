# Start from golang base image
FROM golang:alpine as builder

# Install git.
RUN apk update && apk add --no-cache git tzdata gcc musl-dev

# Working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod tidy

# Copy everythings
COPY . .

# Build the Go app
RUN CGO_ENABLED=1 GOOS=linux go build -mod=readonly -v -o main .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates
# 复制时区数据到镜像中
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai

# 设置时区环境变量为Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. Also copy config yml file
COPY --from=builder /app/main .
COPY --from=builder /app/.env .env

# Expose port 8080 to the outside world
EXPOSE 8000

#Command to run the executable
CMD ["./main"]