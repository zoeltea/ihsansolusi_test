FROM golang:1.21-bullseye AS builder
WORKDIR /app

# Copy only module files first (for layer caching)
COPY go.mod go.sum ./


RUN apt-get update
RUN apt-get install -y curl
RUN curl -v https://proxy.golang.org
RUN go mod download

# Copy build
COPY . .
RUN go build -v -x -o service-account .

# Final stage (Alpine for minimal size)
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/service-account .
ENV APP_PORT="8080" \
    DB_HOST="localhost" \
    DB_PORT="5432" \
    DB_USER="postgres" \
    DB_PASSWORD="root" \
    DB_NAME="postgres" \
    DB_SSLMODE="disable" \
    LOG_LEVEL="info"

CMD ["./service-account"]