# Build container
FROM golang:1.24 AS builder

RUN go version

RUN apt-get update && apt-get upgrade -y && apt-get install -y ca-certificates git zlib1g-dev

COPY . /go/src/github.com/SoundBoardBot/server-counter
WORKDIR /go/src/github.com/SoundBoardBot/server-counter

RUN set -Eeux && \
    go mod download && \
    go mod verify

RUN GOOS=linux GOARCH=amd64 \
    go build \
    -tags=jsoniter \
    -trimpath \
    -o main main.go

# Prod container
FROM ubuntu:latest

RUN apt-get update && apt-get upgrade -y && apt-get install -y ca-certificates curl

COPY --from=builder /go/src/github.com/SoundBoardBot/server-counter/main /srv/main

RUN chmod +x /srv/main

RUN useradd -m container
USER container
WORKDIR /srv/

CMD ["/srv/main"]