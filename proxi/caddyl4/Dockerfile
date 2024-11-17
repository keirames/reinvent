# Build the application from source
FROM golang:1.22.2 AS builder

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

RUN go get main

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-app

# Deploy the application binary into a lean image
FROM caddy:latest

COPY --from=builder /docker-app /usr/bin/caddy
