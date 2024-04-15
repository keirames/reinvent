FROM golang:1.22.2-alpine3.19 AS BuildStage

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /build main.go

# Deploy Stage

FROM alpine:latest

WORKDIR /app

COPY --from=BuildStage /build /build

ENTRYPOINT [ "/build" ]