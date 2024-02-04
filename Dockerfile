FROM golang:1.21 as builder
COPY go.mod go.sum /go/src/lamoda_task/
WORKDIR /go/src/lamoda_task
RUN go mod download
COPY . /go/src/lamoda_task
RUN CGO_ENABLED=0 GOOS=linux go build -a -o build/lamoda_task lamoda_task

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/lamoda_task/build/lamoda_task /usr/bin/lamoda_task
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/lamoda_task"]