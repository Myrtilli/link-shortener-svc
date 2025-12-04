FROM golang:1.20-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/Myrtilli/link-shortener-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/link-shortener-svc /go/src/github.com/Myrtilli/link-shortener-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/link-shortener-svc /usr/local/bin/link-shortener-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["link-shortener-svc"]
