#FROM golang:1.22.2-alpine3.19 AS build-stage
#
#WORKDIR /app
#
#COPY . .
#
#RUN go mod download
#
#RUN go build -o /main main.go
#
#FROM alpine:latest
#
#WORKDIR /
#
#COPY --from=build-stage main .
#
#EXPOSE 6969
#
#ENTRYPOINT ["/main"]

FROM golang:1.22.2-alpine3.19 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY server.go ./

RUN go build -o /my-app server.go

#FROM gcr.io/distroless/base-debian11
FROM alpine:latest

COPY --from=build /my-app /my-app

ENTRYPOINT ["/my-app"]