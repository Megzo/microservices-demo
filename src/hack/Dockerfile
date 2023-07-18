FROM golang:1.20.6-alpine3.18 AS builder
ADD . /src
RUN cd /src && go build -o web-tty-connector

FROM tsl0922/ttyd:1.7.3-alpine as ttyd

FROM alpine:3.18

RUN apk add --no-cache bash

WORKDIR /app

# copy ttyd binary
COPY --from=ttyd /usr/bin/ttyd /usr/bin/ttyd

# copy the Go binary
COPY --from=builder /src/web-tty-connector /app/

EXPOSE 8080
EXPOSE 9000-9999

ENTRYPOINT [ "./web-tty-connector" ]
