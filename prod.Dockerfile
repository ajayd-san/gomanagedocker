# dockerfile for building an image to run the executable via container,
# bypassing local setup entirely.  

FROM golang:1-bullseye as builder

ENV TERM=xterm-256color

RUN apt-get update && apt-get install -y libbtrfs-dev libgpgme-dev libx11-dev

RUN go install github.com/ajayd-san/gomanagedocker@main

FROM debian:bullseye-slim

ENV TERM=xterm-256color

RUN apt-get update && apt-get install -y libbtrfs-dev libgpgme-dev libx11-dev

COPY --from=builder /go/bin/gomanagedocker /app/gmd

ENTRYPOINT ["/app/gmd"]
