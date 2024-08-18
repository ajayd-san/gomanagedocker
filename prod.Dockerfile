FROM golang:1-bullseye as builder

ENV TERM=xterm-256color

RUN apt-get update && apt-get install -y libx11-dev

RUN go install github.com/ajayd-san/gomanagedocker@HEAD

FROM debian:bullseye-slim

ENV TERM=xterm-256color

RUN apt-get update && apt-get install -y libx11-dev

COPY --from=builder /go/bin/gomanagedocker /app/gmd

CMD ["/app/gmd"]
