FROM golang:1.24-bookworm AS build

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go build -o music-utils main.go

FROM debian:12-slim

COPY --from=build /app/music-utils /usr/local/bin/music-utils

RUN apt update && apt install -y ca-certificates
RUN update-ca-certificates

ENTRYPOINT ["/usr/local/bin/music-utils"]
