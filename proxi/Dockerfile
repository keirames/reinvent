# Build the application from source
FROM golang:1.22.2 AS build-stage

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-app

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /docker-app /docker-app

EXPOSE 3210

USER nonroot:nonroot

ENTRYPOINT ["/docker-app"]