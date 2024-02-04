FROM golang:1.21-alpine AS builder

WORKDIR /usr/local/go/src/

ADD . /usr/local/go/src/

RUN go clean --modcache
RUN go build -mod=readonly -o app cmd/main.go

FROM alpine

COPY --from=builder /usr/local/go/src/app /

CMD ["/app"]