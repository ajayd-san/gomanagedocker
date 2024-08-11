# this dockerfile is to build the image the release version of gmd is compiled on.
# I do this cuz I want to build against an older version of glibc so that gmd can be run on older os installs as well.

FROM alpine

WORKDIR /test

COPY . .

RUN echo "alpine"
# FROM golang:1-bullseye

# RUN apt-get update && apt-get install -y libx11-dev
