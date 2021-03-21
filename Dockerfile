FROM golang:1.16-alpine as builder

COPY . /src
WORKDIR /src

ENV CGO_ENABLED 0

RUN go build -v -o /s3-proxy ./cmd/main.go



FROM alpine:3.13

EXPOSE 9080

WORKDIR /srv

COPY --from=builder /s3-proxy /srv/s3-proxy

ENTRYPOINT ["/srv/s3-proxy", "--host=:9080"]

