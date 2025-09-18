# syntax=docker/dockerfile:1
FROM golang:1.18-alpine as builder

ENV GO111MODULE=on

WORKDIR /storage
COPY .. .
RUN apk --no-cache add git alpine-sdk build-base gcc
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o storage main.go

FROM alpine:3.15.4

RUN apk --no-cache add ca-certificates

WORKDIR /root
COPY --from=builder /storage/storage .
EXPOSE 8080

CMD ["./storage"]