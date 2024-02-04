FROM alpine

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.18.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

WORKDIR /migrations

ADD *.sql /migrations/
ADD entrypoint.sh /migrations/

ENTRYPOINT ["/migrations/entrypoint.sh"]