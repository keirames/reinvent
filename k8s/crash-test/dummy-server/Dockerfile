FROM golang:1.22.0-alpine3.19 AS BuildStage

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE 8000

RUN go build -o /main main.go

FROM alpine:latest

WORKDIR /

COPY --from=BuildStage main .

EXPOSE 8000

ENTRYPOINT ["/main"]